package storage

import (
	"time"
)

type Event struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	StartDate   time.Time `json:"startDate" db:"start_date_time"`
	EndDate     time.Time `json:"endDate" db:"end_date_time"`
	UserID      int64     `json:"userId" db:"user_id"`
	Description string    `json:"description" db:"description"`
	RemindFor   int       `json:"remindFor" db:"remind_for"`
}
