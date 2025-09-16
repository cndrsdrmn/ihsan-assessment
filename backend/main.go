package main

import (
	"log/slog"

	backend "github.com/cndrsdrmn/ihsan-assessment/backend/internal"
	"github.com/cndrsdrmn/ihsan-assessment/config"
	"github.com/cndrsdrmn/ihsan-assessment/infrastructure"
	"github.com/cndrsdrmn/ihsan-assessment/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := infrastructure.NewDBConnection()
	if err != nil {
		slog.Error("connection failed to database", slog.String("error", err.Error()))
	}

	handler := backend.NewBlogHandler(repository.NewBlogRepository(db))

	r := gin.Default()
	r.POST("/blogs", handler.Create)
	r.PUT("/blogs/:id", handler.Update)
	r.DELETE("/blogs/:id", handler.Delete)

	r.Run(config.PORT_BACKEND)
}
