package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"web-ssh-terminal/internal/config"
	"web-ssh-terminal/internal/model"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/ssh"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type resizeMessage struct {
	Type string `json:"type"`
	Cols int    `json:"cols"`
	Rows int    `json:"rows"`
}

type TerminalHandler struct {
	cfg *config.Config
}

func NewTerminalHandler(cfg *config.Config) *TerminalHandler {
	return &TerminalHandler{cfg: cfg}
}

// HandleTerminal handles GET /ws/terminal?robotId=xxx
func (h *TerminalHandler) HandleTerminal(c echo.Context) error {
	robotID := c.QueryParam("robotId")

	var robot *model.Robot
	for i := range h.cfg.Robots {
		if h.cfg.Robots[i].ID == robotID {
			robot = &h.cfg.Robots[i]
			break
		}
	}
	if robot == nil {
		return echo.NewHTTPError(http.StatusNotFound, "robot not found")
	}

	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	sshConfig, err := h.buildSSHConfig(robot)
	if err != nil {
		ws.WriteJSON(map[string]string{"type": "error", "message": err.Error()})
		return nil
	}

	addr := fmt.Sprintf("%s:%d", robot.Host, robot.Port)
	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		ws.WriteJSON(map[string]string{"type": "error", "message": "SSH connection failed: " + err.Error()})
		return nil
	}
	defer conn.Close()

	session, err := conn.NewSession()
	if err != nil {
		ws.WriteJSON(map[string]string{"type": "error", "message": "SSH session failed: " + err.Error()})
		return nil
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm-256color", 24, 80, modes); err != nil {
		ws.WriteJSON(map[string]string{"type": "error", "message": "PTY request failed: " + err.Error()})
		return nil
	}

	sshIn, err := session.StdinPipe()
	if err != nil {
		return err
	}
	sshOut, err := session.StdoutPipe()
	if err != nil {
		return err
	}

	if err := session.Shell(); err != nil {
		ws.WriteJSON(map[string]string{"type": "error", "message": "Shell start failed: " + err.Error()})
		return nil
	}

	ws.WriteJSON(map[string]string{"type": "status", "message": "connected"})

	// SSH stdout → WebSocket (Browser)
	done := make(chan struct{})
	go func() {
		defer close(done)
		buf := make([]byte, 4096)
		for {
			n, err := sshOut.Read(buf)
			if err != nil {
				return
			}
			if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// WebSocket (Browser) → SSH stdin
	go func() {
		for {
			_, msg, err := ws.ReadMessage()
			if err != nil {
				session.Close()
				return
			}

			var resize resizeMessage
			if json.Unmarshal(msg, &resize) == nil && resize.Type == "resize" {
				session.WindowChange(resize.Rows, resize.Cols)
				continue
			}

			sshIn.Write(msg)
		}
	}()

	<-done
	return nil
}

func (h *TerminalHandler) buildSSHConfig(robot *model.Robot) (*ssh.ClientConfig, error) {
	sshCfg := &ssh.ClientConfig{
		User:            robot.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	switch robot.AuthType {
	case model.AuthPassword:
		normalizedID := strings.ToUpper(strings.ReplaceAll(robot.ID, "-", "_"))
		envKey := fmt.Sprintf("ROBOT_%s_PASSWORD", normalizedID)
		password := os.Getenv(envKey)
		sshCfg.Auth = []ssh.AuthMethod{
			ssh.Password(password),
			ssh.KeyboardInteractive(func(user, instruction string, questions []string, echos []bool) ([]string, error) {
				answers := make([]string, len(questions))
				for i := range answers {
					answers[i] = password
				}
				return answers, nil
			}),
		}

	case model.AuthPrivateKey:
		keyPath := h.cfg.SSH.PrivateKeyPath
		key, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		sshCfg.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	}

	return sshCfg, nil
}
