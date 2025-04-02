package app

import (
	"context"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type App struct {
	logg logger.Logger
	stor Storage
}

type Storage interface {
	Open(ctx context.Context) error
	Close(ctx context.Context) error

	CreateUser(ctx context.Context, user storage.User) (int64, error)
	UpdateUser(ctx context.Context, user storage.User) error
	DeleteUser(ctx context.Context, userID int64) error

	CreateEvent(ctx context.Context, evt storage.Event) (int64, error)
	UpdateEvent(ctx context.Context, evt storage.Event) error
	DeleteEvent(ctx context.Context, evtID int64) error
	ListEventsForDay(ctx context.Context, startDay time.Time) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error)
}

func New(logger *logger.Logger, storage Storage) *App {
	return &App{
		logg: *logger,
		stor: storage,
	}
}

func (a *App) Open(ctx context.Context) error {
	return a.stor.Open(ctx)
}

func (a *App) Close(ctx context.Context) error {
	return a.stor.Close(ctx)
}

func (a *App) CreateEvent(ctx context.Context, evt storage.Event) (int64, error) {
	// TODO
	return a.stor.CreateEvent(ctx, evt)
}

func (a *App) UpdateEvent(ctx context.Context, evt storage.Event) error {
	// TODO
	return a.stor.UpdateEvent(ctx, evt)
}

func (a *App) DeleteEvent(ctx context.Context, evtID int64) error {
	// TODO
	return a.stor.DeleteEvent(ctx, evtID)
}

func getStartDay(in *time.Time) time.Time {
	out := time.Date(in.Year(), in.Month(), in.Day(), 0, 0, 0, 0, in.Location())

	return out
}

func (a *App) ListEventsForDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	inDate := getStartDay(&startDay)

	return a.stor.ListEventsForDay(ctx, inDate)
}

func (a *App) ListEventsForWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	inDate := getStartDay(&startDay)

	return a.stor.ListEventsForWeek(ctx, inDate)
}

func (a *App) ListEventsForMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	inDate := getStartDay(&startDay)
	return a.stor.ListEventsForMonth(ctx, inDate)
}

func (a *App) CreateUser(ctx context.Context, user storage.User) (int64, error) {
	// TODO
	return a.stor.CreateUser(ctx, user)
}

func (a *App) UpdateUser(ctx context.Context, user storage.User) error {
	// TODO
	return a.stor.UpdateUser(ctx, user)
}

func (a *App) DeleteUser(ctx context.Context, userID int64) error {
	// TODO
	return a.stor.DeleteUser(ctx, userID)
}
