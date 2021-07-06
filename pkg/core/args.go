package core

import "time"

// 当前已截取的房源
var LastTimeData = make(map[string]Room)

// 自如房源地址
var ZiroomURL string

// 钉钉预警
var WebhookUrl string

// 钉钉预警 授权KEY
var WebhookUrlKey = "114"

// 每x分钟进行一次比较
var TaskInterval time.Duration

// 每x秒每页请求限制（强制3秒）
var PageSecondInterval time.Duration = 3

// [测试] 通知x条房源
var TestForNotice = 0

func Check() {
	if len(ZiroomURL) == 0 {
		panic("房源地址未设置。")
	}

	if len(WebhookUrl) == 0 {
		panic("钉钉通知未设置。")
	}

	if len(WebhookUrlKey) == 0 {
		panic("钉钉密钥未设置。")
	}

	if TaskInterval <= 300 {
		TaskInterval = 300
	}
}
