package adapters

import "context"

type IAdapter interface {
	Post(ctx context.Context, url string, headers map[string]string, payload []byte) (response []byte, err error)
}
