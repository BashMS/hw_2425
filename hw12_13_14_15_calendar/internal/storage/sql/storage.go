package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"
	_ "github.com/jackc/pgx/stdlib"
)

type Storage struct { 
	Log logger.Logger
	dsn string
	timeOut time.Duration
	numOpenConns int
	connLifeTime time.Duration
	DB *sql.DB
}

func New(cfg config.DBConf, logg *logger.Logger) *Storage {
	return &Storage{
		Log: *logg,
		dsn: fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.Name),
		timeOut: time.Duration(time.Duration(cfg.Timeout) * time.Second),
		numOpenConns: cfg.NumOpenConns,
		connLifeTime: time.Duration(time.Duration(cfg.ConnLifeTime) * time.Second),
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

func (s *Storage) Close() error {
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
