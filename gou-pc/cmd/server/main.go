package main

import (
	"fmt"
	"gou-pc/internal/api"
	"gou-pc/internal/api/repository"
	"gou-pc/internal/api/service"
	"gou-pc/internal/config"
	"gou-pc/internal/logutil"
	"gou-pc/internal/tcpserver"
	"os"
	"sync"
)

func main() {
	cfg := config.DefaultServerConfig()
	if err := logutil.Init(cfg.LogFile, logutil.DEBUG); err != nil {
		fmt.Printf("Could not open log file: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Starting servers...")

	// Khởi tạo repository với file path từ config
	userRepo := repository.NewFileUserRepository(cfg.UserDBFile)
	clientRepo := repository.NewFileClientRepository(cfg.ClientDBFile)
	// TODO: logRepo nếu cần

	// Khởi tạo service
	userService := service.NewUserService(userRepo)
	clientService := service.NewClientService(clientRepo, userRepo)
	logService := service.NewLogService(cfg.ArchiveFile)
	// TODO: logService nếu cần

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := tcpserver.Start(cfg); err != nil {
			logutil.Error("TCP server error: %v", err)
			os.Exit(1)
		}
	}()
	go func() {
		defer wg.Done()
		api.Start(cfg.APIPort, userService, clientService, logService, clientRepo, cfg.JWTSecret, cfg.JWTExpire)
	}()
	wg.Wait()
}
