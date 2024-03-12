package cardData

import (
	"database/sql"
	"strconv"
)

type CardData struct {
	Id               int
	Name             string
	Name_eng         string
	Code             string
	Rarity           string
	Rarity_abb       string
	Card_type        string
	Color            string
	Color_sub        string
	Level            sql.NullInt16
	Plain_string_eng string
	Plain_string     string
	Expansion        sql.NullString
	Illustrator      string
	Link             string
	Image_link       string
}

func CardDataToString(cardRow CardData) string {
	botMessage := cardRow.Image_link + "\n" +
		// "id: " + strconv.Itoa(cardRow.id) + "\n" +
		"Name_EN: " + cardRow.Name_eng + "\n" +
		"Name_KR: " + cardRow.Name + "\n" +
		"Code: " + cardRow.Code + "\n" +
		"Rarity: " + cardRow.Rarity + "\n" +
		// "rarity_abb: " + cardRow.rarity_abb + "\n" +
		"Card Type: " + cardRow.Card_type + "\n" +
		"Color: " + cardRow.Color + "\n" +
		// "color_sub: " + cardRow.color_sub + "\n" +
		"Level: " + strconv.FormatInt(int64(cardRow.Level.Int16), 10) + "\n" +
		"Card Text_EN: " + cardRow.Plain_string_eng + "\n" +
		"Card Text_KR:: " + cardRow.Plain_string
	// "expansion: " + cardRow.expansion.String + "\n" +
	// "illustrator: " + cardRow.illustrator + "\n" +
	// "link: " + cardRow.link + "\n"

	return botMessage
}
