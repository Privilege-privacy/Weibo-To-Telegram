package main

import (
	"Weibo-To-Telegram/spider"
	"Weibo-To-Telegram/tg"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

func init() {
	_, file := os.Stat("config.toml")
	if file == nil {
		viper.AddConfigPath(".")
	}
	if os.IsNotExist(file) {
		log.Println("未在当前目录找到配置文件 在当前目录创建 Config.toml")
		log.Println("根据要求填写 Config.toml 后运行")
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")
		viper.Set("Tgbotapi", "")
		viper.Set("TgChatid", "")
		viper.Set("Weibo_uid", []int{})
		if err := viper.SafeWriteConfig(); err != nil {
			log.Println("保存配置文件错误", err)
		}
		os.Exit(3)
	}
}

func main() {
	if err := viper.ReadInConfig(); err != nil {
		log.Println("加载配置文件错误", err)
	}
	tg.Token = viper.GetString("tgbotapi")
	tg.Userid = viper.GetString("tgchatid")

	for {
		for _, uid := range viper.GetIntSlice("weibo_uid") {
			spider.Run(uid)
		}
		fmt.Printf("防止 IP 被拉黑 一分钟后下一轮\n")
		time.Sleep(60 * time.Second)
	}
}
