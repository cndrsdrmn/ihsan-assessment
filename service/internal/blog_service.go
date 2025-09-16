package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/cndrsdrmn/ihsan-assessment/config"
	e "github.com/cndrsdrmn/ihsan-assessment/entities"
	pbs "github.com/cndrsdrmn/ihsan-assessment/generated/blog/service/v1"
	pb "github.com/cndrsdrmn/ihsan-assessment/generated/blog/v1"
	r "github.com/cndrsdrmn/ihsan-assessment/repository"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BlogService interface {
	pbs.BlogServiceServer
}

type srv struct {
	pbs.UnimplementedBlogServiceServer
	repo r.BlogRepository
}

// ProcessMessage implements BlogService.
func (s *srv) ProcessMessage(ctx context.Context, req *pbs.ProcessMessageRequest) (*pbs.ProcessMessageResponse, error) {
	data := req.GetData()
	if data == nil {
		return failure(codes.InvalidArgument, "blog data is missing")
	}

	switch {
	case data.GetCreate() != nil:
		return s.forwardToBackend(ctx, http.MethodPost, "/blogs", data.GetCreate())
	case data.GetUpdate() != nil:
		return s.forwardToBackend(ctx, http.MethodPut, fmt.Sprintf("/blogs/%s", data.GetUpdate().GetId()), data.GetUpdate())
	case data.GetDelete() != nil:
		return s.forwardToBackend(ctx, http.MethodDelete, fmt.Sprintf("/blogs/%s", data.GetDelete().GetId()), nil)
	case data.GetRead() != nil:
		record, err := s.repo.Read(data.Read.GetId())
		if err != nil {
			return failure(codes.NotFound, fmt.Sprintf("blog with ID %s not found", data.GetRead().GetId()))
		}
		return success(record)
	default:
		return failure(codes.InvalidArgument, "unsupported transaction type")
	}
}

func (s *srv) forwardToBackend(ctx context.Context, method string, endpoint string, payload interface{}) (*pbs.ProcessMessageResponse, error) {
	var reqBody io.Reader

	if payload != nil {
		body, err := json.Marshal(payload)
		if err != nil {
			return failure(codes.Internal, fmt.Sprintf("failed to encode request: %v", err))
		}
		reqBody = bytes.NewBuffer(body)
	}

	url := config.REST_BACKEND + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return failure(codes.Internal, fmt.Sprintf("failed to create request: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return failure(codes.Unavailable, fmt.Sprintf("backend unreachable: %v", err))
	}
	defer res.Body.Close()

	if res.StatusCode >= http.StatusBadRequest {
		body, _ := io.ReadAll(res.Body)
		return failure(codes.Internal, fmt.Sprintf("backend error: %s, body: %s", res.Status, string(body)))
	}

	if method == http.MethodDelete {
		return success(nil)
	}

	var record e.Blog
	if err := json.NewDecoder(res.Body).Decode(&record); err != nil {
		return failure(codes.Internal, fmt.Sprintf("failed to decode backend response: %v", err))
	}
	return success(&record)
}

func failure(code codes.Code, message string) (*pbs.ProcessMessageResponse, error) {
	return &pbs.ProcessMessageResponse{
		Status: &pb.Status{
			Success: false,
			Code:    code.String(),
			Message: message,
		},
	}, status.Error(code, message)
}

func success(record *e.Blog) (*pbs.ProcessMessageResponse, error) {
	resp := &pbs.ProcessMessageResponse{
		Status: &pb.Status{
			Success: true,
			Code:    codes.OK.String(),
			Message: "Request processed successfully",
		},
	}

	fmt.Println(record)

	if record != nil {
		resp.Data = &pb.Blog{
			Id:        record.ID,
			Title:     record.Title,
			Content:   record.Content,
			CreatedAt: timestamppb.New(record.UpdatedAt),
			UpdatedAt: timestamppb.New(record.UpdatedAt),
		}
	}

	return resp, nil
}

func NewBlogService(repo r.BlogRepository) BlogService {
	return &srv{repo: repo}
}
