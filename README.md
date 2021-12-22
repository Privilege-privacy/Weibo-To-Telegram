# Weibo-To-Telegram
将微博动态实时同步到 Telegram 支持关注多个用户

## 使用 

- Linux x86
     
     从 [Releases](https://github.com/Privilege-privacy/Weibo-To-Telegram/releases/download/main/weibo.zip) 下载 zip 后解压，修改 `conf.json` 文件， 将 `Tgbotapi` 的 Value 值改为从 Botfather 中获取的 API Token ， `Chat_id` 的 value 值 ，可填 Telegram 群组 ID 或 个人账户 ID(可从@getmyid_bot 这个机器人获得)，`Weibo_uid` 的 Value 就是跟在 weibo.com/2201313382 的那一串数字，每个 uid 之间用 `,` 隔开，最后运行 `main` 二进制文件。

- Docker
    

    部署完容器后，进入 `/` 目录，根据 `linux x86` 的方法修改 `conf.json` 后， `./main` 运行 

        docker pull privileges/weibo:latest
     
     
     ### Windows 或其他平台架构，自行 `clone` 代码后交叉编译
