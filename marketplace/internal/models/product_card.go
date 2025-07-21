package models

import (
	"github.com/google/uuid"
	"time"
)

type ProductCard struct {
	ID          uuid.UUID `gorm:"primary_key;type:uuid;gen_random_uuid()" json:"id"`
	UserID      uuid.UUID `gorm:"not null;index" json:"user_id"`
	Title       string    `gorm:"size:100;NOT NULL" json:"title"`
	Description string    `gorm:"size:1000;NOT NULL" json:"description"`
	ImageURL    string    `gorm:"size:511;NOT NULL" json:"image_url"`
	Price       float64   `gorm:"type:decimal(10,2);NOT NULL" json:"price"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
