package main

import (
	"strings"
	"time"
	"ziroom/internal/pkg"
	"ziroom/pkg/platform"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	dingUrl      = kingpin.Flag("dingUrl", "钉钉消息通知接口地址").Short('d').String()
	dingKey      = kingpin.Flag("dingKey", "钉钉消息通知授权KEY（白名单）").Short('k').Default("推送").String()
	taskInterval = kingpin.Flag("taskInterval", "任务周期间隔时长，单位：秒").Short('p').Default("300").Int()
	url          = kingpin.Arg("url", "自如/链家网页版房源请求地址。").Strings()

	//ziroomCommand    = app.Command("ziroom", "请输入自如房源地址，房源搜索地址参考：https://www.ziroom.com/z/，多个地址通过空格分割。")
	//examplesOfZiroom = ziroomCommand.Arg("examplesOfZiroom", "URLS").Required().Strings()
	//
	//lianjiaCommand    = app.Command("lianjia", "请输入链家房源地址，通过空格分离。")
	//examplesOfLianjia = lianjiaCommand.Arg("examplesOfLianjia", "URLS").Required().Strings()
)

func main() {
	kingpin.Parse()
	examples := *url

	if examples == nil {
		panic("请设置自如/链家网页版房源请求地址。")
	}

	runExamples := make([]pkg.AbilityService, 0, 10)

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

	if len(*dingUrl) == 0 {
		panic("钉钉通知未设置。")
	}

	if len(*dingKey) == 0 {
		panic("钉钉密钥未设置。")
	}

	if *taskInterval <= 30 {
		*taskInterval = 30
	}

	if len(runExamples) <= 0 {
		panic("请至少输入一个平台的搜索地址...")
	}

	pkg.BeginToInspect(runExamples, time.Duration(*taskInterval), *dingUrl, *dingKey)

}
