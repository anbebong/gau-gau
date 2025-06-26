package handler

import (
	"gou-pc/internal/api/response"
	"gou-pc/internal/api/service"
	"gou-pc/internal/logutil"

	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	clientService service.ClientService
	otpService    service.OTPService
)

func InjectClientService(s service.ClientService) { clientService = s }
func InjectOTPService(s service.OTPService)       { otpService = s }

func ListClientsHandler(c *gin.Context) {
	logutil.Debug("ListClientsHandler called")
	clients, err := clientService.GetAllClients()
	if err != nil {
		logutil.Debug("ListClientsHandler error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	logutil.Debug("ListClientsHandler success, %d clients", len(clients))
	response.Success(c, clients)
}

func GetClientByAgentIDHandler(c *gin.Context) {
	logutil.Debug("GetClientByAgentIDHandler called")
	agentID := c.Param("agent_id")
	if agentID == "" {
		logutil.Debug("GetClientByAgentIDHandler missing agent_id")
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	client, err := clientService.GetClientByAgentID(agentID)
	if err != nil {
		logutil.Debug("GetClientByAgentIDHandler error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if client == nil {
		logutil.Debug("GetClientByAgentIDHandler not found: %s", agentID)
		response.Error(c, http.StatusNotFound, "client not found")
		return
	}
	logutil.Debug("GetClientByAgentIDHandler success: %s", agentID)
	response.Success(c, client)
}

func GetClientsByUserIDHandler(c *gin.Context) {
	logutil.Debug("GetClientsByUserIDHandler called")
	userID := c.Param("user_id")
	if userID == "" {
		logutil.Debug("GetClientsByUserIDHandler missing user_id")
		response.Error(c, http.StatusBadRequest, "user_id required")
		return
	}
	clients, err := clientService.GetClientsByUserID(userID)
	if err != nil {
		logutil.Debug("GetClientsByUserIDHandler error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	logutil.Debug("GetClientsByUserIDHandler success, %d clients", len(clients))
	response.Success(c, clients)
}

func DeleteClientHandler(c *gin.Context) {
	logutil.Debug("DeleteClientHandler called")
	var req struct {
		AgentID string `json:"agent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logutil.Debug("DeleteClientHandler: invalid request body")
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AgentID == "" {
		logutil.Debug("DeleteClientHandler: missing agent_id")
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	if err := clientService.DeleteClient(req.AgentID); err != nil {
		logutil.Debug("DeleteClientHandler: failed to delete client: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	logutil.Debug("DeleteClientHandler: client deleted successfully: %s", req.AgentID)
	response.Success(c, gin.H{"message": "client deleted successfully"})
}

func HandleAssignUser(c *gin.Context) {
	logutil.Debug("HandleAssignUser called")
	var req struct {
		AgentID  string `json:"agent_id"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logutil.Debug("HandleAssignUser: invalid request body")
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AgentID == "" || req.Username == "" {
		logutil.Debug("HandleAssignUser: missing agent_id or username")
		response.Error(c, http.StatusBadRequest, "agent_id and username required")
		return
	}
	if err := clientService.AssignUserToClient(req.AgentID, req.Username); err != nil {
		logutil.Debug("HandleAssignUser: failed to assign user: %v", err)
		response.Error(c, http.StatusNotFound, err.Error())
		return
	}
	logutil.Debug("HandleAssignUser: user assigned to client successfully: agent_id=%s, username=%s", req.AgentID, req.Username)
	response.Success(c, gin.H{"message": "user assigned to client successfully"})
}

func GetOTPByAgentIDHandler(c *gin.Context) {
	agentID := c.Param("agent_id")
	if agentID == "" {
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	otp, secondsLeft, err := otpService.GetOTPByAgentIDWithExpire(agentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"agent_id": agentID, "otp": otp, "expire_in": secondsLeft})
}

func GetMyOTPHandler(c *gin.Context) {
	agentID := c.Query("agent_id")
	if agentID == "" {
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	otp, secondsLeft, err := otpService.GetOTPByAgentIDWithExpire(agentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, gin.H{"agent_id": agentID, "otp": otp, "expire_in": secondsLeft})
}
