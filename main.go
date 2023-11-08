package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/segunjkf/server/pkg/database/bolt"
	"github.com/segunjkf/server/pkg/server"
)

func main() {
	fmt.Println("Starting my server")

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	b, err := bolt.New(ctx, "./data")
	if err != nil {
		log.Fatal("Failed to start the database: %v", err)
	}

	Address := ":8080"
	r := mux.NewRouter()
	s := server.New(ctx, b)
	println(s)

	r.HandleFunc("/user/create", s.HandleCreateUsers)
	r.HandleFunc("/user/{name}", s.HandleUsers)

	r.HandleFunc("/", s.HandleFuncHome)
	srv := &http.Server{
		Addr:           Address,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Start Server: %v", Address)
	log.Fatal(srv.ListenAndServe())
}
