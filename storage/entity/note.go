package entity

type Note struct {
	Id int64 `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	CreatedAt string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
	FolderId int64 `json:"folder_id"`
}