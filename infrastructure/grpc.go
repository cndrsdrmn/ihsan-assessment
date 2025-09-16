package infrastructure

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func RunGrpcServer(srv *grpc.Server, ctx context.Context, address string) error {
	grp, ctx := errgroup.WithContext(ctx)

	grp.Go(func() error {
		lis, err := net.Listen("tcp", address)
		if err != nil {
			return fmt.Errorf("failed to listen on address %q: %w", address, err)
		}

		slog.Info("starting grpc server", slog.String("address", address))

		if err := srv.Serve(lis); err != nil {
			return fmt.Errorf("failed to serve grpc server: %w", err)
		}

		slog.Info("connected grpc server", slog.String("address", address))

		return nil
	})

	grp.Go(func() error {
		<-ctx.Done()

		srv.GracefulStop()

		return ctx.Err()
	})

	return grp.Wait()

}
