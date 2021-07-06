package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var template = "## [%Key 新房源提醒]  %Title \n\n" +
	"%Image \n\n" +
	"链接：%Url \n\n" +
	"详情：%Desc \n\n" +
	"标签：%Tag \n\n"

func DingToInfo(room Room) bool {
	sendTemplate := template

	sendTemplate = strings.Replace(sendTemplate, "%Key", WebhookUrlKey, -1)
	sendTemplate = strings.Replace(sendTemplate, "%Title", room.Title, -1)

	if room.Url != "" {
		sendTemplate = strings.Replace(sendTemplate, "%Url", room.Url, -1)
	}

	if room.Image != "" {
		sendTemplate = strings.Replace(sendTemplate, "%Image", "![image]("+room.Image+")", -1)
	}

	if room.Desc != nil && len(room.Desc) > 0 {
		sendTemplate = strings.Replace(sendTemplate, "%Desc", strings.Join(room.Desc, "、"), -1)
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

	resp, err := http.Post(WebhookUrl,
		"application/json",
		bytes.NewBuffer(b))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(string(body))
	}
	return true
}
