package memorystorage

import (
	"context"
	"sync"
	"fmt"

	//"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"
)

type Storage struct {
	Log logger.Logger
	dsn string
	client *mongo.Client
	mu sync.RWMutex //nolint:unused
}

// New mongodb://username:password@host:port/database?options
func New(cfg config.Config, logg *logger.Logger) *Storage {
	return &Storage{
		Log: *logg,
		dsn: fmt.Sprintf("mongodb://%s:%s@%s:%d/%s?authSource=admin", cfg.Mongo.User, cfg.Mongo.Pass, cfg.Mongo.Host, cfg.Mongo.Port, cfg.Mongo.Name),
	}
}

// Open 
func (s *Storage) Open(ctx context.Context) error {
	s.Log.Info("Opening mongoDB...")
	clientOptions := options.Client().ApplyURI(s.dsn)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		s.Log.Error("failed to load driver", "Error:", err)
	    return fmt.Errorf("mongo.Connect: %w", err)
	}

	// Пинг сервера для проверки соединения
	err = client.Ping(ctx, nil)
	if err != nil {
		s.Log.Error("failed to connect to mongoDB", "Error:", err)
        return fmt.Errorf("Ping: %w", err)
	}
	s.client = client
	s.Log.Info("Ping mongoDB OK...")

	//db := client.Database(*clientOptions.SRVServiceName)
	
	//collectionEvent := db.Collection("event")

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	s.Log.Info("Closing mongoDB...")
	if err := s.client.Disconnect(ctx); err != nil {
		s.Log.Error("failed to close mongoDB", "Error:", err)
		return fmt.Errorf("Disconnect: %w", err)
	}
	
	s.client = nil
	s.Log.Info("Close DB OK...")
	
	return nil
}
