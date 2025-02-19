package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	port := os.Getenv("SERVER_PORT")
	log.Print("Server starting on port " + port)

	http.HandleFunc("/cotacao", QuoteHandler)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
