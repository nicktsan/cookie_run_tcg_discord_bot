package iBot

import (
	cardModel "discordbot/cookieruntcg_bot/CardData"

	"github.com/bwmarrin/discordgo"
)

type IBot interface {
	//Handler for reacting to discord messages with valid user input
	HandleNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate)

	//Handler for reacting to interactions from discord like menu selections.
	HandleNewInteraction(s *discordgo.Session, i *discordgo.InteractionCreate)

	//Get cards from the database
	SelectCards(query string, queryArgs []interface{}) ([]cardModel.CardData, error)

	//Create a select menu when more than one result is returned from SelectCards
	ListMultipleCards(discord *discordgo.Session, message *discordgo.MessageCreate, cardRows []cardModel.CardData)
}
