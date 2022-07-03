package main

import (
	"os"
	"strconv"
)

func main() {
	cfg := BotConfig{}
	cfg.Token = os.Getenv("TOKEN")
	cfg.Domain = os.Getenv("DOMAIN")
	cfg.AdminId, _ = strconv.ParseInt(os.Getenv("ADMIN_ID"), 10, 64)
	cfg.YummyId, _ = strconv.ParseInt(os.Getenv("YUMMY_ID"), 10, 64)
	cfg.GroupId, _ = strconv.ParseInt(os.Getenv("GROUP_ID"), 10, 64)
	cfg.OrderHourEnd, _ = strconv.Atoi(os.Getenv("ORDER_HOUR_END"))
	cfg.AbbyyUsername = os.Getenv("ABBYY_USERNAME")
	cfg.AbbyyPassword = os.Getenv("ABBYY_PASSWORD")
	dev := os.Getenv("ENV") == "DEV"

	InitDb()
	go RunWeb(dev)
	RunBot(cfg)
}
