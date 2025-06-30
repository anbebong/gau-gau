package handler

import (
	"encoding/json"
	"fmt"
	"gou-pc/internal/api/model"
	"gou-pc/internal/api/response"
	"gou-pc/internal/api/service"
	"gou-pc/internal/logutil"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	userService service.UserService
	jwtSecret   []byte
	jwtExpire   int64 // seconds
)

func InjectUserService(s service.UserService) { userService = s }
func InjectJWTConfig(secret string, expireSeconds int64) {
	jwtSecret = []byte(secret)
	jwtExpire = expireSeconds
}

func LoginHandler(c *gin.Context) {
	logutil.APIDebug("LoginHandler called")
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// Đọc body ra chuỗi để debug
	bodyBytes, _ := c.GetRawData()
	fmt.Printf("[DEBUG] LoginHandler: raw request body: %s\n", string(bodyBytes))
	if err := json.Unmarshal(bodyBytes, &req); err != nil {
		logutil.APIDebug("LoginHandler: invalid request body (unmarshal)")
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	fmt.Printf("[DEBUG] LoginHandler: parsed username: %s, password: %s\n", req.Username, req.Password)
	if req.Username == "" || req.Password == "" {
		logutil.APIDebug("LoginHandler: missing username or password")
		response.Error(c, http.StatusBadRequest, "username and password required")
		return
	}
	user, err := userService.UserGetByUsername(req.Username)
	fmt.Printf("[DEBUG] userService.GetByUsername(%s) => user: %+v, err: %v\n", req.Username, user, err)
	if err != nil || user == nil {
		logutil.APIDebug("LoginHandler: invalid username or password")
		response.Error(c, http.StatusUnauthorized, "invalid username or password")
		return
	}
	fmt.Printf("[DEBUG] LoginHandler: user in DB: username=%s, password=%s\n", user.Username, user.Password)
	if user.Password != req.Password {
		logutil.APIDebug("LoginHandler: password mismatch for user %s", req.Username)
		response.Error(c, http.StatusUnauthorized, "invalid username or password")
		return
	}
	// Sinh JWT thực tế
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Duration(jwtExpire) * time.Second).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "could not generate token")
		return
	}
	// Ẩn trường password trước khi trả về user
	safeUser := gin.H{
		"id":         user.ID,
		"username":   user.Username,
		"full_name":  user.FullName, // Thêm trường full_name nếu có
		"email":      user.Email,
		"role":       user.Role,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	}
	response.Success(c, gin.H{"token": tokenString, "user": safeUser})
}

