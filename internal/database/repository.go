package database

import (
	"context"
)

// DatabaseRepo - contract for our DB calls
type DatabaseRepo interface {
	Ping(ctx context.Context) error
}
