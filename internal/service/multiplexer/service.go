package multiplexer

import (
	"context"
	"fmt"
)

type Result map[string]map[string]interface{}

type service struct {
	client Client
}

func New(opts ...Option) (*service, error) {
	var (
		s   = &service{}
		err error
	)

	for i, configure := range opts {
		if err = configure(s); err != nil {
			return nil, fmt.Errorf("invalid option %d: %w", i, err)
		}
	}

	if s.client == nil {
		return nil, fmt.Errorf("client is nil")
	}

	return s, nil
}

func (s service) Urls(ctx context.Context, urls []string) (Result, error) {
	result := make(map[string]map[string]interface{}, len(urls))
	for _, url := range urls {
		r, errDo := s.client.Get(ctx, url)
		if errDo != nil {
			return nil, fmt.Errorf("can't request %s: %w", url, errDo)
		}
		result[url] = r
	}

	return result, nil
}
