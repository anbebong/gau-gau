package service

import (
	"gou-pc/internal/logcollector"
)

// LogService interface
//go:generate mockgen -source=log_service.go -destination=mock_log_service.go -package=service

type LogService interface {
	GetAllLogs() ([]logcollector.ArchiveLogEntry, error)
	GetLogsByAgentID(agentID string) ([]logcollector.ArchiveLogEntry, error)
}

type logServiceImpl struct {
	archiveFile string
}

func NewLogService(archiveFile string) LogService {
	return &logServiceImpl{archiveFile: archiveFile}
}

func (s *logServiceImpl) GetAllLogs() ([]logcollector.ArchiveLogEntry, error) {
	return logcollector.LoadArchiveLogs(s.archiveFile)
}

func (s *logServiceImpl) GetLogsByAgentID(agentID string) ([]logcollector.ArchiveLogEntry, error) {
	logs, err := logcollector.LoadArchiveLogs(s.archiveFile)
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
