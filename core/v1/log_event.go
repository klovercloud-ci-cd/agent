package v1

import "time"

// LogEvent LogEvent struct
type LogEvent struct {
	ProcessId string    `bson:"process_id"`
	Log       string    `bson:"log"`
	Step      string    `bson:"step"`
	Footmark  string    `bson:"footmark"`
	CreatedAt time.Time `bson:"created_at"`
	Claim     int       `bson:"claim" json:"claim"`
}
