package handlers

import (
	"log"
	"net/http"
	"os"
)

func ModEnv(next http.Handler) http.Handler {
	handler := func(response http.ResponseWriter, request *http.Request) {
		tag := request.URL.Query().Get("tag")
		if tag != "" {
			err := os.Setenv("QueryString", tag)
			if err != nil {
				log.Fatal("Couldn't write to environment")
			}
		}
		next.ServeHTTP(response, request)
	}
	return http.HandlerFunc(handler)
}
