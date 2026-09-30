package models

// WebSocket event type discriminators.
const (
	EventTaskStarted   = "task_started"
	EventLogChunk      = "log_chunk"
	EventTaskCompleted = "task_completed"
	EventTaskAborted   = "task_aborted"
	EventTaskFailed    = "task_failed"
	EventCommitSuccess = "commit_success"
	EventSubscribe     = "subscribe"
)

// WSEvent wraps all bidirectional messages sent over the WebSocket bridge.
type WSEvent struct {
	Event     string `json:"event"`
	Seq       uint64 `json:"seq,omitempty"`
	Timestamp int64  `json:"timestamp"`
	Payload   any    `json:"payload,omitempty"`
}

// TaskStartedPayload notifies the mobile client that an Aider run has begun.
type TaskStartedPayload struct {
	TaskID string `json:"taskId"`
	Prompt string `json:"prompt"`
}

// LogChunkPayload carries a streaming line of stdout/stderr from Aider.
type LogChunkPayload struct {
	Seq  uint64 `json:"seq"`
	Text string `json:"text"`
}

// TaskCompletedPayload delivers finalized diff cards and stats when Aider exits 0.
type TaskCompletedPayload struct {
	TaskID    string     `json:"taskId"`
	DiffCards []DiffCard `json:"diffCards"`
	Stats     DiffStats  `json:"stats"`
}

// TaskAbortedPayload confirms an active task was stopped via user request.
type TaskAbortedPayload struct {
	TaskID string `json:"taskId"`
}

// TaskFailedPayload alerts the client that an agent run crashed or exited non-zero.
type TaskFailedPayload struct {
	TaskID   string `json:"taskId"`
	ExitCode int    `json:"exitCode"`
	Error    string `json:"error"`
}

// CommitSuccessPayload confirms changes were committed to git and optionally pushed.
type CommitSuccessPayload struct {
	CommitSHA string `json:"commitSha"`
	Branch    string `json:"branch"`
}

// SubscribePayload is sent by the phone upon connecting to replay buffered logs since lastSeq.
type SubscribePayload struct {
	LastSeq uint64 `json:"lastSeq"`
}
