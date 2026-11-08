package models

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
)

var expiryOptions = map[string]time.Duration{
	"1h": time.Hour,
	"1d": 24 * time.Hour,
	"7d": 7 * 24 * time.Hour,
}

const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func parseExpires(value string) (*time.Time, bool) {
	d, ok := expiryOptions[value]
	if !ok {
		return nil, false
	}
	t := time.Now().UTC().Add(d)
	return &t, true
}

type PasteModelInterface interface {
	Insert(title string, body string, password string, expires string) (string, error)
	Read(code, password string) (*Paste, error)
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

func (m *PasteModel) Insert(title string, body string, password string, expires string) (string, error) {
	parsedExpiry, ok := parseExpires(expires)
	if !ok {
		return "", errors.New("cannot parse expiry.")
	}

	for range 3 {
		var mysqlErr *mysql.MySQLError
		code, err := GenerateCode()

		if err != nil {
			return "", err
		}

		stmt := `INSERT INTO pastes (code, title, body,created, expires)
	VALUES(?, ?, ?, UTC_TIMESTAMP(), ?)`

		_, err = m.DB.Exec(stmt, code, title, body, parsedExpiry)
		if err != nil {
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				continue
			}
			return "", err
		}
		return code, nil
	}
	return "", errors.New("failed to generate a unique code after 3 attempts")
}

func (m *PasteModel) Read(code string, password string) (*Paste, error) {

	stmt := "SELECT title, body, password_hash, created FROM pastes WHERE (expires IS NULL OR expires > UTC_TIMESTAMP()) AND code = ?"

	p := &Paste{Code: code}
	err := m.DB.QueryRow(stmt, code).Scan(&p.Title, &p.Body, &p.PasswordHash, &p.Created)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	// TODO password protected
	//if p.PasswordHash.Valid {
	//	if password != "" {
	//		return nil, ErrPasswordRequired
	//	}
	//	err := bcrypt.CompareHashAndPassword([]byte(p.PasswordHash.String), []byte(password))
	//	if err != nil {
	//		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
	//			return nil, ErrInvalidCredentials
	//		}
	//		return nil, err
	//	}
	//}

	p.PasswordHash = sql.NullString{}
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
