package pkg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	TgBotApiToken string
	ChatId        int64
	Bot           *tgbotapi.BotAPI
)

func sendMediaGroup(mediaGroup []interface{}, delete bool) error {
	if len(mediaGroup) > 9 {
		bot1, err1 := Bot.SendMediaGroup(tgbotapi.NewMediaGroup(ChatId, mediaGroup[len(mediaGroup)-10:]))
		bot2, err2 := Bot.SendMediaGroup(tgbotapi.NewMediaGroup(ChatId, mediaGroup[:len(mediaGroup)-10]))

		if (err1 != nil || err2 != nil) && delete {
			removeSend(append(bot1, bot2...))
		}
		if err1 != nil {
			return err1
		} else if err2 != nil {
			return err2
		}
		return nil
	}
	bot, err := Bot.SendMediaGroup(tgbotapi.NewMediaGroup(ChatId, mediaGroup))
	if err != nil && delete {
		removeSend(bot)
		return err
	}
	return nil
}

func removeSend(Message []tgbotapi.Message) {
	for _, v := range Message {
		Bot.Send(tgbotapi.NewDeleteMessage(ChatId, v.MessageID))
	}
}
