package pg

import (
	"database/sql"
	"time"

	// Initialise the pq package.
	_ "github.com/lib/pq"

	pb "github.com/golang/protobuf/ptypes"
	commentpb "github.com/tomasbasham/blunderlist-comment/blunderlist_comment_v1"
)

// Storage implements the Repository interface, with an active connection to a
// PostgreSQL instance.
type Storage struct {
	db *sql.DB
}

// NewStorage returns a Storage that implements the Repository interface.
func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// GetComments makes a SQL query to retreive all comments from the database.
func (s *Storage) GetComments(id uint) []*commentpb.CommentResponse {
	rows, err := s.db.Query("SELECT * FROM comments WHERE parent_id=$1 ORDER BY created_at", id)
	if err != nil {
		return []*commentpb.CommentResponse{}
	}

	var comments []*commentpb.CommentResponse

	for rows.Next() {
		comment, err := scan(rows)
		if err != nil {
			continue
		}

		comments = append(comments, comment)
	}

	return comments
}

// GetComment makes a SQL query to retreive a single comment from the database.
func (s *Storage) GetComment(id uint) (*commentpb.CommentResponse, error) {
	row := s.db.QueryRow("SELECT * FROM comments WHERE id=$1", id)

	comment, err := scan(row)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

// CreateComment makes a SQL query to insert a new comment into the database.
func (s *Storage) CreateComment(comment *commentpb.CommentCreateRequest) (*commentpb.CommentResponse, error) {
	var lastInsertID uint

	err := s.db.QueryRow("INSERT INTO comments(parent_id, text, created_at) VALUES($1, $2, $3) RETURNING id", comment.ParentId, comment.Text, time.Now()).Scan(&lastInsertID)
	if err != nil {
		return nil, err
	}

	return s.GetComment(lastInsertID)
}

// UpdateComment makes a SQL query to update a comment already present in the
// database.
func (s *Storage) UpdateComment(comment *commentpb.CommentUpdateRequest) (*commentpb.CommentResponse, error) {
	var lastUpdateID uint

	err := s.db.QueryRow("UPDATE comments SET text=$1, WHERE id=$2 RETURNING id", comment.Text, uint(comment.Id)).Scan(&lastUpdateID)
	if err != nil {
		return nil, err
	}

	return s.GetComment(lastUpdateID)
}

// DeleteComment makes a SQL query to remove a comment from the database.
func (s *Storage) DeleteComment(id uint) error {
	if _, err := s.db.Exec("DELETE FROM comments WHERE id=$1", id); err != nil {
		return err
	}

	return nil
}

// IsAvailable enusues there is an active connection to the database,
// establishing a connection if necessary.
func (s *Storage) IsAvailable() error {
	if err := s.db.Ping(); err != nil {
		return err
	}

	// This is necessary because once a database connection is initiallly
	// established subsequent calls to the Ping method return success even if the
	// database goes down.
	//
	// https://github.com/lib/pq/issues/533
	if _, err := s.db.Exec("SELECT 1"); err != nil {
		return err
	}

	return nil
}

type resultScanner interface {
	Scan(dest ...interface{}) error
}

func scan(r resultScanner) (*commentpb.CommentResponse, error) {
	var id uint
	var parentID uint
	var text string
	var createdAt time.Time

	if err := r.Scan(&id, &parentID, &text, &createdAt); err != nil {
		return nil, err
	}

	timestamp, err := pb.TimestampProto(createdAt)
	if err != nil {
		return nil, err
	}

	return &commentpb.CommentResponse{
		Id:         uint64(id),
		ParentId:   uint64(parentID),
		Text:       text,
		CreateTime: timestamp,
	}, nil
}
