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
	"ziroom/internal/pkg/core"

	"github.com/PuerkitoBio/goquery"
)

type LianJiaImpl struct {
	/**
	 * @Description: 缓存房间
	 */
	cacheRoom map[string]core.Room

	/**
	 * @Description: 搜索房源URL
	 */
	InputURL string

	/**
	 * @Description: 搜索房源请求URL模版
	 */
	requestTemplateURL string
}

func (platform LianJiaImpl) GetPlatform() string {
	return "链家房源"
}

func (platform *LianJiaImpl) Validation() {
	platform.InputURL = strings.Replace(platform.InputURL, "/#contentList", "", -1)
	platform.InputURL = strings.Replace(platform.InputURL, "/?showMore=1", "", -1)

	if platform.InputURL[len(platform.InputURL)-1:] == "/" {
		platform.InputURL = platform.InputURL[0 : len(platform.InputURL)-1]
	}

	// 解析url 并保证没有错误
	curl, err := url.Parse(platform.InputURL)
	if err != nil {
		panic("请求URL错误（解析失败）！" + err.Error())
	}

	rURL := curl.Path
	if rURL[0:1] == "/" {
		rURL = rURL[1:]
	}

	params := strings.Split(rURL, "/")

	// 需满足规则，方便制作模版
	if len(params) < 2 || len(params) > 3 {
		panic("请求URL错误（请选择过滤条件）！" + platform.InputURL)
	} else if len(params) == 2 {
		params = append(params, "pg#")
	} else if len(params) == 3 {
		if strings.HasPrefix(params[2], "pg") {
			panic("链家房源，请输入首页检索地址，首页地址无'pg'参数～")
		}
		params[2] = "pg#" + params[2]
	}

	platform.requestTemplateURL = "https://" + curl.Host + "/" + strings.Join(params, "/")
}

func (platform *LianJiaImpl) TotalPage() int {
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

	totalPage, _ := dom.Find(".content__pg").First().Attr("data-totalpage")

	j, err := strconv.Atoi(totalPage)
	if err == nil {
		return j
	}

	return 1
}

func (platform *LianJiaImpl) ObtainRefreshRooms(page int) []core.Room {

	allRooms := make([]core.Room, 0, 10)

	for i := 1; i <= page; i++ {
		nextRequestUrl := strings.Replace(platform.requestTemplateURL, "#", strconv.Itoa(i), -1)

		req, _ := http.NewRequest("GET", nextRequestUrl, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36")
		resp, err := (&http.Client{}).Do(req)
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

		// 防止给链家找房压力，设置每页请求间隔
		time.Sleep(time.Second * time.Duration(rand.Intn(5)))
	}

	return allRooms
}

func (platform *LianJiaImpl) Calculation(refreshRooms []core.Room) []core.Room {

	if platform.cacheRoom == nil {
		platform.cacheRoom = make(map[string]core.Room)
	}

	notifyRooms := make([]core.Room, 0, 10)

	if refreshRooms != nil && len(refreshRooms) > 0 {
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
	}

	log.Printf("【链家】URL: %s, 缓存房源个数: %d\n", platform.InputURL, len(platform.cacheRoom))

	return notifyRooms
}

func (platform *LianJiaImpl) perRoomInfo(dom *goquery.Document) []core.Room {
	rooms := make([]core.Room, 0, 10)

	// 解析url 并保证没有错误
	curl, err := url.Parse(platform.InputURL)
	if err != nil {
		panic("请求URL错误（解析失败）！" + err.Error())
	}

	dom.Find("div[class=content__list]>div[class=content__list--item]").Each(func(i int, s *goquery.Selection) {
		perRoom := core.Room{}

		// 打开链接
		openUrl, openUrlExists := s.Find("a[class=content__list--item--aside]").Attr("href")
		if openUrlExists {
			perRoom.Url = "https:" + curl.Host + openUrl
		}

		// 首张图片
		image, imageExists := s.Find("a[class=content__list--item--aside]>img").First().Attr("data-src")
		if imageExists {
			perRoom.Image = image
		}

		titleText := s.Find("div[class=content__list--item--main]>p[class=content__list--item--title]>a").First().Text()

		// 房间标题
		perRoom.Title = strings.Replace(titleText, " ", "", -1)
		perRoom.Title = strings.Replace(perRoom.Title, "\n", "", -1)

		// 房间详情
		s.Find("div[class=content__list--item--main]>p[class=content__list--item--des]").Each(func(i int, s *goquery.Selection) {
			text := strings.Replace(s.Text(), " ", "", -1)
			text = strings.Replace(text, "\n", " ", -1)

			split := strings.Split(text, "/")

			for i := 0; i < len(split); i++ {
				perRoom.Desc = append(perRoom.Desc, split[i])
			}
		})

		s.Find("div[class=content__list--item--main]>p[class='content__list--item--bottom oneline']>i").Each(func(i int, s *goquery.Selection) {
			perRoom.Tag = append(perRoom.Tag, s.Text())
		})

		perRoom.Tag = append(perRoom.Tag, s.Find("div[class=content__list--item--main]>span").First().Text())

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
