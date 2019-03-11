package internal

import commentpb "github.com/tomasbasham/blunderlist-comment/blunderlist_comment_v1"

// Repository provides the storage requirements through which to access
// comments.
type Repository interface {
	GetComments(uint) []*commentpb.CommentResponse
	GetComment(uint) (*commentpb.CommentResponse, error)
	CreateComment(*commentpb.CommentCreateRequest) (*commentpb.CommentResponse, error)
	UpdateComment(*commentpb.CommentUpdateRequest) (*commentpb.CommentResponse, error)
	DeleteComment(uint) error

	// Used to check of the storge media is ready to use. If not this method
	// should return an error.
	IsAvailable() error
}

// Store proxies requests to the underlying storage implementation.
type Store struct {
	repo Repository
}

// NewStore returns a Store containing a repository.
func NewStore(repo Repository) *Store {
	return &Store{repo: repo}
}

// GetComments proxies to the repo.
func (s *Store) GetComments(parentID uint) []*commentpb.CommentResponse {
	return s.repo.GetComments(parentID)
}

// GetComment proxies to the repo.
func (s *Store) GetComment(id uint) (*commentpb.CommentResponse, error) {
	return s.repo.GetComment(id)
}

// CreateComment proxies to the repo.
func (s *Store) CreateComment(comment *commentpb.CommentCreateRequest) (*commentpb.CommentResponse, error) {
	return s.repo.CreateComment(comment)
}

// UpdateComment proxies to the repo.
func (s *Store) UpdateComment(comment *commentpb.CommentUpdateRequest) (*commentpb.CommentResponse, error) {
	return s.repo.UpdateComment(comment)
}

// DeleteComment proxies to the repo.
func (s *Store) DeleteComment(id uint) error {
	return s.repo.DeleteComment(id)
}

// IsAvailable proxies to the repo.
func (s *Store) IsAvailable() error {
	return s.repo.IsAvailable()
}
