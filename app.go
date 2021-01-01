package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pawanverma1337/atlan-challenge/db"
	"github.com/pawanverma1337/atlan-challenge/routes"
)

// Start to function to run the app from main.go
func Start() {
	fmt.Println("Start")
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	routes.Run(r)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	db.Connect()
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	defer db.M.Client.Disconnect(ctx)
	db.M.Database = db.M.Client.Database("test")
	db.M.Collection = db.M.Database.Collection("files")

	log.Fatal(srv.ListenAndServe())
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
