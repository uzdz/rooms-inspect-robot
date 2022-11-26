package notice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"ziroom/internal/pkg/core"
)

var dingTemplate = "## [%Key 新房源提醒] %Platform %Title \n\n" +
	"%Image \n\n" +
	"PC链接：[房源地址](%Url) \n\n" +
	"手机链接：[房源地址](%MUrl) \n\n" +
	"详情：%Desc \n\n" +
	"标签：%Tag \n\n"

type DingImpl struct {
	Name string `default:"钉钉"`
}

func (ding *DingImpl) GetName() string {
	return ding.Name
}

func (ding *DingImpl) Send(room core.Room, url, key string) bool {

	if room.Title == "" {
		return false
	}

	sendTemplate := dingTemplate

	sendTemplate = strings.Replace(sendTemplate, "%Platform", room.Platform, -1)
	sendTemplate = strings.Replace(sendTemplate, "%Key", key, -1)
	sendTemplate = strings.Replace(sendTemplate, "%Title", room.Title, -1)

	if room.Url != "" {
		sendTemplate = strings.Replace(sendTemplate, "%Url", room.Url, -1)
	}

	if room.MUrl != "" {
		sendTemplate = strings.Replace(sendTemplate, "%MUrl", room.MUrl, -1)
	} else {
		sendTemplate = strings.Replace(sendTemplate, "%MUrl", room.Url, -1)
	}

	if room.Image != "" {
		sendTemplate = strings.Replace(sendTemplate, "%Image", "![image]("+room.Image+")", -1)
	}

	for i := 0; i < len(room.Desc); i++ {
		room.Desc[i] = strings.ReplaceAll(room.Desc[i], " ", "")
		room.Desc[i] = strings.ReplaceAll(room.Desc[i], "\t", "")
		room.Desc[i] = strings.ReplaceAll(room.Desc[i], "\n", "")
	}

	if room.Desc != nil && len(room.Desc) > 0 {
		sendTemplate = strings.Replace(sendTemplate, "%Desc", strings.Join(room.Desc, "、"), -1)
	}

	for i := 0; i < len(room.Tag); i++ {
		room.Tag[i] = strings.ReplaceAll(room.Tag[i], " ", "")
		room.Tag[i] = strings.ReplaceAll(room.Tag[i], "\t", "")
		room.Tag[i] = strings.ReplaceAll(room.Tag[i], "\n", "")
	}

	if room.Tag != nil && len(room.Tag) > 0 {
		sendTemplate = strings.Replace(sendTemplate, "%Tag", strings.Join(room.Tag, "、"), -1)
	}

	content, data := make(map[string]string), make(map[string]interface{})
	content["title"] = "新房源提醒"
	content["text"] = sendTemplate
	data["msgtype"] = "markdown"
	data["markdown"] = content

	b, _ := json.Marshal(data)

	resp, err := http.Post(url,
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(string(body))
		return false
	}
	return true
}
