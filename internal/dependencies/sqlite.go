package dependencies

import (
	"context"
	"database/sql"

	"github.com/ashtonx86/nybl/internal/supervisor"

	_ "github.com/mattn/go-sqlite3"
)

var _ supervisor.Singleton = (*SQLiteSingleton)(nil)

type SQLiteSingleton struct {
	DB *sql.DB
}

func NewSQLiteSingleton() *SQLiteSingleton {
	return &SQLiteSingleton{}
}

func (s *SQLiteSingleton) Init(ctx context.Context) error {
	db, err := sql.Open("sqlite3", "./dev.db")
	if err != nil {
		return err
	}

	initSqlStmt := `
	CREATE TABLE IF NOT EXISTS accounts (
		id VARCHAR(10) NOT NULL PRIMARY KEY, 
		name VARCHAR(50) NOT NULL,
		emailID TEXT NOT NULL UNIQUE,
		avatarURL TEXT,
		verified BOOLEAN NOT NULL,
		createdAt TEXT,
		updatedAt TEXT
	);

	CREATE TABLE IF NOT EXISTS emailOTPs (
		id VARCHAR(10) NOT NULL PRIMARY KEY,
		code TEXT NOT NULL UNIQUE,
		emailID TEXT NOT NULL UNIQUE,
		accountID TEXT NOT NULL UNIQUE,
		requestedAt TEXT
	);

	CREATE TABLE IF NOT EXISTS spaces (
		id VARCHAR(10) NOT NULL PRIMARY KEY,
		name TEXT NOT NULL,
		about TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS members (
		id VARCHAR(10) NOT NULL PRIMARY KEY,
		spaceID VARCHAR(10) NOT NULL,
		role TEXT NOT NULL DEFAULT 'normie',

		notifyAllMentions BOOLEAN NOT NULL,
		notifyAllMessages BOOLEAN NOT NULL,

		joinedAt TEXT NOT NULL,

		FOREIGN KEY (id) REFERENCES accounts(id) ON DELETE CASCADE,
		FOREIGN KEY (spaceID) REFERENCES spaces(id) ON DELETE CASCADE,

		UNIQUE (id, spaceID) 
	);
	`
	_, err = db.ExecContext(ctx, initSqlStmt)
	if err != nil {
		return err
	}

	s.DB = db
	return nil
}

func (s *SQLiteSingleton) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}

func MustGetSQLite(su *supervisor.Supervisor) *SQLiteSingleton {
	sqlite, ok := supervisor.GetSingletonAs[*SQLiteSingleton](su, "sqlite")

	if !ok {
		panic("Missing required dependency: sqlite (SQLiteSingleton) -- dependency is not registered")
	}

	if sqlite.DB == nil {
		panic("SQLiteSingleton.DB must not be nil")
	}

	return sqlite
}
