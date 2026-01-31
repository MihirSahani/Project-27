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

	FolderStorageManager interface {
		CreateFolder(context.Context, *sql.Tx, *entity.Folder) (*entity.Folder, error)
		DeleteFolder(context.Context, *sql.Tx, int64) (*entity.Folder, error)
		GetNotesInFolder(context.Context, *sql.Tx, int64, int64) ([]*entity.Note, error)
		GetAllFolders(context.Context, *sql.Tx, int64) ([]*entity.Folder, error)
	}

	NoteStorageManager interface {
		CreateNote(context.Context, *sql.Tx, *entity.Note) (*entity.Note, error)
		DeleteNote(context.Context, *sql.Tx, int64) (*entity.Note, error)
		GetNoteByID(context.Context, *sql.Tx, int64) (*entity.Note, error)
		UpdateNote(context.Context, *sql.Tx, *entity.Note) (*entity.Note, error)
	}
}

func NewStorageManager() *StorageManager {
	databaseConnection, err := postgres.CreateConfiguredPostgresStorage()
	if err != nil {
		panic(err)
	}

	return &StorageManager{
		databaseConnection:   databaseConnection,
		UserStorageManager:   postgres.NewPostgresUserStorageManager(),
		FolderStorageManager: postgres.NewPostgresFolderStorageManager(),
		NoteStorageManager:   postgres.NewPostgresNoteStorageManager(),
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
