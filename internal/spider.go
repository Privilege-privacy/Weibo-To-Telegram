package internal

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

func Run(uid int, silent bool) {
	body := doGet(uid)
	//uid 不存在则返回 0
	if gjson.Get(body, "ok").Int() == 1 {
		gjson.Get(body, "data.cards").ForEach(func(key, value gjson.Result) bool {
			name := value.Get("mblog.user.screen_name").String()
			scheme := value.Get("scheme").String()
			content := value.Get("mblog.text").String()
			//是否以在数据库中保存过
			if Check(scheme) <= 0 {
				sendMessage(name, regx(content), scheme)
				if silent {
					fmt.Println(name, regx(content))
				}
				pics := value.Get("mblog.pics").Array()
				for _, pic := range pics {
					sendPhoto(pic.Get("large.url").String())
				}
				Insert(regx(content), scheme)
			}
			return true
		})
	}
	if silent {
		log.Printf("User %d Done\n", uid)
	}
}

func regx(src string) string {
	r, _ := regexp.Compile("<[^>]*>")
	return r.ReplaceAllString(src, "")
}

func doGet(uid int) string {
	url := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?containerid=107603%d", uid)
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		return string(body)
	}
	return ""
}
