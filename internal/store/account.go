package store

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/ashtonx86/nybl/internal/schemas"
	"github.com/google/uuid"
)

type AccountStore interface {
	Create(ctx context.Context, name string, email string) (*schemas.Account, error)
	// Get(ctx context.Context, selectField string, selectValue string) (error, *schemas.Account)

	// PatchVerified(ctx context.Context, accountID string) (error, *schemas.Account)

	// Delete(ctx context.Context, accountID string, emailID string) error
}

type SQLiteAccountStore struct {
	DB *sql.DB
}

func NewSQLiteAccountStore(db *sql.DB) *SQLiteAccountStore {
	log.Print("sqlite account store created, db: ", db)
	return &SQLiteAccountStore{
		DB: db,
	}
}

func (s *SQLiteAccountStore) Create(ctx context.Context, name string, emailID string) (*schemas.Account, error) {
	id := uuid.NewString()

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})

	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	stmt := `
	INSERT INTO accounts (id, name, emailID, verified, createdAt, updatedAt)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	timeNow := time.Now()

	_, err = tx.ExecContext(ctx, stmt, id, name, emailID, false, timeNow, timeNow)
	if err != nil {
		return nil, err
	}

	account := &schemas.Account{
		ID:       id,
		Name:     name,
		EmailID:  emailID,
		Verified: false,

		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	return account, nil
}
