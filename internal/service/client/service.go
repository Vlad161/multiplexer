package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/vlad161/multiplexer/internal/logger"
)

type service struct {
	transport Transport
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

	if s.transport == nil {
		return nil, fmt.Errorf("http transport is nil")
	}

	return s, nil
}

func (s service) Get(ctx context.Context, url string) (map[string]interface{}, error) {
	log := logger.FromContext(ctx)

	req, errNewRequest := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if errNewRequest != nil {
		return nil, errNewRequest
	}

	resp, errDo := s.transport.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer func() {
		if errCloseBody := resp.Body.Close(); errCloseBody != nil {
			log.Error("can't close response body: %v", errCloseBody)
		}
	}()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code is %v", resp.StatusCode)
	}

	var respBody map[string]interface{}
	if errDecode := json.NewDecoder(resp.Body).Decode(&respBody); errDecode != nil {
		return nil, errDecode
	}
	return respBody, nil
}
