package repository

import "gou-pc/internal/logcollector"

// LogRepository interface cho thao tác log
//go:generate mockgen -source=log_repository.go -destination=mock_log_repository.go -package=repository

type LogRepository interface {
	GetAllLogs() ([]logcollector.ArchiveLogEntry, error)
	GetLogsByAgentID(agentID string) ([]logcollector.ArchiveLogEntry, error)
}

type fileLogRepository struct {
	archiveFile string
}

func NewFileLogRepository(archiveFile string) LogRepository {
	return &fileLogRepository{archiveFile: archiveFile}
}

func (r *fileLogRepository) GetAllLogs() ([]logcollector.ArchiveLogEntry, error) {
	return logcollector.LoadArchiveLogs(r.archiveFile)
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
	return result, nil
}
