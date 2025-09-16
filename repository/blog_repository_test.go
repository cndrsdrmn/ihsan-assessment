package repository_test

import (
	"testing"
	"time"

	e "github.com/cndrsdrmn/ihsan-assessment/entities"
	"github.com/cndrsdrmn/ihsan-assessment/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	// Auto migrate schema
	if err := db.AutoMigrate(&e.Blog{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func TestBlogRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := repository.NewBlogRepository(db)

	// --- Create ---
	blog := &e.Blog{
		Title:   "My First Blog",
		Content: "Hello, world!",
	}
	created, err := repo.Create(blog)
	assert.NoError(t, err)
	assert.NotEmpty(t, created.ID)
	assert.Equal(t, "My First Blog", created.Title)
	assert.Equal(t, "Hello, world!", created.Content)
	assert.WithinDuration(t, time.Now(), created.CreatedAt, time.Second*2)

	// --- Read ---
	read, err := repo.Read(created.ID)
	assert.NoError(t, err)
	assert.Equal(t, created.ID, read.ID)
	assert.Equal(t, "My First Blog", read.Title)

	// --- Update ---
	read.Title = "Updated Blog"
	updated, err := repo.Update(read)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Blog", updated.Title)
	assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))

	// --- Delete ---
	err = repo.Delete(updated.ID)
	assert.NoError(t, err)

	_, err = repo.Read(updated.ID)
	assert.Error(t, err) // should return record not found
}
