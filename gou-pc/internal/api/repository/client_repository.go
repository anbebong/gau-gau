package repository

import (
	"gou-pc/internal/agent"
	"gou-pc/internal/logutil"
)

// ClientRepository interface cho thao tác client
//go:generate mockgen -source=client_repository.go -destination=mock_client_repository.go -package=repository

type ClientRepository interface {
	GetAll() ([]agent.ManagedClient, error)
	SaveAll([]agent.ManagedClient) error
	FindByID(clientID string) (*agent.ManagedClient, error)
	FindByAgentID(agentID string) (*agent.ManagedClient, error)
	FindByUserID(userID string) ([]agent.ManagedClient, error)
	GetClientIDByAgentID(agentID string) (string, error)
}

type fileClientRepository struct {
	managerFile string
}

func NewFileClientRepository(managerFile string) ClientRepository {
	return &fileClientRepository{managerFile: managerFile}
}

func (r *fileClientRepository) GetAll() ([]agent.ManagedClient, error) {
	return agent.LoadClients(r.managerFile)
}

func (r *fileClientRepository) SaveAll(clients []agent.ManagedClient) error {
	return agent.SaveClients(r.managerFile, clients)
}

func (r *fileClientRepository) FindByID(clientID string) (*agent.ManagedClient, error) {
	clients, err := agent.LoadClients(r.managerFile)
	if err != nil {
		return nil, err
	}
	for _, c := range clients {
		if c.ClientID == clientID {
			return &c, nil
		}
	}
	return nil, nil
}

func (r *fileClientRepository) FindByAgentID(agentID string) (*agent.ManagedClient, error) {
	clients, err := agent.LoadClients(r.managerFile)
	if err != nil {
		return nil, err
	}
	for _, c := range clients {
		if c.AgentID == agentID {
			return &c, nil
		}
	}
	return nil, nil
}

func (r *fileClientRepository) FindByUserID(userID string) ([]agent.ManagedClient, error) {
	clients, err := agent.LoadClients(r.managerFile)
	if err != nil {
		logutil.Debug("FindByUserID: error loading clients: %v", err)
		return nil, err
	}
	logutil.Debug("FindByUserID: loaded %d clients, looking for userID=%s", len(clients), userID)
	var result []agent.ManagedClient
	for _, c := range clients {
		logutil.Debug("Checking client: agentID=%s, userID=%s", c.AgentID, c.UserID)
		if c.UserID == userID {
			logutil.Debug("Match: agentID=%s", c.AgentID)
			result = append(result, c)
		}
	}
	logutil.Debug("FindByUserID: found %d clients for userID=%s", len(result), userID)
	return result, nil
}

func (r *fileClientRepository) GetClientIDByAgentID(agentID string) (string, error) {
	clients, err := agent.LoadClients(r.managerFile)
	if err != nil {
		return "", err
	}
	for _, c := range clients {
		if c.AgentID == agentID {
			return c.ClientID, nil
		}
	}
	return "", nil
}
