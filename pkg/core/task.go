package core

import (
	"log"
	"time"
)

func TaskForSearch() {
	log.Println("已开启定时任务，等待调度...")

	for {
		t := time.NewTimer(time.Second * TaskInterval)

		<-t.C

		log.Println("开始进行数据爬取...")
		BeginInit(ZiroomURL)
		log.Println("数据爬取结束，等待下次任务调度...")
	}
}
