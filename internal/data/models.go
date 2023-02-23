package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Models struct which wraps the MovieModel
// kind of enveloping
type Models struct {
	Users   UserModel
	Posts   PostModel
	Tokens  TokenModel
	Stories StoryModel
	Story   StoryModel
	Comment CommentModel
}

// NewModels method which returns a Models struct containing the initialized MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Users:   UserModel{DB: db},
		Posts:   PostModel{DB: db},
		Tokens:  TokenModel{DB: db},
		Story:   StoryModel{DB: db},
		Comment: CommentModel{DB: db},
	}
}
