package utils

import (
	"net/url"
	"strings"
)

// 获取指定字符的中间值
func Between(str, starting, ending string) string {
	s := strings.Index(str, starting)
	if s < 0 {
		return ""
	}
	s += len(starting)
	e := strings.Index(str[s:], ending)
	if e < 0 {
		return ""
	}
	return str[s : s+e]
}

// 解析并返回请求URL模版
func ParseZiRoomUrl(requestUrl string) string {

	// 解析url 并保证没有错误
	curl, err := url.Parse(requestUrl)
	if err != nil {
		panic("请求URL错误（解析失败）！" + err.Error())
	}

	end := "/?isOpen=0"
	oParams := strings.Split(requestUrl, "?")
	if len(oParams) >= 2 {
		end = "/?" + oParams[1]
	}

	params := strings.Split(curl.Path, "/")

	// 需满足规则，方便制作模版
	if len(params) < 3 {
		panic("请求URL错误（格式）！" + err.Error())
	}

	details := strings.Split(params[2], "-")

	if strings.Contains(details[len(details)-1], "p") {
		details[len(details)-1] = "p#"
	} else {
		details = append(details, "p#")
	}

	return "https://www.ziroom.com/" + params[1] + "/" + strings.Join(details, "-") + end
}
