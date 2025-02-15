package models

type Event struct {
	UserID   int64  `json:"user_id"`
	AuthorID int64  `json:"author_id"`
	PostID   string `json:"post_id"`
	UserName string `json:"user_name"`
	Action   string `json:"action"`
}
