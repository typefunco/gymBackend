package models

import (
	"encoding/json"
	"time"
)

type Session struct {
	Id             uint
	UserId         int
	TrainerId      int
	StartTime      time.Time
	EndTime        time.Time
	ScheduledTime  time.Time
	WorkoutDetails json.RawMessage
}
