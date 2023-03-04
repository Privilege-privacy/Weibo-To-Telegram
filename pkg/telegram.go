package pkg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"path/filepath"
	"strings"
	"time"
)

var (
	TgBotApiToken string
	ChatId        int64
	Bot           *tgbotapi.BotAPI
)

func SendSeparatelyMessage(url, message string, mediaFiles []string) error {
	var messageInlineKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗点击查看原微博", url),
		),
	)
	var msg = tgbotapi.NewMessage(ChatId, message)
	msg.ReplyMarkup = messageInlineKeyboard
	_, err := Bot.Send(msg)

	for _, file := range mediaFiles {
		ext := filepath.Ext(file)

		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg", ".png":
			Bot.Send(tgbotapi.NewPhoto(ChatId, tgbotapi.FileURL(file)))
			time.Sleep(time.Second)
		default:
			Bot.Send(tgbotapi.NewVideo(ChatId, tgbotapi.FileURL(SavePics(file))))
			time.Sleep(time.Second)
		}
	}

	return err
}

func SendMergeMessage(message string, mediaFiles []string) error {
	mediaGroup := make([]interface{}, 0, len(mediaFiles))

	for _, file := range mediaFiles {
		ext := filepath.Ext(file)

		switch strings.ToLower(ext) {
		case ".jpg", ".jpeg", ".png":
			mediaGroup = append(mediaGroup, tgbotapi.InputMediaPhoto{
				BaseInputMedia: tgbotapi.BaseInputMedia{
					Type:      "photo",
					Media:     tgbotapi.FileURL(file),
					ParseMode: tgbotapi.ModeMarkdown,
				},
			})
		default:
			mediaGroup = append(mediaGroup, tgbotapi.InputMediaVideo{
				BaseInputMedia: tgbotapi.BaseInputMedia{
					Type:      "video",
					Media:     tgbotapi.FilePath(SavePics(file)),
					ParseMode: tgbotapi.ModeMarkdown,
				},
			})
		}
	}

	if mediaItem, ok := mediaGroup[0].(tgbotapi.InputMediaPhoto); ok {
		mediaItem.Caption = message
		mediaGroup[0] = mediaItem
	} else if mediaItem, ok := mediaGroup[0].(tgbotapi.InputMediaVideo); ok {
		mediaItem.Caption = message
		mediaGroup[0] = mediaItem
	}

	if len(mediaGroup) > 9 {
		Bot.SendMediaGroup(tgbotapi.NewMediaGroup(ChatId, mediaGroup[len(mediaGroup[:10]):]))
		time.Sleep(time.Second)
		_, err := Bot.SendMediaGroup(tgbotapi.NewMediaGroup(ChatId, mediaGroup[:10]))
		return err
	}

	_, err := Bot.SendMediaGroup(tgbotapi.NewMediaGroup(ChatId, mediaGroup))
	return err
}

func SendMessage(name, url, content string, mediaFies []string) {
	var err error
	message := fmt.Sprintf("博主: [%s]\n%s\n", name, content)

	if len(mediaFies) == 0 || !MergeMessage {
		err = SendSeparatelyMessage(url, message, mediaFies)
	} else {
		message += fmt.Sprintf("\n[🔗点击查看原微博](%s)", url)
		err = SendMergeMessage(message, mediaFies)
	}

	if err == nil || strings.Contains(err.Error(), `Bad Request: failed to send message #1 with the error message "Wrong file identifier/HTTP URL specified"`) {
		Insert(content, url)
	}

	if SavePicLocal {
		SaveAllPics(mediaFies)
	}

}
