package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"
	"github.com/lucabecci/golang-blog-psql.git/internal/server"
)

func main() {
	port := os.Getenv("PORT") //env var
	srv, err := server.New(port)
	if err != nil {
		log.Fatal(err)
	}

	//start srv
	go srv.Start()

	//If the developer use the Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	//Shutdown
	srv.Close()
}
