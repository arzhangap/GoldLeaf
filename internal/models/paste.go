package models

import (
	"crypto/rand"
	"database/sql"
	"time"
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

func (m *PasteModel) Read(code string) ([]*Paste, error) {
	return nil, nil
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
