package v1

import (
	"github.com/klovercloud-ci-cd/agent/enums"
	"time"
)

type ProcessLifeCycleEvent struct {
	ProcessId string            `bson:"process_id" json :"process_id"`
	Step string                 `bson:"step" json :"step"`
	Status enums.PROCESS_STATUS `bson:"status" json :"status"`
	Next []string               `bson:"next" json :"next"`
	Agent string                `bson:"agent" json :"agent"`
	CreatedAt time.Time         `bson:"created_at" json :"created_at"`
}