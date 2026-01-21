package middleware

import (
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/organiza-aqui/backend/pkg/logger"
	"go.uber.org/zap"
)

// RecoverMiddleware retorna um middleware de recuperação
func RecoverMiddleware() echo.MiddlewareFunc {
	return middleware.Recover()
}

// RequestIDMiddleware adiciona um request ID único a cada requisição
func RequestIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Verificar se já existe um request ID no header
			requestID := c.Request().Header.Get("X-Request-ID")
			if requestID == "" {
				// Gerar novo UUID se não existir
				requestID = uuid.New().String()
			}

			// Adicionar ao context
			c.Set(logger.RequestIDKey(), requestID)

			// Adicionar ao header da resposta
			c.Response().Header().Set("X-Request-ID", requestID)

			return next(c)
		}
	}
}

// StructuredLoggerMiddleware retorna um middleware de logging estruturado com zap
func StructuredLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Obter request ID do context
			requestID := c.Get(logger.RequestIDKey())
			requestIDStr := ""
			if id, ok := requestID.(string); ok {
				requestIDStr = id
			}

			// Criar logger com campos estruturados
			log := logger.New().With(
				zap.String("request_id", requestIDStr),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Request().URL.Path),
				zap.String("ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			)

			// Log da requisição recebida
			log.Info("Request recebida",
				zap.String("query", c.Request().URL.RawQuery),
			)

			// Processar requisição
			err := next(c)

			// Calcular duração
			duration := time.Since(start)

			// Campos de resposta
			fields := []zap.Field{
				zap.Int("status_code", c.Response().Status),
				zap.Duration("duration", duration),
				zap.Int64("size", c.Response().Size),
			}

			// Log da resposta
			if err != nil {
				log.Error("Request finalizada com erro", append(fields, zap.Error(err))...)
				return err
			}

			// Log baseado no status code
			switch {
			case c.Response().Status >= 500:
				log.Error("Request finalizada", fields...)
			case c.Response().Status >= 400:
				log.Warn("Request finalizada", fields...)
			default:
				log.Info("Request finalizada com sucesso", fields...)
			}

			return nil
		}
	}
}

// LoggerMiddleware retorna um middleware de logging (legado, usar StructuredLoggerMiddleware)
func LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.Logger()
}