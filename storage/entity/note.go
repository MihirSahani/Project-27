package entity

type Note struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UserId int64 `json:"user_id"`
	FolderId int64 `json:"folder_id"`
}