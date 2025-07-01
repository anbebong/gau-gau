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

// GetLogsPaged phân trang log từ file archive.log, đảo ngược thứ tự (mới nhất lên đầu)
func GetLogsPaged(archiveFile string, page, pageSize int) ([]ArchiveLogEntry, int, error) {
	logs, err := LoadArchiveLogs(archiveFile)
	if err != nil {
		return nil, 0, err
	}
	// Đảo ngược slice logs
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}
	total := len(logs)
	start := (page - 1) * pageSize
	if start > total {
		return []ArchiveLogEntry{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return logs[start:end], total, nil
}

// RotateLog đổi tên file log hiện tại sang <log>.old, tạo file mới
func RotateLog(archiveFile string) error {
	oldPath := archiveFile + ".old"
	_ = os.Remove(oldPath)
	if err := os.Rename(archiveFile, oldPath); err != nil {
		return err
	}
	f, err := os.Create(archiveFile)
	if err != nil {
		return err
	}
	return f.Close()
}
