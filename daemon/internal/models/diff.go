package models

// FileActionType describes the git working tree status of a file.
type FileActionType string

const (
	ActionModified FileActionType = "modified"
	ActionAdded    FileActionType = "added"
	ActionDeleted  FileActionType = "deleted"
	ActionRenamed  FileActionType = "renamed"
	ActionBinary   FileActionType = "binary"
)

// DiffLineType categorizes a line in a diff hunk.
type DiffLineType string

const (
	LineContext DiffLineType = "context"
	LineAdd     DiffLineType = "add"
	LineDel     DiffLineType = "del"
)

// DiffLine represents a single line in a hunk with structured line numbers for touch anchors.
type DiffLine struct {
	Type      DiffLineType `json:"type"`
	OldLineNo int          `json:"oldLineNo,omitempty"`
	NewLineNo int          `json:"newLineNo,omitempty"`
	Content   string       `json:"content"`
}

// DiffHunk groups a contiguous block of changed and context lines.
type DiffHunk struct {
	Header   string     `json:"header"`
	OldStart int        `json:"oldStart"`
	OldLines int        `json:"oldLines"`
	NewStart int        `json:"newStart"`
	NewLines int        `json:"newLines"`
	Lines    []DiffLine `json:"lines"`
}

// DiffCard represents a changed file rendered as a card in the mobile review deck.
type DiffCard struct {
	Path      string         `json:"path"`
	OldPath   string         `json:"oldPath,omitempty"`
	Status    FileActionType `json:"status"`
	IsBinary  bool           `json:"isBinary"`
	Additions int            `json:"additions"`
	Deletions int            `json:"deletions"`
	Hunks     []DiffHunk     `json:"hunks"`
}

// DiffStats provides change totals for the mobile cockpit header pill.
type DiffStats struct {
	FilesChanged int `json:"filesChanged"`
	Additions    int `json:"additions"`
	Deletions    int `json:"deletions"`
}

// WorkspaceStatus captures the repository git state and active task.
type WorkspaceStatus struct {
	Branch         string   `json:"branch"`
	IsClean        bool     `json:"isClean"`
	ActiveTask     *Task    `json:"activeTask,omitempty"`
	ModifiedFiles  []string `json:"modifiedFiles"`
	StagedFiles    []string `json:"stagedFiles"`
	UntrackedFiles []string `json:"untrackedFiles"`
	HeadSeq        uint64   `json:"headSeq"`
}

// FileEntry represents a file or directory node in the repository explorer.
type FileEntry struct {
	Path     string `json:"path"`
	IsDir    bool   `json:"isDir"`
	IsBinary bool   `json:"isBinary"`
	Size     int64  `json:"size"`
}

// CommitRequest instructs the daemon to commit staged files and optionally push.
type CommitRequest struct {
	Message string `json:"message"`
	Push    bool   `json:"push"`
}

// StageRequest instructs the daemon to stage a file via git add.
type StageRequest struct {
	Path string `json:"path"`
}

// RevertRequest instructs the daemon to discard working tree edits on a file.
type RevertRequest struct {
	Path string `json:"path"`
}
