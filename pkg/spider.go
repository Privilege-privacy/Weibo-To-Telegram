package pkg

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tidwall/gjson"
)

var (
	SendLivePics bool
	SavePicLocal bool
	MergeMessage bool

	client = resty.New().SetBaseURL("https://m.weibo.cn")
	PostCh = make(chan PostQueue)
	logger = slog.New(slog.Default().Handler())
)

type (
	MediaInfo struct {
		imageUrl   string
		url        string
		isOversize bool
		isVideo    bool
		isPhoto    bool
		isLive     bool
	}
	Post struct {
		Raw           string
		Text          string
		Url           string
		BlogID        string
		SkipLivePhoto bool
		MediaInfo     []MediaInfo
	}

	PostQueue struct {
		MediaGroup []interface{}
		Info       Post
	}
)

func Run(uid int) {
	resp, err := client.R().Get("/api/container/getIndex?containerid=107603" + strconv.Itoa(uid))
	body := resp.String()
	if err != nil || gjson.Get(body, "ok").Int() != 1 {
		return
	}
	ParsePost(body)
}

func NewPost(Url, Raw, BlogID string) *Post {
	return &Post{
		Url:    Url,
		Raw:    Raw,
		BlogID: BlogID,
	}
}

func (p *Post) MessageBuilder(Author, Raw, Url string) {
	mergeMessage := MergeMessage && len(p.MediaInfo) != 0

	if p.Text != "" {
		Raw = p.Text
	}

	if mergeMessage {
		Raw = strings.ReplaceAll(Raw, "<br />", "\n")
	}

	if strings.Contains(Raw, "https://weibo.cn/sinaurl?u=") {
		Raw = parseHTML(Raw)
	}

	Raw = removeHTMLTags(Raw)

	Template := fmt.Sprintf("「 #%s 」\n\n %s\n", Author, Raw)
	if mergeMessage {
		Template = fmt.Sprintf("「 #%s 」\n\n *%s*\n", Author, Raw)
		Template += fmt.Sprintf("\n[🔗点击查看原微博](%s)", Url)
	}

	p.Raw = Raw
	p.Text = Template
}

func (p *Post) SendChBuilder(saveImage bool) PostQueue {
	mediaGroup := make([]interface{}, 0, len(p.MediaInfo))
	for _, file := range p.MediaInfo {
		url := file.imageUrl

		if p.SkipLivePhoto && file.isLive {
			url = file.url
		}

		var media tgbotapi.RequestFileData = tgbotapi.FileURL(url)
		if saveImage || SavePicLocal {
			filename := SavePics(url)
			if filename == "" {
				continue
			}
			media = tgbotapi.FilePath(SavePics(file.imageUrl))
		}

		switch {
		case file.isPhoto || p.SkipLivePhoto:
			mediaGroup = append(mediaGroup, tgbotapi.InputMediaPhoto{
				BaseInputMedia: tgbotapi.BaseInputMedia{
					Type:      "photo",
					Media:     media,
					ParseMode: tgbotapi.ModeMarkdown,
				},
			})
		case file.isVideo || file.isLive:
			mediaGroup = append(mediaGroup, tgbotapi.InputMediaVideo{
				BaseInputMedia: tgbotapi.BaseInputMedia{
					Type:      "video",
					Media:     media,
					ParseMode: tgbotapi.ModeMarkdown,
				},
			})
		}
	}

	switch mediaItem := mediaGroup[0].(type) {
	case tgbotapi.InputMediaPhoto:
		mediaItem.Caption = p.Text
		mediaGroup[0] = mediaItem
	case tgbotapi.InputMediaVideo:
		mediaItem.Caption = p.Text
		mediaGroup[0] = mediaItem
	}

	return PostQueue{
		MediaGroup: mediaGroup,
		Info:       *p,
	}
}

func (p *Post) getBlogPostContent() {
	resp, err := client.R().Get("/statuses/show?id=" + p.BlogID)
	if err == nil {
		p.Text = removeHTMLTags(gjson.Get(resp.String(), "data.text").String())
	}
}

func (p *Post) getAllImages() {
	resp, _ := client.R().Get("/statuses/show?id=" + p.BlogID)
	gjson.Get(resp.String(), "data.pics").ForEach(func(key, value gjson.Result) bool {
		p.MediaInfo = append(p.MediaInfo, parseImages(value))
		return true
	})
}

func (p *Post) addImages(list gjson.Result) {
	gjson.Get(list.String(), "@values").ForEach(func(key, value gjson.Result) bool {
		p.MediaInfo = append(p.MediaInfo, parseImages(value))
		return true
	})
}

