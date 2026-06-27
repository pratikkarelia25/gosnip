package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var ErrNotFound = errors.New("short code not found")

type UrlMapping struct {
	ShortCode string     `json:"short_code"`
	LongUrl   string     `json:"long_url"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

type Store struct {
	db *sql.DB
}

func New(connStr string) (*Store, error) {
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func NewFromParts(host, user, password, dbname string) (*Store, error) {
	// Build DSN using key=value format to avoid URL-encoding issues with special chars in passwords.
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=require", host, user, password, dbname)
	return New(connStr)
}

func (s *Store) Migrate() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS url_mappings (
		id BIGSERIAL PRIMARY KEY,
		short_code TEXT NOT NULL UNIQUE,
		long_url TEXT NOT NULL,
		expires_at BIGINT
	);
	`)
	return err
}

// InsertMapping stores a new short code -> long URL mapping. A nil
// expiresAt means the mapping never expires.
func (s *Store) InsertMapping(shortCode, longUrl string, expiresAt *time.Time) error {
	var expiresAtUnix any
	if expiresAt != nil {
		expiresAtUnix = expiresAt.Unix()
	}

	_, err := s.db.Exec(
		`INSERT INTO url_mappings (short_code, long_url, expires_at) VALUES ($1, $2, $3)`,
		shortCode, longUrl, expiresAtUnix,
	)
	return err
}

// ShortCodeTaken reports whether a short code is already in use,
// regardless of whether the mapping has expired.
func (s *Store) ShortCodeTaken(shortCode string) (bool, error) {
	var count int
	err := s.db.QueryRow(`SELECT COUNT(1) FROM url_mappings WHERE short_code = $1`, shortCode).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetActiveLongUrl returns the long URL for a short code, as long as
// it exists and has not expired. ErrNotFound is returned otherwise.
func (s *Store) GetActiveLongUrl(shortCode string) (string, error) {
	var longUrl string
	var expiresAt sql.NullInt64

	row := s.db.QueryRow(`SELECT long_url, expires_at FROM url_mappings WHERE short_code = $1`, shortCode)
	if err := row.Scan(&longUrl, &expiresAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}

	if expiresAt.Valid && time.Now().Unix() >= expiresAt.Int64 {
		return "", ErrNotFound
	}

	return longUrl, nil
}

// FindActiveShortCode returns the short code already mapped to a long
// URL, if one exists and has not expired.
func (s *Store) FindActiveShortCode(longUrl string) (string, error) {
	var shortCode string
	var expiresAt sql.NullInt64

	rows, err := s.db.Query(`SELECT short_code, expires_at FROM url_mappings WHERE long_url = $1`, longUrl)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	now := time.Now().Unix()
	for rows.Next() {
		if err := rows.Scan(&shortCode, &expiresAt); err != nil {
			return "", err
		}
		if !expiresAt.Valid || now < expiresAt.Int64 {
			return shortCode, nil
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	return "", ErrNotFound
}

func (s *Store) Close() error {
	return s.db.Close()
}
