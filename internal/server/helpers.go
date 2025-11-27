package server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const userContextKey = "user"

type AuthContext struct {
	UserID uint
	Role   string
}

func getAuthContext(c *gin.Context) *AuthContext {
	val, ok := c.Get(userContextKey)
	if !ok {
		return nil
	}
	if ctx, ok := val.(*AuthContext); ok {
		return ctx
	}
	return nil
}

func respondError(c *gin.Context, status int, msg string) {
	c.JSON(status, gin.H{"error": msg})
}

func parsePagination(q PaginationQuery) (page, pageSize int, sort string) {
	page = q.Page
	pageSize = q.PageSize
	sort = q.Sort

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 10
	}
	if sort == "" {
		sort = "created_at desc"
	}
	return
}

func sanitizeSort(sort string, allowed map[string]string, fallback string) string {
	if mapped, ok := allowed[sort]; ok {
		return mapped
	}
	if sort == "" {
		return fallback
	}
	return fallback
}

func errorsIs(err error, target error) bool {
	return errors.Is(err, target)
}

func parseUintParam(c *gin.Context, name string) (uint, bool) {
	val := c.Param(name)
	id64, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, "invalid "+name)
		return 0, false
	}
	return uint(id64), true
}
