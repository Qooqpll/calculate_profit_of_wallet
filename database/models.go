package database

import "time"

type ApiKeys struct {
	ID     uint
	ApiKey string
	Time   time.Time
}
