package main

import (
	"log"
	"os"

	"github.com/Lalo64GG/api-gateway/internal/config"
	"github.com/Lalo64GG/api-gateway/internal/server"
	"github.com/joho/godotenv"
)

func main(){
	if err := godotenv.Load(); err != nil{
		log.Printf("Warning: No .env file found, using environment variables only")
	}

	cfg := config.New()

	srv := server.New(cfg)


	log.Printf("API Gateway started: %s", os.Getenv("API_GATEWAY_URL"))
	log.Fatal(srv.Start())
}