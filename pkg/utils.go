package pkg

import (
	"log"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pelletier/go-toml"
)

type Config struct {
	TgBotApiToken string
	TgChatid      int64
	WeiboUid      []int
	MergeMessage  bool
	Interval      int
	SavePicLocal  bool
	SendLivePics  bool
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

	err = toml.NewDecoder(f).Decode(config)
	if err != nil {
		log.Fatal(err)
	}

	return config
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

	// Check if the file already exists
	if _, err := os.Stat(filename); err == nil {
		// File already exists, return the filename without downloading again
		return filename
	}

	// File does not exist, proceed with the download

	_, err = client.R().SetOutput(filename).Get(decodedUrl)
	if err != nil {
		log.Println("下载文件失败: ", err)
		return ""
	}

	return filename
}
