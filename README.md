# Weibo-To-Telegram
将微博动态实时同步到 Telegram 支持关注多个用户

## 使用
1. 从 [releases](https://github.com/Privilege-privacy/Weibo-To-Telegram/releases) 页面下载对应平台的压缩包并解压缩。
   

2. 运行 `./main` 会在当前目录生成 `config.toml` 文件和 `weibo.db` 文件。
   

3. 修改 `config.toml` 配置文件。

| 配置项              | 含义                             | 示例                                       |
|------------------|--------------------------------|------------------------------------------|
| tgbotapitoken    | Telegram Bot Api Token         | 90804:pqwozgkoadsaa...                   |
| tguseridorchatid | 可填写 Telegram 用户 Id 或需要转发的群组 ID | `UserId: 1234586` `GroupId:  -294892475` |
| Weibo_uid        | 微博用户 UID                       | 2201313382                               |

### 配置项

#### TGBotApiToken

> 可从 [Telegram BotFather](https://t.me/botfather) 处创建新 Bot 或选择已有 Bot 获取。

#### TGUserIdORChatId
> 将你创建的 `Bot` 和 [GetMyId_bot](https://t.me/getmyid_bot) 拉进同一个群组内，随便发送一条信息 `GetMyId_bot` 就会输出当前的`Chat ID`和你的`UserID`
>> 填写 `Chat ID ` Bot 就会转发消息到当前群组</br> 填写 `User ID` 就会以私聊的方式转发。

#### Weibo_uid
> `weibo_uid` 的配置项为微博用户 ID，可以在 [微博用户主页](https://weibo.com/u/<your_weibo_uid>) 中查看。</br>

> 例如 **weibo.com/u/2201313382** 这个微博个人主页，那么这个 `weibo_uid` 配置项的值就是 **2201313382**，多个用户之间用`,`分隔。</br> 示例： "Weibo_uid":[2201313382,123456,654321]

### 感谢
#### Idea： [Weibo](https://github.com/cndiandian/weibo)


