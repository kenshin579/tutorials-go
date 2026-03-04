package custom

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// ZapLoggerConfig는 구조화된 로깅 미들웨어 설정
type ZapLoggerConfig struct {
	// Skipper는 미들웨어를 건너뛸 조건을 정의한다
	Skipper middleware.Skipper
}

// DefaultZapLoggerConfig는 기본 설정값
var DefaultZapLoggerConfig = ZapLoggerConfig{
	Skipper: middleware.DefaultSkipper,
}

// ZapLogger는 zap 로거를 사용하는 요청 로깅 미들웨어를 반환한다
func ZapLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return ZapLoggerWithConfig(logger, DefaultZapLoggerConfig)
}

// ZapLoggerWithConfig는 설정을 받아 zap 로깅 미들웨어를 반환한다
func ZapLoggerWithConfig(logger *zap.Logger, config ZapLoggerConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = DefaultZapLoggerConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			start := time.Now()
			req := c.Request()
			res := c.Response()

			err := next(c)

			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("path", req.URL.Path),
				zap.Int("status", res.Status),
				zap.Duration("latency", time.Since(start)),
				zap.String("remote_ip", c.RealIP()),
				zap.String("request_id", res.Header().Get(echo.HeaderXRequestID)),
			}

			if err != nil {
				fields = append(fields, zap.Error(err))
				logger.Error("request", fields...)
			} else {
				logger.Info("request", fields...)
			}

			return err
		}
	}
}
