package cardData

import (
	"database/sql"
	"strconv"
)

type CardData struct {
	// Id               int
	Name_kr  string
	Name_eng string
	Code     string
	Rarity   string
	// Rarity_abb       string
	Card_type string
	Color     string
	// Color_sub        string
	Card_level     sql.NullInt16
	Plain_text_eng string
	Plain_text     string
	// Expansion        sql.NullString
	// Illustrator      string
	// Link             string
	Image_link string
	// Name_eng_lower string
}

func CardDataToString(cardRow CardData) string {
	botMessage := cardRow.Image_link + "\n" +
		// "id: " + strconv.Itoa(cardRow.id) + "\n" +
		"Name_EN: " + cardRow.Name_eng + "\n" +
		"Name_KR: " + cardRow.Name_kr + "\n" +
		"Code: " + cardRow.Code + "\n" +
		"Rarity: " + cardRow.Rarity + "\n" +
		// "rarity_abb: " + cardRow.rarity_abb + "\n" +
		"Card Type: " + cardRow.Card_type + "\n" +
		"Color: " + cardRow.Color + "\n" +
		// "color_sub: " + cardRow.color_sub + "\n" +
		"Level: " + strconv.FormatInt(int64(cardRow.Card_level.Int16), 10) + "\n" +
		"Card Text_EN: " + cardRow.Plain_text_eng + "\n" +
		"Card Text_KR:: " + cardRow.Plain_text
	// "expansion: " + cardRow.expansion.String + "\n" +
	// "illustrator: " + cardRow.illustrator + "\n" +
	// "link: " + cardRow.link + "\n"

	return botMessage
}
