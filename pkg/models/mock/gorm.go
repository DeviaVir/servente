package mock

import (
	"time"

	"gorm.io/gorm"
)

var mockGorm = &gorm.Model{
	ID:        2,
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}
