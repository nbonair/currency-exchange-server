package repo

import (
	"github.com/nbonair/currency-exchange-server/internal/dataaccess/db"
	"github.com/nbonair/currency-exchange-server/internal/database"
)

type UserRepository interface{}

type userRepository struct {
	db      *db.Database
	queries *database.Queries
}

func NewUserRepository(db *db.Database) UserRepository {
	return &userRepository{
		db:      db,
		queries: database.New(db.Pool),
	}
}
