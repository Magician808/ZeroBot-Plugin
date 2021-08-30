package lolicon

import (
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/tidwall/gjson"
)

func querylolicon(keyword string, pid *string, uid *string, title *string, author *string, urls *string) {
	clientsetu := http.Client{}
	req, err := http.NewRequest("GET", "https://api.lolicon.app/setu/v2?r18=1?keyword="+url.QueryEscape(keyword), nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36 Edg/87.0.664.66")
	res, err := clientsetu.Do(req)
	if err != nil {
		return
	}
	data, err := ioutil.ReadAll(res.Body)
	_ = res.Body.Close()
	if err != nil {
		return
	}

	*pid = gjson.ParseBytes(data).Get("data.0.pid").String()

	*uid = gjson.ParseBytes(data).Get("data.0.uid").String()

	*title = gjson.ParseBytes(data).Get("data.0.title").String()

	*author = gjson.ParseBytes(data).Get("data.0.author").String()

	*urls = gjson.ParseBytes(data).Get("data.0.urls.original").String()

	return
}
