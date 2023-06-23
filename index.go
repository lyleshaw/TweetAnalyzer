package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/madebywelch/anthropic-go/pkg/anthropic"
	"github.com/spf13/viper"
	"gopkg.in/antage/eventsource.v1"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
)

type TweetsResponse struct {
	Tweets []struct {
		CreatedAt string `json:"created_at"`
		ID        int64  `json:"id"`
		IDStr     string `json:"id_str"`
		Text      string `json:"text"`
		Truncated bool   `json:"truncated"`
		Entities  struct {
			Hashtags     []interface{} `json:"hashtags"`
			Symbols      []interface{} `json:"symbols"`
			UserMentions []struct {
				ScreenName string `json:"screen_name"`
				Name       string `json:"name"`
				ID         int64  `json:"id"`
				IDStr      string `json:"id_str"`
				Indices    []int  `json:"indices"`
			} `json:"user_mentions"`
			Urls []interface{} `json:"urls"`
		} `json:"entities,omitempty"`
		Source               string `json:"source"`
		InReplyToStatusID    int64  `json:"in_reply_to_status_id"`
		InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
		InReplyToUserID      int64  `json:"in_reply_to_user_id"`
		InReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`
		InReplyToScreenName  string `json:"in_reply_to_screen_name"`
		User                 struct {
			ID          int64       `json:"id"`
			IDStr       string      `json:"id_str"`
			Name        string      `json:"name"`
			ScreenName  string      `json:"screen_name"`
			Location    string      `json:"location"`
			Description string      `json:"description"`
			URL         interface{} `json:"url"`
			Entities    struct {
				Description struct {
					Urls []interface{} `json:"urls"`
				} `json:"description"`
			} `json:"entities"`
			Protected                      bool          `json:"protected"`
			FollowersCount                 int           `json:"followers_count"`
			FastFollowersCount             int           `json:"fast_followers_count"`
			NormalFollowersCount           int           `json:"normal_followers_count"`
			FriendsCount                   int           `json:"friends_count"`
			ListedCount                    int           `json:"listed_count"`
			CreatedAt                      string        `json:"created_at"`
			FavouritesCount                int           `json:"favourites_count"`
			UtcOffset                      interface{}   `json:"utc_offset"`
			TimeZone                       interface{}   `json:"time_zone"`
			GeoEnabled                     bool          `json:"geo_enabled"`
			Verified                       bool          `json:"verified"`
			StatusesCount                  int           `json:"statuses_count"`
			MediaCount                     int           `json:"media_count"`
			Lang                           interface{}   `json:"lang"`
			ContributorsEnabled            bool          `json:"contributors_enabled"`
			IsTranslator                   bool          `json:"is_translator"`
			IsTranslationEnabled           bool          `json:"is_translation_enabled"`
			ProfileBackgroundColor         string        `json:"profile_background_color"`
			ProfileBackgroundImageURL      interface{}   `json:"profile_background_image_url"`
			ProfileBackgroundImageURLHTTPS interface{}   `json:"profile_background_image_url_https"`
			ProfileBackgroundTile          bool          `json:"profile_background_tile"`
			ProfileImageURL                string        `json:"profile_image_url"`
			ProfileImageURLHTTPS           string        `json:"profile_image_url_https"`
			ProfileBannerURL               string        `json:"profile_banner_url"`
			ProfileLinkColor               string        `json:"profile_link_color"`
			ProfileSidebarBorderColor      string        `json:"profile_sidebar_border_color"`
			ProfileSidebarFillColor        string        `json:"profile_sidebar_fill_color"`
			ProfileTextColor               string        `json:"profile_text_color"`
			ProfileUseBackgroundImage      bool          `json:"profile_use_background_image"`
			HasExtendedProfile             bool          `json:"has_extended_profile"`
			DefaultProfile                 bool          `json:"default_profile"`
			DefaultProfileImage            bool          `json:"default_profile_image"`
			PinnedTweetIds                 []int64       `json:"pinned_tweet_ids"`
			PinnedTweetIdsStr              []string      `json:"pinned_tweet_ids_str"`
			HasCustomTimelines             bool          `json:"has_custom_timelines"`
			Following                      interface{}   `json:"following"`
			FollowRequestSent              interface{}   `json:"follow_request_sent"`
			Notifications                  interface{}   `json:"notifications"`
			AdvertiserAccountType          string        `json:"advertiser_account_type"`
			AdvertiserAccountServiceLevels []string      `json:"advertiser_account_service_levels"`
			BusinessProfileState           string        `json:"business_profile_state"`
			TranslatorType                 string        `json:"translator_type"`
			WithheldInCountries            []interface{} `json:"withheld_in_countries"`
			RequireSomeConsent             bool          `json:"require_some_consent"`
		} `json:"user"`
		Geo                  interface{} `json:"geo"`
		Coordinates          interface{} `json:"coordinates"`
		Place                interface{} `json:"place"`
		Contributors         interface{} `json:"contributors"`
		IsQuoteStatus        bool        `json:"is_quote_status"`
		RetweetCount         int         `json:"retweet_count"`
		FavoriteCount        int         `json:"favorite_count"`
		ConversationID       int64       `json:"conversation_id"`
		ConversationIDStr    string      `json:"conversation_id_str"`
		Favorited            bool        `json:"favorited"`
		Retweeted            bool        `json:"retweeted"`
		Lang                 string      `json:"lang"`
		SupplementalLanguage interface{} `json:"supplemental_language"`
		Entities0            struct {
			Hashtags     []interface{} `json:"hashtags"`
			Symbols      []interface{} `json:"symbols"`
			UserMentions []struct {
				ScreenName string `json:"screen_name"`
				Name       string `json:"name"`
				ID         int64  `json:"id"`
				IDStr      string `json:"id_str"`
				Indices    []int  `json:"indices"`
			} `json:"user_mentions"`
			Urls  []interface{} `json:"urls"`
			Media []struct {
				ID            int64  `json:"id"`
				IDStr         string `json:"id_str"`
				Indices       []int  `json:"indices"`
				MediaURL      string `json:"media_url"`
				MediaURLHTTPS string `json:"media_url_https"`
				URL           string `json:"url"`
				DisplayURL    string `json:"display_url"`
				ExpandedURL   string `json:"expanded_url"`
				Type          string `json:"type"`
				OriginalInfo  struct {
					Width      int `json:"width"`
					Height     int `json:"height"`
					FocusRects []struct {
						X int `json:"x"`
						Y int `json:"y"`
						H int `json:"h"`
						W int `json:"w"`
					} `json:"focus_rects"`
				} `json:"original_info"`
				Sizes struct {
					Small struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"small"`
					Thumb struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"thumb"`
					Large struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"large"`
					Medium struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"medium"`
				} `json:"sizes"`
			} `json:"media"`
		} `json:"entities,omitempty"`
		ExtendedEntities struct {
			Media []struct {
				ID            int64  `json:"id"`
				IDStr         string `json:"id_str"`
				Indices       []int  `json:"indices"`
				MediaURL      string `json:"media_url"`
				MediaURLHTTPS string `json:"media_url_https"`
				URL           string `json:"url"`
				DisplayURL    string `json:"display_url"`
				ExpandedURL   string `json:"expanded_url"`
				Type          string `json:"type"`
				OriginalInfo  struct {
					Width      int `json:"width"`
					Height     int `json:"height"`
					FocusRects []struct {
						X int `json:"x"`
						Y int `json:"y"`
						H int `json:"h"`
						W int `json:"w"`
					} `json:"focus_rects"`
				} `json:"original_info"`
				Sizes struct {
					Small struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"small"`
					Thumb struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"thumb"`
					Large struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"large"`
					Medium struct {
						W      int    `json:"w"`
						H      int    `json:"h"`
						Resize string `json:"resize"`
					} `json:"medium"`
				} `json:"sizes"`
				MediaKey string `json:"media_key"`
			} `json:"media"`
		} `json:"extended_entities,omitempty"`
		PossiblySensitive         bool   `json:"possibly_sensitive,omitempty"`
		PossiblySensitiveEditable bool   `json:"possibly_sensitive_editable,omitempty"`
		QuotedStatusID            int64  `json:"quoted_status_id,omitempty"`
		QuotedStatusIDStr         string `json:"quoted_status_id_str,omitempty"`
		QuotedStatus              struct {
			CreatedAt string `json:"created_at"`
			ID        int64  `json:"id"`
			IDStr     string `json:"id_str"`
			Text      string `json:"text"`
			Truncated bool   `json:"truncated"`
			Entities  struct {
				Hashtags     []interface{} `json:"hashtags"`
				Symbols      []interface{} `json:"symbols"`
				UserMentions []struct {
					ScreenName string `json:"screen_name"`
					Name       string `json:"name"`
					ID         int64  `json:"id"`
					IDStr      string `json:"id_str"`
					Indices    []int  `json:"indices"`
				} `json:"user_mentions"`
				Urls []interface{} `json:"urls"`
			} `json:"entities"`
			Source               string `json:"source"`
			InReplyToStatusID    int64  `json:"in_reply_to_status_id"`
			InReplyToStatusIDStr string `json:"in_reply_to_status_id_str"`
			InReplyToUserID      int64  `json:"in_reply_to_user_id"`
			InReplyToUserIDStr   string `json:"in_reply_to_user_id_str"`
			InReplyToScreenName  string `json:"in_reply_to_screen_name"`
			User                 struct {
				ID          int64       `json:"id"`
				IDStr       string      `json:"id_str"`
				Name        string      `json:"name"`
				ScreenName  string      `json:"screen_name"`
				Location    string      `json:"location"`
				Description string      `json:"description"`
				URL         interface{} `json:"url"`
				Entities    struct {
					Description struct {
						Urls []interface{} `json:"urls"`
					} `json:"description"`
				} `json:"entities"`
				Protected                      bool          `json:"protected"`
				FollowersCount                 int           `json:"followers_count"`
				FastFollowersCount             int           `json:"fast_followers_count"`
				NormalFollowersCount           int           `json:"normal_followers_count"`
				FriendsCount                   int           `json:"friends_count"`
				ListedCount                    int           `json:"listed_count"`
				CreatedAt                      string        `json:"created_at"`
				FavouritesCount                int           `json:"favourites_count"`
				UtcOffset                      interface{}   `json:"utc_offset"`
				TimeZone                       interface{}   `json:"time_zone"`
				GeoEnabled                     bool          `json:"geo_enabled"`
				Verified                       bool          `json:"verified"`
				StatusesCount                  int           `json:"statuses_count"`
				MediaCount                     int           `json:"media_count"`
				Lang                           interface{}   `json:"lang"`
				ContributorsEnabled            bool          `json:"contributors_enabled"`
				IsTranslator                   bool          `json:"is_translator"`
				IsTranslationEnabled           bool          `json:"is_translation_enabled"`
				ProfileBackgroundColor         string        `json:"profile_background_color"`
				ProfileBackgroundImageURL      interface{}   `json:"profile_background_image_url"`
				ProfileBackgroundImageURLHTTPS interface{}   `json:"profile_background_image_url_https"`
				ProfileBackgroundTile          bool          `json:"profile_background_tile"`
				ProfileImageURL                string        `json:"profile_image_url"`
				ProfileImageURLHTTPS           string        `json:"profile_image_url_https"`
				ProfileBannerURL               string        `json:"profile_banner_url"`
				ProfileLinkColor               string        `json:"profile_link_color"`
				ProfileSidebarBorderColor      string        `json:"profile_sidebar_border_color"`
				ProfileSidebarFillColor        string        `json:"profile_sidebar_fill_color"`
				ProfileTextColor               string        `json:"profile_text_color"`
				ProfileUseBackgroundImage      bool          `json:"profile_use_background_image"`
				HasExtendedProfile             bool          `json:"has_extended_profile"`
				DefaultProfile                 bool          `json:"default_profile"`
				DefaultProfileImage            bool          `json:"default_profile_image"`
				PinnedTweetIds                 []int64       `json:"pinned_tweet_ids"`
				PinnedTweetIdsStr              []string      `json:"pinned_tweet_ids_str"`
				HasCustomTimelines             bool          `json:"has_custom_timelines"`
				Following                      interface{}   `json:"following"`
				FollowRequestSent              interface{}   `json:"follow_request_sent"`
				Notifications                  interface{}   `json:"notifications"`
				AdvertiserAccountType          string        `json:"advertiser_account_type"`
				AdvertiserAccountServiceLevels []interface{} `json:"advertiser_account_service_levels"`
				BusinessProfileState           string        `json:"business_profile_state"`
				TranslatorType                 string        `json:"translator_type"`
				WithheldInCountries            []interface{} `json:"withheld_in_countries"`
				RequireSomeConsent             bool          `json:"require_some_consent"`
			} `json:"user"`
			Geo                  interface{} `json:"geo"`
			Coordinates          interface{} `json:"coordinates"`
			Place                interface{} `json:"place"`
			Contributors         interface{} `json:"contributors"`
			IsQuoteStatus        bool        `json:"is_quote_status"`
			RetweetCount         int         `json:"retweet_count"`
			FavoriteCount        int         `json:"favorite_count"`
			ConversationID       int64       `json:"conversation_id"`
			ConversationIDStr    string      `json:"conversation_id_str"`
			Favorited            bool        `json:"favorited"`
			Retweeted            bool        `json:"retweeted"`
			Lang                 string      `json:"lang"`
			SupplementalLanguage interface{} `json:"supplemental_language"`
		} `json:"quoted_status,omitempty"`
		SelfThread struct {
			ID    int64  `json:"id"`
			IDStr string `json:"id_str"`
		} `json:"self_thread,omitempty"`
	} `json:"tweets"`
	HasMore bool `json:"hasMore"`
}

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

