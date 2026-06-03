package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/model"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/repository"
	"github.com/kenshin579/tutorials-go/web/sns-login-session/backend/service"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newSvc(t *testing.T) *service.SessionService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("DB: %v", err)
	}
	db.AutoMigrate(&model.User{}, &model.Session{})
	return service.NewSessionService(repository.NewSessionRepository(db), time.Hour)
}

func TestSessionAuth_NoCookie_Rejected(t *testing.T) {
	e := echo.New()
	h := SessionAuth(newSvc(t))(func(c echo.Context) error { return c.String(200, "ok") })
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	c := e.NewContext(req, httptest.NewRecorder())
	if err := h(c); err == nil {
		t.Fatal("쿠키 없는 요청이 통과됨")
	}
}

func TestSessionAuth_ValidCookie_Passes(t *testing.T) {
	svc := newSvc(t)
	sess, _ := svc.Create(9)

	e := echo.New()
	h := SessionAuth(svc)(func(c echo.Context) error {
		if c.Get("user_id").(uint) != 9 {
			t.Errorf("user_id 기대 9")
		}
		return c.String(200, "ok")
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: sess.ID})
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if err := h(c); err != nil {
		t.Fatalf("유효 쿠키가 거부됨: %v", err)
	}
}
