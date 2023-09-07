package pkg

import (
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml"
)

type Config struct {
	MergeMessage  bool
	SavePicLocal  bool
	SendLivePics  bool
	Interval      int
	TgChatid      int64
	WeiboUid      []int
	TgBotApiToken string
}

func CreateConfig() error {
	config := &Config{
		TgBotApiToken: "token",
		TgChatid:      123456,
		WeiboUid:      []int{1, 2, 3},
		MergeMessage:  true,
		Interval:      120,
		SavePicLocal:  false,
		SendLivePics:  true,
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile("config.toml", data, 0o644)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfig() *Config {
	config := &Config{}
	f, err := os.Open("config.toml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	err = toml.NewDecoder(f).Decode(config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func GetVideoLength(input string) int {
	resp, err := http.Get(input)
	if err != nil {
		return 0
	}
	defer resp.Body.Close()
	return int(resp.ContentLength)
}

func parseHTML(input string) string {
	urlMatches := regexp.MustCompile(`<a\s+href="(https://weibo\.cn/sinaurl\?u=.*?)".*?>`).FindAllString(input, -1)
	extractedURL, ok := extractURL(urlMatches[0])
	if ok {
		extractedURL, _ = url.QueryUnescape(extractedURL)
		return strings.ReplaceAll(input, urlMatches[0], extractedURL)
	}

	return input
}

func extractURL(input string) (string, bool) {
	startTag := "https://weibo.cn/sinaurl?u="
	startIndex := strings.Index(input, startTag)

	if startIndex == -1 {
		return "", false
	}

	startIndex += len(startTag)
	endIndex := strings.Index(input[startIndex:], "\"")
	if endIndex == -1 {
		return "", false
	}

	endIndex += startIndex
	return input[startIndex:endIndex] + " ", true
}

func removeHTMLTags(src string) string {
	return strings.TrimSpace(regexp.MustCompile("<[^>]*>").ReplaceAllString(src, ""))
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

	if _, err := os.Stat(filename); err == nil {
		return filename
	}

	_, err = client.R().SetOutput(filename).Get(decodedUrl)
	if err != nil {
		log.Println("下载文件失败: ", err)
		return ""
	}

	return filename
}
