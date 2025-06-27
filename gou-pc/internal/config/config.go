package config

import "time"

// ClientConfig holds all configurable paths and options for the client
type ClientConfig struct {
	LogFile    string        // Đường dẫn file log client
	EventLog   string        // File log sự kiện cần theo dõi
	OffsetFile string        // File lưu offset log
	ConfigFile string        // File lưu client_id, agent_id
	ServerAddr string        // Địa chỉ server
	Interval   time.Duration // Chu kỳ kiểm tra log
}

func DefaultClientConfig() *ClientConfig {
	return &ClientConfig{
		LogFile:    "etc/client.log",
		EventLog:   "etc/event.log",
		OffsetFile: "etc/event.log.offset",
		ConfigFile: "etc/client.json",
		ServerAddr: "192.168.15.12:9000",
		Interval:   2 * time.Second,
	}
}

// ServerConfig holds all configurable paths and options for the server
type ServerConfig struct {
	LogFile      string        // Đường dẫn file log server
	ArchiveFile  string        // File lưu log thu thập từ agent
	ClientDBFile string        // File lưu thông tin client/agent
	UserDBFile   string        // File lưu thông tin user
	ListenAddr   string        // Địa chỉ lắng nghe TCP
	APIPort      string        // Cổng chạy API server
	JWTSecret    string        // Secret key cho JWT
	JWTExpire    time.Duration // Thời gian sống của JWT
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		LogFile:      "etc/server.log",
		ArchiveFile:  "etc/archive.log",
		ClientDBFile: "etc/manager_client.db",
		UserDBFile:   "etc/users.db",
		ListenAddr:   ":9000",
		APIPort:      "8082",
		JWTSecret:    "an-pt-2001",
		JWTExpire:    10 * time.Minute,
	}
}
