package main

import (
	"fmt"
	"net/http"

	"github.com/vlad161/affise_test_task/internal/http/handler"
	"github.com/vlad161/affise_test_task/internal/http/middleware"
	"github.com/vlad161/affise_test_task/internal/logger"
	"github.com/vlad161/affise_test_task/internal/service/multiplexer"
)

const PORT = 8081

func main() {
	log := logger.New()

	mp, errMultiplexer := multiplexer.New(multiplexer.WithDefaultClient())
	if errMultiplexer != nil {
		log.Fatal("can't create multiplexer service: %w", errMultiplexer)
	}

	mux := http.NewServeMux()
	mux.Handle("/multiplexer", handler.NewMultiplexer(mp))

	s := http.Server{
		Addr:           fmt.Sprintf(":%d", PORT),
		Handler:        mux,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
	}
	s.Handler = middleware.Logger(log, s.Handler)

	if errListen := s.ListenAndServe(); errListen != nil {
		log.Fatal("can't listen and serve: %w", errListen)
	}

	// TODO Graceful shutdown
}
