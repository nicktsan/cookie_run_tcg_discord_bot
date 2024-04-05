package bot

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	cardD "discordbot/cookieruntcg_bot/CardData"
	errFunc "discordbot/cookieruntcg_bot/error"

	"github.com/bwmarrin/discordgo"
	"github.com/huandu/go-sqlbuilder"
)

type Bot struct {
	Db                     *sql.DB
	CardSelectMenuCustomID string
	PlaceHolderText        string
}

func NewBot(dbConnection *sql.DB, customId string, placeholder string) *Bot {
	return &Bot{
		Db:                     dbConnection,
		CardSelectMenuCustomID: customId,
		PlaceHolderText:        placeholder,
	}
}

func (bot *Bot) HandleNewMessage(discord *discordgo.Session, message *discordgo.MessageCreate) {

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
		// Prevent the Discord bot from returning hundreds of cards if the user only has "cookie" as their card search parameter,
		// or if the command only has one rune (character)
		if len(split_message) == 2 && (utf8.RuneCountInString(split_message[1]) < 2 || strings.Contains("cookie", strings.ToLower(split_message[1])) || strings.Contains("쿠키", split_message[1])) {
			discord.ChannelMessageSend(message.ChannelID, "Please use more specific search parameters than just "+split_message[1]+".")
		} else {
			joined_message := strings.Join(split_message[1:], " ")
			discord.ChannelMessageSend(message.ChannelID, "Fetching data for "+joined_message+".")
			// Editting user input to allow to search for records that contain the input
			card_name_editted := "%" + strings.Trim(joined_message, "%") + "%"
			//Use an SQL builder to build the query
			sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
			sb.Select("*")
			sb.From("cards")
			sb.Where(
				sb.Or(
					sb.Like("UPPER(name_eng)", strings.ToUpper(card_name_editted)),
					sb.Like("name", card_name_editted),
				),
			)
			sql, args := sb.Build()
			// fmt.Println(sql)
			// fmt.Println(args)
			// sql := "SELECT * FROM cards WHERE UPPER(name_eng) LIKE UPPER($1) OR name LIKE $1"
			cardRows, selectErr := bot.SelectCards(sql, args)
			errFunc.CheckNilErrChannelMessageSend("An error occured while attempting to scan database rows.", selectErr, discord, message.ChannelID)

			//For each card found, make the Discord bot display its data in a message.
			if len(cardRows) == 0 {
				discord.ChannelMessageSend(message.ChannelID, "No data found for "+joined_message+".")
			} else if len(cardRows) == 1 {
				botMessage := cardD.CardDataToString(cardRows[0])
				discord.ChannelMessageSend(message.ChannelID, botMessage)
			} else {
				bot.ListMultipleCards(discord, message, cardRows)
			}
		}
	}
}
func (bot *Bot) HandleNewInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the interaction type is a component interaction.
	if i.Type == discordgo.InteractionMessageComponent {
		// Check if the component type is a select menu interaction.
		if i.MessageComponentData().CustomID == bot.CardSelectMenuCustomID {
			// Get the selected value from the interaction.
			selectedValue := i.MessageComponentData().Values[0]
			//split the selected value by its separator
			split_value := strings.Split(selectedValue, ":")
			//Use an SQL builder to build the query
			sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
			sb.Select("*")
			sb.From("cards")
			sb.Where(
				sb.Equal("UPPER(code)", strings.ToUpper(split_value[1])),
			)
			sql, args := sb.Build()
			// fmt.Println(sql)
			// fmt.Println(args)
			//sql := "SELECT * FROM cards WHERE UPPER(code) = UPPER($1)"
			result, err := bot.SelectCards(sql, args)
			errFunc.CheckNilErrChannelMessageSend("An error occured while attempting to scan database rows.", err, s, i.ChannelID)
			botMessage := cardD.CardDataToString(result[0])
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

func (bot *Bot) SelectCards(query string, queryArgs []interface{}) ([]cardD.CardData, error) {
	// fmt.Println("query: " + query)
	// fmt.Println("queryArgs: " + strings.Join(queryArgs, ", "))
	var rows *sql.Rows
	var err error

	// Create a Context with a timeout.
	queryCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err = bot.Db.QueryContext(queryCtx, query, queryArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// time.Sleep(time.Second * 6) //Comment out when not testing
	// A CardData slice to hold data from returned rows.
	var cardDatas []cardD.CardData
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var card cardD.CardData
		if err := rows.Scan(&card.Id, &card.Name, &card.Name_eng, &card.Code, &card.Rarity,
			&card.Rarity_abb, &card.Card_type, &card.Color, &card.Color_sub, &card.Level, &card.Plain_string_eng,
			&card.Plain_string, &card.Expansion, &card.Illustrator, &card.Link, &card.Image_link); err != nil {
			return cardDatas, err
		}
		cardDatas = append(cardDatas, card)
	}
	if err = rows.Err(); err != nil {
		return cardDatas, err
	}
	return cardDatas, nil
}

func (bot *Bot) ListMultipleCards(discord *discordgo.Session, message *discordgo.MessageCreate, cardRows []cardD.CardData) {
	var cardSelectMenuOptions []discordgo.SelectMenuOption
	// Ensure the Select Menu Value is in the format of [Card Name] / [Card Code] / [Color] / [Rarity] / [Card Type] / [Level]
	for _, cardRow := range cardRows {
		optionLabel := cardRow.Name_eng + " / " + cardRow.Name //+
		// " / " + cardRow.code + " / " + cardRow.Color + " / " + cardRow.Rarity + " / " +
		// 	cardRow.Card_type + " / " + strconv.FormatInt(int64(cardRow.Level.Int16), 10)
		optionValue := "code:" + cardRow.Code
		optionDescription := cardRow.Code + " / " + cardRow.Color + " / " + cardRow.Rarity + " / " +
			cardRow.Card_type + " / " + strconv.FormatInt(int64(cardRow.Level.Int16), 10)
		// fmt.Println("optionValue from ListMultipleCards: " + optionValue)
		// fmt.Println("optionLabel from ListMultipleCards: " + optionLabel)
		// fmt.Println("optionDescription from ListMultipleCards: " + optionDescription)
		var colourEmoji string
		switch strings.ToLower(cardRow.Color) {
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
	// fmt.Println("Formatting select Menu")
	selectMenu := []discordgo.MessageComponent{
		// Type: discordgo.MessageComponentTypeActionRow,
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					// Select menu, as other components, must have a customID, so we set it to this value.
					CustomID:    bot.CardSelectMenuCustomID,
					Placeholder: bot.PlaceHolderText, //"Choose a card.",
					Options:     cardSelectMenuOptions,
				},
			},
		},
	}
	// fmt.Println("Attempting to create select menu")
	discord.ChannelMessageSend(message.ChannelID, "Attempting to create select menu from multiple cards.")
	// Send the select menu in the channel where the command was received.
	_, err := discord.ChannelMessageSendComplex(message.ChannelID, &discordgo.MessageSend{
		Content: "Choose an option:",
		// Components:      []discordgo.MessageComponent{&selectMenu},
		Components: selectMenu,
		// AllowedMentions: &discordgo.MessageAllowedMentions{},
	})
	errFunc.CheckNilErrChannelMessageSend("Error creating card select menu.", err, discord, message.ChannelID)
}
