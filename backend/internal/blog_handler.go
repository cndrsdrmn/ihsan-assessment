package backend

import (
	"fmt"
	"net/http"

	"github.com/cndrsdrmn/ihsan-assessment/entities"
	pb "github.com/cndrsdrmn/ihsan-assessment/generated/blog/v1"
	"github.com/cndrsdrmn/ihsan-assessment/repository"
	"github.com/gin-gonic/gin"
)

type BlogHandler interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type blogHandler struct {
	repo repository.BlogRepository
}

// Create handles POST /blogs
func (b *blogHandler) Create(ctx *gin.Context) {
	var blog entities.Blog
	if err := ctx.ShouldBindJSON(&blog); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("create", blog)

	created, err := b.repo.Create(&blog)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create blog"})
		return
	}

	ctx.JSON(http.StatusCreated, created)
}

// Update handles PUT /blogs/:id
func (b *blogHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	var req pb.UpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("update", req)

	existing, err := b.repo.Read(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}

	for _, path := range req.UpdateMask.Paths {
		switch path {
		case "title":
			existing.Title = req.Blog.Title
		case "content":
			existing.Content = req.Blog.Content
		default:
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported field: " + path})
			return
		}
	}

	updated, err := b.repo.Update(existing)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update blog"})
		return
	}

	ctx.JSON(http.StatusOK, updated)
}

// Delete handles DELETE /blogs/:id
func (b *blogHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	if err := b.repo.Delete(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete blog"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "blog deleted"})
}

func NewBlogHandler(repo repository.BlogRepository) BlogHandler {
	return &blogHandler{repo: repo}
}
