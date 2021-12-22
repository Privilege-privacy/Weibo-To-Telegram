package main

import (
	"Weibo-To-Telegram/spider"
	"Weibo-To-Telegram/tg"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Tgbotapi string  `json:"Tgbotapi"`
	ChatID   string  `json:"Chat_id"`
	WeiboUID []int64 `json:"Weibo_uid"`
}

func main() {
	
	//读取配置文件
	file, _ := os.Open("conf.json")
	defer file.Close()
	decode := json.NewDecoder(file)
	conf := Config{}
	err := decode.Decode(&conf)
	if err != nil {
		fmt.Println(err)
	}

	tg.Token = conf.Tgbotapi
	tg.Userid = conf.ChatID

	for {
		for _, uid := range conf.WeiboUID {
			spider.Run(int(uid))
		}
		fmt.Printf("防止 IP 被拉黑 一分钟后下一轮\n")
		time.Sleep(60 * time.Second)
	}
}
