package handlers

import (
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/go-enterprise-api/internal/database"
	"github.com/yourusername/go-enterprise-api/pkg/response"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{
		db: db,
	}
}

// HealthStatus represents the health status response
type HealthStatus struct {
	Status   string            `json:"status"`
	Version  string            `json:"version"`
	Services map[string]string `json:"services"`
}

// Health returns a simple health check
// @Summary Health check
// @Description Basic health check endpoint
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "healthy",
	})
}

// Ready returns readiness status including dependencies
// @Summary Readiness check
// @Description Check if the service and its dependencies are ready
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 503 {object} response.Response
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	services := make(map[string]string)

	// Check database
	if err := h.db.HealthCheck(); err != nil {
		services["database"] = "unhealthy"
		c.JSON(503, gin.H{
			"status":   "not ready",
			"services": services,
		})
		return
	}
	services["database"] = "healthy"

	response.Success(c, gin.H{
		"status":   "ready",
		"services": services,
	})
}

// Live returns liveness status
// @Summary Liveness check
// @Description Check if the service is alive
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "alive",
	})
}

// Info returns system information
// @Summary System info
// @Description Get system information (admin only)
// @Tags health
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /health/info [get]
func (h *HealthHandler) Info(c *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	response.Success(c, gin.H{
		"go_version":    runtime.Version(),
		"go_os":         runtime.GOOS,
		"go_arch":       runtime.GOARCH,
		"cpu_count":     runtime.NumCPU(),
		"goroutines":    runtime.NumGoroutine(),
		"heap_alloc_mb": memStats.HeapAlloc / 1024 / 1024,
		"sys_mb":        memStats.Sys / 1024 / 1024,
	})
}
