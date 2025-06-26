package logcollector

import (
	"encoding/json"
	"io"
	"os"
)

type ArchiveLogEntry struct {
	Time    string `json:"time"`
	AgentID string `json:"agent_id"`
	Message string `json:"message"`
}

// LoadArchiveLogs đọc toàn bộ log từ file archive.log
func LoadArchiveLogs(archiveFile string) ([]ArchiveLogEntry, error) {
	f, err := os.Open(archiveFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var logs []ArchiveLogEntry
	dec := json.NewDecoder(f)
	for {
		var entry ArchiveLogEntry
		if err := dec.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		logs = append(logs, entry)
	}
	return logs, nil
}
