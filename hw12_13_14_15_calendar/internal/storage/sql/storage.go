package sqlstorage

import (
	"context"
	"fmt"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"  //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage" //nolint:depguard
	_ "github.com/jackc/pgx/stdlib"                                     //nolint:depguard
	"github.com/jmoiron/sqlx"                                           //nolint:depguard
)

type Storage struct {
	Log          logger.Logger
	dsn          string
	timeOut      time.Duration
	numOpenConns int
	connLifeTime time.Duration
	maxIdleConns int
	DB           *sqlx.DB
}

func New(cfg config.Config, logg *logger.Logger) *Storage {
	return &Storage{
		Log: *logg,
		dsn: fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DB.User, cfg.DB.Pass, cfg.DB.Host,
			cfg.DB.Port, cfg.DB.Name),
		timeOut:      time.Duration(cfg.DB.Timeout) * time.Second,
		numOpenConns: cfg.DB.NumOpenConns,
		maxIdleConns: cfg.DB.MaxIdleConns,
		connLifeTime: time.Duration(cfg.DB.ConnLifeTime) * time.Second,
	}
}

func (s *Storage) Open(ctx context.Context) error {
	s.Log.Info("Opening DB...")
	db, err := sqlx.Open("pgx", s.dsn)
	if err != nil {
		s.Log.Error("failed to load driver", "Error:", err)
		return fmt.Errorf("sql.Open: %w", err)
	}

	// (по умолчанию - 0, без ограничений)
	db.SetMaxOpenConns(s.numOpenConns)
	// Макс. число открытых неиспользуемых соединений
	db.SetMaxIdleConns(s.maxIdleConns)
	// Макс. время жизни одного подключения
	db.SetConnMaxLifetime(s.connLifeTime)

	err = db.PingContext(ctx)
	if err != nil {
		s.Log.Error("failed to connect to db", "Error:", err)
		return fmt.Errorf("PingContext: %w", err)
	}
	s.DB = db
	s.Log.Info("Ping DB OK...")
	return nil
}

func (s *Storage) Close(_ context.Context) error {
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

// CreateUser создает запись пользователя.
func (s *Storage) CreateUser(ctx context.Context, user storage.User) (int64, error) {
	ext, err := s.checkUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("checkUser: %w", err)
	}
	if !ext {
		return 0, fmt.Errorf("CreateUser: %w", storage.ErrUserExists)
	}

	user.ID, err = s.genNewID(ctx)
	if err != nil {
		return 0, fmt.Errorf("genNewID: %w", err)
	}

	id, err := s.createUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("createUser: %w", err)
	}

	return id, nil
}

// UpdateUser обновляет данные пользователя.
func (s *Storage) UpdateUser(ctx context.Context, user storage.User) error {
	err := s.updateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("updateUser: %w", err)
	}
	return nil
}

func (s *Storage) DeleteUser(ctx context.Context, userID int64) error {
	err := s.deleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("deleteUser: %w", err)
	}
	return nil
}

// CreateEvent добавляет новое событие.
func (s *Storage) CreateEvent(ctx context.Context, evt storage.Event) (int64, error) {
	ext, err := s.checkEvent(ctx, evt)
	if err != nil {
		return 0, fmt.Errorf("checkEvent: %w", err)
	}
	if !ext {
		return 0, fmt.Errorf("CreateEvent: %w", storage.ErrDateBusy)
	}

	evt.ID, err = s.genNewID(ctx)
	if err != nil {
		return 0, fmt.Errorf("genNewID: %w", err)
	}

	_, err = s.createEvent(ctx, evt)
	if err != nil {
		return 0, fmt.Errorf("createEvent: %w", err)
	}

	return evt.ID, nil
}

// UpdateEvent обновляет событие.
func (s *Storage) UpdateEvent(ctx context.Context, evt storage.Event) error {
	ext, err := s.checkEvent(ctx, evt)
	if err != nil {
		return fmt.Errorf("checkEvent: %w", err)
	}
	if !ext {
		return fmt.Errorf("CreateEvent: %w", storage.ErrDateBusy)
	}

	err = s.updateEvent(ctx, evt)
	if err != nil {
		return fmt.Errorf("updateEvent: %w", err)
	}

	return nil
}

// DeleteEvent удаляет событие.
func (s *Storage) DeleteEvent(ctx context.Context, evtID int64) error {
	err := s.deleteEvent(ctx, evtID)
	if err != nil {
		return fmt.Errorf("deleteEvent: %w", err)
	}
	return nil
}

// ListEventsForDay возвращает события за указанный день.
func (s *Storage) ListEventsForDay(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	endDay := startDay.Add(24 * time.Hour)
	result, err := s.getListEventByPeriod(ctx, startDay, endDay)
	if err != nil {
		return nil, fmt.Errorf("getListEventByPeriod: %w", err)
	}

	return result, nil
}

// ListEventsForWeek возвращает события за неделю с указанного дня.
func (s *Storage) ListEventsForWeek(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	endDay := startDay.Add(7 * 24 * time.Hour)
	result, err := s.getListEventByPeriod(ctx, startDay, endDay)
	if err != nil {
		return nil, fmt.Errorf("getListEventByPeriod: %w", err)
	}

	return result, nil
}

// ListEventsForMonth возвращает события за месяц с указанного дня.
func (s *Storage) ListEventsForMonth(ctx context.Context, startDay time.Time) ([]storage.Event, error) {
	endDay := startDay.Add(31 * 24 * time.Hour)
	result, err := s.getListEventByPeriod(ctx, startDay, endDay)
	if err != nil {
		return nil, fmt.Errorf("getListEventByPeriod: %w", err)
	}

	return result, nil
}
