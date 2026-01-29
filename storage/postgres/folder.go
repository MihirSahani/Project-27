package postgres

import (
	"context"
	"database/sql"

	"github.com/MihirSahani/Project-27/storage/entity"
)

type PostgresFolderStorageManager struct{}

func NewPostgresFolderStorageManager() *PostgresFolderStorageManager {
	return &PostgresFolderStorageManager{}
}

func (p *PostgresFolderStorageManager) CreateFolder(ctx context.Context, db *sql.Tx, folder *entity.Folder) (*entity.Folder, error) {
	query := `
		INSERT INTO folders (name, user_id)
		VALUES ($1, $2)
		RETURNING id, name, user_id, created_at, updated_at
	`
	err := db.QueryRowContext(ctx, query,
		folder.Name,
		folder.UserId,
	).Scan(&folder.Id, &folder.Name, &folder.UserId, &folder.CreatedAt, &folder.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (p *PostgresFolderStorageManager) DeleteFolder(ctx context.Context, db *sql.Tx, id int64) (*entity.Folder, error) {
	query := `
		DELETE FROM folders
		WHERE id = $1
		RETURNING id, name, user_id, created_at, updated_at
	`
	folder := &entity.Folder{}
	err := db.QueryRowContext(ctx, query, id).Scan(&folder.Id, &folder.Name, &folder.UserId, &folder.CreatedAt, &folder.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (p *PostgresFolderStorageManager) GetAllFolders(ctx context.Context, db *sql.Tx, id int64) ([]*entity.Folder, error) {
	query := `
		SELECT id, name, user_id, created_at, updated_at
		FROM folders
		WHERE user_id = $1
	`
	rows, err := db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	folders := []*entity.Folder{}
	for rows.Next() {
		folder := &entity.Folder{}
		err := rows.Scan(
			&folder.Id,
			&folder.Name,
			&folder.UserId,
			&folder.CreatedAt,
			&folder.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		folders = append(folders, folder)
	}
	return folders, nil
}

func (p *PostgresFolderStorageManager) GetNotesInFolder(ctx context.Context, db *sql.Tx, id int64) ([]*entity.Note, error) {
	return nil, nil
}
