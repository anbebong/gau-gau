package handler

import (
	"gou-pc/internal/api/response"
	"gou-pc/internal/api/service"
	"gou-pc/internal/logutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var logService service.LogService

func InjectLogService(s service.LogService) { logService = s }

func GetArchiveLogHandler(c *gin.Context) {
	logutil.APIDebug("GetArchiveLogHandler called")
	agentID := c.Query("agent")
	if agentID != "" {
		logutil.APIDebug("GetArchiveLogHandler: filter by agent_id=%s", agentID)
		logs, err := logService.GetLogsByAgentID(agentID)
		if err != nil {
			logutil.APIDebug("GetArchiveLogHandler error: %v", err)
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		response.Success(c, logs)

	} else {
		logs, err := logService.GetAllLogs()
		if err != nil {
			logutil.APIDebug("GetArchiveLogHandler error: %v", err)
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		logutil.APIDebug("GetArchiveLogHandler success, %d logs", len(logs))
		response.Success(c, logs)
	}
}

func GetMyDeviceLogHandler(c *gin.Context) {
	logutil.APIDebug("GetMyDeviceLogHandler called")
	username, ok := c.Get("username")
	if !ok {
		logutil.APIDebug("GetMyDeviceLogHandler missing username in context")
		response.Error(c, http.StatusUnauthorized, "username not found in context")
		return
	}
	clients, err := clientService.GetAllClients()
	if err != nil {
		logutil.APIDebug("GetMyDeviceLogHandler: error getting all clients: %v", err)
		response.Error(c, http.StatusInternalServerError, "Không lấy được danh sách thiết bị")
		return
	}
	var agentIDs []string
	for _, cl := range clients {
		if cl.Username == username {
			agentIDs = append(agentIDs, cl.AgentID)
		}
	}
	if len(agentIDs) == 0 {
		logutil.APIDebug("GetMyDeviceLogHandler: user %v has no agentIDs", username)
		response.Error(c, http.StatusNotFound, "User chưa được gán thiết bị hoặc không tìm thấy agent_id")
		return
	}
	var allLogs []interface{}
	for _, agentID := range agentIDs {
		logs, err := logService.GetLogsByAgentID(agentID)
		if err == nil && logs != nil {
			logutil.APIDebug("GetMyDeviceLogHandler: found %d logs for agentID=%s", len(logs), agentID)
			for _, l := range logs {
				allLogs = append(allLogs, l)
			}
		}
	}
	logutil.APIDebug("GetMyDeviceLogHandler: total logs returned: %d", len(allLogs))
	response.Success(c, allLogs)
}

func GetLogsPagedHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	agentID := c.Query("agent")
	if agentID != "" {
		logs, total, err := logService.GetLogsPagedByAgentID(agentID, page, pageSize)
		if err != nil {
			response.Error(c, http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"logs":  logs,
			"total": total,
		})
		return
	}
	logs, total, err := logService.GetLogsPaged(page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
	})
}

func GetMyDeviceLogPagedHandler(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		response.Error(c, http.StatusUnauthorized, "username not found in context")
		return
	}
	agentID := c.Query("agent")
	if agentID == "" {
		response.Error(c, http.StatusBadRequest, "agent param required")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	// Kiểm tra agentID có thuộc user không
	clients, err := clientService.GetClientsByUsername(username.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	found := false
	for _, cl := range clients {
		if cl.AgentID == agentID {
			found = true
			break
		}
	}
	if !found {
		response.Error(c, http.StatusForbidden, "agent_id not assigned to user")
		return
	}
	logs, total, err := logService.GetLogsPagedByAgentID(agentID, page, pageSize)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"logs":  logs,
		"total": total,
	})
}
