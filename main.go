package main

import (
	"strings"
	"time"
	"ziroom/internal/pkg"
	"ziroom/internal/pkg/core"
	notice2 "ziroom/internal/pkg/notice"
	"ziroom/pkg/platform"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	notice      = kingpin.Flag("notice", "消息通知平台：ding（钉钉）、fs（飞书）").Short('p').Default("ding").String()
	noticeUrl      = kingpin.Flag("noticeUrl", "消息通知接口地址").Short('u').String()
	noticeKey      = kingpin.Flag("noticeKey", "消息通知授权KEY（白名单）").Short('k').Default("推送").String()
	taskInterval = kingpin.Flag("taskInterval", "任务周期间隔时长，单位：秒").Short('t').Default("300").Int()
	url          = kingpin.Arg("url", "自如/链家网页版房源请求地址，支持录入多地址，多个地址通过`空格`分隔。").Strings()
)

func main() {
	kingpin.Parse()
	examples := *url

	if examples == nil {
		panic("请设置自如/链家网页版房源请求地址。")
	}

	runExamples := make([]core.AbilityService, 0, 10)

	for i := 0; i < len(examples); i++ {
		value := examples[i]

		if strings.Contains(value, "ziroom") {
			example := &platform.ZIRoomImpl{
				InputURL: value,
			}

			// 生成请求模版
			example.Validation()
			runExamples = append(runExamples, example)
		} else if strings.Contains(value, "lianjia") {
			example := &platform.LianJiaImpl{
				InputURL: value,
			}

			// 生成请求模版
			example.Validation()
			runExamples = append(runExamples, example)
		} else {
			panic("存在非自如/链家房源搜索地址，请检查～")
		}
	}

	// 消息通知平台
	var noticePlatform core.NoticeService = nil
	if *notice == "ding" {
		noticePlatform = &notice2.DingImpl{}
	} else if *notice == "fs" {
		noticePlatform = &notice2.FeishuImpl{}
	}

	if noticePlatform == nil {
		panic("通知平台暂未支持，请提Issues～")
	}

	if len(*noticeUrl) == 0 {
		panic(noticePlatform.GetName() + "消息通知平台地址未设置。")
	}

	if len(*noticeKey) == 0 {
		panic(noticePlatform.GetName() + "消息通知平台密钥未设置。")
	}

	if *taskInterval <= 30 {
		*taskInterval = 30
	}

	if len(runExamples) <= 0 {
		panic("请至少输入一个平台的搜索地址...")
	}

	pkg.BeginToInspect(runExamples, noticePlatform, time.Duration(*taskInterval), *noticeUrl, *noticeKey)

}
