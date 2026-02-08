package version

import (
	"encoding/json"
	"os"
)

// Build-injected variables (set via -ldflags -X at build time).
var (
	CommitHash string
	CommitTime string
	BuildTime  string
)

// Info holds VCS and build metadata for JSON output.
type Info struct {
	CommitHash string `json:"commit_hash"`
	CommitTime string `json:"commit_time"`
	BuildTime  string `json:"build_time"`
}

// Get returns the current version info, using "unknown" for any unset field.
func Get() Info {
	commitHash := CommitHash
	if commitHash == "" {
		commitHash = "unknown"
	}
	commitTime := CommitTime
	if commitTime == "" {
		commitTime = "unknown"
	}
	buildTime := BuildTime
	if buildTime == "" {
		buildTime = "unknown"
	}
	return Info{
		CommitHash: commitHash,
		CommitTime: commitTime,
		BuildTime:  buildTime,
	}
}

// WriteJSON marshals the version info as indented JSON to os.Stdout.
func WriteJSON() error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(Get())
}
