package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/priyxansh/rudder/daemon/internal/models"
)

// CreateTask inserts a new agent execution record when a prompt is dispatched.
func (d *DB) CreateTask(ctx context.Context, task *models.Task) error {
	query := `
	INSERT INTO tasks (id, repo_path, prompt, status, started_at, completed_at, stats_adds, stats_dels)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := d.conn.ExecContext(ctx, query,
		task.ID,
		task.RepoPath,
		task.Prompt,
		string(task.Status),
		task.StartedAt,
		task.CompletedAt,
		task.StatsAdds,
		task.StatsDels,
	)
	if err != nil {
		return fmt.Errorf("insert task: %w", err)
	}
	return nil
}

// UpdateTaskStatus records lifecycle transitions (e.g. RUNNING -> REVIEW or ABORTED).
// It saves the completion timestamp and final diff stats (+adds, -dels). Returns sql.ErrNoRows if task doesn't exist.
func (d *DB) UpdateTaskStatus(ctx context.Context, id string, status models.TaskStatus, completedAt *time.Time, statsAdds, statsDels int) error {
	query := `
	UPDATE tasks
	SET status = ?, completed_at = ?, stats_adds = ?, stats_dels = ?
	WHERE id = ?
	`
	res, err := d.conn.ExecContext(ctx, query, string(status), completedAt, statsAdds, statsDels, id)
	if err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetTask looks up a single task by ID. Used during reconnects and status checks.
func (d *DB) GetTask(ctx context.Context, id string) (*models.Task, error) {
	query := `
	SELECT id, repo_path, prompt, status, started_at, completed_at, stats_adds, stats_dels
	FROM tasks
	WHERE id = ?
	`
	row := d.conn.QueryRowContext(ctx, query, id)

	var task models.Task
	var statusStr string
	var completedAt sql.NullTime

	err := row.Scan(
		&task.ID,
		&task.RepoPath,
		&task.Prompt,
		&statusStr,
		&task.StartedAt,
		&completedAt,
		&task.StatsAdds,
		&task.StatsDels,
	)
	if err != nil {
		return nil, err
	}

	task.Status = models.TaskStatus(statusStr)
	if completedAt.Valid {
		task.CompletedAt = &completedAt.Time
	}

	return &task, nil
}

// GetRecentTasks fetches the last N tasks for the mobile timeline feed.
// If repoPath is provided, filters tasks to that specific project folder.
func (d *DB) GetRecentTasks(ctx context.Context, repoPath string, limit int) ([]models.Task, error) {
	if limit <= 0 {
		limit = 20
	}

	var query string
	var rows *sql.Rows
	var err error

	if repoPath != "" {
		query = `
		SELECT id, repo_path, prompt, status, started_at, completed_at, stats_adds, stats_dels
		FROM tasks
		WHERE repo_path = ?
		ORDER BY started_at DESC
		LIMIT ?
		`
		rows, err = d.conn.QueryContext(ctx, query, repoPath, limit)
	} else {
		query = `
		SELECT id, repo_path, prompt, status, started_at, completed_at, stats_adds, stats_dels
		FROM tasks
		ORDER BY started_at DESC
		LIMIT ?
		`
		rows, err = d.conn.QueryContext(ctx, query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("query recent tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		var statusStr string
		var completedAt sql.NullTime

		if err := rows.Scan(
			&task.ID,
			&task.RepoPath,
			&task.Prompt,
			&statusStr,
			&task.StartedAt,
			&completedAt,
			&task.StatsAdds,
			&task.StatsDels,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		task.Status = models.TaskStatus(statusStr)
		if completedAt.Valid {
			task.CompletedAt = &completedAt.Time
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

// SaveMessage appends a conversational turn (user prompt or agent reply) to the task history.
// An optional diff_json snapshot is stored to preserve what the code looked like at that moment.
func (d *DB) SaveMessage(ctx context.Context, msg *models.Message) error {
	query := `
	INSERT INTO messages (id, task_id, role, content, diff_json, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err := d.conn.ExecContext(ctx, query,
		msg.ID,
		msg.TaskID,
		string(msg.Role),
		msg.Content,
		msg.DiffJSON,
		msg.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}
	return nil
}

// GetTaskMessages returns all messages for a given task in ascending chronological order.
// Used by the mobile app to populate the multi-turn chat timeline.
func (d *DB) GetTaskMessages(ctx context.Context, taskID string) ([]models.Message, error) {
	query := `
	SELECT id, task_id, role, content, diff_json, created_at
	FROM messages
	WHERE task_id = ?
	ORDER BY created_at ASC
	`
	rows, err := d.conn.QueryContext(ctx, query, taskID)
	if err != nil {
		return nil, fmt.Errorf("query task messages: %w", err)
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		var roleStr string
		var diffJSON sql.NullString

		if err := rows.Scan(
			&msg.ID,
			&msg.TaskID,
			&roleStr,
			&msg.Content,
			&diffJSON,
			&msg.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan message: %w", err)
		}

		msg.Role = models.MessageRole(roleStr)
		if diffJSON.Valid {
			msg.DiffJSON = &diffJSON.String
		}

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
