# Weibo-To-Telegram
将微博动态实时同步到 Telegram 支持关注多个用户

## 使用
1. 从 [releases](https://github.com/Privilege-privacy/Weibo-To-Telegram/releases) 页面下载对应平台的压缩包并解压缩。
2. 运行 `./main` 会在当前目录生成 `config.toml` 文件和 `weibo.db` 文件。
3. 填写 `config.toml` 中的 `tgbotapi`,`tgchatid`,`weibo_uid` 配置项 `Value` 值。
   1. `tgbotapi` 的配置项为 Telegram Bot API Token，可从 [Telegram Botfather](https://t.me/botfather) 获取。
   2. `tgchatid` 的配置可以填写 Telegram 群聊 ID 或者 Telegram 用户 ID。
      1. 如果填写的是群聊 ID，则需要在 Telegram 群聊中加入 Bot，否则无法收到消息。  
         获取当前群组 ID 可以把你的 Telegram Bot 和 [getmyid_bot](https://t.me/getmyid_bot) 拉进同一个群组，然后随便发一条信息 getmyid_bot 就会输出当前群组 id 和 Telegram 用户 ID。
      4. 填 Telegram 群聊 ID 时，Bot 就会转发消息到指定的群组。
      5. 写 Telegram 用户 ID 时，Bot 就会转发消息到指定的用户。
   3. `weibo_uid` 的配置项为微博用户 ID，可以在 [微博用户主页](https://weibo.com/u/<weibo_uid>) 中查看。
      1. 例如 weibo.com/u/2201313382 这个微博个人主页，那么这个 `weibo_uid` 配置项的值就是 2201313382，多个用户之间用`,`分隔。  
         示例： "Weibo_uid":[2201313382,123123123]
4. 修改完配置后，运行 `./main` 即可开始转发。