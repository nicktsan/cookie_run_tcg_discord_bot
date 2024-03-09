package repository

import (
	"database/sql"
	"discordbot/cookieruntcg_bot/repo/cardRepo"
)

// Repositories contains all the repo structs
type Repositories struct {
	CardRepo *cardRepo.CardRepo
}

// InitRepositories should be called in main.go
func InitRepositories(db *sql.DB) *Repositories {
	cardRepo := cardRepo.NewCardRepo(db)
	return &Repositories{CardRepo: cardRepo}
}
