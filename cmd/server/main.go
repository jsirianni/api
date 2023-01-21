package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jsirianni/server/logging"
	"github.com/jsirianni/server/server"
	"go.uber.org/zap"
)

func main() {
	// Configure a Zap logger, which will be used for
	// unified logging between your business logic
	// logging and the Gin server's request logging.
	logger, err := logging.New(logging.DebugLevel)
	if err != nil {
		fmt.Printf("failed to configure logger: %s\n", err)
		os.Exit(1)
	}

	// Create the server with the logger and options.
	s, err := server.New(
		logger,
		server.WithBindAddress("", 8000),
		server.WithMemoryStore(true),
	)
	if err != nil {
		logger.Sugar().Errorf("failed to initialize server: %s", err)
		os.Exit(1)
	}

	// Configure a context which will be cancled by SIGINT and
	// SIGTERM signals. This will allow for graceful shutdown.
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		cancel()
	}()

	// Start the server within a goroutine.
	var startErr error
	go func() {
		defer cancel()
		logger.Sugar().Infof("starting server on %s", s.Addr())
		startErr = s.Start()
	}()

	// Pause here until the context is cancled. This can happen
	// via 'ctrl+c', 'kill <pid>', or 'systemctl stop <service name>'.
	<-ctx.Done()

	// If there was an error during startup or runtime, log it here.
	if startErr != nil {
		logger.Error("runtime error", zap.Error(startErr))
		os.Exit(1)
	}

	// Stop the server with a 60 second timeout to allow for existing
	// requests to finish.
	logger.Info("stopping server")
	if err := s.Stop(time.Second * 60); err != nil {
		logger.Error("failed to stop server gracefully", zap.Error(err))
		os.Exit(1)
	}

	// Exit with 0 if there were no runtime or shutdown errors.
	os.Exit(0)
}
