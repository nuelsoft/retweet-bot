package main

import (
	twitter "github.com/dghubble/gologin/twitter"
	oauth "github.com/dghubble/oauth1"
	twitterAuthorization "github.com/dghubble/oauth1/twitter"
	"log"
	"net/http"
	"os"
	handlers "retweet-bot/handlers"
)

const defaultPort = "9090"

func main() {

	//err := godotenv.Load()

	//if err != nil {
	//	log.Fatal(err)
	//}

	var (
		ConsumerKey    = os.Getenv("ConsumerKey")
		ConsumerSecret = os.Getenv("ConsumerSecret")
	)

	if len(ConsumerKey) == 0 || len(ConsumerSecret) == 0 {
		log.Fatal("couldn't retrieve consumerKey and/or consumerSecret in application environment")
	}

	config := &oauth.Config{ConsumerKey: ConsumerKey, ConsumerSecret: ConsumerSecret,
		Endpoint: twitterAuthorization.AuthorizeEndpoint,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()
	mux.Handle("/retweet", twitter.CallbackHandler(config, handlers.Bot(), nil))
	mux.Handle("/with", handlers.ModEnv(twitter.LoginHandler(config, nil)))
	mux.Handle("/", twitter.LoginHandler(config, nil))

	log.Printf("Server listening on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
