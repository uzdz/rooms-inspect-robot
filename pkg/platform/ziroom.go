package platform

import (
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"ziroom/internal/pkg"
	"ziroom/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

type ZIRoomImpl struct {
	/**
	 * @Description: 缓存房间
	 */
	cacheRoom map[string]pkg.Room

	/**
	 * @Description: 搜索房源URL
	 */
	InputURL string

	/**
	 * @Description: 搜索房源请求URL模版
	 */
	requestTemplateURL string
}

func (platform ZIRoomImpl) GetPlatform() string {
	return "自如房源"
}

func (platform *ZIRoomImpl) Validation() {
	// 解析url 并保证没有错误
	curl, err := url.Parse(platform.InputURL)
	if err != nil {
		panic("请求URL错误（解析失败）！" + err.Error())
	}

	end := "/?isOpen=0"
	oParams := strings.Split(platform.InputURL, "?")
	if len(oParams) >= 2 {
		end = "/?" + oParams[1]
	}

	params := strings.Split(curl.Path, "/")

	// 需满足规则，方便制作模版
	if len(params) < 3 {
		panic("请求URL错误（格式）！" + platform.InputURL)
	}

	details := strings.Split(params[2], "-")

	if strings.Contains(details[len(details)-1], "p") {
		details[len(details)-1] = "p#"
	} else {
		details = append(details, "p#")
	}

	platform.requestTemplateURL = "https://" + curl.Host + "/" + params[1] + "/" + strings.Join(details, "-") + end
}

func (platform *ZIRoomImpl) TotalPage() int {
	resp, err := http.Get(platform.InputURL)
	if err != nil {
		log.Printf("此次网络请求异常，等待下次调度～")
		return 0
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// 解析起始页内容
	dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if err != nil {
		log.Printf("此次起始页解析异常，等待下次调度～")
		return 0
	}

	totalPage := "1"

	dom.Find(".page>span").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "共") && strings.Contains(s.Text(), "页") {
			totalPage = utils.Between(s.Text(), "共", "页")
		}
	})

	j, err := strconv.Atoi(totalPage)
	if err == nil {
		return j
	}

	return 1
}

func (platform *ZIRoomImpl) ObtainRefreshRooms(page int) []pkg.Room {

	allRooms := make([]pkg.Room, 0, 10)

	for i := 1; i <= page; i++ {
		nextRequestUrl := strings.Replace(platform.requestTemplateURL, "#", strconv.Itoa(i), -1)

		resp, err := http.Get(nextRequestUrl)
		if err != nil {
			// 此页此次循环不计算
			log.Printf("分页网络请求异常，放弃失败分页～")
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)

		dom, err := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
		if err != nil {
			log.Printf("DOM分析失败，放弃失败分页～")
			continue
		}

		allRooms = append(allRooms, platform.perRoomInfo(dom)...)
		resp.Body.Close()

		// 防止给自如找房压力，设置每页请求间隔
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	}

	return allRooms
}

func (platform *ZIRoomImpl) Calculation(refreshRooms []pkg.Room) []pkg.Room {

	if platform.cacheRoom == nil {
		platform.cacheRoom = make(map[string]pkg.Room)
	}

	notifyRooms := make([]pkg.Room, 0, 10)

	if refreshRooms == nil || len(refreshRooms) <= 0 {
		return notifyRooms
	}

	begin := false
	if platform.cacheRoom == nil || len(platform.cacheRoom) <= 0 {
		begin = true
	}

	for i := 0; i < len(refreshRooms); i++ {
		if begin == false {
			if _, ok := platform.cacheRoom[refreshRooms[i].Url]; ok {
				// 存在，不预警
			} else {
				refreshRooms[i].Platform = platform.GetPlatform()
				notifyRooms = append(notifyRooms, refreshRooms[i])
			}
		}

		platform.cacheRoom[refreshRooms[i].Url] = refreshRooms[i]
	}

	return notifyRooms
}

func (platform *ZIRoomImpl) perRoomInfo(dom *goquery.Document) []pkg.Room {
	rooms := make([]pkg.Room, 0, 10)

	dom.Find("div[class=Z_list]>div[class=Z_list-box]>div[class=item]").Each(func(i int, s *goquery.Selection) {
		perRoom := pkg.Room{}

		picContent := s.Find("div[class=pic-box]>a").First()

		// 打开链接
		openUrl, openUrlExists := picContent.Attr("href")
		if openUrlExists {
			perRoom.Url = "https:" + openUrl

			splitOfMobile := strings.Split(perRoom.Url, "/")
			if len(splitOfMobile) > 0 {
				numOfRoom := splitOfMobile[len(splitOfMobile)-1]
				if numOfRoom != "" {
					num := strings.Split(numOfRoom, ".")
					if len(num) > 0 {
						perRoom.MUrl = "https://m.ziroom.com/bj/x/999.html?inv_no=" + num[0]
					}
				}
			}
		}

		// 首张图片
		image, imageExists := picContent.Find("img").Attr("data-original")
		if imageExists {
			perRoom.Image = "https:" + image
		}

		infoContent := s.Find("div[class=info-box]")

		// 房间名称
		perRoom.Title = infoContent.Children().First().Children().First().Text()

		// 房间详情
		infoContent.Find("div[class=desc]>div").Each(func(i int, s *goquery.Selection) {
			perRoom.Desc = append(perRoom.Desc, s.Text())
		})

		infoContent.Find("div[class=tag]>span").Each(func(i int, s *goquery.Selection) {
			perRoom.Tag = append(perRoom.Tag, s.Text())
		})

		perRoom.Tag = append(perRoom.Tag, infoContent.Find("div[class=tip-info]").First().Text())

		//jsonBytes, err := json.Marshal(perRoom)
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//fmt.Println(string(jsonBytes))
		rooms = append(rooms, perRoom)

	})

	return rooms
}
