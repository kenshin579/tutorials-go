package handler

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"web-ssh-terminal/internal/config"
	"web-ssh-terminal/internal/model"

	"github.com/labstack/echo/v4"
)

type RobotHandler struct {
	cfg *config.Config
}

func NewRobotHandler(cfg *config.Config) *RobotHandler {
	return &RobotHandler{cfg: cfg}
}

// ListRobots handles GET /api/robots
func (h *RobotHandler) ListRobots(c echo.Context) error {
	robots := make([]model.Robot, len(h.cfg.Robots))
	copy(robots, h.cfg.Robots)

	for i := range robots {
		robots[i].IsOnline = checkSSHPort(robots[i].Host, robots[i].Port)
	}

	return c.JSON(http.StatusOK, robots)
}

func checkSSHPort(host string, port int) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}
