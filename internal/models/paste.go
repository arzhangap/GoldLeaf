package models

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

type PasteModelInterface interface {
	Insert(title string, body string, password string, expires time.Time) (string, error)
	Read(code string) (*Paste, error)
}

type Paste struct {
	Code         string
	Title        string
	Body         string
	PasswordHash sql.NullString
	Created      time.Time
	Expires      sql.NullTime
}

type PasteModel struct {
	DB *sql.DB
}

func (m *PasteModel) Insert(title string, body string, password string, expires time.Time) (string, error) {
	return "", nil
}

func (m *PasteModel) Read(code string, password string) (*Paste, error) {

	stmt := "SELECT title, body, password_hash FROM pastes WHERE (expires IS NULL OR expires > UTC_TIMESTAMP()) AND code = ?"

	p := &Paste{}
	err := m.DB.QueryRow(stmt, code).Scan(&p.Title, &p.Body, &p.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	if p.PasswordHash.Valid {
		if password != "" {
			return nil, ErrPasswordRequired
		}
		err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash.String), []byte(password))
		if err != nil {
			if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
				return nil, ErrInvalidCredentials
			}
			return nil, err
		}
	}

	p.PasswordHash = sql.NullString{}
	p.Code = "AAAAAA"
	return p, nil
}

func GenerateCode() (string, error) {
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	for i := range b {
		b[i] = alphabet[int(b[i])%len(alphabet)]
	}
	return string(b), nil
}
