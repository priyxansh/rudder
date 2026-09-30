package models

import "time"

// TaskStatus defines the lifecycle state of an agent execution run.
type TaskStatus string

const (
	TaskStatusIdle     TaskStatus = "IDLE"
	TaskStatusRunning  TaskStatus = "RUNNING"
	TaskStatusReview   TaskStatus = "REVIEW"
	TaskStatusResolved TaskStatus = "RESOLVED"
	TaskStatusAborted  TaskStatus = "ABORTED"
	TaskStatusFailed   TaskStatus = "FAILED"
)

// Task represents an agent execution session and its uncommitted diff stats.
type Task struct {
	ID          string     `json:"id"`
	RepoPath    string     `json:"repoPath"`
	Prompt      string     `json:"prompt"`
	Status      TaskStatus `json:"status"`
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt,omitempty"`
	StatsAdds   int        `json:"statsAdds"`
	StatsDels   int        `json:"statsDels"`
}

// MessageRole distinguishes the author of a conversation transcript item.
type MessageRole string

const (
	MessageRoleUser      MessageRole = "user"
	MessageRoleAssistant MessageRole = "assistant"
	MessageRoleSystem    MessageRole = "system"
)

// Message is a conversational turn or transcript entry linked to a task.
type Message struct {
	ID        string      `json:"id"`
	TaskID    string      `json:"taskId"`
	Role      MessageRole `json:"role"`
	Content   string      `json:"content"`
	DiffJSON  *string     `json:"diffJson,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
}

// CreateTaskRequest is sent from mobile to trigger an Aider run with context files.
type CreateTaskRequest struct {
	Prompt string   `json:"prompt"`
	Files  []string `json:"files,omitempty"`
}

// CritiqueRequest sends targeted feedback anchored to a specific file line.
type CritiqueRequest struct {
	TaskID        string `json:"taskId,omitempty"`
	FilePath      string `json:"filePath"`
	Line          int    `json:"line"`
	AnchorSnippet string `json:"anchorSnippet"`
	Critique      string `json:"critique"`
}

// AbortTaskRequest cancels an in-progress agent run.
type AbortTaskRequest struct {
	TaskID string `json:"taskId,omitempty"`
}
