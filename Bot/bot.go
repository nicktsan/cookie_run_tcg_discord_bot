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

type CookieData struct {
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

func connectPostGres() (*sql.DB, error) {
	conn, err := sql.Open("postgres", ConnectionStr)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	fmt.Println("Connection successful")
	return conn, nil
}

func selectEN(name_eng string, conn *sql.DB) ([]CookieData, error) {
	query := "SELECT * FROM cards WHERE UPPER(name_eng) LIKE UPPER('%" + name_eng + "%')"
	rows, err := conn.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	// A CookieData slice to hold data from returned rows.
	var cookieDatas []CookieData
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var cookie CookieData
		if err := rows.Scan(&cookie.id, &cookie.name, &cookie.name_eng, &cookie.code, &cookie.rarity,
			&cookie.rarity_abb, &cookie.card_type, &cookie.color, &cookie.color_sub, &cookie.level, &cookie.plain_string_eng,
			&cookie.plain_string, &cookie.expansion, &cookie.illustrator, &cookie.link, &cookie.image_link); err != nil {
			return cookieDatas, err
		}
		cookieDatas = append(cookieDatas, cookie)
	}
	if err = rows.Err(); err != nil {
		return cookieDatas, err
	}
	return cookieDatas, nil
}

func displayCookieData(discord *discordgo.Session, message *discordgo.MessageCreate, cookieRow CookieData) {
	botMessage := cookieRow.image_link + "\n" +
		"id: " + strconv.Itoa(cookieRow.id) + "\n" +
		"name: " + cookieRow.name + "\n" +
		"name_eng: " + cookieRow.name_eng + "\n" +
		"code: " + cookieRow.code + "\n" +
		"rarity: " + cookieRow.rarity + "\n" +
		"rarity_abb: " + cookieRow.rarity_abb + "\n" +
		"card_type: " + cookieRow.card_type + "\n" +
		"color: " + cookieRow.color + "\n" +
		"color_sub: " + cookieRow.color_sub + "\n" +
		"level: " + strconv.Itoa(cookieRow.level) + "\n" +
		"plain_string_eng: " + cookieRow.plain_string_eng + "\n" +
		"plain_string: " + cookieRow.plain_string + "\n" +
		"expansion: " + cookieRow.expansion + "\n" +
		"illustrator: " + cookieRow.illustrator + "\n" +
		"link: " + cookieRow.link + "\n"
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
		conn, err := connectPostGres()
		if err != nil {
			discord.ChannelMessageSend(message.ChannelID, "An error occured while attempting to connect to the database")
		} else {
			cookieRows, selectErr := selectEN(joined_message, conn)
			if selectErr != nil {
				log.Fatal(selectErr)
				discord.ChannelMessageSend(message.ChannelID, "An error occured while attempting to Scan Rows")
			}
			for _, cookieRow := range cookieRows {
				displayCookieData(discord, message, cookieRow)
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
