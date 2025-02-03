package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/SobolevTim/finance_bot/internal/database"
)

func (b *Bot) EvryDayMessage(service *database.Service) {
	for {
		users, err := service.GetUsersWitchNotify()
		if err != nil {
			log.Printf("ERROR: %v", err)
		}
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 20, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}
		time.Sleep(next.Sub(now))
		for _, user := range users {
			message := fmt.Sprintf("%s, не забудь внести траты за сегодняшний день!", user.Username)
			go b.sendMessage(user.TelegramID, message)
		}
	}
}
