package v1

import "time"

type LogEvent struct {
	ProcessId string    `bson:"process_id"`
	Log       string    `bson:"log"`
	Step      string    `bson:"step"`
	CreatedAt time.Time `bson:"created_at"`
}

