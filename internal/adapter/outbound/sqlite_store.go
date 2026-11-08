package outbound

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/ahmad-banikamali/internal/core/domain"
	"github.com/ahmad-banikamali/internal/core/port"
)

var _ port.LinkStore = (*SQLiteStore)(nil)

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1) // SQLite allows one writer; avoids "database is locked"

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS links (
		code       TEXT PRIMARY KEY,
		url        TEXT NOT NULL,
		created_at INTEGER NOT NULL
	)`)
	if err != nil {
		db.Close()
		return nil, err
	}
	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) Close() error { return s.db.Close() }

func (s *SQLiteStore) Save(ctx context.Context, l domain.Link) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO links (code, url, created_at) VALUES (?, ?, ?)`,
		l.Code, l.URL, l.CreatedAt.Unix())
	if err != nil {
		// Translate the adapter-specific error into the domain error.
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return domain.ErrCodeTaken
		}
		return err
	}
	return nil
}

func (s *SQLiteStore) FindByCode(ctx context.Context, code string) (domain.Link, error) {
	var l domain.Link
	var created int64
	err := s.db.QueryRowContext(ctx,
		`SELECT code, url, created_at FROM links WHERE code = ?`, code).
		Scan(&l.Code, &l.URL, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Link{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.Link{}, err
	}
	l.CreatedAt = time.Unix(created, 0)
	return l, nil
}
