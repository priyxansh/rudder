package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// DB wraps the embedded SQLite database connection for task and chat history.
type DB struct {
	conn *sql.DB
}

// Open initializes SQLite at the given path (defaulting to ~/.rudder/history.db) and applies migrations.
func Open(path string) (*DB, error) {
	resolvedPath, err := resolvePath(path)
	if err != nil {
		return nil, fmt.Errorf("resolve db path: %w", err)
	}

	dir := filepath.Dir(resolvedPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, fmt.Errorf("create db directory %s: %w", dir, err)
	}

	dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)", resolvedPath)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("run migrations: %w", err)
	}

	return db, nil
}

// Close shuts down the SQLite connection pool.
func (d *DB) Close() error {
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}

// migrate creates the tasks and messages tables and indexes if missing.
func (d *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS tasks (
		id           TEXT PRIMARY KEY,
		repo_path    TEXT NOT NULL,
		prompt       TEXT NOT NULL,
		status       TEXT NOT NULL,
		started_at   DATETIME NOT NULL,
		completed_at DATETIME,
		stats_adds   INTEGER DEFAULT 0,
		stats_dels   INTEGER DEFAULT 0
	);

	CREATE TABLE IF NOT EXISTS messages (
		id         TEXT PRIMARY KEY,
		task_id    TEXT REFERENCES tasks(id) ON DELETE CASCADE,
		role       TEXT NOT NULL,
		content    TEXT NOT NULL,
		diff_json  TEXT,
		created_at DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_repo    ON tasks(repo_path);
	CREATE INDEX IF NOT EXISTS idx_messages_task ON messages(task_id);
	`
	_, err := d.conn.Exec(schema)
	return err
}

// resolvePath expands home directory tildes and resolves relative paths to absolute paths.
func resolvePath(path string) (string, error) {
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".rudder", "history.db"), nil
	}

	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, path[2:]), nil
	}

	return filepath.Abs(path)
}
