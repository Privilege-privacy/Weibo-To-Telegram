package main

import (
	"Weibo-To-Telegram/internal"
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
		log.Println("未在当前目录找到配置文件 将在当前目录创建 Config.toml")
		viper.SetConfigName("config")
		viper.SetConfigType("toml")
		viper.AddConfigPath(".")

		viper.Set("TgBotApiToken", "")
		viper.Set("TgChatid", 0)
		viper.Set("Weibo_uid", []int{})
		viper.Set("MergeMessage", true)
		viper.Set("Interval", 120)
		viper.Set("SavePicLocal", false)
		viper.Set("SendLivePics", true)

		if err := viper.SafeWriteConfig(); err != nil {
			log.Fatal("保存配置文件失败", err)
		}
		log.Fatal("根据要求填写 Config.toml 后运行")
	}
}

func main() {
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("加载配置文件错误", err)
	}

	interval := viper.GetDuration("Interval")

	for {
		for _, uid := range viper.GetIntSlice("weibo_uid") {
			internal.Run(uid)
			time.Sleep(3 * time.Second)
		}
		time.Sleep(time.Second * interval)
	}
}
