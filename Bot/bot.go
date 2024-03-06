package bot

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

var BotToken string
var ConnectionStr string

type CardData struct {
	id               int
	name             string
	name_eng         string
	code             string
	rarity           string
	rarity_abb       string
	card_type        string
	color            string
	color_sub        string
	level            int
	plain_string_eng string
	plain_string     string
	expansion        string
	illustrator      string
	link             string
	image_link       string
}

func checkNilErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func Run() {
	// user=postgres.uexpudztesdujzrmclis password=[YOUR-PASSWORD] host=aws-0-ca-central-1.pooler.supabase.com port=5432 dbname=postgres
	// create a discord session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	// add a event handler
	discord.AddHandler(newMessage)

	//Send intents to discord servers
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// open session
	err = discord.Open()
	checkNilErr(err)

	defer discord.Close() // close session, after function termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-c

}

// Connect to the PostGreSQL database
func connectPostGres() (*sql.DB, error) {
	conn, err := sql.Open("postgres", ConnectionStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Connection successful")
	return conn, nil
}

// Searches the cards database using English parameters.
func selectEN(name_eng string, conn *sql.DB) ([]CardData, error) {
	// Editting user input to allow to search for records that contain the input
	name_eng_editted := "%" + strings.Trim(name_eng, "%") + "%"
	// Use parameters to prevent SQL injection attacks
	query := "SELECT * FROM cards WHERE UPPER(name_eng) LIKE UPPER($1)"
	rows, err := conn.Query(query, name_eng_editted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// A CardData slice to hold data from returned rows.
	var cardDatas []CardData
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var card CardData
		if err := rows.Scan(&card.id, &card.name, &card.name_eng, &card.code, &card.rarity,
			&card.rarity_abb, &card.card_type, &card.color, &card.color_sub, &card.level, &card.plain_string_eng,
			&card.plain_string, &card.expansion, &card.illustrator, &card.link, &card.image_link); err != nil {
			return cardDatas, err
		}
		cardDatas = append(cardDatas, card)
	}
	if err = rows.Err(); err != nil {
		return cardDatas, err
	}
	return cardDatas, nil
}

// Makes the Discord bot display the data it fetched via Discord message.
func displayCardData(discord *discordgo.Session, message *discordgo.MessageCreate, cardRow CardData) {
	botMessage := cardRow.image_link + "\n" +
		"id: " + strconv.Itoa(cardRow.id) + "\n" +
		"name: " + cardRow.name + "\n" +
		"name_eng: " + cardRow.name_eng + "\n" +
		"code: " + cardRow.code + "\n" +
		"rarity: " + cardRow.rarity + "\n" +
		"rarity_abb: " + cardRow.rarity_abb + "\n" +
		"card_type: " + cardRow.card_type + "\n" +
		"color: " + cardRow.color + "\n" +
		"color_sub: " + cardRow.color_sub + "\n" +
		"level: " + strconv.Itoa(cardRow.level) + "\n" +
		"plain_string_eng: " + cardRow.plain_string_eng + "\n" +
		"plain_string: " + cardRow.plain_string + "\n" +
		"expansion: " + cardRow.expansion + "\n" +
		"illustrator: " + cardRow.illustrator + "\n" +
		"link: " + cardRow.link + "\n"
	discord.ChannelMessageSend(message.ChannelID, botMessage)
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

		//Connect to PostGreSQL database
		conn, err := connectPostGres()
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "An error occured while attempting to connect to the database")
		} else {
			//Fetch Card data
			cardRows, selectErr := selectEN(joined_message, conn)
			if selectErr != nil {
				log.Fatal(selectErr)
				discord.ChannelMessageSend(message.ChannelID, "An error occured while attempting to Scan Rows")
			}
			//For each card found, make the Discord bot display its data in a message.
			for _, cardRow := range cardRows {
				displayCardData(discord, message, cardRow)
			}
		}
		fmt.Println("Closing connection")
		conn.Close()
		fmt.Println("DB connection closed")
	} else if split_message[0] == "!fetchKR" && len(split_message) > 1 {
		joined_message := strings.Join(split_message[1:], " ")
		discord.ChannelMessageSend(message.ChannelID, "Fetching data for "+joined_message)
		//todo implement korean searching
	}
}
