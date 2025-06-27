package main

import (
	"database/sql"
	"fmt"
	"gou-pc/internal/agent"
	"gou-pc/internal/api"
	"gou-pc/internal/api/repository"
	"gou-pc/internal/api/service"
	"gou-pc/internal/config"
	"gou-pc/internal/logutil"
	"gou-pc/internal/tcpserver"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

func InitAgentDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}
	// Tạo bảng managed_clients nếu chưa có
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS managed_clients (
		client_id TEXT PRIMARY KEY,
		agent_id TEXT UNIQUE,
		hardware_id TEXT,
		user_name TEXT,
		last_seen TEXT,
		online INTEGER
	)`)
	if err != nil {
		return nil, err
	}
	// Tạo bảng users nếu chưa có (đủ các trường)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id TEXT PRIMARY KEY,
		username TEXT UNIQUE,
		password TEXT,
		email TEXT,
		full_name TEXT,
		role TEXT,
		created_at TEXT,
		updated_at TEXT
	)`)
	if err != nil {
		return nil, err
	}
	// Tạo user admin mặc định nếu chưa có
	row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", "admin")
	var count int
	err = row.Scan(&count)
	if err != nil {
		return nil, err
	}
	if count == 0 {
		adminID := uuid.NewString()
		fmt.Printf("[DEBUG] Creating default admin user with id: %s, username: admin, password: 1\n", adminID)
		_, err = db.Exec(`INSERT INTO users (id, username, password, email, full_name, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			adminID, "admin", "1", "admin@example.com", "ADMIN", "admin", "2024-01-01T00:00:00Z", "2024-01-01T00:00:00Z")
		if err != nil {
			fmt.Printf("[DEBUG] Failed to create admin user: %v\n", err)
			return nil, err
		}
	} else {
		fmt.Println("[DEBUG] Admin user already exists in DB")
	}
	agent.SetDB(db)
	return db, nil
}

func main() {
	cfg := config.DefaultServerConfig()
	if err := logutil.Init(cfg.LogFile, logutil.DEBUG); err != nil {
		fmt.Printf("Could not open log file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Starting servers...")

	db, err := InitAgentDB(cfg.ClientDBFile)
	if err != nil {
		fmt.Printf("Could not open agent DB: %v\n", err)
		os.Exit(1)
	}

	// Khởi tạo repository với SQLite
	userRepo := repository.NewSQLiteUserRepository(db)
	clientRepo := repository.NewSQLiteClientRepository(db)
	// TODO: logRepo nếu cần

	// Khởi tạo service
	logService := service.NewLogService(cfg.ArchiveFile)
	userService := service.NewUserService(userRepo)
	clientService := service.NewClientService(clientRepo, userRepo)
	// TODO: logService nếu cần

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if err := tcpserver.Start(cfg); err != nil {
			logutil.Error("TCP server error: %v", err)
			os.Exit(1)
		}
	}()
	go func() {
		defer wg.Done()
		fmt.Printf("Serving API at http://localhost:%s/\n", cfg.APIPort)
		api.Start(cfg.APIPort, userService, clientService, logService, clientRepo, cfg.JWTSecret, cfg.JWTExpire)
	}()
	// Static web server
	go func() {
		fmt.Println("Serving static web at http://localhost:8080/")
		err := http.ListenAndServe(":8080", http.FileServer(http.Dir("web")))
		if err != nil {
			logutil.Error("Static web server error: %v", err)
			os.Exit(1)
		}
	}()
	wg.Wait()
}
