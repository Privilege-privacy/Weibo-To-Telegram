package main

import (
	"Weibo-To-Telegram/internal"
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

var Silent *bool

func init() {
	Silent = flag.Bool("s", true, "默认开启消息打印 false 关闭")
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
		viper.Set("tgbotapitoken", "")
		viper.Set("tguseridorchatid", 0)
		viper.Set("Weibo_uid", []int{})
		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatal("保存配置文件失败", err)
		}
		os.Exit(3)
	}
}

func main() {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("加载配置文件错误", err)
	}
	flag.Parse()

	internal.TgBotApiToken = viper.GetString("tgbotapitoken")
	internal.TgUseridORChatId = viper.GetInt64("tguseridorchatid")
	Silents := *Silent

	for {
		for _, uid := range viper.GetIntSlice("weibo_uid") {
			internal.Run(uid, Silents)
		}
		if Silents {
			fmt.Printf("防止 IP 被拉黑 一分钟后下一轮\n")
		}
		time.Sleep(time.Minute)
	}
}
