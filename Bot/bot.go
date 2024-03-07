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
var postGres *sql.DB
var postGresErr error

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
	level            sql.NullInt16
	plain_string_eng string
	plain_string     string
	expansion        sql.NullString
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
	postGres, postGresErr = connectPostGres()
	checkNilErr(postGresErr)

	// add a event handler
	discord.AddHandler(newMessage)

	//Send intents to discord servers
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// open session
	err = discord.Open()
	checkNilErr(err)

	defer discord.Close()  // close session, after function termination
	defer postGres.Close() // close postGres connection after functin termination

	// keep bot running untill there is NO os interruption (ctrl + C)
	fmt.Println("Bot running....")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-c
	fmt.Println("DB connection closed.")
}

// Connect to the PostGreSQL database
func connectPostGres() (*sql.DB, error) {
	conn, err := sql.Open("postgres", ConnectionStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("DB Connection successful.")
	return conn, nil
}

// Searches the cards database using English parameters.
func selectCards(language string, card_name string, conn *sql.DB) ([]CardData, error) {
	// Editting user input to allow to search for records that contain the input
	card_name_editted := "%" + strings.Trim(card_name, "%") + "%"
	// Use parameters to prevent SQL injection attacks
	var query string
	query = "SELECT * FROM cards WHERE UPPER(name_eng) LIKE UPPER($1)"
	if language == "!fetchKR" {
		query = "SELECT * FROM cards WHERE name LIKE $1"
	}
	// fmt.Println("query: " + query)
	rows, err := conn.Query(query, card_name_editted)
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
		"level: " + strconv.FormatInt(int64(cardRow.level.Int16), 10) + "\n" +
		"plain_string_eng: " + cardRow.plain_string_eng + "\n" +
		"plain_string: " + cardRow.plain_string + "\n" +
		"expansion: " + cardRow.expansion.String + "\n" +
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
		helpMessage := "`!fetchEN [card name in English]`\n" +
			"Searches for a card by its English name. \n" +
			"`!fetchKR [card name in Korean]`\n" +
			"Searches for a card by its Korean name."
		discord.ChannelMessageSend(message.ChannelID, helpMessage)
	} else if trimmed_string == "!fetchEN" || trimmed_string == "!fetchKR" {
		discord.ChannelMessageSend(message.ChannelID, "Command missing card argument.")
	} else if (split_message[0] == "!fetchEN" || split_message[0] == "!fetchKR") && len(split_message) > 1 {
		// Prevent the Discord bot from returning hundreds of cards if the user only has "cookie" as their card search parameter.
		if len(split_message) == 2 && (strings.Contains("cookie", strings.ToLower(split_message[1])) || strings.Contains("쿠키", split_message[1])) {
			discord.ChannelMessageSend(message.ChannelID, "Please use more specific search parameters than just "+split_message[1]+".")
		} else {
			joined_message := strings.Join(split_message[1:], " ")
			discord.ChannelMessageSend(message.ChannelID, "Fetching data for "+joined_message+".")

			//Fetch Card data
			cardRows, selectErr := selectCards(split_message[0], joined_message, postGres)
			if selectErr != nil {
				log.Fatal(selectErr)
				discord.ChannelMessageSend(message.ChannelID, "An error occured while attempting to scan database rows.")
			}
			//For each card found, make the Discord bot display its data in a message.
			if len(cardRows) == 0 {
				discord.ChannelMessageSend(message.ChannelID, "No data found for "+joined_message+".")
			} else {
				for _, cardRow := range cardRows {
					displayCardData(discord, message, cardRow)
				}
			}
		}
	}
}
