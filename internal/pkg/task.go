package pkg

import (
	"time"
)

func BeginToInspect(examples []AbilityService, taskInterval time.Duration, WebHookUrl, WebHookUrlKey string) {
	for {
		t := time.NewTimer(time.Second * 1)

		<-t.C

		for i := 0; i < len(examples); i++ {
			runSearchExample(examples[i], WebHookUrl, WebHookUrlKey)
		}
	}
}

/**
 * @Description: 运行每个房源搜索实例
 * @param example 搜索相关信息
 */
func runSearchExample(example AbilityService, WebHookUrl, WebHookUrlKey string) {

	totalPageNum := example.TotalPage()

	if totalPageNum == 0 {
		return
	}

	refreshRooms := example.ObtainRefreshRooms(totalPageNum)

	if refreshRooms == nil || len(refreshRooms) == 0 {
		return
	}

	valueNotifyRooms := example.Calculation(refreshRooms)

	// 新房源发送钉钉通知
	for i := 0; i < len(valueNotifyRooms); i++ {
		DingNotify(valueNotifyRooms[i], WebHookUrl, WebHookUrlKey)
	}
}
