package data

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Repositories struct {
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
		GetAll(title string, genres []string, filters Filters) ([]*Movie, Metadata, error)
	}
	Users interface {
		Insert(user *User) error
		GetByEmail(email string) (*User, error)
		Update(user *User) error
		GetForToken(tokenScope string, tokenPlaintext string) (*User, error)
	}
	Tokens interface {
		New(userId int64, ttl time.Duration, scope string) (*Token, error)
		Insert(token *Token) error
		DeleteAllForUser(scope string, userId int64) error
	}
	Permissions interface {
		GetAllForUser(userID int64) (Permissions, error)
		AddForUser(userID int64, codes ...string) error
	}
}

func NewRepositories(db *sql.DB) Repositories {
	return Repositories{
		Movies:      MovieRepository{DB: db},
		Users:       UserRepository{DB: db},
		Tokens:      TokenRepository{DB: db},
		Permissions: PermissionsRepository{DB: db},
	}
}

func NewMockRepositories(db *sql.DB) Repositories {
	return Repositories{
		Movies:      MockMovieRepository{},
		Users:       MockUserRepository{},
		Tokens:      MockTokenRepository{},
		Permissions: MockPermissionsRepository{},
	}
}
