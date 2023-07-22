package twitter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/dghubble/oauth1"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

type authorizer struct{}

func (a *authorizer) Add(_ *http.Request) {}

type TwitterClient struct {
	Client *http.Client
	Ctx    context.Context
}

const tweetURL = "https://api.twitter.com/2/tweets"

func NewTwitterClient() *TwitterClient {
	log.SetLevel(log.DebugLevel)
	ctx := context.Background()
	//config := oauth1.NewConfig("m7BQ7z9dT5jcbONP514vI4tlj", "LuHOQlOah5NgB1QhJdMB0jEt1u86R24ww4PCpSDoZQgVjGyeZD")
	//token := oauth1.NewToken("1679372641835307011-gXqEDSKuY2HrIyJomwYDetFdJs2mg7", "WW7UtOXjhoafWfNo4aCaV9PWHSrMsyJsKYlopwRCvDmFq")

	consumerKey := viper.GetString("consumerKey")
	consumerSecret := viper.GetString("consumerSecret")
	tokenKey := viper.GetString("token")
	tokenSecret := viper.GetString("tokenSecret")

	//fmt.Println(consumerKey, consumerSecret, tokenKey, tokenSecret)
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(tokenKey, tokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	return &TwitterClient{
		Ctx:    ctx,
		Client: httpClient,
	}
}

func (t *TwitterClient) PostTwitter(text string) {

	tweetData := map[string]string{
		"text": text,
	}

	tweetPayload, err := json.Marshal(tweetData)
	if err != nil {
		fmt.Println("==", err)
	}

	// Send a POST request to the Twitter API to create the tweet
	req, err := http.NewRequest("POST", tweetURL, bytes.NewBuffer(tweetPayload))
	if err != nil {
		fmt.Printf("Error sending tweet1: %v", err)
		os.Exit(1)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		fmt.Printf("Error occured while sending tweet: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusCreated {
		errorResponse := struct {
			Errors []struct {
				Message string `json:"message"`
			} `json:"errors"`
		}{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			fmt.Printf("Error sending tweet: Twitter API responded with status code %v", resp.StatusCode)
			//os.Exit(1)
		}
		fmt.Printf("Error sending twee2t2: %v\n", errorResponse.Errors)
		//os.Exit(1)
	}
}

//func (t *twitterClient) DeleteTwitter(Id string) {
//
//	res, err := t.Client.DeleteTweet(t.Ctx, Id)
//	if err != nil {
//		log.WithError(err).Fatal("Could not publish tweet")
//	}
//	fmt.Println(res)
//}
