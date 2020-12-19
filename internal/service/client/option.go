package client

import (
	"net/http"
	"time"
)

const (
	DefaultTimeout = 1 * time.Second
)

type Option func(*service) error

func WithTransport(t Transport) Option {
	return func(s *service) error {
		s.transport = t
		return nil
	}
}

func WithDefaultTransport() Option {
	return func(s *service) error {
		s.transport = &http.Client{
			Timeout: DefaultTimeout,
		}
		return nil
	}
}
