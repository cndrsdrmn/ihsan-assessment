package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"buf.build/go/protovalidate"
	"github.com/cndrsdrmn/ihsan-assessment/config"
	pbs "github.com/cndrsdrmn/ihsan-assessment/generated/blog/service/v1"
	"github.com/cndrsdrmn/ihsan-assessment/infrastructure"
	repo "github.com/cndrsdrmn/ihsan-assessment/repository"
	srvc "github.com/cndrsdrmn/ihsan-assessment/service/internal"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	db, err := infrastructure.NewDBConnection()
	if err != nil {
		slog.Error("connection failed to database", slog.String("error", err.Error()))
	}

	srv, err := server(db)
	if err != nil {
		slog.Error("cannot initialize service", slog.String("error", err.Error()))
	}

	if err := infrastructure.RunGrpcServer(srv, ctx, config.GRPC_SERVICE); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("error running application", slog.String("error", err.Error()))
	}

	slog.Info("closing server")
}

func server(db *gorm.DB) (*grpc.Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, errors.New("failed to create protovalidate validator")
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
	)

	reflection.Register(srv)

	pbs.RegisterBlogServiceServer(srv, srvc.NewBlogService(
		repo.NewBlogRepository(db),
	))

	return srv, nil
}
