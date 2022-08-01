package v1

import "time"

// LogEvent LogEvent struct
type LogEvent struct {
	ProcessId string    `bson:"process_id" json:"process_id"`
	Log       string    `bson:"log" json:"log"`
	Step      string    `bson:"step" json:"step"`
	Footmark  string    `bson:"footmark" json:"footmark"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	Claim     int       `bson:"claim" json:"claim"`
}
