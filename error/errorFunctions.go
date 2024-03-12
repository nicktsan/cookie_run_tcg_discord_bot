package errorFunctions

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func CheckNilErrPanic(customMessage string, e error) {
	if e != nil {
		log.Println(customMessage)
		log.Panic(e)
	}
}

func CheckNilErrPrintln(customMessage string, e error) {
	if e != nil {
		log.Println(customMessage)
		log.Println(e)
	}
}

func CheckNilErrChannelMessageSend(customMessage string, e error, discord *discordgo.Session, ChannelID string) {
	if e != nil {
		log.Println(e)
		discord.ChannelMessageSend(ChannelID, "An error occured while attempting to scan database rows.")
	}
}
