package cardRepo

import (
	"database/sql"
	"strconv"
	"strings"

	cardModel "discordbot/cookieruntcg_bot/repo/cardModel"
)

type CardRepo struct {
	Db *sql.DB
}

func NewCardRepo(db *sql.DB) *CardRepo {
	return &CardRepo{Db: db}
}

// Searches the cards database using English parameters.
func (repo *CardRepo) SelectCards(language string, card_name string, query string) ([]cardModel.CardData, error) {
	// Editting user input to allow to search for records that contain the input
	card_name_editted := "%" + strings.Trim(card_name, "%") + "%"
	// Use parameters to prevent SQL injection attacks
	// var query string
	// query = "SELECT * FROM cards WHERE UPPER(name_eng) LIKE UPPER($1)"
	// if language == "!fetchKR" {
	// 	query = "SELECT * FROM cards WHERE name LIKE $1"
	// }
	// fmt.Println("query: " + query)
	// rows, err := conn.Query(query, card_name_editted)
	rows, err := repo.Db.Query(query, card_name_editted)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	// A CardData slice to hold data from returned rows.
	var cardDatas []cardModel.CardData
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var card cardModel.CardData
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

func (repo *CardRepo) FormatCardData(cardRow cardModel.CardData) string {
	botMessage := cardRow.Image_link + "\n" +
		"id: " + strconv.Itoa(cardRow.Id) + "\n" +
		"name: " + cardRow.Name + "\n" +
		"name_eng: " + cardRow.Name_eng + "\n" +
		"code: " + cardRow.Code + "\n" +
		"rarity: " + cardRow.Rarity + "\n" +
		"rarity_abb: " + cardRow.Rarity_abb + "\n" +
		"card_type: " + cardRow.Card_type + "\n" +
		"color: " + cardRow.Color + "\n" +
		"color_sub: " + cardRow.Color_sub + "\n" +
		"level: " + strconv.FormatInt(int64(cardRow.Level.Int16), 10) + "\n" +
		"plain_string_eng: " + cardRow.Plain_string_eng + "\n" +
		"plain_string: " + cardRow.Plain_string + "\n" +
		"expansion: " + cardRow.Expansion.String + "\n" +
		"illustrator: " + cardRow.Illustrator + "\n" +
		"link: " + cardRow.Link + "\n"
	return botMessage
}
