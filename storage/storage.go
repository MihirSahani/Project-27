package storage

import (
	"context"
	"database/sql"

	"github.com/MihirSahani/Project-27/storage/entity"
	"github.com/MihirSahani/Project-27/storage/postgres"
)


type StorageManager struct {
	databaseConnection interface {
		Close() error
		GetDb() *sql.DB

	}
	UserStorageManager interface {
		CreateUser(context.Context, *sql.Tx, *entity.User) (*entity.User, error)
		GetUserByEmail(context.Context, *sql.Tx, string) (*entity.User, error)
		GetUserByID(context.Context, *sql.Tx, int64) (*entity.User, error)
		UpdateUser(context.Context, *sql.Tx, *entity.User) (*entity.User, error)
		DeleteUser(context.Context, *sql.Tx, int64) (*entity.User, error)
		ActivateUser(context.Context, *sql.Tx, int64) (*entity.User, error)
	}
}

func NewStorageManager() *StorageManager {
	databaseConnection, err := postgres.CreateConfiguredPostgresStorage()
	if err != nil {
		panic(err)
	}
	
	return &StorageManager{
		databaseConnection: databaseConnection,
		UserStorageManager: postgres.NewPostgresUserStorageManager(),
	}
}

func (sm *StorageManager) Close() error {
	return sm.databaseConnection.Close()
}

func (sm *StorageManager) WithTx(ctx context.Context, fn func(context.Context, *sql.Tx) (any, error)) (any, error) {
	tx, err := sm.databaseConnection.GetDb().BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	data, err := fn(ctx, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return data, nil
}