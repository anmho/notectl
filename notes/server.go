package notes

import (
	"context"
	"database/sql"
	"fmt"
	pb "github.com/anmho/notectl/gen/proto/notes"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"time"
)

type Service struct {
	pb.UnsafeNoteServiceServer
	db *sql.DB
}

// Ensure implements the interface
var _ pb.NoteServiceServer = Service{}

func NewService(db *sql.DB) Service {
	return Service{
		db: db,
	}
}

func (s Service) CreateNote(ctx context.Context, r *pb.CreateNoteRequest) (*pb.Note, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("opening conn: %w", err)
	}
	defer conn.Close()
	stmt, err := conn.PrepareContext(ctx,
		`INSERT INTO notes (
	              id, title, content, createdAt, updatedAt
			   ) VALUES ($1, $2, $3, $4, $5)
			   RETURNING *;`)
	if err != nil {
		return nil, fmt.Errorf("preparing statement: %w", err)
	}

	id := uuid.NewString()
	rows, err := stmt.QueryContext(ctx, id, r.Title, r.Content, time.Now(), time.Now())
	if err != nil {
		return nil, fmt.Errorf("executing insert: %w", err)
	}
	defer rows.Close()
	note := pb.Note{}
	if rows.Next() {
		err := scanNote(rows, &note)
		if err != nil {
			return nil, fmt.Errorf("scanning note struct: %w", err)
		}
	}
	return &note, nil
}

func (s Service) GetNote(ctx context.Context, r *pb.GetNoteRequest) (*pb.Note, error) {
	log.Println("getting note with id", r.Id)
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating conn to get note: %w", err)
	}
	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, title, content, createdAt, updatedAt 
				FROM notes
				WHERE id = $1
				;`)
	if err != nil {
		return nil, fmt.Errorf("preparing query: %w", err)
	}
	rows, err := stmt.QueryContext(ctx, r.Id)
	if err != nil {
		return nil, fmt.Errorf("querying note id=%s: %w", r.Id, err)
	}

	note := pb.Note{}

	if rows.Next() {
		err := scanNote(rows, &note)
		if err != nil {
			return nil, fmt.Errorf("scanning note row id=%s: %w", r.Id, err)
		}
	}

	return &note, nil
}

func (s Service) ListNotes(ctx context.Context, empty *emptypb.Empty) (*pb.NoteList, error) {
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "getting db connection: %w", err)
	}

	stmt, err := conn.PrepareContext(ctx,
		`SELECT id, title, content, createdAt, updatedAt 
				FROM notes`,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "preparing stmt: %w", err)
	}
	rows, err := stmt.QueryContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "executing query: %w", err)
	}

	var notes pb.NoteList
	for rows.Next() {
		var note pb.Note
		err := scanNote(rows, &note)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "scanning note: %w", err)
		}
		notes.Notes = append(notes.Notes, &note)
	}

	return &notes, nil
}

func (s Service) DeleteNote(ctx context.Context, request *pb.DeleteNoteRequest) (*emptypb.Empty, error) {
	log.Println("deleting note with id", request.Id)
	conn, err := s.db.Conn(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "connecting db: %w", err)
	}
	stmt, err := conn.PrepareContext(ctx, "DELETE FROM notes WHERE id = $1")
	if err != nil {
		return nil, err
	}

	exec, err := stmt.Exec(request.Id)
	if err != nil {
		return nil, err
	}
	_, err = exec.RowsAffected()
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func scanNote(rows *sql.Rows, note *pb.Note) error {
	createdAt := time.Time{}
	updatedAt := time.Time{}
	err := rows.Scan(&note.Id, &note.Title, &note.Content, &createdAt, &updatedAt)
	if err != nil {
		return err
	}
	note.CreatedAt = timestamppb.New(createdAt)
	note.UpdatedAt = timestamppb.New(updatedAt)
	return nil
}
