package grpc

import (
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	commentpb "github.com/tomasbasham/blunderlist-comment/blunderlist_comment_v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/tomasbasham/blunderlist-comment/internal"
)

// Server implements the Comment server.
type Server struct {
	logger *log.Logger
	store  internal.Repository
}

// NewServer returns a Server that implements the Comment service.
func NewServer(logger *log.Logger, store internal.Repository) *Server {
	return &Server{
		logger: logger,
		store:  store,
	}
}

// Serve initialises a gRPC server and has it accept and handle incoming
// requests.
func (s *Server) Serve(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	commentpb.RegisterCommentServer(server, s)
	healthpb.RegisterHealthServer(server, s)

	// Register reflection service on gRPC server for debugging.
	reflection.Register(server)

	s.logger.Printf("server started on %s", lis.Addr())
	return server.Serve(lis)
}
