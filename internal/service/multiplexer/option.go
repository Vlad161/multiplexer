package multiplexer

import "github.com/vlad161/multiplexer/internal/service/client"

type Option func(*service) error

func WithClient(cl Client) Option {
	return func(s *service) error {
		s.client = cl
		return nil
	}
}

func WithDefaultClient() Option {
	return func(s *service) error {
		cl, err := client.New(client.WithDefaultHttpClient())
		if err != nil {
			return err
		}

		s.client = cl
		return nil
	}
}
