package sandbox

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type runRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
}

func RunSandbox(c *gin.Context) {
	var req runRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code and language are required"})
		return
	}

	run, ok := runners[req.Language]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported language"})
		return
	}

	output, err := run(c.Request.Context(), req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": output})
}
