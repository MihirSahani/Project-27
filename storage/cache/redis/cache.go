package redis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/MihirSahani/Project-27/storage/entity"
	"github.com/go-redis/redis/v8"
)

var (
	CACHE_DISABLED = errors.New("cache disabled")
)

type RedisCacheManager struct {
	client *redis.Client
	config *RedisConfig
}

func NewRedisCacheManager() *RedisCacheManager {
	config := NewRedisConfig()
	return &RedisCacheManager{
		config: config,
		client: redis.NewClient(&redis.Options{
			Addr:     config.Address,
			Password: config.Password,
			DB:       config.Db,
		}),
	}
}

func (r *RedisCacheManager) get(key string) (string, error) {
	return r.client.Get(r.client.Context(), key).Result()
}

func (r *RedisCacheManager) set(key string, value []byte) error {
	return r.client.Set(r.client.Context(), key, value, 0).Err()
}

func (r *RedisCacheManager) delete(key string) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	return r.client.Del(r.client.Context(), key).Err()
}

func (r *RedisCacheManager) flush() error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	return r.client.FlushDB(r.client.Context()).Err()
}

func (r *RedisCacheManager) Close() error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	if err := r.flush(); err != nil {
		return err
	}
	return r.client.Close()
}

func (r *RedisCacheManager) Ping() error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	return r.client.Ping(r.client.Context()).Err()
}

func (r *RedisCacheManager) GetUser(id int64) (*entity.User, error) {
	if !r.config.enabled {
		fmt.Printf("CACHE MISS!")
		return nil, CACHE_DISABLED
	}
	data, err := r.get(fmt.Sprintf("user-%v", id))
	if err != nil {
		fmt.Printf("CACHE MISS!")
		return nil, err
	}
	var user entity.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RedisCacheManager) SetUser(user *entity.User) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	jsonData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.set(fmt.Sprintf("user-%v", user.Id), jsonData)
}

func (r *RedisCacheManager) DeleteUser(id int64) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	return r.delete(fmt.Sprintf("user-%v", id))
}

func (r *RedisCacheManager) GetFolder(id int64) (*entity.Folder, error) {
	if !r.config.enabled {
		return nil, CACHE_DISABLED
	}
	data, err := r.get(fmt.Sprintf("folder-%v", id))
	if err != nil {
		return nil, err
	}
	var folder entity.Folder
	err = json.Unmarshal([]byte(data), &folder)
	if err != nil {
		return nil, err
	}
	return &folder, nil
}

func (r *RedisCacheManager) SetFolder(folder *entity.Folder) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	jsonData, err := json.Marshal(folder)
	if err != nil {
		return err
	}
	return r.set(fmt.Sprintf("folder-%v", folder.Id), jsonData)
}

func (r *RedisCacheManager) DeleteFolder(id int64) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	return r.delete(fmt.Sprintf("folder-%v", id))
}

func (r *RedisCacheManager) GetNote(id int64) (*entity.Note, error) {
	if !r.config.enabled {
		return nil, CACHE_DISABLED
	}
	data, err := r.get(fmt.Sprintf("note-%v", id))
	if err != nil {
		return nil, err
	}
	var note entity.Note
	err = json.Unmarshal([]byte(data), &note)
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (r *RedisCacheManager) SetNote(note *entity.Note) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	jsonData, err := json.Marshal(note)
	if err != nil {
		return err
	}
	return r.set(fmt.Sprintf("note-%v", note.Id), jsonData)
}

func (r *RedisCacheManager) DeleteNote(id int64) error {
	if r.config.enabled {
		return CACHE_DISABLED
	}
	return r.delete(fmt.Sprintf("note-%v", id))
}

func (r *RedisCacheManager) GetAllFolders(userId int64) ([]*entity.Folder, error) {
	if !r.config.enabled {
		return nil, CACHE_DISABLED
	}
	data, err := r.get(fmt.Sprintf("folders-%v", userId))
	if err != nil {
		return nil, err
	}
	var folders []*entity.Folder
	err = json.Unmarshal([]byte(data), &folders)
	if err != nil {
		return nil, err
	}
	return folders, nil
}

func (r *RedisCacheManager) SetAllFolders(folders []*entity.Folder, userId int64) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	jsonData, err := json.Marshal(folders)
	if err != nil {
		return err
	}
	return r.set(fmt.Sprintf("folders-%v", userId), jsonData)
}

func (r *RedisCacheManager) DeleteAllFolders(userId int64) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	return r.delete(fmt.Sprintf("folders-%v", userId))
}

func (r *RedisCacheManager) GetNotesInFolder(folderId int64) ([]*entity.Note, error) {
	if !r.config.enabled {
		return nil, CACHE_DISABLED
	}
	data, err := r.get(fmt.Sprintf("notes-%v", folderId))
	if err != nil {
		return nil, err
	}
	var notes []*entity.Note
	err = json.Unmarshal([]byte(data), &notes)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (r *RedisCacheManager) SetNotesInFolder(notes []*entity.Note) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	jsonData, err := json.Marshal(notes)
	if err != nil {
		return err
	}
	return r.set(fmt.Sprintf("notes-%v", notes[0].FolderId), jsonData)
}

func (r *RedisCacheManager) DeleteNotesInFolder(folderId int64) error {
	if !r.config.enabled {
		return CACHE_DISABLED
	}
	return r.delete(fmt.Sprintf("notes-%v", folderId))
}