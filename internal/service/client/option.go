package client

import (
	"net/http"
	"time"
)

const (
	DefaultTimeout = 1 * time.Second
)

type Option func(*service) error

func WithHttpClient(cl *http.Client) Option {
	return func(s *service) error {
		s.http = cl
		return nil
	}
}

func WithDefaultHttpClient() Option {
	return func(s *service) error {
		s.http = &http.Client{
			Timeout: DefaultTimeout,
		}
		return nil
	}
}
