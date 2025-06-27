package agent

import (
	"database/sql"
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
	UserName   string     `json:"user_name"`
	LastSeen   string     `json:"last_seen"` // ISO8601 string
	Online     bool       `json:"online"`
}

var (
	mu          sync.Mutex
	nextAgentID = 1
	db          *sql.DB // Đảm bảo đã khởi tạo ở nơi khác
)

func LoadClients() ([]ManagedClient, error) {
	rows, err := db.Query("SELECT client_id, agent_id, hardware_id, user_id, last_seen, online FROM managed_clients")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clients []ManagedClient
	for rows.Next() {
		var c ManagedClient
		var hardwareID string
		var onlineInt int
		err := rows.Scan(&c.ClientID, &c.AgentID, &hardwareID, &c.UserName, &c.LastSeen, &onlineInt)
		c.DeviceInfo.HardwareID = hardwareID
		c.Online = onlineInt == 1
		if err == nil {
			clients = append(clients, c)
		}
	}
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

func FindClientByDevice(hardwareID string) (*ManagedClient, error) {
	row := db.QueryRow("SELECT client_id, agent_id, user_id, last_seen, online FROM managed_clients WHERE hardware_id = ?", hardwareID)
	var c ManagedClient
	var onlineInt int
	c.DeviceInfo.HardwareID = hardwareID
	err := row.Scan(&c.ClientID, &c.AgentID, &c.UserName, &c.LastSeen, &onlineInt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	c.Online = onlineInt == 1
	return &c, nil
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

// SaveClient lưu một client mới vào DB
func SaveClient(c ManagedClient) error {
	_, err := db.Exec(`INSERT OR REPLACE INTO managed_clients
		(client_id, agent_id, hardware_id, user_name, last_seen, online)
		VALUES (?, ?, ?, ?, ?, ?)`,
		c.ClientID, c.AgentID, c.DeviceInfo.HardwareID, c.UserName, c.LastSeen, boolToInt(c.Online))
	return err
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// UpdateClientStatus cập nhật trạng thái online và last_seen cho agent
func UpdateClientStatus(agentID string, online bool, lastSeen string) error {
	_, err := db.Exec(`UPDATE managed_clients SET online=?, last_seen=? WHERE agent_id=?`, boolToInt(online), lastSeen, agentID)
	return err
}

// SetDB cho phép gán biến db toàn cục từ bên ngoài
func SetDB(newDB *sql.DB) {
	db = newDB
}
