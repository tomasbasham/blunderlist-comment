package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	//details "google.golang.org/genproto/googleapis/rpc/errordetails"

	commentpb "github.com/tomasbasham/blunderlist-comment/blunderlist_comment_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// ListComments streams all comments from the store.
func (s *Server) ListComments(in *commentpb.CommentListRequest, stream commentpb.Comment_ListCommentsServer) error {
	for _, comment := range s.store.GetComments(uint(in.ParentId)) {
		if err := stream.Send(comment); err != nil {
			return status.Error(codes.Unknown, "Unknown streaming error")
		}
	}

	return nil
}

// GetComment returns a single comment from the store.
func (s *Server) GetComment(ctx context.Context, in *commentpb.CommentQuery) (*commentpb.CommentResponse, error) {
	comment, err := s.store.GetComment(uint(in.Id))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "comment not found with id %d", in.Id)
	}

	return comment, nil
}

// CreateComment creates and returns a new comment.
func (s *Server) CreateComment(ctx context.Context, in *commentpb.CommentCreateRequest) (*commentpb.CommentResponse, error) {
	comment, err := s.store.CreateComment(in)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, "failed to create comment")
	}

	return comment, nil
}

// UpdateComment updates and returns an existing comment.
func (s *Server) UpdateComment(ctx context.Context, in *commentpb.CommentUpdateRequest) (*commentpb.CommentResponse, error) {
	comment, err := s.store.UpdateComment(in)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "comment not found with id %d", in.Id)
	}

	return comment, nil
}

// DeleteComment deletes a comment from the store.
func (s *Server) DeleteComment(ctx context.Context, in *commentpb.CommentQuery) (*empty.Empty, error) {
	if err := s.store.DeleteComment(uint(in.Id)); err != nil {
		return &empty.Empty{}, status.Errorf(codes.NotFound, "comment not found with if %d", in.Id)
	}

	return &empty.Empty{}, nil
}

// Check is used by clients to know when the service is ready.
func (s *Server) Check(context.Context, *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	if err := s.store.IsAvailable(); err != nil {
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}, nil
	}

	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

// Watch is used by clients to receive updates when the service status changes.
// Here it has no meaningful implementation just to satisfy the interface.
func (s *Server) Watch(*healthpb.HealthCheckRequest, healthpb.Health_WatchServer) error {
	return nil
}
