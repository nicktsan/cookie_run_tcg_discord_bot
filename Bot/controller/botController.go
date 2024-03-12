package botController

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	botInterface "discordbot/cookieruntcg_bot/Bot/interface"
	errFunc "discordbot/cookieruntcg_bot/error"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

type BotController struct {
	ibot     botInterface.IBot
	BotToken string
}

func NewBotController(interF botInterface.IBot, bToken string) *BotController {
	return &BotController{
		ibot:     interF,
		BotToken: bToken,
	}
}

func (bController *BotController) Run() {
	// create a discord session
	discord, err := discordgo.New("Bot " + bController.BotToken)
	errFunc.CheckNilErrPanic("Error occured while attempting to create a new Discord session.", err)

	// add a event handler
	discord.AddHandler(bController.ibot.HandleNewMessage)
	discord.AddHandler(bController.ibot.HandleNewInteraction)

	//Send intents to discord servers
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// open session
	discordOpenErr := discord.Open()
	errFunc.CheckNilErrPanic("Error occured while attempting to open a websocket connection to Discord.", discordOpenErr)
	defer discord.Close() // close session, after function termination
	// keep bot running untill there is NO os interruption (ctrl + C)
	log.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-c
}
