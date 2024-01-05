package types

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var Users = make(map[int]*User)

type User struct {
	TelegramUser *tgbotapi.User
}
