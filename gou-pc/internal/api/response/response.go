package response

import "github.com/gin-gonic/gin"

// Chuẩn hóa response thành công
func Success(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"success": true, "data": data})
}

// Chuẩn hóa response lỗi
func Error(c *gin.Context, code int, msg string) {
	c.JSON(code, gin.H{"success": false, "error": msg})
}
