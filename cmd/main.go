package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vlad161/multiplexer/internal/http/handler"
	"github.com/vlad161/multiplexer/internal/http/middleware"
	"github.com/vlad161/multiplexer/internal/logger"
	"github.com/vlad161/multiplexer/internal/service/multiplexer"
)

const (
	port                 = 8081
	multiplexerRateLimit = 100
)

func main() {
	ctx := context.Background()
	log := logger.New()

	mp, errMultiplexer := multiplexer.New(multiplexer.WithDefaultClient())
	if errMultiplexer != nil {
		log.Fatal("can't create multiplexer service: %v", errMultiplexer)
	}

	mux := http.NewServeMux()
	mux.Handle("/multiplexer", middleware.RateLimiter(multiplexerRateLimit, handler.NewMultiplexer(mp)))

	s := http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        mux,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	s.Handler = middleware.Logger(log, s.Handler)

	go func() {
		if errListen := s.ListenAndServe(); errListen != nil && errListen != http.ErrServerClosed {
			log.Fatal("can't listen and serve: %v", errListen)
		}
	}()

	// Waiting OS signals or context cancellation
	wait(ctx, log)

	ctxShutdown, cancelCtxShutdown := context.WithTimeout(ctx, 5*time.Second)
	defer cancelCtxShutdown()

	if errShutdown := s.Shutdown(ctxShutdown); errShutdown != nil {
		log.Error("shutdown error: %v", errShutdown)
	}
}

func wait(ctx context.Context, log logger.Logger) {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-osSignals:
	case <-ctx.Done():
		log.Error(ctx.Err().Error())
	}

	log.Info("termination signal received")
}
