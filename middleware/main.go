package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"

	"buf.build/go/protovalidate"
	"github.com/cndrsdrmn/ihsan-assessment/config"
	pbm "github.com/cndrsdrmn/ihsan-assessment/generated/blog/middleware/v1"
	pbs "github.com/cndrsdrmn/ihsan-assessment/generated/blog/service/v1"
	"github.com/cndrsdrmn/ihsan-assessment/infrastructure"
	midd "github.com/cndrsdrmn/ihsan-assessment/middleware/internal"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	conn, err := grpc.NewClient(config.GRPC_SERVICE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("cannot connect to grpc server", slog.String("error", err.Error()))
	}
	defer conn.Close()

	srv, err := server(conn)
	if err != nil {
		slog.Error("cannot initialize service", slog.String("error", err.Error()))
	}

	if err := infrastructure.RunGrpcServer(srv, ctx, config.GRPC_MIDDLEWARE); err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("error running application", slog.String("error", err.Error()))
	}

	slog.Info("closing server")
}

func server(conn *grpc.ClientConn) (*grpc.Server, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, errors.New("failed to create protovalidate validator")
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(protovalidate_middleware.UnaryServerInterceptor(validator)),
	)

	reflection.Register(srv)

	pbm.RegisterBlogMiddlewareServiceServer(srv, midd.NewBlogMiddleware(
		pbs.NewBlogServiceClient(conn),
	))

	return srv, nil
}
