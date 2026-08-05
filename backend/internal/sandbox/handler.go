package sandbox

import (
	"net/http"

	"github.com/anssuy/code-colosseum/backend/internal/httpx"

	"github.com/gin-gonic/gin"
)

type runRequest struct {
	Code     string `json:"code" binding:"required"`
	Language string `json:"language" binding:"required"`
}

func RunSandbox(c *gin.Context) {
	var req runRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpx.WriteError(c, http.StatusBadRequest, "code and language are required")
		return
	}

	run, ok := runners[req.Language]
	if !ok {
		httpx.WriteError(c, http.StatusBadRequest, "unsupported language")
		return
	}

	output, err := run(c.Request.Context(), req.Code)
	if err != nil {
		httpx.WriteError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"output": output})
}
