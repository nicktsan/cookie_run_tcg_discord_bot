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
		log.Panic(e)
	}
}

func Run() {
	// create a discord session
	discord, err := discordgo.New("Bot " + BotToken)
	checkNilErr(err)

	postGres, postGresErr = connectPostGres()
	checkNilErr(postGresErr)
	defer postGres.Close() // close postGres connection after functin termination
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
		log.Panic(err)
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

func listMultipleCards(discord *discordgo.Session, message *discordgo.MessageCreate, cardRows []CardData) {
	var cardSelectMenuOptions []discordgo.SelectMenuOption
	// Ensure the Select Menu Value is in the format of [Card Name] / [Card Code] / [Color] / [Rarity] / [Card Type] / [Level]
	for _, cardRow := range cardRows {
		optionLabel := cardRow.name_eng + " / " + cardRow.name //+
		// " / " + cardRow.code + " / " + cardRow.color + " / " + cardRow.rarity + " / " +
		// 	cardRow.card_type + " / " + strconv.FormatInt(int64(cardRow.level.Int16), 10)
		optionValue := cardRow.name_eng + " / " + cardRow.name +
			" / " + cardRow.code + " / " + cardRow.color + " / " + cardRow.rarity + " / " +
			cardRow.card_type + " / " + strconv.FormatInt(int64(cardRow.level.Int16), 10)
		optionDescription := cardRow.code + " / " + cardRow.color + " / " + cardRow.rarity + " / " +
			cardRow.card_type + " / " + strconv.FormatInt(int64(cardRow.level.Int16), 10)
		fmt.Println("optionValue from listMultipleCards: " + optionValue)
		fmt.Println("optionLabel from listMultipleCards: " + optionLabel)
		fmt.Println("optionDescription from listMultipleCards: " + optionDescription)
		var colourEmoji string
		switch strings.ToLower(cardRow.color) {
		case "red":
			colourEmoji = "🔴"
		case "yellow":
			colourEmoji = "🟡"
		case "blue":
			colourEmoji = "🔵"
		case "purple":
			colourEmoji = "🟣"
		case "green":
			colourEmoji = "🟢"
		}
		cardOption := discordgo.SelectMenuOption{
			Label: optionLabel,
			Value: optionValue,
			Emoji: discordgo.ComponentEmoji{
				Name: colourEmoji,
			},
			Description: optionDescription,
			Default:     false,
		}
		cardSelectMenuOptions = append(cardSelectMenuOptions, cardOption)
	}
	fmt.Println("Formatting select Menu")
	selectMenu := []discordgo.MessageComponent{
		// Type: discordgo.MessageComponentTypeActionRow,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					// Select menu, as other components, must have a customID, so we set it to this value.
					CustomID:    "select",
					Placeholder: "Choose a card.",
					Options:     cardSelectMenuOptions,
					// Options: []discordgo.SelectMenuOption{
					// 	{
					// 		Label: "Go",
					// 		// As with components, this things must have their own unique "id" to identify which is which.
					// 		// In this case such id is Value field.
					// 		Value: "go",
					// 		Emoji: discordgo.ComponentEmoji{
					// 			Name: "🦦",
					// 		},
					// 		// You can also make it a default option, but in this case we won't.
					// 		Default:     false,
					// 		Description: "Go programming language",
					// 	},
					// 	{
					// 		Label: "JS",
					// 		Value: "js",
					// 		Emoji: discordgo.ComponentEmoji{
					// 			Name: "🟨",
					// 		},
					// 		Description: "JavaScript programming language",
					// 	},
					// 	{
					// 		Label: "Python",
					// 		Value: "py",
					// 		Emoji: discordgo.ComponentEmoji{
					// 			Name: "🐍",
					// 		},
					// 		Description: "Python programming language",
					// 	},
					// },
				},
			},
		},
	}
	fmt.Println("Attempting to create select menu")
	discord.ChannelMessageSend(message.ChannelID, "Attempting to create select menu")
	// Send the select menu in the channel where the command was received.
	_, err := discord.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content: "Choose an option:",
		// Components:      []discordgo.MessageComponent{&selectMenu},
		Components: selectMenu,
		// AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	if err != nil {
		fmt.Println("Error sending message: ", err)
	}
}

// func CreateSelectMenu(s *discordgo.Session, i *discordgo.InteractionCreate) {

// 	// 	Also, if you are goning to implement this fucntion the string format for the each listup cards should be something like this down below.
// 	// Brave Cookie / ST4-011 / Blue / Common / COOKIE / 3
// 	// Which is in order of [Card Name] / [Card Code] / [Color] / [Rarity] / [Card Type] / [Level]
// 	var response *discordgo.InteractionResponse
// 	switch i.ApplicationCommandData().Options[0].Name {
// 	case "single":
// 		response = &discordgo.InteractionResponse{
// 			Type: discordgo.InteractionResponseChannelMessageWithSource,
// 			Data: &discordgo.InteractionResponseData{
// 				Content: "Now let's take a look on selects. This is single item select menu.",
// 				Flags:   discordgo.MessageFlagsEphemeral,
// 				Components: []discordgo.MessageComponent{
// 					discordgo.ActionsRow{
// 						Components: []discordgo.MessageComponent{
// 							discordgo.SelectMenu{
// 								// Select menu, as other components, must have a customID, so we set it to this value.
// 								CustomID:    "select",
// 								Placeholder: "Choose your favorite programming language 👇",
// 								Options: []discordgo.SelectMenuOption{
// 									{
// 										Label: "Go",
// 										// As with components, this things must have their own unique "id" to identify which is which.
// 										// In this case such id is Value field.
// 										Value: "go",
// 										Emoji: discordgo.ComponentEmoji{
// 											Name: "🦦",
// 										},
// 										// You can also make it a default option, but in this case we won't.
// 										Default:     false,
// 										Description: "Go programming language",
// 									},
// 									{
// 										Label: "JS",
// 										Value: "js",
// 										Emoji: discordgo.ComponentEmoji{
// 											Name: "🟨",
// 										},
// 										Description: "JavaScript programming language",
// 									},
// 									{
// 										Label: "Python",
// 										Value: "py",
// 										Emoji: discordgo.ComponentEmoji{
// 											Name: "🐍",
// 										},
// 										Description: "Python programming language",
// 									},
// 								},
// 							},
// 						},
// 					},
// 				},
// 			},
// 		}
// 	}
// 	err := s.InteractionRespond(i.Interaction, response)
// 	if err != nil {
// 		panic(err)
// 	}
// }

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
				log.Panic(selectErr)
				discord.ChannelMessageSend(message.ChannelID, "An error occured while attempting to scan database rows.")
			}
			//For each card found, make the Discord bot display its data in a message.
			if len(cardRows) == 0 {
				discord.ChannelMessageSend(message.ChannelID, "No data found for "+joined_message+".")
			} else if len(cardRows) == 1 {
				displayCardData(discord, message, cardRows[0])
			} else {
				// for _, cardRow := range cardRows {
				// 	displayCardData(discord, message, cardRow)
				// }
				listMultipleCards(discord, message, cardRows)
			}
		}
	}
}
