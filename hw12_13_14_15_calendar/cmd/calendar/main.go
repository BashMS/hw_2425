package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	config "github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config"                //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	internalhttp "github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/server/http"     //nolint:depguard
	memorystorage "github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	sqlstorage "github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/storage/sql"       //nolint:depguard
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.txt", "Path to configuration file")
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) > 0 && args[0] == "version" {
		printVersion()
		return
	}

	if len(configFile) == 0 {
		panic("No configuration file specified")
	}

	// получим конфигурацию приложения
	cfg := config.NewConfig(configFile)
	logg := logger.New(cfg.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	// инициализируем хранилище
	var storage app.Storage
	switch cfg.Source {
	case "postgres":
		storage = sqlstorage.New(cfg, logg)
	case "memory":
		storage = memorystorage.New(cfg, logg)
	default:
		panic("No storage configuration source specified")
	}

	calendar := app.New(logg, storage)
	if err := storage.Open(ctx); err != nil {
		panic(fmt.Sprintf("storage.Open: %s", err.Error()))
	}
	defer func() { storage.Close(ctx) }()

	server := internalhttp.NewServer(logg, cfg, calendar)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
