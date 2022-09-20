package repository

import (
	"context"
	"fmt"

	"github.com/alpine-hodler/gidari/internal/storage"
	"github.com/alpine-hodler/gidari/proto"
)

// Generic is the interface for the generic service.
type Generic interface {
	storage.Storage
	storage.Tx

	Transact(fn func(ctx context.Context, repo Generic) error)
}

// GenericService is the implementation of the Generic service.
type GenericService struct {
	storage.Storage
	storage.Tx
}

// New returns a new Generic service.
func New(ctx context.Context, dns string) (Generic, error) {
	stg, err := storage.New(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("failed to construct storage: %v", err)
	}
	return &GenericService{stg, nil}, nil
}

// NewTx returns a new Generic service with an initialized transaction object that can be used to commit or rollback
// storage operations made by the repository layer.
func NewTx(ctx context.Context, dns string) (Generic, error) {
	stg, err := storage.New(ctx, dns)
	if err != nil {
		return nil, fmt.Errorf("failed to construct storage: %v", err)
	}

	tx, err := stg.StartTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %v", err)
	}
	return &GenericService{stg, tx}, nil

}

// Transact is a helper function that wraps a function in a transaction and commits or rolls back the transaction. If
// svc is not a transaction, the function will be executed without executing.
func (svc *GenericService) Transact(fn func(ctx context.Context, repo Generic) error) {
	svc.Tx.Send(func(ctx context.Context, stg storage.Storage) error {
		err := fn(ctx, svc)
		if err != nil {
			return fmt.Errorf("error executing transaction: %v", err)
		}
		return nil
	})
}

// Truncate truncates a table.
func (svc *GenericService) Truncate(ctx context.Context, req *proto.TruncateRequest) (*proto.TruncateResponse, error) {
	rsp, err := svc.Storage.Truncate(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error truncating table: %v", err)
	}
	return rsp, nil
}
