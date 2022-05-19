package v1

import "github.com/klovercloud-ci-cd/agent/enums"

type KubeEventMessage struct {
	Body   interface{}   `json:"body"`
	Header MessageHeader `json:"header"`
}

type MessageHeader struct {
	Offset  int               `json:"offset"`
	Command enums.Command     `json:"command"`
	Extras  map[string]string `json:"extras"`
}
