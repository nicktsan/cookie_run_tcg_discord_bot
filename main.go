package main

import (
	bot "discordbot/cookieruntcg_bot/Bot"

	"os"
)

func main() {
	//assign environment variable BOT_TOKEN to botToken
	botToken := os.Getenv("BOT_TOKEN")
	connectionStr := os.Getenv("CONNECTION_STR")
	bot.BotToken = botToken
	bot.ConnectionStr = connectionStr
	bot.Run() // call the run function of bot/bot.go
}
