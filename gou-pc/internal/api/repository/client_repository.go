package repository

import (
	"database/sql"
	"gou-pc/internal/agent"
	"gou-pc/internal/logutil"
)

// ClientRepository interface cho thao tác client
//go:generate mockgen -source=client_repository.go -destination=mock_client_repository.go -package=repository

type ClientRepository interface {
	ClientGetAll() ([]agent.ManagedClient, error)
	ClientCreate(client *agent.ManagedClient) error
	ClientUpdate(client *agent.ManagedClient) error
	ClientDeleteByID(clientID string) error
	ClientFindByID(clientID string) (*agent.ManagedClient, error)
	ClientFindByAgentID(agentID string) (*agent.ManagedClient, error)
	ClientFindByUserID(userID string) ([]agent.ManagedClient, error)
	ClientGetClientIDByAgentID(agentID string) (string, error)
}

type sqliteClientRepository struct {
	db *sql.DB
}

func NewSQLiteClientRepository(db *sql.DB) ClientRepository {
	return &sqliteClientRepository{db: db}
}

func (r *sqliteClientRepository) ClientGetAll() ([]agent.ManagedClient, error) {
	rows, err := r.db.Query(`SELECT client_id, agent_id, hardware_id, host_name, ip_address, mac_address, user_name, last_seen, online FROM managed_clients`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clients []agent.ManagedClient
	for rows.Next() {
		var c agent.ManagedClient
		var hardwareID, hostName, ipAddress, macAddress string
		var onlineInt int
		err := rows.Scan(&c.ClientID, &c.AgentID, &hardwareID, &hostName, &ipAddress, &macAddress, &c.UserName, &c.LastSeen, &onlineInt)
		c.DeviceInfo.HardwareID = hardwareID
		c.DeviceInfo.HostName = hostName
		c.DeviceInfo.IPAddress = ipAddress
		c.DeviceInfo.MacAddress = macAddress
		c.Online = onlineInt == 1
		if err == nil {
			clients = append(clients, c)
		}
	}
	return clients, nil
}

func (r *sqliteClientRepository) ClientCreate(client *agent.ManagedClient) error {
	if client == nil {
		return sql.ErrNoRows
	}
	if client.ClientID == "" {
		return sql.ErrNoRows
	}
	if client.AgentID == "" {
		return sql.ErrNoRows
	}
	if client.DeviceInfo.HardwareID == "" {
		return sql.ErrNoRows
	}
	_, err := r.db.Exec(`INSERT INTO managed_clients (client_id, agent_id, hardware_id, host_name, ip_address, mac_address, user_name, last_seen, online) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		client.ClientID, client.AgentID, client.DeviceInfo.HardwareID, client.DeviceInfo.HostName, client.DeviceInfo.IPAddress, client.DeviceInfo.MacAddress, client.UserName, client.LastSeen, boolToInt(client.Online))
	if err != nil {
		logutil.APIDebug("ClientRepository.Create: failed to create client %s: %v", client.ClientID, err)
		return err
	}
	logutil.APIDebug("ClientRepository.Create: created client %s", client.ClientID)
	return nil
}

func (r *sqliteClientRepository) ClientUpdate(client *agent.ManagedClient) error {
	if client == nil {
		return sql.ErrNoRows
	}
	if client.ClientID == "" {
		return sql.ErrNoRows
	}
	// Lấy client hiện tại từ DB
	existing, err := r.ClientFindByID(client.ClientID)
	if err != nil {
		return err
	}
	if existing == nil {
		return sql.ErrNoRows
	}
	fields := []string{}
	args := []interface{}{}
	if client.AgentID != "" && client.AgentID != existing.AgentID {
		fields = append(fields, "agent_id=?")
		args = append(args, client.AgentID)
	}
	if client.DeviceInfo.HardwareID != "" && client.DeviceInfo.HardwareID != existing.DeviceInfo.HardwareID {
		fields = append(fields, "hardware_id=?")
		args = append(args, client.DeviceInfo.HardwareID)
	}
	if client.UserName != "" && client.UserName != existing.UserName {
		fields = append(fields, "user_name=?")
		args = append(args, client.UserName)
	}
	if client.LastSeen != "" && client.LastSeen != existing.LastSeen {
		fields = append(fields, "last_seen=?")
		args = append(args, client.LastSeen)
	}
	if client.Online != existing.Online {
		fields = append(fields, "online=?")
		args = append(args, boolToInt(client.Online))
	}
	if len(fields) == 0 {
		return nil // Không có gì để update
	}
	args = append(args, client.ClientID)
	query := "UPDATE managed_clients SET " + joinFields(fields) + " WHERE client_id=?"
	_, err = r.db.Exec(query, args...)
	if err != nil {
		logutil.APIDebug("ClientRepository.Update: failed to update client %s: %v", client.ClientID, err)
		return err
	}
	logutil.APIDebug("ClientRepository.Update: updated client %s", client.ClientID)
	return nil
}

func (r *sqliteClientRepository) ClientDeleteByID(clientID string) error {
	if clientID == "" {
		return sql.ErrNoRows
	}
	_, err := r.db.Exec(`DELETE FROM managed_clients WHERE client_id=?`, clientID)
	if err != nil {
		logutil.APIDebug("ClientRepository.DeleteByID: failed to delete client %s: %v", clientID, err)
		return err
	}
	logutil.APIDebug("ClientRepository.DeleteByID: deleted client %s", clientID)
	return nil
}

func (r *sqliteClientRepository) ClientFindByID(clientID string) (*agent.ManagedClient, error) {
	if clientID == "" {
		return nil, sql.ErrNoRows
	}
	row := r.db.QueryRow(`SELECT client_id, agent_id, hardware_id, host_name, ip_address, mac_address, user_name, last_seen, online FROM managed_clients WHERE client_id=?`, clientID)
	var c agent.ManagedClient
	var hardwareID, hostName, ipAddress, macAddress string
	var onlineInt int
	err := row.Scan(&c.ClientID, &c.AgentID, &hardwareID, &hostName, &ipAddress, &macAddress, &c.UserName, &c.LastSeen, &onlineInt)
	c.DeviceInfo.HardwareID = hardwareID
	c.DeviceInfo.HostName = hostName
	c.DeviceInfo.IPAddress = ipAddress
	c.DeviceInfo.MacAddress = macAddress
	c.Online = onlineInt == 1
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *sqliteClientRepository) ClientFindByAgentID(agentID string) (*agent.ManagedClient, error) {
	if agentID == "" {
		return nil, sql.ErrNoRows
	}
	row := r.db.QueryRow(`SELECT client_id, agent_id, hardware_id, host_name, ip_address, mac_address, user_name, last_seen, online FROM managed_clients WHERE agent_id=?`, agentID)
	var c agent.ManagedClient
	var hardwareID, hostName, ipAddress, macAddress string
	var onlineInt int
	err := row.Scan(&c.ClientID, &c.AgentID, &hardwareID, &hostName, &ipAddress, &macAddress, &c.UserName, &c.LastSeen, &onlineInt)
	c.DeviceInfo.HardwareID = hardwareID
	c.DeviceInfo.HostName = hostName
	c.DeviceInfo.IPAddress = ipAddress
	c.DeviceInfo.MacAddress = macAddress
	c.Online = onlineInt == 1
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *sqliteClientRepository) ClientFindByUserID(userID string) ([]agent.ManagedClient, error) {
	if userID == "" {
		return nil, sql.ErrNoRows
	}
	rows, err := r.db.Query(`SELECT client_id, agent_id, hardware_id, host_name, ip_address, mac_address, user_name, last_seen, online FROM managed_clients WHERE user_name=?`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var clients []agent.ManagedClient
	for rows.Next() {
		var c agent.ManagedClient
		var hardwareID, hostName, ipAddress, macAddress string
		var onlineInt int
		err := rows.Scan(&c.ClientID, &c.AgentID, &hardwareID, &hostName, &ipAddress, &macAddress, &c.UserName, &c.LastSeen, &onlineInt)
		c.DeviceInfo.HardwareID = hardwareID
		c.DeviceInfo.HostName = hostName
		c.DeviceInfo.IPAddress = ipAddress
		c.DeviceInfo.MacAddress = macAddress
		c.Online = onlineInt == 1
		if err == nil {
			clients = append(clients, c)
		}
	}
	return clients, nil
}

func (r *sqliteClientRepository) ClientGetClientIDByAgentID(agentID string) (string, error) {
	if agentID == "" {
		return "", sql.ErrNoRows
	}
	row := r.db.QueryRow(`SELECT client_id FROM managed_clients WHERE agent_id=?`, agentID)
	var clientID string
	err := row.Scan(&clientID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return clientID, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
