package abstractions

import (
	"context"
)

type IUnitOfWorkFactory interface {
	NewUnitOfWork(ctx context.Context) (IUnitOfWork, error)
}
