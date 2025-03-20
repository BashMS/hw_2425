package memorystorage

import (
	"context"
	"sync"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	Log logger.Logger
	Events map[int64]storage.Event
	Users  map[int64]storage.User
	mu sync.RWMutex //nolint:unused
}

// New mongodb://username:password@host:port/database?options
func New(cfg config.Config, logg *logger.Logger) *Storage {
	return &Storage{
		Events: make(map[int64]storage.Event),
		Users: make(map[int64]storage.User),
		Log: *logg,
	}
}

// Open 
func (s *Storage) Open(ctx context.Context) error {
	
	return nil
}

