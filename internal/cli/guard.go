package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var readOnly bool

type auditRecord struct {
	Timestamp string `json:"timestamp"`
	Command   string `json:"command"`
	Action    string `json:"action"`
	Target    string `json:"target"`
	Result    string `json:"result"`
}

func isReadOnly() bool {
	return readOnly || os.Getenv("S1_READONLY") == "1"
}

func guard(w io.Writer, command, action, target string, yes bool, fn func() error) error {
	if isReadOnly() {
		fmt.Fprintf(w, "Read-only mode: would %s\n", action)
		return nil
	}
	if !yes {
		fmt.Fprintf(w, "Would %s. Pass --yes to apply.\n", action)
		return nil
	}
	err := fn()
	result := "ok"
	if err != nil {
		result = err.Error()
	}
	_ = writeAuditRecord(auditRecord{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Command:   command,
		Action:    action,
		Target:    target,
		Result:    result,
	})
	return err
}

func auditLogPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".s1ctl", "audit.jsonl")
}

func writeAuditRecord(rec auditRecord) error {
	path := auditLogPath()
	if path == "" {
		fmt.Fprintln(os.Stderr, "Warning: could not determine home directory, audit log skipped")
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	return json.NewEncoder(f).Encode(rec)
}
