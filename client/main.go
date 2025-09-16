package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/cndrsdrmn/ihsan-assessment/config"
	middlewarev1 "github.com/cndrsdrmn/ihsan-assessment/generated/blog/middleware/v1"
	blogv1 "github.com/cndrsdrmn/ihsan-assessment/generated/blog/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
)

func main() {
	conn, err := grpc.NewClient(config.GRPC_MIDDLEWARE, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("cannot connect to service")
	}
	defer conn.Close()

	server := middlewarev1.NewBlogMiddlewareServiceClient(conn)

	created, err := server.SendMessage(context.Background(), &middlewarev1.SendMessageRequest{
		TrxType: "create",
		Data: &blogv1.BlogRequest{
			Create: &blogv1.CreateRequest{
				Blog: &blogv1.Blog{
					Title:   "Lorem Ipsum",
					Content: "Hello World This Is From Client.",
				},
			},
		},
	})

	fmt.Println(created)

	if err != nil {
		slog.Error("failed to create transaction", slog.String("error", err.Error()))
	}

	updated, err := server.SendMessage(context.Background(), &middlewarev1.SendMessageRequest{
		TrxType: "update",
		Data: &blogv1.BlogRequest{
			Update: &blogv1.UpdateRequest{
				Id: created.GetData().GetId(),
				Blog: &blogv1.Blog{
					Title: "Foo Bar",
				},
				UpdateMask: &fieldmaskpb.FieldMask{
					Paths: []string{"title"},
				},
			},
		},
	})

	if err != nil {
		slog.Error("failed to update transaction", slog.String("error", err.Error()))
	}

	fmt.Println(updated)
}
