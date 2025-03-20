package storage

import (
	"time"
)

type Event struct {
	ID          int64
	Name        string
	StartDate   time.Time
	EndDate     time.Time
	UserID      int64
	Description string
	RemindFor   int
}
