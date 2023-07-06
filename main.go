package main

import (
	"log"
	"os"
	"time"

	"github.com/Privilege-privacy/Weibo-To-Telegram/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
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

func main() {
	var config Config
	viper.AddConfigPath(".")
	if _, file := os.Stat("config.toml"); os.IsNotExist(file) {
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.SetDefault("TgBotApiToken", "")
		viper.SetDefault("TgChatid", 0)
		viper.SetDefault("WeiboUid", []int{})
		viper.SetDefault("MergeMessage", true)
		viper.SetDefault("Interval", 120)
		viper.SetDefault("SavePicLocal", false)
		viper.SetDefault("SendLivePics", true)

		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatal("保存配置文件失败", err)
		}
		log.Fatal("根据要求填写 Config.toml 后运行")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("加载配置文件错误", err)
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal("解析配置文件错误", err)
	}

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

	for {
		for _, uid := range config.WeiboUid {
			pkg.Run(uid)
			time.Sleep(3 * time.Second)
		}
		time.Sleep(time.Duration(config.Interval))
	}
}
