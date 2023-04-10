package batches

import (
	"context"
	"time"

	"github.com/sumup/sii-certification/internal/entities"
)

type IDomain interface {
	Generate(ctx context.Context, startDate, endDate time.Time) ([]entities.Batch, error)
	Send(ctx context.Context) ([]entities.Batch, error)
	Auth(ctx context.Context) (string, string, string, string, error)
	GetBatchesByDay(ctx context.Context, day string, page, pageSize int) ([]entities.Batch, int64, error)
}

type IGateway interface {
	SendMany(ctx context.Context, batch []entities.Batch) (bool, string, []entities.Batch, error)
	GetSeed(ctx context.Context) (string, string, error)
	GetToken(ctx context.Context, seed string) (string, string, error)
}
