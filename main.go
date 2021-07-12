package main

import (
	"time"
	"ziroom/pkg/core"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	ding         = kingpin.Flag("ding", "钉钉消息通知接口地址").Short('d').String()
	dingKey      = kingpin.Flag("dingKey", "钉钉消息通知授权KEY（白名单）").Short('k').Default("自如").String()
	taskInterval = kingpin.Flag("taskInterval", "任务周期间隔时长（最少5分钟），单位：秒").Short('p').Default("300").Int()
	//TestForNotice = kingpin.Flag("TestForNotice", "[测试] 每次调度推送t个房源通知").Default("0").Short('t').Int()
	url = kingpin.Arg("url", "自如网页版请求链接").String()
)

func main() {
	kingpin.Parse()

	core.WebhookUrl = *ding
	core.WebhookUrlKey = *dingKey
	core.TaskInterval = time.Duration(*taskInterval)
	core.ZiroomURL = *url
	//core.TestForNotice = *TestForNotice

	// 数据检查
	core.Check()

	// 数据初始化
	core.FirstDataInit()

	// 定时任务
	core.TaskForSearch()
}
