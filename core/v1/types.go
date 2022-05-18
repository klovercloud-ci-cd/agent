package v1

import "github.com/klovercloud-ci-cd/agent/enums"

type KafkaMessage struct {
	Body   interface{}        `json:"body"`
	Header KafkaMessageHeader `json:"header"`
}

type KafkaMessageHeader struct {
	Offset  int               `json:"offset"`
	Command enums.Command     `json:"command"`
	Extras  map[string]string `json:"extras"`
}
