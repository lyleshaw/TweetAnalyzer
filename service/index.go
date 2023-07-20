package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/madebywelch/anthropic-go/pkg/anthropic"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"gopkg.in/antage/eventsource.v1"
)

type TweetsResponse struct {
	Modules []struct {
		Status struct {
			Data struct {
				CreatedAt string `json:"created_at"`
				Text      string `json:"full_text"`
				ID        int64  `json:"id"`
				Entities  struct {
					UserMentions []struct {
						ScreenName string `json:"screen_name"`
					} `json:"user_mentions"`
				} `json:"entities,omitempty"`
				User struct {
					Name            string `json:"name"`
					ScreenName      string `json:"screen_name"`
					Description     string `json:"description"`
					FollowersCount  int    `json:"followers_count"`
					FriendsCount    int    `json:"friends_count"`
					CreatedAt       string `json:"created_at"`
					FavouritesCount int    `json:"favourites_count"`
					StatusesCount   int    `json:"statuses_count"`
				} `json:"user"`
			} `json:"data"`
		} `json:"status"`
	} `json:"modules"`
}

type GuestTokenResponse struct {
	Response struct {
		GuestToken string `json:"guest_token"`
	}
	RateLimit       int
	Expire int64
}


var BearerToken = "Bearer AAAAAAAAAAAAAAAAAAAAAFQODgEAAAAAVHTp76lzh3rFzcHbmHVvQxYYpTw%3DckAlMINMjmCwxUcaXbAN4XqJVdgMJaHqNOFgPMK0zN1qLqLQCF"
var GuestTokenHandle *GuestTokenResponse

func CompletionWithoutSessionWithStreamByClaude(client *anthropic.Client, prompt string, callBack anthropic.StreamCallback) error {
	_, err := client.Complete(&anthropic.CompletionRequest{
		Prompt:            fmt.Sprintf("\n\nHuman: %s\n\nAssistant:", prompt),
		Model:             anthropic.ClaudeInstantV1_100k,
		MaxTokensToSample: 100000,
		StopSequences:     []string{"\r", "Human:"},
		Stream:            true,
		Temperature:       0,
	}, callBack)
	if err != nil {
		fmt.Printf("CompletionStream error: %v\n", err)
		//return err
	}
	return nil
}

func CompletionWithoutSessionWithStreamByOpenAI(ctx context.Context, client *openai.Client, prompt string) (*openai.ChatCompletionStream, error) {
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 1000,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}
	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return nil, err
	}
	return stream, nil
}

func GetClaudeClient() (*anthropic.Client, error) {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	if err := viper.BindEnv("ANTHROPIC_API_KEY"); err != nil {
		log.Fatal(err)
	}
	AnthropicApiKey := viper.GetString("ANTHROPIC_API_KEY")
	return anthropic.NewClient(AnthropicApiKey)
}
func GetOpenAIClient() (context.Context, *openai.Client, error) {
	ctx := context.Background()
	config := openai.DefaultConfig("sb-4d2ed1a575db4bd4e1df1c377252a6ae")
	config.BaseURL = "https://api.openai-sb.com/v1"
	fmt.Printf("base url: %s\n", config.BaseURL)
	return ctx, openai.NewClientWithConfig(config), nil
}
func getTweetPrompt(response TweetsResponse) (string, string) {
	prompt := ""
	maxId := ""
	for _, tweet := range response.Modules {
		var tweetPrompt string
		tweetPrompt += "这条推文的作者是" + tweet.Status.Data.User.Name + "@" + tweet.Status.Data.User.ScreenName + "。\n"
		tweetPrompt += "这条推文的内容是：{{" + tweet.Status.Data.Text + "}}\n"
		if len(tweet.Status.Data.Entities.UserMentions) > 0 {
			tweetPrompt += "这条推文中提到了："
			for _, mention := range tweet.Status.Data.Entities.UserMentions {
				tweetPrompt += "@" + mention.ScreenName + " "
			}
			tweetPrompt += "\n"
		}
		tweetPrompt += "这条推文的发布时间是" + tweet.Status.Data.CreatedAt + "。\n"
		prompt += tweetPrompt
		maxId = strconv.FormatInt(tweet.Status.Data.ID, 10)
	}
	return prompt, maxId
}

func getPromptFromTweetResponse(response TweetsResponse) string {
	prompt := "你是一个专业的心理咨询师。你的工作是从一个人所发表的推文里专业而详细的分析其性格并分点给出依据，下面是你所需要分析的推主的一些信息：\n\n"
	prompt += "这个推主的名字是" + response.Modules[0].Status.Data.User.Name + "；" +
		"ID 是" + response.Modules[0].Status.Data.User.ScreenName + "；" +
		"自我描述是{{" + response.Modules[0].Status.Data.User.Description + "}}；" +
		"拥有" + strconv.Itoa(response.Modules[0].Status.Data.User.FollowersCount) + "个粉丝；" +
		"关注了" + strconv.Itoa(response.Modules[0].Status.Data.User.FriendsCount) + "个人；" +
		"发表了" + strconv.Itoa(response.Modules[0].Status.Data.User.StatusesCount) + "条推文；" +
		"喜欢了" + strconv.Itoa(response.Modules[0].Status.Data.User.FavouritesCount) + "条推文；" +
		"注册于" + response.Modules[0].Status.Data.User.CreatedAt + "。\n"

	tweetPrompt, maxId := getTweetPrompt(response)
	prompt += tweetPrompt
	fmt.Println(maxId)
	//这里默认是前端回传了所有的需要的推文

	prompt += "\n\n请你根据以上信息，分析这个推主的性格特点，并给出你的分析依据(即所引用的推文原文)。要求写出 500 字以上的分析内容，必须从 10 点以上论述，并在最后从多个维度总结推主是什么样的人。"

	return prompt
}

