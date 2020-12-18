package multiplexer

import "context"

type Client interface {
	Get(ctx context.Context, url string) (map[string]interface{}, error)
}
