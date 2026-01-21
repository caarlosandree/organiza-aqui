package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/organiza-aqui/backend/docs"
	"github.com/organiza-aqui/backend/internal/config"
	"github.com/organiza-aqui/backend/internal/handler"
	appMiddleware "github.com/organiza-aqui/backend/internal/middleware"
	"github.com/organiza-aqui/backend/internal/repository"
	"github.com/organiza-aqui/backend/internal/service"
	"github.com/organiza-aqui/backend/pkg/logger"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// @title Organiza Aqui API
// @version 1.0
// @description API do sistema Organiza Aqui
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Carregar configurações primeiro para obter o ambiente
	cfg, err := config.LoadConfig()
	if err != nil {
		// Se não conseguir carregar config, usar logger básico
		panic(err)
	}

	// Inicializar logger com base no ambiente
	if err := logger.Init(cfg.Environment); err != nil {
		panic(err)
	}
	defer logger.Sync()

	log := logger.New()
	log.Info("Logger inicializado",
		zap.String("environment", cfg.Environment),
	)

	// Conectar ao banco de dados
	db, err := config.NewDB(cfg)
	if err != nil {
		logger.Fatal("Erro ao conectar ao banco de dados", err,
			zap.String("host", cfg.DB.Host),
			zap.String("port", cfg.DB.Port),
			zap.String("database", cfg.DB.Name),
		)
	}
	defer db.Close()

	log.Info("Conectado ao banco de dados com sucesso",
		zap.String("host", cfg.DB.Host),
		zap.String("port", cfg.DB.Port),
		zap.String("database", cfg.DB.Name),
	)

	// Executar migrações automaticamente
	if err := config.RunMigrations(db, "./migrations"); err != nil {
		logger.Fatal("Erro ao executar migrações", err)
	}
	log.Info("Migrações executadas com sucesso")

	// Conectar ao Redis
	redisClient, err := config.NewRedis(cfg)
	if err != nil {
		logger.Fatal("Erro ao conectar ao Redis", err,
			zap.String("host", cfg.Redis.Host),
			zap.String("port", cfg.Redis.Port),
			zap.Int("database", cfg.Redis.Database),
		)
	}
	defer redisClient.Close()

	log.Info("Conectado ao Redis com sucesso",
		zap.String("host", cfg.Redis.Host),
		zap.String("port", cfg.Redis.Port),
		zap.Int("database", cfg.Redis.Database),
	)

	// Inicializar Echo
	e := echo.New()

	// Middlewares (ordem importa!)
	// 1. Request ID primeiro para rastreamento
	e.Use(appMiddleware.RequestIDMiddleware())
	// 2. Logging estruturado
	e.Use(appMiddleware.StructuredLoggerMiddleware())
	// 3. Recovery
	e.Use(appMiddleware.RecoverMiddleware())
	// 4. CORS
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-Request-ID"},
	}))
	// 5. Rate Limiting
	rateLimitConfig := appMiddleware.DefaultRateLimitConfig(redisClient)
	rateLimitConfig.RequestsPerMinute = cfg.RateLimit.RequestsPerMinute
	rateLimitConfig.Enabled = cfg.RateLimit.Enabled
	e.Use(appMiddleware.RateLimitMiddleware(rateLimitConfig))

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Swagger documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Inicializar camadas
	repositories := repository.NewRepositories(db, redisClient)
	services := service.NewServices(db, repositories, cfg, redisClient)
	handlers := handler.NewHandlers(services)

	// Setup de rotas
	handler.SetupRoutes(e, handlers)

	// Configurar cron job para sincronização semanal de bancos (domingo às 02:00)
	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("0 0 2 * * 0", func() {
		log.Info("Executando sincronização semanal de bancos")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		if err := services.Bank.SyncBanks(ctx); err != nil {
			logger.Error("Erro ao sincronizar bancos no cron job", err)
		} else {
			log.Info("Sincronização semanal de bancos concluída com sucesso")
		}
	})
	if err != nil {
		logger.Fatal("Erro ao configurar cron job", err)
	}
	c.Start()
	log.Info("Cron job de sincronização de bancos configurado (domingo às 02:00)")

	// Executar sincronização inicial de bancos
	go func() {
		log.Info("Executando sincronização inicial de bancos")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		if err := services.Bank.SyncBanks(ctx); err != nil {
			logger.Error("Erro na sincronização inicial de bancos", err)
		} else {
			log.Info("Sincronização inicial de bancos concluída com sucesso")
		}
	}()

	// Iniciar servidor
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Erro ao iniciar servidor", err,
				zap.String("port", port),
				zap.String("host", cfg.Server.Host),
			)
		}
	}()

	log.Info("Servidor iniciado",
		zap.String("port", port),
		zap.String("host", cfg.Server.Host),
		zap.String("environment", cfg.Environment),
	)

	// Aguardar sinal de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("Sinal de encerramento recebido, iniciando graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		logger.Fatal("Erro ao encerrar servidor", err)
	}

	log.Info("Servidor encerrado com sucesso")
}