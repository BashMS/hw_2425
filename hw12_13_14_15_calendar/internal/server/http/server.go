package internalhttp

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	"github.com/BashMS/hw_2425/hw12_13_14_15_calendar/internal/logger" //nolint:depguard
)

type Server struct {
	Port         string
	server       *http.Server
	mux          *http.ServeMux
	log          logger.Logger
	fileLog      *os.File
	readTimeout  time.Duration
	writeTimeout time.Duration
}

type MyHandler struct{}

var FileLog *log.Logger

func (h *MyHandler) Hello(w http.ResponseWriter, r *http.Request) { //nolint:revive
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

type Application interface { // TODO
}

func NewServer(logger *logger.Logger, cfg config.Config, app Application) *Server {
	_ = app // TODO
	f, err := os.OpenFile(cfg.Logger.LogFile, os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		logger.Error("Error create api log file", "Error:", err)
	}
	FileLog = log.New(f, "> ", 0)

	handler := &MyHandler{}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", loggingMiddleware(handler.Hello))

	return &Server{
		Port:         fmt.Sprintf(":%d", cfg.Server.Port),
		log:          *logger,
		mux:          mux,
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
