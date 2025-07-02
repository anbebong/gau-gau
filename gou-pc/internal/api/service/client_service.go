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
	AssignUserToClientByAgentID(agentID, username string) error
	AssignUserToClientByClientID(clientID, username string) error
	GetAllClients() ([]ManagedClientResponse, error)
	GetClientByAgentID(agentID string) (*ManagedClientResponse, error)
	GetClientByID(clientID string) (*ManagedClientResponse, error)
	DeleteClientByAgentID(agentID string) error
	DeleteClientByClientID(clientID string) error
	GetClientsByUsername(username string) ([]ManagedClientResponse, error)
}

// Response struct trả về client cho API, có Username thay vì user_id
type ManagedClientResponse struct {
	ClientID   string           `json:"client_id"`
	AgentID    string           `json:"agent_id"`
	DeviceInfo agent.DeviceInfo `json:"device_info"`
	Username   string           `json:"username"`
	LastSeen   string           `json:"last_seen"` // ISO8601 string
	Online     bool             `json:"online"`
}

type clientServiceImpl struct {
	repo     repository.ClientRepository
	userRepo repository.UserRepository
}

func NewClientService(repo repository.ClientRepository, userRepo repository.UserRepository) ClientService {
	return &clientServiceImpl{repo: repo, userRepo: userRepo}
}

// Gán user cho client bằng username (mapping sang user_id ở service)
func (s *clientServiceImpl) AssignUserToClientByAgentID(agentID, username string) error {
	if agentID == "" {
		return errors.New("agentID is required")
	}
	if username == "" {
		return errors.New("username is required")
	}
	logutil.APIDebug("ClientService.AssignUserToClientByAgentID called with agentID=%s, username=%s", agentID, username)
	client, err := s.repo.ClientFindByAgentID(agentID)
	if err != nil || client == nil {
		return errors.New("client not found")
	}
	client.UserName = username
	return s.repo.ClientUpdate(client)
}

func (s *clientServiceImpl) AssignUserToClientByClientID(clientID, username string) error {
	if clientID == "" {
		return errors.New("clientID is required")
	}
	if username == "" {
		return errors.New("username is required")
	}
	logutil.APIDebug("ClientService.AssignUserToClientByClientID called with clientID=%s, username=%s", clientID, username)
	client, err := s.repo.ClientFindByID(clientID)
	if err != nil || client == nil {
		return errors.New("client not found")
	}
	client.UserName = username
	return s.repo.ClientUpdate(client)
}

func (s *clientServiceImpl) GetAllClients() ([]ManagedClientResponse, error) {
	logutil.APIDebug("ClientService.GetAllClients called")
	clients, err := s.repo.ClientGetAll()
	if err != nil {
		return nil, err
	}
	var resp []ManagedClientResponse
	for _, c := range clients {
		resp = append(resp, ManagedClientResponse{
			ClientID:   c.ClientID,
			AgentID:    c.AgentID,
			DeviceInfo: c.DeviceInfo,
			Username:   c.UserName,
			LastSeen:   c.LastSeen,
			Online:     c.Online,
		})
	}
	return resp, nil
}

func (s *clientServiceImpl) GetClientByAgentID(agentID string) (*ManagedClientResponse, error) {
	if agentID == "" {
		return nil, errors.New("agentID is required")
	}
	logutil.APIDebug("ClientService.GetClientByAgentID called with agentID=%s", agentID)
	c, err := s.repo.ClientFindByAgentID(agentID)
	if err != nil || c == nil {
		return nil, err
	}
	r := ManagedClientResponse{
		ClientID:   c.ClientID,
		AgentID:    c.AgentID,
		DeviceInfo: c.DeviceInfo,
		Username:   c.UserName,
		LastSeen:   c.LastSeen,
		Online:     c.Online,
	}
	return &r, nil
}

func (s *clientServiceImpl) GetClientByID(clientID string) (*ManagedClientResponse, error) {
	if clientID == "" {
		return nil, errors.New("clientID is required")
	}
	logutil.APIDebug("ClientService.GetClientByID called with clientID=%s", clientID)
	c, err := s.repo.ClientFindByID(clientID)
	if err != nil || c == nil {
		return nil, err
	}
	r := ManagedClientResponse{
		ClientID:   c.ClientID,
		AgentID:    c.AgentID,
		DeviceInfo: c.DeviceInfo,
		Username:   c.UserName,
		LastSeen:   c.LastSeen,
		Online:     c.Online,
	}
	return &r, nil
}

func (s *clientServiceImpl) DeleteClientByAgentID(agentID string) error {
	if agentID == "" {
		return errors.New("agentID is required")
	}
	client, err := s.repo.ClientFindByAgentID(agentID)
	if err != nil || client == nil {
		return errors.New("client not found")
	}
	return s.repo.ClientDeleteByID(client.ClientID)
}

func (s *clientServiceImpl) DeleteClientByClientID(clientID string) error {
	if clientID == "" {
		return errors.New("clientID is required")
	}
	client, err := s.repo.ClientFindByID(clientID)
	if err != nil || client == nil {
		return errors.New("client not found")
	}
	return s.repo.ClientDeleteByID(clientID)
}

func (s *clientServiceImpl) GetClientsByUsername(username string) ([]ManagedClientResponse, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	logutil.APIDebug("ClientService.GetClientsByUsername called with username=%s", username)
	clients, err := s.repo.ClientGetAll()
	if err != nil {
		return nil, err
	}
	var resp []ManagedClientResponse
	for _, c := range clients {
		if c.UserName == username {
			resp = append(resp, ManagedClientResponse{
				ClientID:   c.ClientID,
				AgentID:    c.AgentID,
				DeviceInfo: c.DeviceInfo,
				Username:   c.UserName,
				LastSeen:   c.LastSeen,
				Online:     c.Online,
			})
		}
	}
	return resp, nil
}
