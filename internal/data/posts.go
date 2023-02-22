package data

import "time"

type Post struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-,omitempty"`
	Title       string
	Updated     time.Time
	Description string
	Author      *User
}
