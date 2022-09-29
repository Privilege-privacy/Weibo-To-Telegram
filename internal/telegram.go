package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
	"path"
	"strings"
)

func SendSeparatelyMessage(author, content, url string, pics ...string) {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TgBotApiToken"))
	if err != nil {
		log.Println("NewBotAPI ERR", err)
	}
	var messageInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗点击查看原微博", url),
		),
	)
	var msg = tgbotapi.NewMessage(viper.GetInt64("TgChatid"), fmt.Sprintf("博主：[%s]\n%s", author, content))
	msg.ReplyMarkup = messageInlineKeyboard
	_, err = bot.Send(msg)

	for _, x := range pics {
		photo := tgbotapi.NewPhoto(viper.GetInt64("TgChatid"), tgbotapi.FileURL(x))
		bot.Send(photo)
	}
}

func SendMediaGroupMessage(name, content, scheme string, pics ...string) {
	if pics == nil {
		SendSeparatelyMessage(name, content, scheme)
		return
	}

	bot, err := tgbotapi.NewBotAPI(viper.GetString("TgBotApiToken"))
	if err != nil {
		log.Println("NewBotAPI ERR", err)
	}

	var listMediaVideo []interface{}
	for i, value := range pics {
		if i == 0 {
			temp := tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(value))
			temp.Caption = fmt.Sprintf("博主: [%s]\n%s\n\n[🔗点击查看原微博](%s)", name, content, scheme)
			temp.ParseMode = "Markdown"
			listMediaVideo = append(listMediaVideo, temp)
		} else {
			listMediaVideo = append(listMediaVideo, tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(value)))
		}
	}

	msg := tgbotapi.NewMediaGroup(viper.GetInt64("TgChatid"), listMediaVideo)
	bot.SendMediaGroup(msg)
}

func SendVideoGroupMessage(name, content, scheme string, pics ...string) {
	bot, err := tgbotapi.NewBotAPI(viper.GetString("TgBotApiToken"))
	if err != nil {
		log.Println("NewBotAPI ERR", err)
	}

	var listMediaVideo []interface{}
	for i, value := range pics {
		switch x := Filter(value).(type) {
		case tgbotapi.InputMediaPhoto:
			if i == 0 {
				x.Caption = fmt.Sprintf("博主: [%s]\n%s\n\n[🔗点击查看原微博](%s)", name, content, scheme)
				x.ParseMode = "Markdown"
				listMediaVideo = append(listMediaVideo, x)
			} else {
				listMediaVideo = append(listMediaVideo, x)
			}

		case tgbotapi.InputMediaVideo:
			if i == 0 {
				x.Caption = fmt.Sprintf("博主: [%s]\n%s\n\n[🔗点击查看原微博](%s)", name, content, scheme)
				x.ParseMode = "Markdown"
				listMediaVideo = append(listMediaVideo, x)
			} else {
				listMediaVideo = append(listMediaVideo, x)
			}
		}

	}

	bot.SendMediaGroup(tgbotapi.NewMediaGroup(viper.GetInt64("TgChatid"), listMediaVideo))
}

func Filter(url string) interface{} {
	if !strings.Contains(path.Ext(url), "jpg") {
		return tgbotapi.NewInputMediaVideo(tgbotapi.FilePath(SavePics(url)))
	} else {
		return tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(url))
	}
}
