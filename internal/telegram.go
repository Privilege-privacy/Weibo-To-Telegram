package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

var (
	TgBotApiToken    string
	TgUseridORChatId int64
)

func sendPhoto(url string) {
	bot, err := tgbotapi.NewBotAPI(TgBotApiToken)
	if err != nil {
		log.Println(err)
	}
	photo := tgbotapi.NewPhoto(TgUseridORChatId, tgbotapi.FileURL(url))
	if _, err = bot.Send(photo); err != nil {
		log.Println("Send Photo Err !", err)
	}
}

func sendMessage(author, content, url string) {
	var messageInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗点击查看原微博", url),
		),
	)
	bot, err := tgbotapi.NewBotAPI(TgBotApiToken)
	if err != nil {
		log.Println(err)
	}
	var msg = tgbotapi.NewMessage(TgUseridORChatId, fmt.Sprintf("博主：[%s]\n%s", author, content))
	msg.ReplyMarkup = messageInlineKeyboard
	if _, err := bot.Send(msg); err != nil {
		log.Println("Send Message Err!", err)
	}
}
