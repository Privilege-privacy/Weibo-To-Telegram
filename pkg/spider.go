package pkg

import (
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/hashicorp/go-getter"
	"github.com/tidwall/gjson"
)

var (
	SendLivePics bool
	SavePicLocal bool
	MergeMessage bool

	client = resty.New().SetBaseURL("https://m.weibo.cn")
	PostCh = make(chan PostQueue)
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

func NewPost(Url, Raw, BlogID string) *Post {
	return &Post{
		Url:    Url,
		Raw:    Raw,
		BlogID: BlogID,
	}
}

func (p *Post) MessageBuilder(Author, Raw, Url string) {
	if p.Text != "" {
		Raw = p.Text
	}

	Template := fmt.Sprintf("「 #%s 」\n\n %s\n", Author, Raw)

	if MergeMessage && len(p.MediaInfo) != 0 {
		Template = fmt.Sprintf("「 #%s 」\n\n *%s*\n", Author, Raw)
		Template += fmt.Sprintf("\n[🔗点击查看原微博](%s)", Url)
	}

	p.Text = Template
}

func (p *Post) SeedChBuilder(saveImage bool) PostQueue {
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
	log.Println(p.Raw, p.Url)

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

func Run(uid int) {
	resp, err := client.R().Get("/api/container/getIndex?containerid=107603" + strconv.Itoa(uid))
	if err != nil || gjson.Get(resp.String(), "ok").Int() != 1 {
		return
	}
	ParsePost(resp.String())
}

func ParsePost(jsonData string) {
	gjson.Get(jsonData, "data.cards").ForEach(func(_, value gjson.Result) bool {
		Url := value.Get("scheme").String()
		if ExistsInDB(Url) {
			return true
		}

		Author := value.Get("mblog.user.screen_name").String()
		Raw := removeHTMLTags(value.Get("mblog.text").String())
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

		PostCh <- posts.SeedChBuilder(false)

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

		log.Println(Info.Raw, Info.Url)

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
								log.Println("重发失败: ", ch.Info.Raw, ch.Info.Url)
							}
						}(ch.Info.SeedChBuilder(true))
					}
				}
			}(post.Info.SeedChBuilder(true))
		}

		if ok := InsertDB(Info.Raw, Info.Url); !ok {
			log.Println("插入失败: ", Info.Url)
		}

		if SendCount >= 30 {
			<-ticker.C
			SendCount = 0
		}
	}
}

func removeHTMLTags(src string) string {
	return strings.TrimSpace(regexp.MustCompile("<[^>]*>").ReplaceAllString(src, ""))
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

func SavePics(schema string) string {
	downloadPath := "download/"
	_, err := os.Stat(downloadPath)
	if os.IsNotExist(err) {
		if err := os.Mkdir(downloadPath, os.ModePerm); err != nil {
			log.Fatal("创建 download 文件夹失败, 检查当前目录下的文件夹权限", err)
		}
	}

	decodedUrl, _ := url.QueryUnescape(schema)
	filename := filepath.Join(downloadPath, path.Base(decodedUrl))

	// Check if the file already exists
	if _, err := os.Stat(filename); err == nil {
		// File already exists, return the filename without downloading again
		return filename
	}

	// File does not exist, proceed with the download

	if err := getter.GetFile(filename, decodedUrl); err != nil {
		log.Println("下载文件失败: ", err)
		return ""
	}

	// Check image size and compress if necessary
	if getImageSize(filename) > 10*1024*1024 {
		compressAndReplaceImage(filename, 70)
	}

	return filename
}

func compressAndReplaceImage(imagePath string, compressionQuality int) error {
	inputFile, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("无法打开输入图片文件: %v", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		return fmt.Errorf("无法解码图片文件: %v", err)
	}

	tempOutputPath := imagePath + ".temp"
	outputFile, err := os.Create(tempOutputPath)
	if err != nil {
		return fmt.Errorf("无法创建临时输出图片文件: %v", err)
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, img, &jpeg.Options{Quality: compressionQuality})
	if err != nil {
		return fmt.Errorf("无法压缩图片: %v", err)
	}

	outputFile.Close()
	err = os.Rename(tempOutputPath, imagePath)
	if err != nil {
		return fmt.Errorf("无法替代原始图片文件: %v", err)
	}

	return nil
}

func getImageSize(filePath string) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0
	}
	if !fileInfo.Mode().IsRegular() {
		return 0
	}
	return fileInfo.Size()
}
