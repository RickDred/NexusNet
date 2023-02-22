package data

import "time"

type User struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Role        string    `json:"role"`
	Password    password  `json:"-"`
	Activated   bool      `json:"activated"`
	Description string    `json:"description"`
}

type password struct {
	plaintext *string
	hash      []byte
}
