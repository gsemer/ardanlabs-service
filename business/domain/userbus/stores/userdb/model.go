package userdb

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type dbUser struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        []string       `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	Department   sql.NullString `db:"department"`
	Enabled      bool           `db:"enabled"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}
