package bot

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var BotToken string

func checkNilErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Run() {

	// create a session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)

	//Send intents to discord servers
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// open session
	err = discord.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-c

}

func newMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

	/* prevent bot responding to its own message
	this is achived by looking into the message author id
	if message.author.id is same as bot.author.id then just return
	*/
	if message.Author.ID == discord.State.User.ID {
		return
	}

	// Remove leading and trailing whitespace from the discord message
	trimmed_string := strings.TrimSpace(message.Content)

	// The strings.Fields() function will split on all whitespace, and exclude it from the final result.
	//This is useful if you don’t care about the kind of whitespace, for example, tabs, spaces, and newlines all
	//count as whitespace.
	split_message := strings.Fields(trimmed_string)

	if trimmed_string == "!help" {
		discord.ChannelMessageSend(message.ChannelID, "Hello World😃")
	} else if trimmed_string == "!fetchEN" {
		discord.ChannelMessageSend(message.ChannelID, "Command missing card argument")
	} else if split_message[0] == "!fetchEN" && len(split_message) > 1 {
		joined_message := strings.Join(split_message[1:], " ")
		discord.ChannelMessageSend(message.ChannelID, "Fetching data for "+joined_message)
	} else if split_message[0] == "!fetchKR" && len(split_message) > 1 {
		joined_message := strings.Join(split_message[1:], " ")
		discord.ChannelMessageSend(message.ChannelID, "Fetching data for "+joined_message)
	}
}
