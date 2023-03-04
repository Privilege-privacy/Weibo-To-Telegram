package pkg

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"strings"
)

var (
	SendLivePics bool
	SavePicLocal bool
	MergeMessage bool
)

func Run(uid int) {
	body := doGet(uid)

	if gjson.Get(body, "ok").Int() != 1 {
		return
	}

	gjson.Get(body, "data.cards").ForEach(func(key, value gjson.Result) bool {
		name := value.Get("mblog.user.screen_name").String()
		url := value.Get("scheme").String()
		content := regx(value.Get("mblog.text").String())
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

func regx(src string) string {
	re := regexp.MustCompile("<[^>]*>")
	return strings.TrimSpace(re.ReplaceAllString(src, ""))
}

func doGet(uid int) string {
	url := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?containerid=107603%d", uid)
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return string(body)
	}
	return ""
}

func GetFullContent(bid string) string {
	var url strings.Builder
	url.WriteString("https://m.weibo.cn/statuses/show?id=")
	url.WriteString(bid)

	resp, err := http.Get(url.String())
	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return regx(gjson.Get(string(body), "data.text").String())
	}
	return ""
}

func GetFullPics(bid string) []string {
	var url strings.Builder
	url.WriteString("https://m.weibo.cn/statuses/show?id=")
	url.WriteString(bid)

	resp, err := http.Get(url.String())
	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return GetListPics(gjson.Get(string(body), "data.pics").Array())
	}
	return nil
}

func GetListPics(list []gjson.Result) []string {
	var temp = make([]string, 0, len(list))
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
		err := os.Mkdir("download/", os.ModePerm)
		if err != nil {
			log.Fatal("创建 download 文件夹失败, 检查当前目录下的文件夹权限", err)
		}
	}

	schema, _ = url.QueryUnescape(schema)

	var filename strings.Builder
	filename.WriteString("download/")
	filename.WriteString(path.Base(schema))

	resp, err := http.Get(schema)
	if err != nil {
		log.Println("图片", filename.String(), "下载失败")
		return ""
	}
	defer resp.Body.Close()

	file, err := os.Create(filename.String())
	if err != nil {
		log.Println("文件创建失败")
		return ""
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err == nil {
		return filename.String()
	}
	return ""
}

func SaveAllPics(pics []string) {
	for _, pic := range pics {
		var filename strings.Builder
		filename.WriteString("download/")
		filename.WriteString(path.Base(pic))

		_, err := os.Stat(filename.String())
		if os.IsNotExist(err) {
			SavePics(pic)
		}
	}
}
