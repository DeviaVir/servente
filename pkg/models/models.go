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

// Service model definition of a service
type Service struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// User model definition of a user
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Active         bool
}
