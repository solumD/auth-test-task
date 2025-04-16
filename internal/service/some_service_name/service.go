package someservicename

import (
	"context"

	"github.com/solumD/auth-test-task/internal/client/db"
	"github.com/solumD/auth-test-task/internal/repository"
	"github.com/solumD/auth-test-task/internal/service"
)

type srv struct {
	someRepository repository.SomeRepository
	txManager      db.TxManager
}

// New returns new service object
func New(someRepository repository.SomeRepository, txManager db.TxManager) service.SomeService {
	return &srv{
		someRepository: someRepository,
		txManager:      txManager,
	}
}

// SomeMethod ...
func (s *srv) SomeMethod(_ context.Context, _ ...interface{}) (interface{}, error) {
	// some business logic
	return struct{}{}, nil
}
