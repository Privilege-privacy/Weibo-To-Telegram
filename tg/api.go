package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var Token string = ""
var Userid string = ""

type sendPhotoReqBody struct {
	ChatID string `json:"chat_id"`
	Photo  string `json:"photo"`
}

type sendMessageReqBody struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

type sendReply struct {
	ChatID      string `json:"chat_id"`
	Text        string `json:"text"`
	ReplyMarkup struct {
		InlineKeyboard [][]struct {
			Text string `json:"text"`
			URL  string `json:"url"`
		} `json:"inline_keyboard"`
	} `json:"reply_markup"`
}

func SendPhoto(jpgurl string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto", Token)
	reqbody := &sendPhotoReqBody{
		ChatID: Userid,
		Photo:  fmt.Sprintf(jpgurl),
	}
	reqbytes, err := json.Marshal(reqbody)
	_, err = http.Post(url, "application/json", bytes.NewBuffer(reqbytes))
	if err != nil {
		println(err)
	}
}

func SendMessage(text string) {

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", Token)
	reqbody := &sendMessageReqBody{
		ChatID: Userid,
		Text:   fmt.Sprintf(text),
	}
	reqbytes, err := json.Marshal(reqbody)
	_, err = http.Post(url, "application/json", bytes.NewBuffer(reqbytes))
	if err != nil {
		println(err)
	}
}

func SendMessageReply(text string, author string, Weibourl string) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", Token)
	reqbody := &sendReply{
		ChatID: Userid,
		Text:   fmt.Sprintf("博主：[%s]\n%s", author, text),
		ReplyMarkup: struct {
			InlineKeyboard [][]struct {
				Text string "json:\"text\""
				URL  string "json:\"url\""
			} "json:\"inline_keyboard\""
		}{
			[][]struct {
				Text string "json:\"text\""
				URL  string "json:\"url\""
			}{{struct {
				Text string "json:\"text\""
				URL  string "json:\"url\""
			}{
				Text: fmt.Sprintf("🔗点击查看原微博"),
				URL:  Weibourl,
			}}},
		},
	}
	reqbytes, err := json.Marshal(reqbody)
	_, err = http.Post(url, "application/json", bytes.NewBuffer(reqbytes))
	if err != nil {
		println(err)
	}
}
