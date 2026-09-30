package db

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/priyxansh/rudder/daemon/internal/models"
)

func TestDB_StoreOperations(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := Open(dbPath)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	taskID := uuid.New().String()
	now := time.Now().Truncate(time.Second)

	task := &models.Task{
		ID:        taskID,
		RepoPath:  "/test/repo",
		Prompt:    "Add user authentication",
		Status:    models.TaskStatusRunning,
		StartedAt: now,
		StatsAdds: 0,
		StatsDels: 0,
	}

	if err := database.CreateTask(ctx, task); err != nil {
		t.Fatalf("CreateTask failed: %v", err)
	}

	fetched, err := database.GetTask(ctx, taskID)
	if err != nil {
		t.Fatalf("GetTask failed: %v", err)
	}
	if fetched.ID != taskID || fetched.Status != models.TaskStatusRunning {
		t.Errorf("GetTask mismatched fields: %+v", fetched)
	}

	completedAt := now.Add(2 * time.Minute)
	if err := database.UpdateTaskStatus(ctx, taskID, models.TaskStatusReview, &completedAt, 42, 5); err != nil {
		t.Fatalf("UpdateTaskStatus failed: %v", err)
	}

	updated, err := database.GetTask(ctx, taskID)
	if err != nil {
		t.Fatalf("GetTask after update failed: %v", err)
	}
	if updated.Status != models.TaskStatusReview {
		t.Errorf("expected status %s, got %s", models.TaskStatusReview, updated.Status)
	}
	if updated.StatsAdds != 42 || updated.StatsDels != 5 {
		t.Errorf("expected stats 42/5, got %d/%d", updated.StatsAdds, updated.StatsDels)
	}
	if updated.CompletedAt == nil || !updated.CompletedAt.Equal(completedAt) {
		t.Errorf("expected completedAt %v, got %v", completedAt, updated.CompletedAt)
	}

	diffSnapshot := `[{"path":"auth.go"}]`
	msg := &models.Message{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		Role:      models.MessageRoleAssistant,
		Content:   "Created auth logic",
		DiffJSON:  &diffSnapshot,
		CreatedAt: now.Add(1 * time.Minute),
	}

	if err := database.SaveMessage(ctx, msg); err != nil {
		t.Fatalf("SaveMessage failed: %v", err)
	}

	messages, err := database.GetTaskMessages(ctx, taskID)
	if err != nil {
		t.Fatalf("GetTaskMessages failed: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
	if messages[0].Content != "Created auth logic" || *messages[0].DiffJSON != diffSnapshot {
		t.Errorf("message content mismatch: %+v", messages[0])
	}

	recent, err := database.GetRecentTasks(ctx, "/test/repo", 10)
	if err != nil {
		t.Fatalf("GetRecentTasks failed: %v", err)
	}
	if len(recent) != 1 {
		t.Fatalf("expected 1 recent task, got %d", len(recent))
	}
}
