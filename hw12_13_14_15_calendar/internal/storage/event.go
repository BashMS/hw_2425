package storage

import (
	"time"
)

type Event struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	StartDate   time.Time `db:"start_date_time"`
	EndDate     time.Time `db:"end_date_time"`
	UserID      int64     `db:"user_id"`
	Description string    `db:"description"`
	RemindFor   int       `db:"remind_for"`
}
