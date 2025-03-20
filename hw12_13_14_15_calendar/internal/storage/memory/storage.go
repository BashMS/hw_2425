package memorystorage

import (
	"context"
	"sync"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
)

type Storage struct {
	Log logger.Logger
	Events map[int64]storage.Event
	Users  map[int64]storage.User
	mu sync.RWMutex //nolint:unused
}

// New.
func New(cfg config.Config, logg *logger.Logger) *Storage {
	return &Storage{
		Events: make(map[int64]storage.Event),
		Users: make(map[int64]storage.User),
		Log: *logg,
	}
}

// Open.
func (s *Storage) Open(ctx context.Context) error {
	
	return nil
}

