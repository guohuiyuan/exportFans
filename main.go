package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

var (
	followURL = "https://api.bilibili.com/x/relation/followers?vmid=%v&ps=%v&pn=%v"
	defaultPs = 50
)

func getFans(mid, cookie string, ps, pn int) (fanList []string) {
	fanList = make([]string, 0, 1000)
	url := fmt.Sprintf(followURL, mid, ps, pn)
	fmt.Println(url)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "apifox/1.0.0 (https://www.apifox.cn)")
	req.Header.Add("Cookie", cookie)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	gjson.ParseBytes(body).Get("data.list").ForEach(func(_, v gjson.Result) bool {
		fanList = append(fanList, v.Get("mid").String())
		return true
	})
	return
}

func main() {
	cookie := flag.String("c", "SESSDATA=6e815a8a%2C1681290824%2C7abe2%2Aa1", "Set AccessToken of WSClient.")
	// 直接写死 URL 时，请更改下面第二个参数
	mid := flag.String("m", "1561377116", "Set Url of WSClient.")
	// 默认昵称
	flag.Parse()
	i := 1
	fanList := make([]string, 0, 1000)
	for {
		list := getFans(*mid, *cookie, defaultPs, i)
		if len(list) == 0 {
			break
		}
		fanList = append(fanList, list...)
		i++
	}
	_ = os.WriteFile("fanList.txt", []byte(strings.Join(fanList, "\n")), 0666)
}
