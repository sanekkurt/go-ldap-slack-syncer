package types

import (
	"time"
)

const (
	UserEnable  UserAction = "enable"
	UserDisable UserAction = "disable"
)

type UserAction string

type User struct {
	SlackID string
	Mail    string
	Action  UserAction
	Date    time.Time
}
