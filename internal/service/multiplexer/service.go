package multiplexer

import (
	"context"
	"fmt"
	"sync"
)

const (
	semaphoreMaxGoroutines = 4
)

type Result map[string]map[string]interface{}

type keyValue struct {
	key   string
	value map[string]interface{}
}

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
	done := make(chan struct{})
	defer close(done)

	kvChan, errChan := s.visitUrls(ctx, urls, done)

	for range urls {
		select {
		case kv := <-kvChan:
			result[kv.key] = kv.value
		case err := <-errChan:
			done <- struct{}{}
			return nil, err
		case <-ctx.Done():
			done <- struct{}{}
			return nil, ctx.Err()
		}
	}

	return result, nil
}

func (s service) visitUrls(ctx context.Context, urls []string, done <-chan struct{}) (kv <-chan keyValue, err <-chan error) {
	sem := make(chan struct{}, semaphoreMaxGoroutines)
	kvChan := make(chan keyValue)
	errChan := make(chan error)

	go func() {
		var wg sync.WaitGroup

		for _, u := range urls {
			sem <- struct{}{}
			wg.Add(1)
			go func(url string) {
				r, errDo := s.client.Get(ctx, url)

				select {
				case <-done:
					// Ignore result if done was closed
				default:
					if errDo != nil {
						errChan <- errDo
					} else {
						kvChan <- keyValue{url, r}
					}
				}

				<-sem
				wg.Done()
			}(u)

			// Abort cycle if done was called
			select {
			case <-done:
				return
			default:
				// do nothing
			}
		}

		go func() {
			wg.Wait()
			close(sem)
			close(kvChan)
			close(errChan)
		}()

	}()

	return kvChan, errChan
}