func GetClaudeClient() (*anthropic.Client, error) {
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()
	if err := viper.BindEnv("ANTHROPIC_API_KEY"); err != nil {
		log.Fatal(err)
	}
	AnthropicApiKey := viper.GetString("ANTHROPIC_API_KEY")
	return anthropic.NewClient(AnthropicApiKey)
}

func sendRequestToGetTweets(twitterId string, maxId *string) *TweetsResponse {
	url := "https://twitter-virtual-scroller.vercel.app/api/user_timeline/" + twitterId
	if maxId != nil {
		url += "?max_id=" + *maxId
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	var tweetResponse TweetsResponse
	err = json.Unmarshal(body, &tweetResponse)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return &tweetResponse
}

func getTweetPrompt(response TweetsResponse) (string, string) {
	prompt := ""
	maxId := ""
	for _, tweet := range response.Tweets {
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

func getPromptFromTweetResponse(w gin.ResponseWriter, response TweetsResponse) string {
	prompt := "你是一个专业的心理咨询师。你的工作是从一个人所发表的推文里专业而详细的分析其性格并分点给出依据，下面是你所需要分析的推主的一些信息：\n\n"
	prompt += "这个推主的名字是" + response.Tweets[0].User.Name + "；" +
		"ID 是" + response.Tweets[0].User.ScreenName + "；" +
		"自我描述是{{" + response.Tweets[0].User.Description + "}}；" +
		"拥有" + strconv.Itoa(response.Tweets[0].User.FollowersCount) + "个粉丝；" +
		"关注了" + strconv.Itoa(response.Tweets[0].User.FriendsCount) + "个人；" +
		"发表了" + strconv.Itoa(response.Tweets[0].User.StatusesCount) + "条推文；" +
		"喜欢了" + strconv.Itoa(response.Tweets[0].User.FavouritesCount) + "条推文；" +
		"注册于" + response.Tweets[0].User.CreatedAt + "。\n"

	hasMore := response.HasMore
	tryTimes := 0

	for hasMore {
		tweetPrompt, maxId := getTweetPrompt(response)
		prompt += tweetPrompt
		response = *sendRequestToGetTweets(response.Tweets[0].User.ScreenName, &maxId)
		hasMore = response.HasMore
		tryTimes++
		_, _ = fmt.Fprintf(w, "第 %d 次抓取推文，本次抓取 %d 条\n\n", tryTimes, len(response.Tweets))
		w.Flush()
		fmt.Printf("第 %d 次抓取推文，本次抓取 %d 条\n\n", tryTimes, len(response.Tweets))
		if tryTimes > 10 {
			break
		}
	}

	prompt += "\n\n请你根据以上信息，分析这个推主的性格特点，并给出你的分析依据(即所引用的推文原文)。要求写出 500 字以上的分析内容，必须从 10 点以上论述，并在最后从多个维度总结推主是什么样的人。"

	return prompt
}

func getTweetAnalysis(c *gin.Context) {
	w := c.Writer

	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_, _ = w.(http.Flusher)

	// get tweet id from query string
	twitterId := c.Query("twitter_id")

	// get tweets from twitter api
	tweetResponse := sendRequestToGetTweets(twitterId, nil)

	if len(tweetResponse.Tweets) == 0 {
		c.JSON(http.StatusOK, gin.H{"error": "No tweets found"})
		return
	}

	// get prompt from tweets
	prompt := getPromptFromTweetResponse(w, *tweetResponse)

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
	port := ":" + os.Getenv("PORT")
	r.Run(port)
}
