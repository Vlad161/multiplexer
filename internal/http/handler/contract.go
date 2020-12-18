package handler

import (
	"context"

	"github.com/vlad161/affise_test_task/internal/service/multiplexer"
)

type Multiplexer interface {
	Urls(ctx context.Context, urls []string) (multiplexer.Result, error)
}
