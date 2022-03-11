package alert

import (
	// "fmt"
	// "log"

	// "github.com/go-resty/resty/v2"

	robot_url "spider/utils/config/robot_url"
	post "spider/utils/http"
)

func Alert(content string) {
	url := robot_url.Robot_url
	data := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]interface{}{
			"content":        content,
			"mentioned_list": "",
		}}
	post.Post(url, data)
}
