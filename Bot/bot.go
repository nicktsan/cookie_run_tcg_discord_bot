package bot

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

// todo: eliminate stinky global variables
var BotToken string
var ConnectionStr string

const CardSelectMenuCustomID = "cardSelectMenu"

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
	discord.AddHandler(HandleNewMessage)
	discord.AddHandler(HandleNewInteraction)
	//Send intents to discord servers
	discord.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	// open session
	err = discord.Open()
	checkNilErr(err)
	defer discord.Close() // close session, after function termination
	// keep bot running untill there is NO os interruption (ctrl + C)
	log.Println("Bot running....")
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
	log.Println("DB Connection successful.")
	return conn, nil
}

// Searches the cards database using English parameters.
func selectCards(query string, conn *sql.DB, queryArgs ...string) ([]CardData, error) {
	// log.Println("query: " + query)
	// log.Println("queryArgs: " + strings.Join(queryArgs, ", "))
	var rows *sql.Rows
	var err error

	rows, err = conn.Query(query, queryArgs[0])

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
		optionValue := "code:" + cardRow.code
		optionDescription := cardRow.code + " / " + cardRow.color + " / " + cardRow.rarity + " / " +
			cardRow.card_type + " / " + strconv.FormatInt(int64(cardRow.level.Int16), 10)
		// log.Println("optionValue from listMultipleCards: " + optionValue)
		// log.Println("optionLabel from listMultipleCards: " + optionLabel)
		// log.Println("optionDescription from listMultipleCards: " + optionDescription)
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
	// log.Println("Formatting select Menu")
	selectMenu := []discordgo.MessageComponent{
		// Type: discordgo.MessageComponentTypeActionRow,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					// Select menu, as other components, must have a customID, so we set it to this value.
					CustomID:    CardSelectMenuCustomID,
					Placeholder: "Choose a card.",
					Options:     cardSelectMenuOptions,
				},
			},
		},
	}
	// log.Println("Attempting to create select menu")
	discord.ChannelMessageSend(message.ChannelID, "Attempting to create select menu from multiple cards.")
	// Send the select menu in the channel where the command was received.
	_, err := discord.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content: "Choose an option:",
		// Components:      []discordgo.MessageComponent{&selectMenu},
		Components: selectMenu,
		// AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	if err != nil {
		log.Println("Error creating card select menu: ", err)
		discord.ChannelMessageSend(message.ChannelID, "Error occured while attempting to create select menu.")
	}
}

// Makes the Discord bot display the data it fetched via Discord message.
func CardDataToString(cardRow CardData) string {
	botMessage := cardRow.image_link + "\n" +
		// "id: " + strconv.Itoa(cardRow.id) + "\n" +
		"Name_EN: " + cardRow.name_eng + "\n" +
		"Name_KR: " + cardRow.name + "\n" +
		"Code: " + cardRow.code + "\n" +
		"Rarity: " + cardRow.rarity + "\n" +
		// "rarity_abb: " + cardRow.rarity_abb + "\n" +
		"Card Type: " + cardRow.card_type + "\n" +
		"Color: " + cardRow.color + "\n" +
		// "color_sub: " + cardRow.color_sub + "\n" +
		"Level: " + strconv.FormatInt(int64(cardRow.level.Int16), 10) + "\n" +
		"Card Text_EN: " + cardRow.plain_string_eng + "\n" +
		"Card Text_KR:: " + cardRow.plain_string
	// "expansion: " + cardRow.expansion.String + "\n" +
	// "illustrator: " + cardRow.illustrator + "\n" +
	// "link: " + cardRow.link + "\n"

	return botMessage
}

func HandleNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

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
		helpMessage := "`!fetch [card name in English or Korean]`\n" +
			"Searches for a card by its English or Korean name."
		discord.ChannelMessageSend(message.ChannelID, helpMessage)
	} else if trimmed_string == "!fetch" {
		discord.ChannelMessageSend(message.ChannelID, "Command missing card argument.")
	} else if split_message[0] == "!fetch" && len(split_message) > 1 {
		// Prevent the Discord bot from returning hundreds of cards if the user only has "cookie" as their card search parameter.
		if len(split_message) == 2 && (strings.Contains("cookie", strings.ToLower(split_message[1])) || strings.Contains("쿠키", split_message[1])) {
			discord.ChannelMessageSend(message.ChannelID, "Please use more specific search parameters than just "+split_message[1]+".")
		} else {
			joined_message := strings.Join(split_message[1:], " ")
			discord.ChannelMessageSend(message.ChannelID, "Fetching data for "+joined_message+".")
			// Editting user input to allow to search for records that contain the input
			card_name_editted := "%" + strings.Trim(joined_message, "%") + "%"
			query := "SELECT * FROM cards WHERE UPPER(name_eng) LIKE UPPER($1) OR name LIKE $1"
			cardRows, selectErr := selectCards(query, postGres, card_name_editted)
			if selectErr != nil {
				log.Println(selectErr)
				discord.ChannelMessageSend(message.ChannelID, "An error occured while attempting to scan database rows.")
			}
			//For each card found, make the Discord bot display its data in a message.
			if len(cardRows) == 0 {
				discord.ChannelMessageSend(message.ChannelID, "No data found for "+joined_message+".")
			} else if len(cardRows) == 1 {
				botMessage := CardDataToString(cardRows[0])
				discord.ChannelMessageSend(message.ChannelID, botMessage)
			} else {
				listMultipleCards(discord, message, cardRows)
			}
		}
	}
}

func HandleNewInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the interaction type is a component interaction.
	if i.Type == discordgo.InteractionMessageComponent {
		// Check if the component type is a select menu interaction.
		if i.MessageComponentData().CustomID == CardSelectMenuCustomID {
			// Get the selected value from the interaction.
			selectedValue := i.MessageComponentData().Values[0]
			//split the selected value by its separator
			split_value := strings.Split(selectedValue, ":")
			query := "SELECT * FROM cards WHERE UPPER(code) = UPPER($1)"
			result, err := selectCards(query, postGres, split_value[1])
			if err != nil {
				s.ChannelMessageSend(i.ChannelID, "An error occured while attempting to scan database rows.")
				log.Println(err)
			}
			botMessage := CardDataToString(result[0])
			// Reply to the user with the selected value.
			_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: botMessage,
				},
			})
		}
	}
}
