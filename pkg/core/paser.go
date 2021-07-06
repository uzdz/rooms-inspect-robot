package core

import (
	"github.com/PuerkitoBio/goquery"
)

func PerRoomInfo(dom *goquery.Document) []Room {
	rooms := make([]Room, 0, 10)

	dom.Find("div[class=Z_list-box]>div[class=item]").Each(func(i int, s *goquery.Selection) {
		perRoom := Room{}

		picContent := s.Find("div[class=pic-box]>a").First()

		// 打开链接
		openUrl, openUrlExists := picContent.Attr("href")
		if openUrlExists {
			perRoom.Url = "https:" + openUrl
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
