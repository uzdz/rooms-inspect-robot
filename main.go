package main

import (
	"os"
	"time"
	"ziroom/internal/pkg"
	"ziroom/pkg/platform"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app          = kingpin.New("robot", "自动化获取【自如】/【链家】新房源机器人🤖️")
	dingUrl      = app.Flag("dingUrl", "钉钉消息通知接口地址").Short('d').String()
	dingKey      = app.Flag("dingKey", "钉钉消息通知授权KEY（白名单）").Short('k').Default("推送").String()
	taskInterval = app.Flag("taskInterval", "任务周期间隔时长（最少5分钟），单位：秒").Short('p').Default("300").Int()

	ziroomCommand    = app.Command("ziroom", "请输入自如房源地址，房源搜索地址参考：https://www.ziroom.com/z/，多个地址通过空格分割。")
	examplesOfZiroom = ziroomCommand.Arg("examplesOfZiroom", "URLS").Required().Strings()

	lianjiaCommand    = app.Command("lianjia", "请输入链家房源地址，通过空格分离。")
	examplesOfLianjia = lianjiaCommand.Arg("examplesOfLianjia", "URLS").Required().Strings()
)

func main() {

	runExamples := make([]pkg.AbilityService, 0, 10)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case ziroomCommand.FullCommand():
		if examplesOfZiroom != nil {
			zm := *examplesOfZiroom

			for i := 0; i < len(zm); i++ {
				example := &platform.ZIRoomImpl{
					InputURL: zm[i],
				}

				// 生成请求模版
				example.Validation()
				runExamples = append(runExamples, example)
			}
		}

	case lianjiaCommand.FullCommand():
		//lj := *examplesOfLianjia
		//
		//for i := 0; i < len(lj); i++ {
		//	example := &platform.ZIRoomImpl{
		//		InputURL: lj[i],
		//	}
		//	runExamples = append(runExamples, example)
		//}
	}

	if len(*dingUrl) == 0 {
		panic("钉钉通知未设置。")
	}

	if len(*dingKey) == 0 {
		panic("钉钉密钥未设置。")
	}

	if *taskInterval <= 300 {
		*taskInterval = 300
	}

	if len(runExamples) <= 0 {
		panic("请至少输入一个平台的搜索地址...")
	}

	pkg.BeginToSearch(runExamples, time.Duration(*taskInterval), *dingUrl, *dingKey)
}
