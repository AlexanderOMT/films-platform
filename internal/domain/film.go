package domain

import (
	"time"
)

type Film struct {
	Id       int
	User_id  int
	Title    string
	Director string
	Release  time.Time
}
