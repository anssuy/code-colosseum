package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func ParseUUIDParam(c *gin.Context, name string) (pgtype.UUID, bool) {
	var id pgtype.UUID
	if err := id.Scan(c.Param(name)); err != nil {
		WriteError(c, http.StatusBadRequest, "invalid "+name)
		return pgtype.UUID{}, false
	}
	return id, true
}
