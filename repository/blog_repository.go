package repository

import (
	e "github.com/cndrsdrmn/ihsan-assessment/entities"
	"gorm.io/gorm"
)

type BlogRepository interface {
	Create(*e.Blog) (*e.Blog, error)
	Read(string) (*e.Blog, error)
	Update(*e.Blog) (*e.Blog, error)
	Delete(string) error
}

type blog struct {
	db *gorm.DB
}

// Create inserts a new blog and returns the full record
func (b *blog) Create(blog *e.Blog) (*e.Blog, error) {
	if err := b.db.Create(blog).Error; err != nil {
		return nil, err
	}
	// re-fetch to populate default values (timestamps, ID, etc.)
	return b.Read(blog.ID)
}

// Read fetches a blog by ID
func (b *blog) Read(id string) (*e.Blog, error) {
	var blog = e.Blog{ID: id}
	if err := b.db.First(&blog).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

// Update performs a full update and returns the updated record
func (b *blog) Update(blog *e.Blog) (*e.Blog, error) {
	if err := b.db.Save(blog).Error; err != nil {
		return nil, err
	}
	return b.Read(blog.ID)
}

// Delete deletes a blog by ID
func (b *blog) Delete(id string) error {
	return b.db.Delete(&e.Blog{}, "id = ?", id).Error
}

func NewBlogRepository(db *gorm.DB) BlogRepository {
	return &blog{db: db}
}
