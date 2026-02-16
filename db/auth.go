package db

import (
	"crypto/rand"
	"encoding/hex"
	"time"
)

func (db *DB) CreateSession(duration time.Duration) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	token := hex.EncodeToString(bytes)
	expires := time.Now().Add(duration)

	_, err := db.Conn.Exec(
		`INSERT INTO sessions (token, expires_at) VALUES (?, ?)`,
		token, expires,
	)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (db *DB) ValidateSession(token string) bool {
	var count int
	err := db.Conn.QueryRow(
		`SELECT COUNT(*) FROM sessions
		 WHERE token = ? AND expires_at > ?`,
		token, time.Now(),
	).Scan(&count)
	return err == nil && count > 0
}

func (db *DB) DeleteSession(token string) error {
	_, err := db.Conn.Exec(
		`DELETE FROM sessions WHERE token = ?`, token,
	)
	return err
}

func (db *DB) PurgeExpiredSessions() error {
	_, err := db.Conn.Exec(
		`DELETE FROM sessions WHERE expires_at < ?`, time.Now(),
	)
	return err
}
