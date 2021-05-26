package server

import (
	"context"
	"sync"
)

type Server interface {
	// execute blocking
	Serve(ctx context.Context, wg *sync.WaitGroup) error
	// execute non-blocked
	Run(ctx context.Context, wg *sync.WaitGroup) error
	Stop() error
	GracefulStop() error
}
