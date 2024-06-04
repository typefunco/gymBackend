package models

import (
	"time"
)

type Notifications struct {
	Id         uint
	UserId     int
	SessionId  int
	NotifyTime time.Time
	Message    string
	Sent       bool
}
