package middleware

import (
	"context"
	"fmt"

	pbm "github.com/cndrsdrmn/ihsan-assessment/generated/blog/middleware/v1"
	pbs "github.com/cndrsdrmn/ihsan-assessment/generated/blog/service/v1"
)

type BlogMiddleware interface {
	pbm.BlogMiddlewareServiceServer
}

type midd struct {
	pbm.UnimplementedBlogMiddlewareServiceServer
	service pbs.BlogServiceClient
}

func (m *midd) SendMessage(ctx context.Context, req *pbm.SendMessageRequest) (*pbm.SendMessageResponse, error) {
	fmt.Println("req", req)

	res, err := m.service.ProcessMessage(ctx, &pbs.ProcessMessageRequest{
		Data: req.GetData(),
	})

	return &pbm.SendMessageResponse{
		Status: res.GetStatus(),
		Data:   res.GetData(),
	}, err
}

func NewBlogMiddleware(service pbs.BlogServiceClient) BlogMiddleware {
	return &midd{service: service}
}
