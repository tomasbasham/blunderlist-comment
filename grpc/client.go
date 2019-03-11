package grpc

//go:generate protoc -I ./proto --go_out=plugins=grpc:. comment.proto

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	commentpb "github.com/tomasbasham/blunderlist-comment/blunderlist_comment_v1"
	"google.golang.org/grpc"

	"github.com/tomasbasham/grpc-service-go/option"
	transport "github.com/tomasbasham/grpc-service-go/transport/grpc"
)

// Client is a wrapper type around a client connection to an RPC server and a
// concrete instance of a `CommentClient` type.
type Client struct {
	conn   *grpc.ClientConn
	client commentpb.CommentClient
}

// NewClientWithTarget creates a new client connection to an RPC server at a
// specific endpoint.
func NewClientWithTarget(ctx context.Context, target string) (*Client, error) {
	return NewClient(ctx, option.WithEndpoint(target))
}

// NewClient creates a new client connectino to an RPC server accepting
// multiple connection options.
func NewClient(ctx context.Context, opts ...option.ClientOption) (*Client, error) {
	conn, err := transport.DialInsecure(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn:   conn,
		client: commentpb.NewCommentClient(conn),
	}, nil
}

// Conn returns a client connection to an RPC server.
func (c *Client) Conn() *grpc.ClientConn {
	return c.conn
}

// Close closes the client connection to an RPC server. If the connection has
// not been made then this method does nothing.
func (c *Client) Close() {
	c.conn.Close()
}

// ListComments delegates to the client, returning a type able to fetch
// comments from a stream.
func (c *Client) ListComments(ctx context.Context, in *commentpb.CommentListRequest, opts ...grpc.CallOption) (commentpb.Comment_ListCommentsClient, error) {
	return c.client.ListComments(ctx, in, opts...)
}

// GetComment delegates to the client, returning a comment with the given id.
// If the comment could not be found then an error is returned.
func (c *Client) GetComment(ctx context.Context, in *commentpb.CommentQuery, opts ...grpc.CallOption) (*commentpb.CommentResponse, error) {
	return c.client.GetComment(ctx, in, opts...)
}

// CreateComment delegates to the client, returning a persisted comment with
// the given information. If the comment could not be persisted then an error
// is returned.
func (c *Client) CreateComment(ctx context.Context, in *commentpb.CommentCreateRequest, opts ...grpc.CallOption) (*commentpb.CommentResponse, error) {
	return c.client.CreateComment(ctx, in, opts...)
}

// UpdateComment delegates to the client, returning a persisted comment with
// the given information. If the comment could not be found, or was unable to
// be persisted then an error is returned.
func (c *Client) UpdateComment(ctx context.Context, in *commentpb.CommentUpdateRequest, opts ...grpc.CallOption) (*commentpb.CommentResponse, error) {
	return c.client.UpdateComment(ctx, in, opts...)
}

// DeleteComment delegates to the client, removing the comment with the given
// id. If the comment could not be found, or it could not be deleted then as
// error is returned.
func (c *Client) DeleteComment(ctx context.Context, in *commentpb.CommentQuery, opts ...grpc.CallOption) (*empty.Empty, error) {
	return c.client.DeleteComment(ctx, in, opts...)
}
