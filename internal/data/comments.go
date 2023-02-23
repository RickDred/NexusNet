package data

import "time"

type Comment struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	AuthorID  int64     `json:"author_id"`
	PostID    int64     `json:"post_id"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}
