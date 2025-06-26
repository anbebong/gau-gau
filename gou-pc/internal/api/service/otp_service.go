package service

import (
	"errors"
	"gou-pc/internal/api/repository"
	"gou-pc/internal/crypto"
)

type OTPService interface {
	GetOTPByAgentID(agentID string) (string, error)
	GetOTPByClientID(clientID string) (string, error)
	GetOTPByAgentIDWithExpire(agentID string) (string, int, error)
	GetOTPByClientIDWithExpire(clientID string) (string, int, error)
}

type otpServiceImpl struct {
	repo repository.ClientRepository
	// mu   sync.RWMutex
}

func NewOTPService(repo repository.ClientRepository) OTPService {
	return &otpServiceImpl{repo: repo}
}

func (s *otpServiceImpl) GetOTPByAgentID(agentID string) (string, error) {
	clientID, err := s.repo.GetClientIDByAgentID(agentID)
	if err != nil || clientID == "" {
		return "", errors.New("agent not found")
	}
	return crypto.GetTOTPByClientID(clientID)
}

func (s *otpServiceImpl) GetOTPByClientID(clientID string) (string, error) {
	return crypto.GetTOTPByClientID(clientID)
}

func (s *otpServiceImpl) GetOTPByAgentIDWithExpire(agentID string) (string, int, error) {
	clientID, err := s.repo.GetClientIDByAgentID(agentID)
	if err != nil || clientID == "" {
		return "", 0, errors.New("agent not found")
	}
	return crypto.GetTOTPWithExpireByClientID(clientID)
}

func (s *otpServiceImpl) GetOTPByClientIDWithExpire(clientID string) (string, int, error) {
	return crypto.GetTOTPWithExpireByClientID(clientID)
}
