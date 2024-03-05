package main

import (
	bot "discordbot/cookieruntcg_bot/Bot"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	//load environment variables from local.env
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatal("Error")
	}
	//assign environment variable BOT_TOKEN to botToken
	botToken := os.Getenv("BOT_TOKEN")
	bot.BotToken = botToken
	bot.Run() // call the run function of bot/bot.go
}
