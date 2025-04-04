package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/app"    //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
	"github.com/gorilla/mux"                                           //nolint:depguard
)

type Server struct {
	Port         string
	server       *http.Server
	mux          *mux.Router
	log          logger.Logger
	fileLog      *os.File
	readTimeout  time.Duration
	writeTimeout time.Duration
}

var FileLog *log.Logger

func NewServer(logger *logger.Logger, cfg config.Config, app app.App) *Server {
	f, err := os.OpenFile(cfg.Logger.LogFile, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		logger.Error("Error create api log file", "Error:", err)
	}
	FileLog = log.New(f, "> ", 0)

	handler := &MyHandler{
		App: app,
	}

	gmux := mux.NewRouter()
	gmux.HandleFunc("/hello", loggingMiddleware(handler.Hello)).Methods(http.MethodGet)
	gmux.HandleFunc("/users", loggingMiddleware(handler.CreateUser)).Methods(http.MethodPost)
	gmux.HandleFunc("/users/{userID}", loggingMiddleware(handler.UpdateUser)).Methods(http.MethodPut)
	gmux.HandleFunc("/users/{userID}", loggingMiddleware(handler.DeleteUser)).Methods(http.MethodDelete)
	gmux.HandleFunc("/events", loggingMiddleware(handler.CreateEvent)).Methods(http.MethodPost)
	gmux.HandleFunc("/events/{eventID}", loggingMiddleware(handler.UpdateEvent)).Methods(http.MethodPut)
	gmux.HandleFunc("/events/{eventID}", loggingMiddleware(handler.DeleteEvent)).Methods(http.MethodDelete)
	gmux.HandleFunc("/events/listByDay", loggingMiddleware(handler.ListEventsForDay)).Methods(http.MethodGet)
	gmux.HandleFunc("/events/listByWeek", loggingMiddleware(handler.ListEventsForWeek)).Methods(http.MethodGet)
	gmux.HandleFunc("/events/listByMonth", loggingMiddleware(handler.ListEventsForMonth)).Methods(http.MethodGet)

	return &Server{
		Port:         fmt.Sprintf(":%d", cfg.Server.Port),
		log:          *logger,
		mux:          gmux,
		fileLog:      f,
		readTimeout:  time.Duration(cfg.Server.ReadTimeOut) * time.Second,
		writeTimeout: time.Duration(cfg.Server.WriteTimeOut) * time.Second,
	}
}

func (s *Server) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:         s.Port,
		Handler:      s.mux,
		ReadTimeout:  s.readTimeout,
		WriteTimeout: s.writeTimeout,
	}
	s.server = server
	s.log.Info("Server is starting... ", "addr", s.server.Addr)
	defer s.fileLog.Close()

	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("server.ListenAndServe: %w", err)
	}
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("server.Shutdown: %w", err)
	}
	s.log.Info("Server is Shutdown... ")
	return nil
}

// TODO
