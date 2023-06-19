package pkg

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
)

var (
	SendLivePics bool
	SavePicLocal bool
	MergeMessage bool
	client       = resty.New().SetBaseURL("https://m.weibo.cn").R()
)

func Run(uid int) {
	resp, err := client.Get("/api/container/getIndex?containerid=107603" + strconv.Itoa(uid))
	if err != nil {	return }

	if gjson.Get(resp.String(), "ok").Int() != 1 { return }

	gjson.Get(resp.String(), "data.cards").ForEach(func(key, value gjson.Result) bool {
		name := value.Get("mblog.user.screen_name").String()
		url := value.Get("scheme").String()
		content := removeHTMLTags(value.Get("mblog.text").String())
		pics := GetListPics(value.Get("mblog.pics").Array())

		if Check(url) != 0 {
			return true
		}

		if strings.Contains(content, "全文") {
			content = GetFullContent(value.Get("mblog.bid").String())
		}

		if value.Get("mblog.pic_num").Int() > 9 {
			pics = GetFullPics(value.Get("mblog.bid").String())
		}

		log.Println(name, content, url)

		if SendLivePics && value.Get("mblog.pics.#.videoSrc").Exists() {
			pics = GetLivePics(value.Get("mblog.pics").Array())
			SendMessage(name, url, content, pics)
			return true
		}

		SendMessage(name, url, content, pics)

		return true
	})
}

func removeHTMLTags(src string) string {
	return strings.TrimSpace(regexp.MustCompile("<[^>]*>").ReplaceAllString(src, ""))
}

func GetFullContent(bid string) string {
	resp, err := client.Get("/statuses/show?id=" + bid)
	if err == nil {
		return removeHTMLTags(gjson.Get(resp.String(), "data.text").String())
	}
	return ""
}

func GetFullPics(bid string) (slice []string) {
	resp, _ := client.Get("/statuses/show?id=" + bid)
	gjson.Get(resp.String(), "data.pics").ForEach(func(key, value gjson.Result) bool {
		slice = append(slice, value.Get("large.url").String())
		return true
	})
	return nil
}

func GetListPics(list []gjson.Result) []string {
	temp := make([]string, 0, len(list))
	for _, result := range list {
		temp = append(temp, result.Get("large.url").String())
	}
	return temp
}

func GetLivePics(list []gjson.Result) (temp []string) {
	for _, result := range list {
		if result.Get("videoSrc").Exists() {
			temp = append(temp, result.Get("videoSrc").String())
		} else {
			temp = append(temp, result.Get("large.url").String())
		}
	}
	return
}

func SavePics(schema string) string {
	_, err := os.Stat("download/")
	if os.IsNotExist(err) {
		if err := os.Mkdir("download/", os.ModePerm); err != nil {
			log.Fatal("创建 download 文件夹失败, 检查当前目录下的文件夹权限", err)
		}
	}

	schema, _ = url.QueryUnescape(schema)
	filename := filepath.Join("download", path.Base(schema))

	resp, err := http.Get(schema)
	if err != nil {
		log.Println("图片", filename, "下载失败")
		return ""
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		log.Println("文件创建失败")
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err == nil {
		return filename
	}
	return ""
}

func SaveAllPics(pics []string) {
	for _, pic := range pics {
		filename := filepath.Join("download", path.Base(pic))
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			SavePics(pic)
		}
	}
}
