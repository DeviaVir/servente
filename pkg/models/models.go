package models

import (
	"errors"

	"gorm.io/gorm"
)

var (
	ErrNoRecord            = errors.New("models: no matching record found")
	ErrInvalidCredentials  = errors.New("models: invalid credentials")
	ErrDuplicateEmail      = errors.New("models: duplicate email")
	ErrDuplicateIdentifier = errors.New("models: duplicate service identifier")

	StatusInDevelopment  = 1
	StatusBeta           = 2
	StatusProduction     = 3
	StatusMaintain       = 4
	StatusDeprecated     = 5
	StatusDecommissioned = 6
)

// Service model definition of a service
type Service struct {
	gorm.Model
	Identifier  string `gorm:"type:varchar(100) unique"`
	Title       string
	Description string
	Attributes  string
	Status      int
}

// User model definition of a user
type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"unique"`
	HashedPassword []byte `gorm:"type:char(60)"`
	Active         bool
}

// Attribute model definition of available attributes
type Attribute struct {
	gorm.Model
	Key    string
	Title  string
	Active bool
}

// Setting model definition of configured settings
type Setting struct {
	gorm.Model
	Key   string
	Value string
}

// AuditLog model definition of logs
type AuditLog struct {
	gorm.Model
	User    int
	Message string
}
