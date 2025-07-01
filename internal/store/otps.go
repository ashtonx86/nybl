package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	"github.com/ashtonx86/nybl/internal/schemas"
	"github.com/google/uuid"
)

type OTPStore interface {
	Create(ctx context.Context, emailID string, accountID string) (*schemas.OTP, error)
	Verify(ctx context.Context, code string) (*schemas.OTP, bool, error)

	DeleteByID(ctx context.Context, id string) error
}

type SQLiteOTPStore struct {
	DB *sql.DB
}

func NewSQLiteOTPStore(db *sql.DB) *SQLiteOTPStore {
	return &SQLiteOTPStore{
		DB: db,
	}
}

func (s *SQLiteOTPStore) Create(ctx context.Context, emailID string, accountID string) (*schemas.OTP, error) {
	id := uuid.NewString()

	query := `
	SELECT id FROM emailOTPs WHERE accountID = ? AND emailID = ?
	`
	row := s.DB.QueryRowContext(ctx, query, accountID, emailID)

	if err := row.Err(); err != nil {
		if err != sql.ErrNoRows {
			return nil, err // return if it's not ErrNoRows error -- we only need to verify if a row with similar accountID already exists.
		}
	}

	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
	})
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	if err != nil {
		return nil, fmt.Errorf("tx init failed: %v", err)
	}

	otp, err := generateOTPCode()
	if err != nil {
		return nil, fmt.Errorf("otp generation failed: %v", err)
	}

	timeNow := time.Now()
	stmt := `
	INSERT INTO emailOTPs (id, code, emailID, accountID, requestedAt)
	VALUES (?, ?, ?, ?, ?)
	`
	_, err = tx.ExecContext(ctx, stmt, id, otp, emailID, accountID, timeNow)
	if err != nil {
		return nil, fmt.Errorf("tx execution failed: %v", err)
	}

	newOTP := &schemas.OTP{
		ID:          id,
		Code:        otp,
		AccountID:   accountID,
		EmailID:     emailID,
		RequestedAt: timeNow,
	}
	return newOTP, err
}

func (s *SQLiteOTPStore) Verify(ctx context.Context, code string) (*schemas.OTP, bool, error) {
	rows, err := s.DB.QueryContext(ctx, "SELECT * FROM emailOTPs WHERE code = ?", code)

	if err != nil {
		return nil, false, fmt.Errorf("query failed: %v", err)
	}

	for rows.Next() {
		var otp schemas.OTP

		if err := rows.Scan(&otp.ID, &otp.Code, &otp.EmailID, &otp.AccountID, &otp.RequestedAt); err != nil {
			delErr := s.DeleteByID(ctx, otp.ID)
			if delErr != nil {
				return nil, false, fmt.Errorf("failed to delete: %v", delErr)
			}
			return nil, false, fmt.Errorf("scanning rows failed: %v", err)
		} else {
			return &otp, true, nil
		}
	}

	if err := rows.Err(); err != nil {
		return nil, false, err
	}
	return nil, false, nil
}

func (s *SQLiteOTPStore) DeleteByID(ctx context.Context, id string) error {
	tx, err := s.DB.BeginTx(ctx, &sql.TxOptions{})
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM emailOTPs WHERE id = ?", id)
	return err
}

func generateOTPCode() (string, error) {
	digits := make([]byte, 6)
	if _, err := rand.Read(digits); err != nil {
		return "", err
	}
	for i := 0; i < 6; i++ {
		digits[i] = uint8(48 + (digits[i] % 10))
	}

	return string(digits), nil
}
