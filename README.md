# Weibo-To-Telegram
将微博动态实时同步到 Telegram 支持关注多个用户

## 使用 

从 [Releases](https://github.com/Privilege-privacy/Weibo-To-Telegram/releases) 下载对应平台的压缩包，解压并运行二进制文件后，会在当前目录生成`config.toml` `weibo.db`文件，
修改 `config.toml` 配置文件，`tgbotapi` 的 `Value`值填写在 [Botfather](https://t.me/botfather) 获取到的 `API TOKEN`, `tgchatid` 的 `Value` 值可以填写在 [getmyid_bot](https://t.me/getmyid_bot) 获取到的`User_id ` , 或者把 [getmyid_bot](https://t.me/getmyid_bot) 和你创建的 `tgbot` 拉进同一个群组后随便发一条信息，就会获得当前的`Chat_id` , 如果填写的是 `User_id` 运行时 `Bot` 就会以私聊的方式发送信息，填写 `Chat_id` 运行时 `Bot` 就会转发信息到群组，`weibo_uid` 的 `Value` 值就是微博个人主页后面的那一串数字 `weibo.com/2201313382` 多个用户用 `,` 隔开，修改完成后，运行二进制文件就会开始转发。