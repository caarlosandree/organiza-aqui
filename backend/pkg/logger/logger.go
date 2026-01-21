package logger

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log        *zap.Logger
	requestIDKey = "request_id"
)

// RequestIDKey retorna a chave usada para armazenar o request ID no context
func RequestIDKey() string {
	return requestIDKey
}

// Init inicializa o logger baseado no ambiente
func Init(environment string) error {
	var config zap.Config
	var err error

	if environment == "production" || environment == "prod" {
		// Configuração para produção: JSON estruturado
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		config.EncoderConfig.LevelKey = "level"
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.StacktraceKey = "stacktrace"
		config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		
		// Adicionar campos padrão estruturados
		config.InitialFields = map[string]interface{}{
			"service": "organiza-aqui-backend",
			"environment": environment,
		}
	} else {
		// Configuração para desenvolvimento: logs coloridos e legíveis
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		config.EncoderConfig.CallerKey = "caller"
		config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "level"
		config.EncoderConfig.StacktraceKey = "stacktrace"
		
		// Desabilitar stacktrace em desenvolvimento para logs mais limpos
		config.Development = true
		config.DisableStacktrace = false
	}

	// Configurar nível de log baseado em variável de ambiente
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		var level zapcore.Level
		if err := level.UnmarshalText([]byte(logLevel)); err == nil {
			config.Level = zap.NewAtomicLevelAt(level)
		}
	}

	log, err = config.Build(
		zap.AddCaller(),
		zap.AddCallerSkip(0),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	return nil
}

// New cria uma nova instância do logger
func New() *zap.Logger {
	if log == nil {
		// Fallback: inicializar com desenvolvimento se não foi inicializado
		if err := Init("development"); err != nil {
			panic(err)
		}
	}
	return log
}

// WithRequestID adiciona o request ID ao logger a partir do context
func WithRequestID(ctx context.Context) *zap.Logger {
	if log == nil {
		log = New()
	}
	
	if requestID := ctx.Value(requestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return log.With(zap.String("request_id", id))
		}
	}
	return log
}

// WithFields cria um logger com campos adicionais
func WithFields(fields ...zap.Field) *zap.Logger {
	if log == nil {
		log = New()
	}
	return log.With(fields...)
}

// Info loga uma mensagem de informação
func Info(msg string, fields ...zap.Field) {
	if log == nil {
		log = New()
	}
	log.Info(msg, fields...)
}

// InfoWithContext loga uma mensagem de informação com request ID do context
func InfoWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Info(msg, fields...)
}

// Error loga uma mensagem de erro
func Error(msg string, err error, fields ...zap.Field) {
	if log == nil {
		log = New()
	}
	allFields := append(fields, zap.Error(err))
	log.Error(msg, allFields...)
}

// ErrorWithContext loga uma mensagem de erro com request ID do context
func ErrorWithContext(ctx context.Context, msg string, err error, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	WithRequestID(ctx).Error(msg, allFields...)
}

// Fatal loga uma mensagem fatal e encerra o programa
func Fatal(msg string, err error, fields ...zap.Field) {
	if log == nil {
		log = New()
	}
	allFields := append(fields, zap.Error(err))
	log.Fatal(msg, allFields...)
}

// FatalWithContext loga uma mensagem fatal com request ID do context e encerra o programa
func FatalWithContext(ctx context.Context, msg string, err error, fields ...zap.Field) {
	allFields := append(fields, zap.Error(err))
	WithRequestID(ctx).Fatal(msg, allFields...)
}

// Debug loga uma mensagem de debug
func Debug(msg string, fields ...zap.Field) {
	if log == nil {
		log = New()
	}
	log.Debug(msg, fields...)
}

// DebugWithContext loga uma mensagem de debug com request ID do context
func DebugWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Debug(msg, fields...)
}

// Warn loga uma mensagem de aviso
func Warn(msg string, fields ...zap.Field) {
	if log == nil {
		log = New()
	}
	log.Warn(msg, fields...)
}

// WarnWithContext loga uma mensagem de aviso com request ID do context
func WarnWithContext(ctx context.Context, msg string, fields ...zap.Field) {
	WithRequestID(ctx).Warn(msg, fields...)
}

// Sync sincroniza o logger (deve ser chamado antes de encerrar a aplicação)
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}

// Helper functions para campos comuns

// WithRequestIDField cria um campo zap com request ID
func WithRequestIDField(requestID string) zap.Field {
	return zap.String("request_id", requestID)
}

// WithUserID cria um campo zap com user ID
func WithUserID(userID interface{}) zap.Field {
	return zap.Any("user_id", userID)
}

// WithDuration cria um campo zap com duração
func WithDuration(duration time.Duration) zap.Field {
	return zap.Duration("duration", duration)
}

// WithMethod cria um campo zap com método HTTP
func WithMethod(method string) zap.Field {
	return zap.String("method", method)
}

// WithPath cria um campo zap com path HTTP
func WithPath(path string) zap.Field {
	return zap.String("path", path)
}

// WithStatusCode cria um campo zap com status code HTTP
func WithStatusCode(statusCode int) zap.Field {
	return zap.Int("status_code", statusCode)
}

// WithIP cria um campo zap com IP
func WithIP(ip string) zap.Field {
	return zap.String("ip", ip)
}