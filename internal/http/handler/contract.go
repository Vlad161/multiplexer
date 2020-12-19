package handler

import (
	"context"

	"github.com/vlad161/multiplexer/internal/service/multiplexer"
)

type Multiplexer interface {
	Urls(ctx context.Context, urls []string) (multiplexer.Result, error)
}
