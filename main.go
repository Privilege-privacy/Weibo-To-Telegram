package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/Privilege-privacy/Weibo-To-Telegram/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func init() {
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		if err := pkg.CreateConfig(); err != nil {
			log.Fatalln("创建 Config.toml 失败:", err)
		}
		log.Fatalln("根据要求填写 Config.toml 后运行")
	}
}

func main() {
	config := pkg.LoadConfig()

	pkg.TgBotApiToken = config.TgBotApiToken
	pkg.ChatId = config.TgChatid
	pkg.SendLivePics = config.SendLivePics
	pkg.SavePicLocal = config.SavePicLocal
	pkg.MergeMessage = config.MergeMessage

	interval := time.Duration(config.Interval) * time.Second
	WeiboUid := config.WeiboUid

	bot, err := tgbotapi.NewBotAPI(pkg.TgBotApiToken)
	if err != nil {
		log.Fatal("连接 Telegram 失败", err)
	}
	pkg.Bot = bot

	post := make(chan pkg.PostQueue)
	var wg sync.WaitGroup

	wg.Add(2)
	go pkg.SendPosts(post)
	go func() {
		for range time.Tick(interval) {
			for _, uid := range WeiboUid {
				pkg.Run(uid, post)
			}
		}
	}()

	wg.Wait()
}
