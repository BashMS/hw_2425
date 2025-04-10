package memorystorage

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Storage struct {
	Log    logger.Logger
	Events map[int64]storage.Event
	Users  map[int64]storage.User
	mu     sync.RWMutex
}

// New.
func New(_ config.Config, logg *logger.Logger) *Storage {
	return &Storage{
		Events: make(map[int64]storage.Event),
		Users:  make(map[int64]storage.User),
		Log:    *logg,
	}
}

func (s *Storage) Open(_ context.Context) error {
	s.Log.Info("Opening Memory storage...")
	return nil
}

func (s *Storage) Close(_ context.Context) error {
	s.Log.Info("Close Memory storage...")
	return nil
}

func (s *Storage) checkUser(userAddr string) bool {
	for _, user := range s.Users {
		if userAddr == user.Address {
			return true
		}
	}
	return false
}

func (s *Storage) checkEvent(evt storage.Event) bool {
	for _, event := range s.Events {
		if evt.StartDate == event.StartDate && evt.UserID == event.UserID && evt.ID != event.ID {
			return true
		}
	}
	return false
}

func (s *Storage) CreateUser(_ context.Context, user storage.User) (int64, error) {
	if s.checkUser(strings.ToLower(user.Address)) {
		return 0, fmt.Errorf("CreateUser: %w", storage.ErrUserExists)
	}

	s.mu.Lock()
	id := int64(len(s.Users) + 1)
	user.ID = id
	s.Users[id] = user
	s.mu.Unlock()

	return id, nil
}

func (s *Storage) UpdateUser(_ context.Context, user storage.User) error {
	if _, ok := s.Users[user.ID]; !ok {
		return fmt.Errorf("UpdateUser: %w", storage.ErrUserNotExists)
	}

	s.mu.Lock()
	s.Users[user.ID] = user
	s.mu.Unlock()

	return nil
}

func (s *Storage) DeleteUser(_ context.Context, userID int64) error {
	if _, ok := s.Users[userID]; !ok {
		return fmt.Errorf("DeleteUser: %w", storage.ErrUserNotExists)
	}

	s.mu.Lock()
	delete(s.Users, userID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) CreateEvent(_ context.Context, evt storage.Event) (int64, error) {
	if s.checkEvent(evt) {
		return 0, fmt.Errorf("CreateEvent: %w", storage.ErrDateBusy)
	}

	s.mu.Lock()
	id := int64(len(s.Events) + 1)
	evt.ID = id
	s.Events[id] = evt
	s.mu.Unlock()

	return id, nil
}

func (s *Storage) UpdateEvent(_ context.Context, evt storage.Event) error {
	if _, ok := s.Events[evt.ID]; !ok {
		return fmt.Errorf("UpdateEvent: %w", storage.ErrEventNotExists)
	}
	if s.checkEvent(evt) {
		return fmt.Errorf("UpdateEvent: %w", storage.ErrDateBusy)
	}

	s.mu.Lock()
	s.Events[evt.ID] = evt
	s.mu.Unlock()

	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, evtID int64) error {
	if _, ok := s.Events[evtID]; !ok {
		return fmt.Errorf("DeleteEvent: %w", storage.ErrEventNotExists)
	}

	s.mu.Lock()
	delete(s.Events, evtID)
	s.mu.Unlock()

	return nil
}

// ListEventsForDay возвращает события за указанный день.
func (s *Storage) ListEventsForDay(_ context.Context, startDay time.Time) ([]storage.Event, error) {
	var result []storage.Event
	endDay := startDay.Add(24 * time.Hour)

	s.mu.Lock()
	for _, evt := range s.Events {
		if (evt.StartDate.Equal(startDay) || evt.StartDate.After(startDay)) && evt.StartDate.Before(endDay) {
			result = append(result, evt)
		}
	}
	s.mu.Unlock()

	return result, nil
}

// ListEventsForWeek возвращает события за неделю с указанной даты.
func (s *Storage) ListEventsForWeek(_ context.Context, startDay time.Time) ([]storage.Event, error) {
	var result []storage.Event
	endDay := startDay.Add(7 * 24 * time.Hour)

	s.mu.Lock()
	for _, evt := range s.Events {
		if (evt.StartDate.Equal(startDay) || evt.StartDate.After(startDay)) && evt.StartDate.Before(endDay) {
			result = append(result, evt)
		}
	}
	s.mu.Unlock()

	return result, nil
}

// ListEventsForMonth возвращает события за месяц с указанной даты.
func (s *Storage) ListEventsForMonth(_ context.Context, startDay time.Time) ([]storage.Event, error) {
	var result []storage.Event
	endDay := startDay.Add(31 * 24 * time.Hour)

	s.mu.Lock()
	for _, evt := range s.Events {
		if (evt.StartDate.Equal(startDay) || evt.StartDate.After(startDay)) && evt.StartDate.Before(endDay) {
			result = append(result, evt)
		}
	}
	s.mu.Unlock()

	return result, nil
}
