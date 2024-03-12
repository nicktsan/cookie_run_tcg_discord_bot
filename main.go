package main

import (

	// ibot "discordbot/cookieruntcg_bot/Bot/interface"
	controller "discordbot/cookieruntcg_bot/Bot/controller"
	db "discordbot/cookieruntcg_bot/Bot/database"
	botL "discordbot/cookieruntcg_bot/Bot/logic"
	"os"
)

func main() {
	//assign environment variable BOT_TOKEN to botToken
	botToken := os.Getenv("BOT_TOKEN")
	connectionStr := os.Getenv("CONNECTION_STR")

	//Open a connection to a database
	dbConn := db.ConnectDB("postgres", connectionStr)
	defer dbConn.Close() // close dbConn connection after functin termination

	//Instantiate a new bot logic object
	bot := botL.NewBot(dbConn, "cardSelectMenu", "Choose a card.")

	//Instantiate a new bot controller object
	BotController := controller.NewBotController(bot, botToken)
	BotController.Run() // Tell the controller to run the bot
}
