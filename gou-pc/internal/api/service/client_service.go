package service

import (
	"errors"
	"gou-pc/internal/agent"
	"gou-pc/internal/api/repository"
	"gou-pc/internal/logutil"
)

// ClientService interface định nghĩa các hàm thao tác với client, chỉ dùng user_id
//go:generate mockgen -source=client_service.go -destination=mock_client_service.go -package=service

type ClientService interface {
	AssignUserToClient(agentID, username string) error // nhận username, mapping ở service
	GetAllClients() ([]ManagedClientResponse, error)
	GetClientByAgentID(agentID string) (*ManagedClientResponse, error)
	GetClientsByUserID(userID string) ([]ManagedClientResponse, error)
	GetAgentIDsByUserID(userID string) ([]string, error)
	DeleteClient(agentID string) error // thêm hàm xóa client
}

// Response struct trả về client cho API, có Username thay vì user_id
type ManagedClientResponse struct {
	ClientID   string           `json:"client_id"`
	AgentID    string           `json:"agent_id"`
	DeviceInfo agent.DeviceInfo `json:"device_info"`
	Username   string           `json:"username"`
}

type clientServiceImpl struct {
	repo     repository.ClientRepository
	userRepo repository.UserRepository
}

func NewClientService(repo repository.ClientRepository, userRepo repository.UserRepository) ClientService {
	return &clientServiceImpl{repo: repo, userRepo: userRepo}
}

// Gán user cho client bằng username (mapping sang user_id ở service)
func (s *clientServiceImpl) AssignUserToClient(agentID, username string) error {
	logutil.Debug("ClientService.AssignUserToClient called with agentID=%s, username=%s", agentID, username)
	user, err := s.userRepo.FindByUsername(username)
	if err != nil || user == nil {
		logutil.Debug("ClientService.AssignUserToClient: username not found: %s", username)
		return errors.New("username not found")
	}
	clients, err := s.repo.GetAll()
	if err != nil {
		logutil.Debug("ClientService.AssignUserToClient repo.GetAll error: %v", err)
		return err
	}
	found := false
	for i, c := range clients {
		if c.AgentID == agentID {
			clients[i].UserID = user.ID
			found = true
			break
		}
	}
	if !found {
		logutil.Debug("ClientService.AssignUserToClient: client not found agentID=%s", agentID)
		return errors.New("client not found")
	}
	err = s.repo.SaveAll(clients)
	if err != nil {
		logutil.Debug("ClientService.AssignUserToClient SaveAll error: %v", err)
	}
	return err
}

// Helper: mapping user_id sang username
func (s *clientServiceImpl) mapClientToResponse(c agent.ManagedClient) ManagedClientResponse {
	username := ""
	if c.UserID != "" {
		if user, err := s.userRepo.FindByID(c.UserID); err == nil && user != nil {
			username = user.Username
		}
	}
	return ManagedClientResponse{
		ClientID:   c.ClientID,
		AgentID:    c.AgentID,
		DeviceInfo: c.DeviceInfo,
		Username:   username,
	}
}

func (s *clientServiceImpl) GetAllClients() ([]ManagedClientResponse, error) {
	logutil.Debug("ClientService.GetAllClients called")
	clients, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	var resp []ManagedClientResponse
	for _, c := range clients {
		resp = append(resp, s.mapClientToResponse(c))
	}
	return resp, nil
}

func (s *clientServiceImpl) GetClientByAgentID(agentID string) (*ManagedClientResponse, error) {
	logutil.Debug("ClientService.GetClientByAgentID called with agentID=%s", agentID)
	c, err := s.repo.FindByAgentID(agentID)
	if err != nil || c == nil {
		return nil, err
	}
	r := s.mapClientToResponse(*c)
	return &r, nil
}

func (s *clientServiceImpl) GetClientsByUserID(userID string) ([]ManagedClientResponse, error) {
	logutil.Debug("ClientService.GetClientsByUserID called with userID=%s", userID)
	clients, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	var resp []ManagedClientResponse
	for _, c := range clients {
		resp = append(resp, s.mapClientToResponse(c))
	}
	return resp, nil
}

func (s *clientServiceImpl) GetAgentIDsByUserID(userID string) ([]string, error) {
	logutil.Debug("ClientService.GetAgentIDsByUserID called with userID=%s", userID)
	clients, err := s.repo.FindByUserID(userID)
	if err != nil {
		logutil.Debug("ClientService.GetAgentIDsByUserID repo.FindByUserID error: %v", err)
		return nil, err
	}
	var agentIDs []string
	for _, c := range clients {
		if c.AgentID != "" {
			agentIDs = append(agentIDs, c.AgentID)
		}
	}
	return agentIDs, nil
}

func (s *clientServiceImpl) DeleteClient(agentID string) error {
	logutil.Debug("ClientService.DeleteClient called with agentID=%s", agentID)
	clients, err := s.repo.GetAll()
	if err != nil {
		return err
	}
	newClients := make([]agent.ManagedClient, 0, len(clients))
	deleted := false
	for _, c := range clients {
		if c.AgentID == agentID {
			deleted = true
			continue
		}
		newClients = append(newClients, c)
	}
	if !deleted {
		logutil.Debug("ClientService.DeleteClient: client not found agentID=%s", agentID)
		return errors.New("client not found")
	}
	return s.repo.SaveAll(newClients)
}
