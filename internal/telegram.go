package internal

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"path"
	"strings"
	"time"
)

var (
	TgBotApiToken string
	TgChatid      int64
	Bot           *tgbotapi.BotAPI
)

func SendSeparatelyMessage(author, content, scheme string, pics ...string) {
	var messageInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗点击查看原微博", scheme),
		),
	)
	var msg = tgbotapi.NewMessage(TgChatid, fmt.Sprintf("博主：[%s]\n%s", author, content))
	msg.ReplyMarkup = messageInlineKeyboard
	_, err := Bot.Send(msg)

	for _, x := range pics {
		photo := tgbotapi.NewPhoto(TgChatid, tgbotapi.FileURL(x))
		time.Sleep(1 * time.Second)
		Bot.Send(photo)
	}
	if err == nil || err.Error() == wrongfileType.Error() {
		Insert(content, scheme)
	}
}

func SendMediaGroupMessage(name, content, scheme string, pics ...string) {
	if pics == nil {
		SendSeparatelyMessage(name, content, scheme)
		return
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

	if len(listMediaVideo) > 9 {
		Bot.SendMediaGroup(tgbotapi.NewMediaGroup(TgChatid, listMediaVideo[len(listMediaVideo[:10]):]))
		time.Sleep(1 * time.Second)
		_, err := Bot.SendMediaGroup(tgbotapi.NewMediaGroup(TgChatid, listMediaVideo[:10]))

		if err == nil || err.Error() == wrongfileType.Error() {
			Insert(content, scheme)
		}
		return
	}

	msg := tgbotapi.NewMediaGroup(TgChatid, listMediaVideo)
	_, err := Bot.SendMediaGroup(msg)

	if err == nil || err.Error() == wrongfileType.Error() {
		Insert(content, scheme)
	}
}

func SendVideoGroupMessage(name, content, scheme string, pics ...string) {
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

	if len(listMediaVideo) > 9 {
		Bot.SendMediaGroup(tgbotapi.NewMediaGroup(TgChatid, listMediaVideo[len(listMediaVideo[:10]):]))
		time.Sleep(2 * time.Second)
		Bot.SendMediaGroup(tgbotapi.NewMediaGroup(TgChatid, listMediaVideo[:10]))
		return
	}

	_, err := Bot.SendMediaGroup(tgbotapi.NewMediaGroup(TgChatid, listMediaVideo))

	if err == nil || err.Error() == wrongfileType.Error() {
		Insert(content, scheme)
	}
}

func Filter(url string) interface{} {
	if !strings.Contains(path.Ext(url), "jpg") {
		return tgbotapi.NewInputMediaVideo(tgbotapi.FilePath(SavePics(url)))
	} else {
		return tgbotapi.NewInputMediaPhoto(tgbotapi.FileURL(url))
	}
}
