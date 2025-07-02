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

func HandleListClients(c *gin.Context) {
	logutil.APIDebug("HandleListClients called")
	clients, err := clientService.GetAllClients()
	if err != nil {
		logutil.APIDebug("HandleListClients error: %v", err)
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	logutil.APIDebug("HandleListClients success, %d clients", len(clients))
	response.Success(c, clients)
}

func HandleGetClientByAgentID(c *gin.Context) {
	agentID := c.Param("agent_id")
	if agentID == "" {
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	client, err := clientService.GetClientByAgentID(agentID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if client == nil {
		response.Error(c, http.StatusNotFound, "client not found")
		return
	}
	response.Success(c, client)
}

func HandleGetClientByID(c *gin.Context) {
	clientID := c.Param("client_id")
	if clientID == "" {
		response.Error(c, http.StatusBadRequest, "client_id required")
		return
	}
	client, err := clientService.GetClientByID(clientID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if client == nil {
		response.Error(c, http.StatusNotFound, "client not found")
		return
	}
	response.Success(c, client)
}

func HandleDeleteClientByAgentID(c *gin.Context) {
	var req struct {
		AgentID string `json:"agent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AgentID == "" {
		response.Error(c, http.StatusBadRequest, "agent_id required")
		return
	}
	if err := clientService.DeleteClientByAgentID(req.AgentID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "client deleted successfully"})
}

func HandleDeleteClientByClientID(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ClientID == "" {
		response.Error(c, http.StatusBadRequest, "client_id required")
		return
	}
	if err := clientService.DeleteClientByClientID(req.ClientID); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "client deleted successfully"})
}

func HandleAssignUserToClientByAgentID(c *gin.Context) {
	var req struct {
		AgentID  string `json:"agent_id"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.AgentID == "" || req.Username == "" {
		response.Error(c, http.StatusBadRequest, "agent_id and username required")
		return
	}
	if err := clientService.AssignUserToClientByAgentID(req.AgentID, req.Username); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "user assigned to client successfully"})
}

func HandleAssignUserToClientByClientID(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id"`
		Username string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ClientID == "" || req.Username == "" {
		response.Error(c, http.StatusBadRequest, "client_id and username required")
		return
	}
	if err := clientService.AssignUserToClientByClientID(req.ClientID, req.Username); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
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

func HandleListMyClients(c *gin.Context) {
	username, ok := c.Get("username")
	if !ok {
		response.Error(c, http.StatusUnauthorized, "username not found in context")
		return
	}
	clients, err := clientService.GetClientsByUsername(username.(string))
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, clients)
}
