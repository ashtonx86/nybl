package store

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ashtonx86/nybl/internal/diabeticerrors"
	"github.com/ashtonx86/nybl/internal/schemas"
	"github.com/google/uuid"
	"github.com/mattn/go-sqlite3"
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
	return &SQLiteAccountStore{
		DB: db,
	}
}
func (s *SQLiteAccountStore) Create(ctx context.Context, name string, emailID string) (account *schemas.Account, err error) {
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

	var sqliteErr sqlite3.Error
	if err != nil && errors.As(err, &sqliteErr) {
		if sqliteErr.Code == sqlite3.ErrConstraint && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return nil, diabeticerrors.AlreadyExistsError("sqlite", err)
		}
		return nil, err
	}

	account = &schemas.Account{
		ID:        id,
		Name:      name,
		EmailID:   emailID,
		Verified:  false,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
	}
	return account, nil
}
