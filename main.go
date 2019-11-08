package main

import (
	"devfestaba19_twitter_bot/utils"
	"fmt"
	"github.com/ChimeraCoder/anaconda"
	oauth12 "github.com/dghubble/gologin/oauth1"
	twitter2 "github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitter3 "github.com/dghubble/oauth1/twitter"
	"log"
	"net/http"
	"net/url"
	"os"
)
const defaultPort = "9090"

var (
	config = &oauth1.Config{ConsumerKey: utils.ConsumerKey, ConsumerSecret: utils.ConsumerSecret,
		Endpoint: twitter3.AuthorizeEndpoint,
	}
)

func success() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {

		ctx := req.Context()
		utils.AccessToken, utils.AccessSecret, _ = oauth12.AccessTokenFromContext(ctx)
		v := url.Values{}
		v.Set("count", utils.Count)
		v.Set("since_id", utils.BenchID)
		api := anaconda.NewTwitterApiWithCredentials(utils.AccessToken, utils.AccessSecret, utils.ConsumerKey, utils.ConsumerSecret)
		search, _ := api.GetSearch(utils.Query, v)

			for i := 0; i < len(search.Statuses); i++ {
				if _, rerr := api.Retweet(search.Statuses[i].Id, false); rerr != nil {
					log.Printf("couldn't retweet, because %s", rerr.Error())
				}
				fmt.Printf("tweet id => %s \n tweet => %s", search.Statuses[i].IdStr, search.Statuses[i].Text)
			}
			if _, werr := w.Write([]byte("Retweeted all.")); werr != nil{
				fmt.Print("Couldn't write to client")
			}
		}

	return http.HandlerFunc(fn)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	//
	//utils.ConsumerKey = os.Getenv("ConsumerKey")
	//utils.ConsumerSecret = os.Getenv("ConsumerSecret")

	mux := http.NewServeMux()
	mux.Handle("/retweet-all", twitter2.LoginHandler(config, nil))
	mux.Handle("/callback", twitter2.CallbackHandler(config, success(), nil))
	mux.Handle("/", twitter2.LoginHandler(config, nil))
	//
	//srv := &http.Server{Handler: mux, Addr: ":7070", ReadTimeout: 20 * time.Second, WriteTimeout: 20 * time.Second}

	log.Fatal(http.ListenAndServe(":"+port, mux))

	//if err := srv.ListenAndServe(); err != nil {
	//	log.Fatalf("couldn't start server %s", err.Error())
	//}
	//
}