func (p *Post) SendMessage() {
	messageInlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("🔗点击查看原微博", p.Url),
		),
	)

	logger.LogAttrs(context.Background(), slog.LevelInfo, "SendMsg",
		slog.String("微博正文", p.Raw),
		slog.String("微博链接", p.Url))

	msg := tgbotapi.NewMessage(ChatId, p.Text)
	msg.ReplyMarkup = messageInlineKeyboard
	msg.ParseMode = tgbotapi.ModeHTML
	_, _ = Bot.Send(msg)

	if len(p.MediaInfo) != 0 {
		return
	}

	for _, file := range p.MediaInfo {
		switch {
		case file.isVideo:
			Bot.Send(tgbotapi.NewVideo(ChatId, tgbotapi.FileURL(SavePics(file.imageUrl))))
			time.Sleep(time.Second)
		case file.isPhoto:
			Bot.Send(tgbotapi.NewPhoto(ChatId, tgbotapi.FileURL(SavePics(file.imageUrl))))
		}
	}
}

func ParsePost(jsonData string) {
	gjson.Get(jsonData, "data.cards").ForEach(func(_, value gjson.Result) bool {
		Url := value.Get("scheme").String()
		if ExistsInDB(Url) {
			return true
		}

		Author := value.Get("mblog.user.screen_name").String()
		Raw := value.Get("mblog.text").String()
		ImgNum := value.Get("mblog.pic_num").Int()
		BlogID := value.Get("mblog.id").String()

		VideoExist := value.Get("mblog.page_info.urls").Exists()

		posts := NewPost(Url, Raw, BlogID)

		if ImgNum > 9 {
			posts.getAllImages()
		} else {
			posts.addImages(value.Get("mblog.pics"))
		}

		if strings.Contains(Raw, "全文") {
			posts.getBlogPostContent()
		}

		if VideoExist {
			url := value.Get("mblog.page_info.urls|@values|0").String()
			resp, _ := client.R().Get(url)
			if resp.Size() < 50*1024*1024 {
				posts.MediaInfo = append(posts.MediaInfo, MediaInfo{imageUrl: url, isVideo: true})
			}
		}

		posts.MessageBuilder(Author, Raw, Url)

		if !MergeMessage || len(posts.MediaInfo) == 0 {
			posts.SendMessage()
			InsertDB(posts.Raw, posts.Url)
			return true
		}

		PostCh <- posts.SendChBuilder(false)

		return true
	})
}

func SendPosts() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	SendCount := 0
	for post := range PostCh {
		Info := post.Info
		MediaGroup := post.MediaGroup

		logger.LogAttrs(context.Background(), slog.LevelInfo, "SendPosts",
			slog.String("微博正文", Info.Raw),
			slog.String("微博链接", Info.Url))

		SendCount += len(MediaGroup)

		// 检查第一次是否成功发送媒体组
		if err := sendMediaGroup(MediaGroup, true); err != nil {
			func(ch PostQueue) {
				// 如果没有成功发送，Double Check URL是否存在于数据库中
				if !ExistsInDB(ch.Info.Url) {
					// 如果不存在，则将图片保存到本地后，再次尝试发送
					if err := sendMediaGroup(ch.MediaGroup, true); err != nil {
						// 如果仍然失败，则标记SkipLivePhoto标志，本次将不会发送 LivePhoto，并使用触发器再次尝试
						ch.Info.SkipLivePhoto = true
						func(ch PostQueue) {
							// 有概率遇到 Wrong file identifier/http url specified 错误，但是仍然成功发送的情况
							if err := sendMediaGroup(ch.MediaGroup, false); err != nil {
								logger.LogAttrs(context.Background(), slog.LevelWarn, "Retry SendPosts Failed",
									slog.String("微博正文", ch.Info.Raw),
									slog.String("微博链接", ch.Info.Url),
								)
							}
						}(ch.Info.SendChBuilder(true))
					}
				}
			}(post.Info.SendChBuilder(true))
		}

		if ok := InsertDB(Info.Raw, Info.Url); !ok {
			logger.LogAttrs(context.Background(), slog.LevelWarn, "写入数据库失败",
				slog.String("URL", Info.Url),
			)
		}

		if SendCount >= 30 {
			<-ticker.C
			SendCount = 0
		}
	}
}

func parseImages(value gjson.Result) MediaInfo {
	if value.Get("videoSrc").Exists() && SendLivePics {
		return MediaInfo{
			imageUrl: value.Get("videoSrc").String(),
			url:      value.Get("large.url").String(),
			isLive:   true,
		}
	}

	url := value.Get("large.url").String()
	width, height := value.Get("large.geo.width").Int(), value.Get("large.geo.height").Int()

	if width+height > 10000 || width/height > 20 {
		return MediaInfo{imageUrl: url, isOversize: true}
	}

	return MediaInfo{imageUrl: url, isPhoto: true}
}
