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

var fsTemplate = "{\"config\":{\"wide_screen_mode\":true},\"elements\":[{\"fields\":[{\"is_short\":true,\"text\":{\"content\":\"**平台：**\\n%Platform\",\"tag\":\"lark_md\"}},{\"is_short\":false,\"text\":{\"content\":\"\",\"tag\":\"lark_md\"}},{\"is_short\":true,\"text\":{\"content\":\"**详情：**\\n%Desc\",\"tag\":\"lark_md\"}},{\"is_short\":false,\"text\":{\"content\":\"\",\"tag\":\"lark_md\"}},{\"is_short\":true,\"text\":{\"content\":\"**标签：**\\n%Tag\",\"tag\":\"lark_md\"}}],\"tag\":\"div\"},{\"tag\":\"hr\"},{\"actions\":[{\"tag\":\"button\",\"text\":{\"content\":\"PC链接\",\"tag\":\"lark_md\"},\"url\":\"https://applink.feishu.cn/client/web_url/open?mode=sidebar-semi&url=%Url\",\"type\":\"primary\"},{\"tag\":\"button\",\"text\":{\"content\":\"Mobile链接\",\"tag\":\"lark_md\"},\"url\":\"https://applink.feishu.cn/client/web_url/open?mode=sidebar-semi&url=%MUrl\",\"type\":\"primary\"}],\"tag\":\"action\"}],\"header\":{\"template\":\"green\",\"title\":{\"content\":\"[%Key 新房源提醒] %Title\",\"tag\":\"plain_text\"}}}"

type FeishuImpl struct {

	Name string `default:"飞书"`

}

func (feishu *FeishuImpl) GetName() string {
	return feishu.Name
}

func (feishu *FeishuImpl) Send(room core.Room, url, key string) bool {

	if room.Title == "" {
		return false
	}

	sendTemplate := fsTemplate

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

	content, data := make(map[string]interface{}), make(map[string]interface{})
	err := json.Unmarshal([]byte(sendTemplate), &content)
	if err != nil {
		fmt.Println(err)
		return false
	}

	data["msg_type"] = "interactive"
	data["card"] = content

	b, _ := json.Marshal(data)

	resp, err := http.Post(url,
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
