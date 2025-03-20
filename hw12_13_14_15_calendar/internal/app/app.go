package app

import (
	"context"
	"time"
	
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"
    sqlstorage "github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage/sql" //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type App struct {
	logg logger.Logger
	stor Storage
}

type Storage interface {
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

func New(logger *logger.Logger, storage sqlstorage.Storage) *App {
	return &App{
		logg: *logger,
		stor: &storage,
	}
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
func (a *App) ListEventsForDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	// TODO
	return a.stor.ListEventsForDay(ctx, startDay)
}
func (a *App) ListEventsForWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	// TODO
	return a.stor.ListEventsForWeek(ctx, startDay)
}
func (a *App) ListEventsForMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	// TODO
	return a.stor.ListEventsForMonth(ctx, startDay)
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