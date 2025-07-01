package repository

import (
	"gou-pc/internal/logcollector"
)

// LogRepository interface cho thao tác log
//go:generate mockgen -source=log_repository.go -destination=mock_log_repository.go -package=repository

type LogRepository interface {
	GetAllLogs() ([]logcollector.ArchiveLogEntry, error)
	GetLogsByAgentID(agentID string) ([]logcollector.ArchiveLogEntry, error)
	GetLogsPaged(page, pageSize int) ([]logcollector.ArchiveLogEntry, int, error)
	GetLogsPagedByAgentID(agentID string, page, pageSize int) ([]logcollector.ArchiveLogEntry, int, error)
	RotateLog() error
}

type fileLogRepository struct {
	archiveFile string
}

func NewFileLogRepository(archiveFile string) LogRepository {
	return &fileLogRepository{archiveFile: archiveFile}
}

func (r *fileLogRepository) GetAllLogs() ([]logcollector.ArchiveLogEntry, error) {
	logs, err := logcollector.LoadArchiveLogs(r.archiveFile)
	if err != nil {
		return nil, err
	}
	// Đảo ngược slice logs để log mới nhất lên đầu
	for i, j := 0, len(logs)-1; i < j; i, j = i+1, j-1 {
		logs[i], logs[j] = logs[j], logs[i]
	}
	return logs, nil
}

func (r *fileLogRepository) GetLogsByAgentID(agentID string) ([]logcollector.ArchiveLogEntry, error) {
	logs, err := logcollector.LoadArchiveLogs(r.archiveFile)
	if err != nil {
		return nil, err
	}
	var result []logcollector.ArchiveLogEntry
	for _, l := range logs {
		if l.AgentID == agentID {
			result = append(result, l)
		}
	}
	// Đảo ngược slice result để log mới nhất lên đầu
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}
	return result, nil
}

// Phân trang log: trả về logs theo trang (page, pageSize)
func (r *fileLogRepository) GetLogsPaged(page, pageSize int) ([]logcollector.ArchiveLogEntry, int, error) {
	return logcollector.GetLogsPaged(r.archiveFile, page, pageSize)
}

func (r *fileLogRepository) GetLogsPagedByAgentID(agentID string, page, pageSize int) ([]logcollector.ArchiveLogEntry, int, error) {
	logs, err := logcollector.LoadArchiveLogs(r.archiveFile)
	if err != nil {
		return nil, 0, err
	}
	var filtered []logcollector.ArchiveLogEntry
	for _, l := range logs {
		if l.AgentID == agentID {
			filtered = append(filtered, l)
		}
	}
	total := len(filtered)
	start := (page - 1) * pageSize
	if start > total {
		return []logcollector.ArchiveLogEntry{}, total, nil
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}

// Rotate log: đổi tên file log hiện tại sang <log>.old, tạo file mới
func (r *fileLogRepository) RotateLog() error {
	return logcollector.RotateLog(r.archiveFile)
}
