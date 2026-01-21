package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/organiza-aqui/backend/pkg/response"
	"github.com/redis/go-redis/v9"
)

// RateLimitConfig contém configurações do rate limiting
type RateLimitConfig struct {
	RedisClient          *redis.Client
	RequestsPerMinute   int
	Enabled              bool
	SkipSuccessfulRequests bool
	SkipFailedRequests   bool
}

// DefaultRateLimitConfig retorna configuração padrão
func DefaultRateLimitConfig(redisClient *redis.Client) RateLimitConfig {
	return RateLimitConfig{
		RedisClient:          redisClient,
		RequestsPerMinute:   60,
		Enabled:              true,
		SkipSuccessfulRequests: false,
		SkipFailedRequests:   false,
	}
}

// RateLimitMiddleware retorna um middleware de rate limiting usando Redis
func RateLimitMiddleware(config RateLimitConfig) echo.MiddlewareFunc {
	if !config.Enabled {
		// Se desabilitado, retornar middleware vazio
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return next
		}
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Identificar cliente (IP ou UserID)
			identifier := getClientIdentifier(c)
			
			// Chave Redis baseada no identificador
			key := fmt.Sprintf("ratelimit:%s", identifier)
			
			// Janela de tempo (1 minuto)
			window := time.Minute
			
			// Obter contador atual
			ctx := context.Background()
			count, err := config.RedisClient.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				// Erro ao acessar Redis - permitir requisição (degradação graciosa)
				return next(c)
			}
			
			// Se não existe, criar com TTL
			if err == redis.Nil {
				err = config.RedisClient.Set(ctx, key, 1, window).Err()
				if err != nil {
					// Erro ao criar - permitir requisição
					return next(c)
				}
				count = 1
			} else {
				// Incrementar contador
				newCount, err := config.RedisClient.Incr(ctx, key).Result()
				if err != nil {
					// Erro ao incrementar - permitir requisição
					return next(c)
				}
				count = int(newCount)
				
				// Se é a primeira requisição após criação, definir TTL
				if count == 2 {
					config.RedisClient.Expire(ctx, key, window)
				}
			}
			
			// Calcular limites restantes
			remaining := config.RequestsPerMinute - count
			if remaining < 0 {
				remaining = 0
			}
			
			// Adicionar headers de rate limit
			c.Response().Header().Set("X-RateLimit-Limit", strconv.Itoa(config.RequestsPerMinute))
			c.Response().Header().Set("X-RateLimit-Remaining", strconv.Itoa(remaining))
			c.Response().Header().Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(window).Unix(), 10))
			
			// Verificar se excedeu o limite
			if count > config.RequestsPerMinute {
				return c.JSON(http.StatusTooManyRequests, response.ErrorResponse{
					Error:   "too_many_requests",
					Message: fmt.Sprintf("Limite de requisições excedido. Máximo: %d por minuto", config.RequestsPerMinute),
				})
			}
			
			// Processar requisição
			err = next(c)
			
			// Opcionalmente, não contar requisições bem-sucedidas ou com falha
			if config.SkipSuccessfulRequests && c.Response().Status < 400 {
				// Decrementar contador se for bem-sucedida
				config.RedisClient.Decr(ctx, key)
			}
			if config.SkipFailedRequests && c.Response().Status >= 400 {
				// Decrementar contador se falhou
				config.RedisClient.Decr(ctx, key)
			}
			
			return err
		}
	}
}

// getClientIdentifier retorna identificador único do cliente (IP ou UserID)
func getClientIdentifier(c echo.Context) string {
	// Tentar obter user_id do contexto (se autenticado)
	if userID, ok := c.Get("user_id").(string); ok && userID != "" {
		return fmt.Sprintf("user:%s", userID)
	}
	
	// Caso contrário, usar IP
	return fmt.Sprintf("ip:%s", c.RealIP())
}
