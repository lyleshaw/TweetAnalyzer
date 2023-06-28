package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/gin-gonic/gin"
	"github.com/madebywelch/anthropic-go/pkg/anthropic"
	"github.com/sashabaranov/go-openai"
	"github.com/spf13/viper"
	"gopkg.in/antage/eventsource.v1"
)

var twi_cli *twitter.Client = setTwitterClient()

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

func getTweetPrompt(response []twitter.Tweet) (string, string) {
	prompt := ""
	maxId := ""
	for _, tweet := range response {
		var tweetPrompt string
		tweetPrompt += "这条推文的作者是" + tweet.User.Name + "@" + tweet.User.ScreenName + "。\n"
		tweetPrompt += "这条推文的内容是：{{" + tweet.Text + "}}\n"
		if len(tweet.Entities.UserMentions) > 0 {
			tweetPrompt += "这条推文中提到了："
			for _, mention := range tweet.Entities.UserMentions {
				tweetPrompt += "@" + mention.ScreenName + " "
			}
			tweetPrompt += "\n"
		}
		tweetPrompt += "这条推文的发布时间是" + tweet.CreatedAt + "。\n"
		prompt += tweetPrompt
		maxId = strconv.FormatInt(tweet.ID, 10)
	}
	return prompt, maxId
}

func getPromptFromTweetResponse(response []twitter.Tweet) string {
	prompt := "你是一个专业的心理咨询师。你的工作是从一个人所发表的推文里专业而详细的分析其性格并分点给出依据，下面是你所需要分析的推主的一些信息：\n\n"
	prompt += "这个推主的名字是" + response[0].User.Name + "；" +
		"ID 是" + response[0].User.ScreenName + "；" +
		"自我描述是{{" + response[0].User.Description + "}}；" +
		"拥有" + strconv.Itoa(response[0].User.FollowersCount) + "个粉丝；" +
		"关注了" + strconv.Itoa(response[0].User.FriendsCount) + "个人；" +
		"发表了" + strconv.Itoa(response[0].User.StatusesCount) + "条推文；" +
		"喜欢了" + strconv.Itoa(response[0].User.FavouritesCount) + "条推文；" +
		"注册于" + response[0].User.CreatedAt + "。\n"

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

func getTweetAnalysis(c *gin.Context) {
	w := c.Writer
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, _ = w.(http.Flusher)
	twitterId := c.Query("twitter_id")
	count := c.Query("count")
	client := c.Query("client")
	keyword := c.Query("keyword")
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	if twitterId == "" {
		c.JSON(http.StatusOK, gin.H{"error": "Wrong Params"})
		return
	}
	if client == "" {
		client = "claude"
	}
	fmt.Fprintf(w, "\n\n抓取数据中...本次抓取量%s条\n\n", count)
	tweetResponse, err := getTweets(twitterId, count, keyword, startDate, endDate, twi_cli)
	if err != nil || len(tweetResponse) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No tweets found"})
		return
	}
	prompt := getPromptFromTweetResponse(tweetResponse)
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

func getTweets(twitterId string, count string, keyword string, startDate string, endDate string, client *twitter.Client) ([]twitter.Tweet, error) {
	cnt, err := strconv.Atoi(count)
	if err != nil {
		return nil, err
	}
	fetchedTweets, err := getAmountTweets(twitterId, cnt, client)
	if err != nil {
		return nil, err
	}
	return filterTweets(fetchedTweets, keyword, startDate, endDate)
}

func getAmountTweets(twitterId string, count int, client *twitter.Client) ([]twitter.Tweet, error) {
	var allTweets []twitter.Tweet
	if count > 200 {
		var maxId int64
		for count > 0 {
			var cnt int
			if count > 200 {
				cnt = 200
			} else {
				cnt = count
			}
			tweets, err := getTimeline(twitterId, cnt, maxId, client)
			if err != nil {
				return nil, err
			}
			allTweets = append(allTweets, tweets...)
			count -= cnt
		}
	} else {
		var err error
		allTweets, err = getTimeline(twitterId, count, 0, client)
		if err != nil {
			return nil, err
		}
	}
	return allTweets, nil
}

func getTimeline(twitterId string, count int, maxId int64, client *twitter.Client) ([]twitter.Tweet, error) {
	var userTimelineParams twitter.UserTimelineParams
	userTimelineParams.ScreenName = twitterId
	userTimelineParams.Count = count
	userTimelineParams.MaxID = maxId
	tweets, resp, err := client.Timelines.UserTimeline(&userTimelineParams)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("Error fetching the tweets")
		return nil, err
	}
	return tweets, err
}

func filterTweets(allTweets []twitter.Tweet, keyword string, startDate string, endDate string) ([]twitter.Tweet, error) {
	var filteredDateTweets []twitter.Tweet
	if startDate != "" && endDate != "" {
		start, err := time.Parse("2006-01-02", startDate)
		if err != nil {
			return nil, err
		}
		end, err := time.Parse("2006-01-02", endDate)
		if err != nil {
			return nil, err
		}
		for _, tweet := range allTweets {
			tweetDate, _ := tweet.CreatedAtTime()
			if tweetDate.After(start) && tweetDate.Before(end) {
				filteredDateTweets = append(filteredDateTweets, tweet)
			}
		}
		fmt.Printf("Successfully filtered time, rest tweets: %d\n", len(filteredDateTweets))
	}
	if keyword == "" {
		return filteredDateTweets, nil
	}
	var filteredKeywordTweets []twitter.Tweet
	for _, tweet := range filteredDateTweets {
		if strings.Contains(tweet.Text, keyword) {
			filteredKeywordTweets = append(filteredKeywordTweets, tweet)
		}
	}
	return filteredKeywordTweets, nil
}

func setTwitterClient() *twitter.Client {
	config := oauth1.NewConfig(os.Getenv("TW_CONSUMER_KEY"), os.Getenv("TW_CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TW_ACCESS_TOKEN"), os.Getenv("TW_ACCESS_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)
	client := twitter.NewClient(httpClient)
	return client
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

	r.GET("/api/get_tweet_analysis", getTweetAnalysis)
	r.GET("/ping", ping)
	//port := ":" + os.Getenv("PORT")
	port := ":8080"
	r.Run(port)
}
