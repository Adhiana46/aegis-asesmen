package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	_defaultAddr            = ":80"
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	engine          *echo.Echo
	notify          chan error
	shutdownTimeout time.Duration
	addr            string

	// additionals props
	name    string
	version string
}

func New(
	opts ...Option,
) *Server {
	srv := Server{
		engine:          echo.New(),
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		addr:            _defaultAddr,
	}

	for _, opt := range opts {
		opt(&srv)
	}

	return &srv
}

func (s *Server) Start() {
	slog.Info(fmt.Sprintf("Running http server at %s", s.addr))

	go func() {
		s.notify <- s.engine.Start(s.addr)
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.engine.Shutdown(ctx)
}

func (s *Server) GetEngine() *echo.Echo {
	return s.engine
}

// Set Handlers/Routes
func (s *Server) SetHandlers(f func(e *echo.Echo) error) error {
	return f(s.engine)
}
