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

	StatusInDevelopment  = 1 // active dev
	StatusBeta           = 2 // testing
	StatusProduction     = 3 // prod and active dev
	StatusMaintain       = 4 // prod no longer developed/extended
	StatusDeprecated     = 5 // legacy, to be replaced
	StatusDecommissioned = 6 // it used to be a thing but it's gone now

	TypeLink    = "link"
	TypeWebhook = "webhook"
	TypeSelect  = "select"
)

// @NOTE: add new models to web/main.go AutoMigrate to have them automatically created

// Organization model definition for organizations
type Organization struct {
	gorm.Model
	Identifier string `gorm:"type:varchar(100) unique"`
	Name       string
	Active     bool
	Users      []*User `gorm:"many2many:user_organizations;"`
	Services   []*Service
	Settings   []*Setting
	Attribute  []*Attribute `gorm:"foreignKey:OrganizationID"`
	AuditLogs  []*AuditLog
}

// Service model definition of a service
type Service struct {
	gorm.Model
	Identifier     string `gorm:"type:varchar(100) unique"`
	Title          string
	Description    string
	Status         int
	Attributes     []*Attribute `gorm:"foreignKey:ServiceID"`
	Organization   Organization
	OrganizationID int
}

// User model definition of a user
type User struct {
	gorm.Model
	Name           string
	Email          string `gorm:"unique"`
	HashedPassword []byte `gorm:"type:char(60)"`
	Active         bool
	Organizations  []*Organization `gorm:"many2many:user_organizations;"`
}

// Attribute model definition of available attributes
type Attribute struct {
	gorm.Model
	Value          string
	Active         bool
	Setting        Setting
	SettingID      int
	Organization   *Organization
	OrganizationID int
	Service        *Service
	ServiceID      int
}

// Setting model definition of configured settings
type Setting struct {
	gorm.Model
	Key            string // unique id
	Title          string // nice name
	Type           string // link, webhook, select, etc
	Scope          string // "organization" or "service"
	Organization   *Organization
	OrganizationID int
}

// AuditLog model definition of logs
type AuditLog struct {
	gorm.Model
	User           int
	Message        string
	Organization   *Organization
	OrganizationID int
}
