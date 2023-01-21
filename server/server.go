package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jsirianni/server/store"
	"go.uber.org/zap"
)

const (
	// DefaultTimeout is the default read and write timeout.
	DefaultTimeout = time.Second * 15
)

// Option is a function that configures a Server option.
type Option func(*Server) error

// WithBindAddress configures the bind address. Returns an error if
// the address or port are invalid.
func WithBindAddress(address string, port uint) Option {
	return func(s *Server) error {
		// Allow empty ip address. When not set, the http server
		// will bind to all addresses.
		if address != "" && net.ParseIP(address) == nil {
			return fmt.Errorf("failed to parse '%s' as an IP address", address)
		}

		if port > 65535 {
			return fmt.Errorf("invalid TCP port passed, '%d'. port must be between 0 and 65535", port)
		}

		addr := net.JoinHostPort(address, strconv.Itoa(int(port)))
		s.server.Addr = addr

		return nil
	}
}

// WithMemoryStore configures the store interface with
// an in memory storage backend. If seed is true, the memory
// store will be seeded with test data.
func WithMemoryStore(seed bool) Option {
	return func(s *Server) error {
		if seed {
			s.store = store.NewTestingMemory()
		} else {
			s.store = store.NewMemory()
		}
		return nil
	}
}

// New takes one or more Option functions and returns a Server
// configured with those options. Returns an error if any errors
// are encountered.
func New(logger *zap.Logger, ops ...Option) (*Server, error) {
	if logger == nil {
		return nil, errors.New("a zap logger is required")
	}

	// Gin runs in debug mode by default, but we always want
	// release unless otherwise specified.
	switch os.Getenv("GIN_MODE") {
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "release", "":
		gin.SetMode(gin.ReleaseMode)
	}

	s := &Server{
		logger: logger,
	}

	// TODO(jsirianni): Add timeout options to option functions
	// in order to allow the user to override default values.
	s.server.ReadTimeout = DefaultTimeout
	s.server.ReadHeaderTimeout = DefaultTimeout
	s.server.WriteTimeout = DefaultTimeout

	s.Router = gin.New()
	s.Router.Use(ginzap.Ginzap(logger, "", false))
	s.Router.Use(ginzap.RecoveryWithZap(logger, true))

	for _, op := range ops {
		if err := op(s); err != nil {
			return nil, err
		}
	}

	if s.store == nil {
		return nil, fmt.Errorf("server must be configured with a storage backend")
	}

	return s, nil
}

// Server is an http server for serving the http api.
type Server struct {
	// Router is the handler for the embedded http.Server's
	// Handler interface. Can be configured using
	// Option functions.
	Router *gin.Engine

	logger *zap.Logger
	server http.Server
	store  store.Store
}

// Start starts the server with net/http's ListenAndServer
// method. Runtime errors are returned and should be handled
// by the caller.
func (s *Server) Start() error {
	s.addRoutes()
	s.server.Handler = s.Router
	return s.server.ListenAndServe()
}

// Stop gracefully stops the server with a timeout. Active connections
// will be allowed to finish their requests within the configured timeout
// duration.
func (s *Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop server: %s", err)
	}

	return nil
}

// Addr returns the server's address with the form
// 'host:port'.
func (s *Server) Addr() string {
	return s.server.Addr
}

func (s *Server) addRoutes() {
	s.Router.GET("/health", healthHandler)

	// /v1/accounts requests
	v1 := s.Router.Group("/v1/accounts")
	v1.POST(":account/validate", s.checkSubscriptionHandler)
	v1.GET(":account", s.accountHandler)
	v1.GET(":account/devices", s.devicesHandler)
	v1.GET(":account/devices/:device", s.deviceHandler)
	v1.PUT(":account/device", s.registerDeviceHandler)
}
