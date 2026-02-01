package cache

import "github.com/MihirSahani/Project-27/storage/entity"

type CacheManager interface {
	Close() error
	Ping() error

	GetUser(int64) (*entity.User, error)
	GetFolder(int64) (*entity.Folder, error)
	GetNote(int64) (*entity.Note, error)
	GetAllFolders(int64) ([]*entity.Folder, error)
	GetNotesInFolder(int64) ([]*entity.Note, error)

	SetUser(*entity.User) error
	SetFolder(*entity.Folder) error
	SetNote(*entity.Note) error
	SetAllFolders([]*entity.Folder, int64) error
	SetNotesInFolder([]*entity.Note) error

	DeleteUser(int64) error
	DeleteFolder(int64) error
	DeleteNote( int64) error
	DeleteAllFolders(int64) error
	DeleteNotesInFolder(int64) error
}