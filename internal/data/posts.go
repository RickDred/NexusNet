package data

import "time"

type Post struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-,omitempty"`
	Title       string    `json:"title"`
	UpdatedAt   time.Time `json:"updated_at"`
	Description string    `json:"description"`
	AuthorID    int64     `json:"author_id"`
}
