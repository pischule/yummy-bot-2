package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := BotConfig{}
	cfg.Token = os.Getenv("TOKEN")
	cfg.Domain = os.Getenv("DOMAIN")
	cfg.AdminId, _ = strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	cfg.YummyId, _ = strconv.ParseInt(os.Getenv("YUMMY_ID"), 10, 64)
	cfg.GroupId, _ = strconv.ParseInt(os.Getenv("GROUP_ID"), 10, 64)
	debug := os.Getenv("DEBUG") == "true"

	InitDb()
	go RunWeb(debug)
	RunBot(cfg)
}
