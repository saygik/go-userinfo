package api

import (
	"math"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/saygik/go-userinfo/internal/state"
)

// AdminOverview отдаёт сводную информацию для панели администратора.
// Доступен только пользователям с ролью SysAdmin.
func (h *Handler) AdminOverview(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	if !perms.IsSysAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "у вас нет прав для просмотра панели администратора"})
		return
	}

	overview, err := h.uc.GetAdminOverview()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно получить данные панели администратора",
			"error":   err.Error(),
		})
		return
	}

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	started := state.StartedAt()
	now := time.Now()

	overview.System.Version = buildInfo.Version
	overview.System.Commit = buildInfo.Commit
	overview.System.Built = buildInfo.Built
	overview.System.GoVersion = runtime.Version()
	overview.System.Platform = runtime.GOOS + "/" + runtime.GOARCH
	overview.System.NumCPU = runtime.NumCPU()
	overview.System.GOMAXPROCS = runtime.GOMAXPROCS(0)
	overview.System.Goroutines = runtime.NumGoroutine()
	overview.System.MemoryAllocMB = roundTo1(float64(mem.Alloc) / 1024 / 1024)
	overview.System.MemorySysMB = roundTo1(float64(mem.Sys) / 1024 / 1024)
	overview.System.StartedAt = started.UTC().Format(time.RFC3339)
	overview.System.ServerTime = now.UTC().Format(time.RFC3339)
	overview.System.UptimeSeconds = int64(now.Sub(started).Seconds())

	c.JSON(http.StatusOK, gin.H{"data": overview})
}

func roundTo1(value float64) float64 {
	return math.Round(value*10) / 10
}

// AdminDeleteUserRights удаляет все роли, доступы к доменам и дополнительные
// разделы пользователя. Доступен только пользователям с ролью SysAdmin.
func (h *Handler) AdminDeleteUserRights(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	if !perms.IsSysAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "у вас нет прав для изменения прав пользователей"})
		return
	}

	user := c.Param("username")
	if err := h.uc.DeleteAllUserRights(perms, user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Невозможно удалить права пользователя",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": "OK"})
}

// AdminPermissionAudit отдаёт последние записи аудита изменений прав.
// Доступен только пользователям с ролью SysAdmin.
func (h *Handler) AdminPermissionAudit(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	if !perms.IsSysAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "у вас нет прав для просмотра аудита"})
		return
	}

	limit := 200
	if raw := c.Query("limit"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value > 0 {
			limit = value
		}
	}

	records, err := h.uc.GetPermissionAudit(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно получить аудит изменений прав",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}

// AdminUserSoftwareAudit отдаёт последние записи аудита изменений записей
// «пользователь — система». Доступен только пользователям с ролью SysAdmin.
func (h *Handler) AdminUserSoftwareAudit(c *gin.Context) {
	perms, ok := h.getPerms(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "Сначала войдите в систему"})
		return
	}
	if !perms.IsSysAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "у вас нет прав для просмотра аудита"})
		return
	}

	limit := 200
	if raw := c.Query("limit"); raw != "" {
		if value, err := strconv.Atoi(raw); err == nil && value > 0 {
			limit = value
		}
	}

	records, err := h.uc.GetUserSoftwareAudit(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Невозможно получить аудит изменений записей «пользователь — система»",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": records})
}
