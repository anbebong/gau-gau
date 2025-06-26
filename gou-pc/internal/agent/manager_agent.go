package agent

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

type ManagedClient struct {
	ClientID   string     `json:"client_id"`
	AgentID    string     `json:"agent_id"`
	DeviceInfo DeviceInfo `json:"device_info"`
	UserID     string     `json:"user_id"` // user sở hữu client này
}

var (
	mu          sync.Mutex
	nextAgentID = 1
)

func LoadClients(managerFile string) ([]ManagedClient, error) {
	mu.Lock()
	defer mu.Unlock()
	f, err := os.ReadFile(managerFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("[WARN] Client DB file '%s' not found, returning empty list\n", managerFile)
			return []ManagedClient{}, nil
		}
		return nil, err
	}
	var clients []ManagedClient
	_ = json.Unmarshal(f, &clients)
	// Tìm agent_id lớn nhất để cập nhật nextAgentID
	next := 1
	for _, c := range clients {
		if len(c.AgentID) == 3 {
			if n, err := strconv.Atoi(c.AgentID); err == nil && n >= next {
				next = n + 1
			}
		}
	}
	nextAgentID = next
	return clients, nil
}

func SaveClients(managerFile string, clients []ManagedClient) error {
	mu.Lock()
	defer mu.Unlock()
	b, _ := json.MarshalIndent(clients, "", "  ")
	return os.WriteFile(managerFile, b, 0644)
}

func FindClientByDevice(info DeviceInfo, clients []ManagedClient) *ManagedClient {
	for _, c := range clients {
		if c.DeviceInfo.HardwareID == info.HardwareID {
			return &c
		}
	}
	return nil
}

func GenAgentID() string {
	mu.Lock()
	id := fmt.Sprintf("%03d", nextAgentID)
	nextAgentID++
	mu.Unlock()
	return id
}

func GenClientID() string {
	return uuid.NewString()
}
