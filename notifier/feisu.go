package notifier

import (
	"encoding/json"
	"fmt"
)

type feishuPayload struct {
	MsgType string               `json:"msg_type"`
	Content feishuPayloadContent `json:"content"`
}

type feishuPayloadContent struct {
	Text string `json:"text"`
}

func NewFeishu(base *Base) *Webhook {
	return &Webhook{
		Base:        *base,
		Service:     "feishu",
		method:      "POST",
		contentType: "application/json",
		buildBody: func(title, message string) ([]byte, error) {
			payload := feishuPayload{
				MsgType: "text",
				Content: feishuPayloadContent{
					Text: fmt.Sprintf("%s\n\n%s", title, message),
				},
			}

			return json.Marshal(payload)
		},
	}
}