func CreateUserHandler(c *gin.Context) {
	logutil.APIDebug("CreateUserHandler called")
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logutil.APIDebug("CreateUserHandler: invalid request body")
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" || req.Password == "" {
		logutil.APIDebug("CreateUserHandler: missing username or password")
		response.Error(c, http.StatusBadRequest, "username and password required")
		return
	}
	now := time.Now().UTC().Format(time.RFC3339)
	user := &model.User{
		Username:  req.Username,
		Password:  req.Password, // Đã dùng Password, không phải PasswordHash
		FullName:  req.FullName,
		Email:     req.Email,
		Role:      "user",
		CreatedAt: now,
		UpdatedAt: now,
	}
	// Sinh UUID nếu chưa có
	if user.ID == "" {
		user.ID = generateUUID()
	}
	logutil.APIDebug("CreateUserHandler: creating user %+v", user)
	if err := userService.UserCreate(user); err != nil {
		logutil.APIDebug("CreateUserHandler: failed to create user: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	// Lấy lại user vừa tạo từ DB để đảm bảo trả về đúng dữ liệu đã lưu
	createdUser, err := userService.UserGetByUsername(user.Username)
	if err != nil || createdUser == nil {
		logutil.APIDebug("CreateUserHandler: failed to fetch created user: %v", err)
		response.Error(c, http.StatusInternalServerError, "failed to fetch created user")
		return
	}
	logutil.APIDebug("CreateUserHandler: user created successfully: %+v", createdUser)
	safeUser := gin.H{
		"id":         createdUser.ID,
		"username":   createdUser.Username,
		"full_name":  createdUser.FullName,
		"email":      createdUser.Email,
		"role":       createdUser.Role,
		"created_at": createdUser.CreatedAt,
		"updated_at": createdUser.UpdatedAt,
	}
	response.Success(c, gin.H{"user": safeUser})
}

// Helper sinh UUID
func generateUUID() string {
	return uuid.NewString()
}

func ChangePasswordHandler(c *gin.Context) {
	var req struct {
		UserID      string `json:"username"`
		NewPassword string `json:"new_password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == "" || req.NewPassword == "" {
		response.Error(c, http.StatusBadRequest, "username and new_password required")
		return
	}
	// TODO: hash password trước khi lưu
	user, err := userService.UserGetByUsername(req.UserID)
	if err != nil || user == nil {
		response.Error(c, http.StatusBadRequest, "user not found")
		return
	}
	user.Password = req.NewPassword
	if err := userService.UserUpdate(user); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "password changed successfully"})
}

func UpdateUserHandler(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" {
		response.Error(c, http.StatusBadRequest, "username required")
		return
	}
	if req.FullName == "" {
		response.Error(c, http.StatusBadRequest, "full_name required")
		return
	}
	if err := userService.UserUpdate(&req); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "user updated successfully"})
}

func ListUsersHandler(c *gin.Context) {
	users, err := userService.UserGetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	// Ẩn trường password, liệt kê đầy đủ user, trả về đúng thứ tự trường
	type safeUser struct {
		ID        string `json:"id"`
		Username  string `json:"username"`
		FullName  string `json:"full_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
	var safeUsers []safeUser
	for _, u := range users {
		safeUsers = append(safeUsers, safeUser{
			ID:        u.ID,
			Username:  u.Username,
			FullName:  u.FullName,
			Email:     u.Email,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	// Sắp xếp theo username tăng dần, nhưng không loại bỏ user nào
	sort.Slice(safeUsers, func(i, j int) bool {
		return safeUsers[i].Username < safeUsers[j].Username
	})
	response.Success(c, safeUsers)
}

func UpdateUserInfoHandler(c *gin.Context) {
	logutil.APIDebug("UpdateUserInfoHandler called")
	var req struct {
		Username string `json:"username"`
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logutil.APIDebug("UpdateUserInfoHandler: invalid request body")
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Username == "" {
		logutil.APIDebug("UpdateUserInfoHandler: missing username")
		response.Error(c, http.StatusBadRequest, "username required")
		return
	}
	user, err := userService.UserGetByUsername(req.Username)
	if err != nil || user == nil {
		logutil.APIDebug("UpdateUserInfoHandler: user not found: %v", err)
		response.Error(c, http.StatusBadRequest, "user not found")
		return
	}
	updated := false
	if req.FullName != "" && req.FullName != user.FullName {
		user.FullName = req.FullName
		updated = true
	}
	if req.Email != "" && req.Email != user.Email {
		user.Email = req.Email
		updated = true
	}
	if !updated {
		logutil.APIDebug("UpdateUserInfoHandler: nothing to update for user %s", req.Username)
		response.Success(c, gin.H{"message": "no changes"})
		return
	}
	if err := userService.UserUpdate(user); err != nil {
		logutil.APIDebug("UpdateUserInfoHandler: failed to update user: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	logutil.APIDebug("UpdateUserInfoHandler: user info updated successfully: %s", req.Username)
	response.Success(c, gin.H{"message": "user info updated successfully"})
}

func DeleteUserHandler(c *gin.Context) {
	logutil.APIDebug("DeleteUserHandler called")
	var req struct {
		UserID string `json:"username"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logutil.APIDebug("DeleteUserHandler: invalid request body")
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.UserID == "" {
		logutil.APIDebug("DeleteUserHandler: missing username")
		response.Error(c, http.StatusBadRequest, "username required")
		return
	}
	if err := userService.UserDeleteByUsername(req.UserID); err != nil {
		logutil.APIDebug("DeleteUserHandler: failed to delete user: %v", err)
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	logutil.APIDebug("DeleteUserHandler: user deleted successfully: %s", req.UserID)
	response.Success(c, gin.H{"message": "user deleted successfully"})
}