func getStreamFromClaude(c *gin.Context, prompt string) {
	w := c.Writer
	completion := ""
	var callback anthropic.StreamCallback = func(resp *anthropic.CompletionResponse) error {
		completion = resp.Completion
		_, _ = fmt.Fprintf(w, "Data:\n%s\n\n", completion)
		w.Flush()
		fmt.Printf("%s\n\n", completion)
		return nil
	}
	client, _ := GetClaudeClient()
	_ = CompletionWithoutSessionWithStreamByClaude(client, prompt, callback)
}
func getStreamFromOpenAI(c *gin.Context, prompt string) {
	w := c.Writer
	ctx, clientOpenai, err := GetOpenAIClient()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	resp, err := CompletionWithoutSessionWithStreamByOpenAI(ctx, clientOpenai, prompt)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}
	finalResp := ""
	for {
		response, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}
		finalResp += response.Choices[0].Delta.Content
		_, _ = fmt.Fprintf(w, "Data:\n%s\n\n", finalResp)
		w.Flush()
	}
}

// 改为post发送数据
func getTweetAnalysis(c *gin.Context) {
	w := c.Writer

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, _ = w.(http.Flusher)
	twitterId := c.Query("twitter_id")
	count := c.Query("count")
	maxId := c.Query("max_id")
	client := c.Query("client")
	if twitterId == "" {
		c.JSON(http.StatusOK, gin.H{"error": "Wrong Params"})
		return
	}
	if client == "" {
		client = "claude"
	}
	if count == "" {
		count = "100"
	}
	fmt.Fprintf(w, "\n\n抓取数据中...本次抓取量%s条\n\n", count)
	tweetResponse, err := getTweeterTimeline(twitterId, count, maxId)
	if err != nil || len(tweetResponse.Modules) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No tweets found"})
		return
	}
	prompt := getPromptFromTweetResponse(*tweetResponse)
	fmt.Print(prompt)

	_, _ = fmt.Fprintf(w, "\n\n请求 AI 中...一分钟还没有结果请重试 orz\n\n")
	w.Flush()
	if client == "claude" {
		getStreamFromClaude(c, prompt)
	}
	if client == "openai" {
		getStreamFromOpenAI(c, prompt)
	}
}

func getGuestToken(BearerToken string) (*GuestTokenResponse, error) {
	req, err := http.NewRequest("POST", "https://api.twitter.com/1.1/guest/activate.json", nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", BearerToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var tmpGuestTokenResponse GuestTokenResponse
	if err = json.Unmarshal(body, &tmpGuestTokenResponse.Response); err != nil {
		fmt.Println(err)
		return nil, err
	}
	tmpGuestTokenResponse.RateLimit = 300
	tmpGuestTokenResponse.Expire = (time.Now()).Unix() + 900
	return &tmpGuestTokenResponse, err
}

// 查询用户timeline
func getTweeterTimeline(twitterId string, count string, maxId string) (*TweetsResponse, error) {
	url := "https://api.twitter.com/1.1/search/universal.json?count=" + count + "&modules=status&result_type=recent&pc=false&ui_lang=en-US&cards_platform=Web-13&include_entities=1&include_user_entities=1&include_cards=1&send_error_codes=1&tweet_mode=extended&include_ext_alt_text=true&include_reply_count=true"
	queryString := "&q=from%3A" + twitterId
	if maxId != "" {
		queryString += " max_id%3A" + maxId
	}
	url += queryString
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", BearerToken)

	
	if (GuestTokenHandle == nil || GuestTokenHandle.RateLimit <= 1 || ((time.Now()).Unix() - GuestTokenHandle.Expire) >= 0) {
		GuestTokenHandle, err = getGuestToken(BearerToken)
		if err != nil {
			return nil, err
		}
	}
	req.Header.Add("x-guest-token", GuestTokenHandle.Response.GuestToken)
	req.Header.Add("cookie", "gt=" + GuestTokenHandle.Response.GuestToken + ";")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	GuestTokenHandle.RateLimit--
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var tweetResponse TweetsResponse
	if err = json.Unmarshal(body, &tweetResponse); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &tweetResponse, err
}

func getTweetDetails(c *gin.Context) {
	w := c.Writer

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	_, _ = w.(http.Flusher)
	twitterId := c.Query("twitter_id")
	count := c.Query("count")

	if count == "" {
		count = "100"
	}

	tweetResponse, err := getTweeterTimeline(twitterId, count, "")
	if err != nil || len(tweetResponse.Modules) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No tweets found"})
		return
	}
	c.JSON(http.StatusOK, tweetResponse)
	return
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func main() {
	es := eventsource.New(nil, nil)
	defer es.Close()
	r := gin.Default()
	r.Use(Cors())
	r.POST("/api/get_tweet_analysis", getTweetAnalysis)
	r.GET("/api/get_tweet_details", getTweetDetails)
	r.GET("/ping", ping)
	//port := ":" + os.Getenv("PORT")
	port := ":8080"
	r.Run(port)
}
