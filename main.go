package main

import (
	"devfestaba19_twitter_bot/utils"
	"fmt"
	"github.com/dghubble/go-twitter/twitter"
	oauth12 "github.com/dghubble/gologin/oauth1"
	twitter2 "github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitter3 "github.com/dghubble/oauth1/twitter"
	"log"
	"net/http"
	"time"
)

var (
	config = &oauth1.Config{ConsumerKey: utils.ConsumerKey, ConsumerSecret: utils.ConsumerSecret,
		Endpoint: twitter3.AuthorizeEndpoint,
	}
)

func success() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		utils.AccessToken, utils.AccessSecret, _ = oauth12.AccessTokenFromContext(ctx)
		token := oauth1.Token{Token: utils.AccessToken, TokenSecret: utils.AccessSecret}
		httpClient := config.Client(oauth1.NoContext, &token)

		client := twitter.NewClient(httpClient)

		if search, _, err := client.Search.Tweets(&twitter.SearchTweetParams{
			Query: utils.Query,
		}); err != nil {
			log.Printf("couldn't complete search because %s", err.Error())
		} else {
			for i := 0; i < len(search.Statuses); i++ {
				if _, _ ,rerr := client.Statuses.Retweet(search.Statuses[i].ID, nil); rerr != nil {
					log.Printf("couldn't retweet, because %s", rerr.Error())
				}
				fmt.Printf("tweet id => %s \n tweet => %s", search.Statuses[i].IDStr, search.Statuses[i].Text)
			}
			if _, werr := w.Write([]byte("Retweeted all.")); werr != nil{
				fmt.Print("Couldn't write to client")
			}
		}
	}
	return http.HandlerFunc(fn)
}

func main() {

	mux := http.NewServeMux()
	mux.Handle("/retweet-all", twitter2.LoginHandler(config, nil))
	mux.Handle("/callback", twitter2.CallbackHandler(config, success(), nil))
	mux.Handle("/", twitter2.LoginHandler(config, nil))

	srv := &http.Server{Handler: mux, Addr: ":7070", ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("couldn't start server %s", err.Error())
	}

}
