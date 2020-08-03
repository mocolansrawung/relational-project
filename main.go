package main

import (
	"fmt"
	"log"
	"os"
	config "github.com/joho/godotenv"

)

func init() {
	err := config.Load(".env")
	if err != nil {
		log.Printf("can't load .env file")
		os.Exit(2)
	}
	cfgenv := os.Getenv("ENV")
	log.Printf("environment ENV=%s", cfgenv)
}


func main() {
	fmt.Println("Hello World")
}
