package core

import (
	"strconv"
	"strings"
	"ziroom/pkg/utils"

	"github.com/PuerkitoBio/goquery"
)

// 返回请求的总页数
func page(dom *goquery.Document) int {

	totalPage := "1"

	dom.Find("#page>span").Each(func(i int, s *goquery.Selection) {
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
