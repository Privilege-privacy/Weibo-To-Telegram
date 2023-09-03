package main

import (
	"log"
	"os"
	"time"

	"github.com/Privilege-privacy/Weibo-To-Telegram/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Config struct {
	TgBotApiToken string
	TgChatid      int64
	WeiboUid      []int
	MergeMessage  bool
	Interval      int
	SavePicLocal  bool
	SendLivePics  bool
}

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

	bot, err := tgbotapi.NewBotAPI(pkg.TgBotApiToken)
	if err != nil {
		log.Fatal("连接 Telegram 失败", err)
	}
	pkg.Bot = bot

	go pkg.SendPosts()

	for {
		for _, uid := range config.WeiboUid {
			pkg.Run(uid)
			time.Sleep(3 * time.Second)
		}
		time.Sleep(time.Duration(config.Interval) * time.Second)
	}
}
