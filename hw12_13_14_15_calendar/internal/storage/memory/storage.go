package memorystorage

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

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
func New(logg *logger.Logger) *Storage {
	return &Storage{
		Events: make(map[int64]storage.Event),
		Users:  make(map[int64]storage.User),
		Log:    *logg,
	}
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

func (s *Storage) CreateUser(ctx context.Context, user storage.User) (int64, error) { //nolint:revive
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

func (s *Storage) UpdateUser(ctx context.Context, user storage.User) error { //nolint:revive
	if _, ok := s.Users[user.ID]; !ok {
		return fmt.Errorf("UpdateUser: %w", storage.ErrUserNotExists)
	}

	s.mu.Lock()
	s.Users[user.ID] = user
	s.mu.Unlock()

	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int64) error { //nolint:revive
	if _, ok := s.Users[userID]; !ok {
		return fmt.Errorf("DeleteUser: %w", storage.ErrUserNotExists)
	}

	s.mu.Lock()
	delete(s.Users, userID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, evt storage.Event) (int64, error) { //nolint:revive
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

func (s *Storage) UpdateEvent(ctx context.Context, evt storage.Event) error { //nolint:revive
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

func (s *Storage) DeleteEvent(ctx context.Context, evtID int64) error { //nolint:revive
	if _, ok := s.Events[evtID]; !ok {
		return fmt.Errorf("DeleteEvent: %w", storage.ErrEventNotExists)
	}

	s.mu.Lock()
	delete(s.Events, evtID)
	s.mu.Unlock()

	return nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) { //nolint:revive
	// TODO
	return []storage.Event{}, nil
}

func (s *Storage) ListEventsForWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) { //nolint:revive
	// TODO
	return []storage.Event{}, nil
}

func (s *Storage) ListEventsForMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) { //nolint:revive
	// TODO
	return []storage.Event{}, nil
}
