package main

import (
	"log"
	"net/http"
	"os"
	"ticket-tracker/api/router"
	mongodb "ticket-tracker/config"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
	port := ":" + os.Getenv("SERVER_PORT")
	r := router.New(mongodb.DB)
	s := &http.Server{
		Addr:    port,
		Handler: r,
	}
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Print(err)
	}

}
