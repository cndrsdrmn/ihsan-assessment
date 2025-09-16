package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Blog struct {
	ID        string `gorm:"type:uuid;primaryKey" json:"id"`
	Title     string `gorm:"index" json:"title"`
	Content   string `json:"content"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Blog) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New().String()
	return
}
