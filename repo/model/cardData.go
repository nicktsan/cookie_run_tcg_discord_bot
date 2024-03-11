package cardModel

import "database/sql"

// type CardData struct {
// 	id               int
// 	name             string
// 	name_eng         string
// 	code             string
// 	rarity           string
// 	rarity_abb       string
// 	card_type        string
// 	color            string
// 	color_sub        string
// 	level            sql.NullInt16
// 	plain_string_eng string
// 	plain_string     string
// 	expansion        sql.NullString
// 	illustrator      string
// 	link             string
// 	image_link       string
// }
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

type CardRepository interface {
	SelectCards(language string, card_name string, query string) (CardData, error)
	FormatCardData(cardRow CardData) string
}
