package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)

type Game struct {
	ID			int
	Title		string
	ImageLink	string
	Description string
	Players		int
}

type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	ImageLink	   string
	Bio 		   *string
}

type GameOwnership struct {
	ID             	int
	UserID 			int
	GameID			int
}

type Review struct {
	ID             		int
	ReviewerID 			int
	ReviewedID			int
	ReviewText			string
}
