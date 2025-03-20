package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct { 
	Log logger.Logger
	dsn string
	timeOut time.Duration
	numOpenConns int
	connLifeTime time.Duration
	DB *sql.DB
}

func New(cfg config.Config, logg *logger.Logger) *Storage {
	return &Storage{
		Log: *logg,
		dsn: fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DB.User, cfg.DB.Pass, cfg.DB.Host, cfg.DB.Port, cfg.DB.Name),
		timeOut: time.Duration(time.Duration(cfg.DB.Timeout) * time.Second),
		numOpenConns: cfg.DB.NumOpenConns,
		connLifeTime: time.Duration(time.Duration(cfg.DB.ConnLifeTime) * time.Second),
	}
}

func (s *Storage) Open(ctx context.Context) error {
	s.Log.Info("Opening DB...")
	db, err := sql.Open("pgx", s.dsn) 
    if err != nil {
        s.Log.Error("failed to load driver", "Error:", err)
	    return fmt.Errorf("sql.Open: %w", err)
    }

	// (по умолчанию - 0, без ограничений)
    db.SetMaxOpenConns(s.numOpenConns)
    // Макс. число открытых неиспользуемых соединений
    //db.SetMaxIdleConns(n int)
    // Макс. время жизни одного подключения
    db.SetConnMaxLifetime(s.connLifeTime)
    // Макс. время ожидания подключения в пуле
    //db.SetConnMaxIdleTime(d time.Duration)

    err = db.PingContext(ctx)
    if err != nil {
		s.Log.Error("failed to connect to db", "Error:", err)
        return fmt.Errorf("PingContext: %w", err)
    }
	s.DB = db
	s.Log.Info("Ping DB OK...")
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.Log.Info("Closing DB...")
	if s.DB == nil {
		return nil
	}
	err := s.DB.Close()
	if err != nil {
		s.Log.Error("failed to close db", "Error:", err)
		return fmt.Errorf("db.Close: %w", err)
	}
	s.DB = nil
	s.Log.Info("Close DB OK...")
	return nil
}

func (s *Storage) CreateUser(ctx context.Context, user storage.User) (int64, error) {
	return 1, nil
}
func (s *Storage) UpdateUser(ctx context.Context, user storage.User) error {
	return nil
}
func (s *Storage) DeleteUser(ctx context.Context, userId int64) error {
	return nil
}

func (s *Storage) CreateEvent(ctx context.Context, evt storage.Event) (int64, error) {

	return 1, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, evt storage.Event) error {
	return nil
}
func (s *Storage) DeleteEvent(ctx context.Context, evtId int64) error {
	return nil
}
func (s *Storage) ListEventsForDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
func (s *Storage) ListEventsForWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}
func (s *Storage) ListEventsForMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	return []storage.Event{}, nil
}