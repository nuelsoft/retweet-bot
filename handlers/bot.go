package handlers

import (
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	oauth "github.com/dghubble/gologin/oauth1"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Bot() http.Handler {

	handler := func(response http.ResponseWriter, request *http.Request) {
		var (
			ConsumerKey    = os.Getenv("ConsumerKey")
			ConsumerSecret = os.Getenv("ConsumerSecret")
			AccessToken    = ""
			AccessSecret   = ""
			QueryString = os.Getenv("QueryString")
			Count = 100000000

		)

		if len(ConsumerKey) == 0 || len(ConsumerSecret) == 0 {
			log.Fatal("couldn't retrieve consumerKey and/or consumerSecret in application environment")
		}

		ctx := request.Context()
		AccessToken, AccessSecret, _ = oauth.AccessTokenFromContext(ctx)
		v := url.Values{}
		v.Set("count", string(Count))

		//v.Set("since_id", utils.BenchID

		api := anaconda.NewTwitterApiWithCredentials(AccessToken, AccessSecret, ConsumerKey, ConsumerSecret)
		search, _ := api.GetSearch(QueryString, v)

		for i := 0; i < len(search.Statuses); i++ {
			if _, retweetError := api.Retweet(search.Statuses[i].Id, false); retweetError != nil {
				log.Printf("couldn't retweet, because %s", retweetError.Error())
			}
			fmt.Printf("tweet id => %s \n tweet => %s", search.Statuses[i].IdStr, search.Statuses[i].Text)
		}
		if _, err := response.Write([]byte("Hi. Job Successful! Retweeted all.")); err != nil {
			fmt.Print("Couldn't write to client")
		}

	}
	return http.HandlerFunc(handler)
}
