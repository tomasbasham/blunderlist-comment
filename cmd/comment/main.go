package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/tomasbasham/blunderlist-comment/internal"
	"github.com/tomasbasham/blunderlist-comment/internal/grpc"
	"github.com/tomasbasham/blunderlist-comment/internal/storage/pg"
)

func main() {
	logger := log.New(os.Stdout, "", log.Lshortfile)

	connectionString := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("BLUNDERLIST_COMMENT_GCLOUD_SQLPROXY_SERVICE_HOST"),
		os.Getenv("BLUNDERLIST_COMMENT_GCLOUD_SQLPROXY_SERVICE_PORT"),
		os.Getenv("PGDATABASE"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		logger.Fatalf("unable to connect to the database: %s", err)
	}

	defer db.Close()

	// Setup dependencies for the service.
	storage := pg.NewStorage(db)
	service := internal.NewStore(storage)

	commentService := grpc.NewServer(logger, service)
	commentServicePort := os.Getenv("BLUNDERLIST_COMMENT_SERVICE_PORT")

	if err := commentService.Serve(commentServicePort); err != nil {
		logger.Fatalf("failed to serve: %v", err)
	}
}
