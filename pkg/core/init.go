package core

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"ziroom/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

func FirstDataInit() {
	log.Println("进行首次数据初始化...")
	BeginInit(ZiroomURL)
	log.Println("首次数据初始化完成，等待任务调度...")
}

func BeginInit(url string) {
	// 验证请求URL是否合法，并返回请求URL模版
	patternUrl := utils.ParseZiRoomUrl(url)

	// 请求起始页，准备总页数数据
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("此次网络请求异常，等待下次调度～")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// 解析起始页内容
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Printf("此次起始页解析异常，等待下次调度～")
		return
	}

	// 返回请求总页数
	totalPage := page(dom)

	allRooms := make([]Room, 0, 10)

	for i := 1; i <= totalPage; i++ {
		requestUrl := strings.Replace(patternUrl, "#", strconv.Itoa(i), -1)

		rooms := perPage(requestUrl)

		if rooms != nil {
			allRooms = append(allRooms, rooms...)
		}

		// 防止给ziroom压力，设置每页请求间隔
		time.Sleep(time.Second * PageSecondInterval)
	}

	Calculation(allRooms)
}

func perPage(url string) []Room {

	resp, err := http.Get(url)
	if err != nil {
		// 此页此次循环不计算
		log.Printf("分页网络请求异常，放弃失败分页～")
		return nil
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Printf("DOM分析失败，放弃失败分页～")
		return nil
	}

	return PerRoomInfo(dom)
}
