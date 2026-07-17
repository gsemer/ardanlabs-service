package userdb

import (
	"database/sql"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gsemer/ardanlabs-service/business/domain/userbus"
)

type user struct {
	ID           uuid.UUID      `db:"user_id"`
	Name         string         `db:"name"`
	Email        string         `db:"email"`
	Roles        string         `db:"roles"`
	PasswordHash []byte         `db:"password_hash"`
	Department   sql.NullString `db:"department"`
	Enabled      bool           `db:"enabled"`
	DateCreated  time.Time      `db:"date_created"`
	DateUpdated  time.Time      `db:"date_updated"`
}

func toDBUser(usr userbus.User) user {
	roleNames := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roleNames[i] = role.Name()
	}

	return user{
		ID:           usr.ID,
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        strings.Join(roleNames, ","),
		PasswordHash: usr.PasswordHash,
		Department: sql.NullString{
			String: usr.Department,
			Valid:  usr.Department != "",
		},
		Enabled:     usr.Enabled,
		DateCreated: usr.DateCreated.UTC(),
		DateUpdated: usr.DateUpdated.UTC(),
	}
}

func toBusUser(dbUsr user) (userbus.User, error) {
	addr := mail.Address{
		Address: dbUsr.Email,
	}

	roleNames := []string{}
	if dbUsr.Roles != "" {
		roleNames = strings.Split(dbUsr.Roles, ",")
	}

	roles := make([]userbus.Role, len(roleNames))
	for i, value := range roleNames {
		role, err := userbus.ParseRole(strings.TrimSpace(value))
		if err != nil {
			return userbus.User{}, fmt.Errorf("parse role %q: %w", value, err)
		}
		roles[i] = role
	}

	return userbus.User{
		ID:           dbUsr.ID,
		Name:         dbUsr.Name,
		Email:        addr,
		Roles:        roles,
		PasswordHash: dbUsr.PasswordHash,
		Enabled:      dbUsr.Enabled,
		Department:   dbUsr.Department.String,
		DateCreated:  dbUsr.DateCreated.In(time.Local),
		DateUpdated:  dbUsr.DateUpdated.In(time.Local),
	}, nil
}

func toBusUsers(dbUsers []user) ([]userbus.User, error) {
	users := make([]userbus.User, len(dbUsers))

	for i, dbUsr := range dbUsers {
		u, err := toBusUser(dbUsr)
		if err != nil {
			return nil, err
		}
		users[i] = u
	}

	return users, nil
}
