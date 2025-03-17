package storage

import "time"

type Event struct {
	ID int64
    Name string
    StartDate time.Time
    EndDate time.Time
    UserId int64
    Description string
    RemindFor int
}

type EventRepo interface {
	CreateEvent(evt Event) (int64, error)
    UpdateEvent(evt Event) error;
    DeleteEvent(evtId int64) error;
    ListEventsForDay(startDay time.Time) ([]Event, error);
    ListEventsForWeek(startDay time.Time) ([]Event, error);
    ListEventsForMonth(startDay time.Time) ([]Event, error);
}