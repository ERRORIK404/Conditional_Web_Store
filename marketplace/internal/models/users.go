package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Login          string    `gorm:"uniqueIndex;not null;size:50" json:"login"`
	HashedPassword string    `gorm:"not null;size:250"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
