package broker

import (
	"context"
)

type Publisher interface {
	PublishMessage(ctx context.Context, body []byte) error
}

type Consumer[T any] interface {
	Consume(ctx context.Context) (<-chan T, error)
}
