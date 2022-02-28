package spider

import (
	"Weibo-To-Telegram/db"
	"Weibo-To-Telegram/tg"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type Weibo struct {
	Ok   int
	Data struct {
		Cards []struct {
			Scheme string
			Mblog  struct {
				Id         string
				Created_at string
				Text       string
				User       struct {
					Screen_name string
				}
				Pics []struct {
					Url   string
					Large struct {
						Url string
					}
				}
			}
		}
	}
}

func Run(uid int) {
	url := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?containerid=107603%d", uid)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		println(err)
	}
	var res Weibo
	err = json.Unmarshal(body, &res)
	//uid 不存在则返回 0
	if res.Ok == 1 {
		for _, item := range res.Data.Cards {
			//检测这条微博是否在数据库中已保存
			if db.Check(item.Scheme) <= 0 {
				weibocontent := item.Mblog.Text
				weiboPhoto := item.Mblog.Pics
				tg.SendMessageReply(regx(weibocontent), item.Mblog.User.Screen_name, item.Scheme)
				println(regx(weibocontent))
				for _, url := range weiboPhoto {
					tg.SendPhoto(url.Large.Url)
				}
				db.Insert(regx(weibocontent), item.Scheme)
			}
		}
	} else {
		fmt.Printf("uid 错误或炸号 %d \n", uid)
	}
	fmt.Printf("User %d Done\n", uid)
}

func regx(src string) string {
	re, _ := regexp.Compile("<[^>]*>")
	src = re.ReplaceAllString(src, "")
	return src
}
