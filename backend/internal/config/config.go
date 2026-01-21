package config

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

// Config contém todas as configurações da aplicação
type Config struct {
	DB struct {
		Host     string `mapstructure:"DB_HOST_DEV"`
		Port     string `mapstructure:"DB_PORT_DEV"`
		User     string `mapstructure:"DB_USER_DEV"`
		Password string `mapstructure:"DB_PASSWORD_DEV"`
		Name     string `mapstructure:"DB_NAME_DEV"`
	} `mapstructure:"DB"`

	Redis struct {
		Host     string `mapstructure:"REDIS_HOST_DEV"`
		Port     string `mapstructure:"REDIS_PORT_DEV"`
		Database int    `mapstructure:"REDIS_DATABASE_DEV"`
	} `mapstructure:"REDIS"`

	Server struct {
		Port string `mapstructure:"SERVER_PORT"`
		Host string `mapstructure:"SERVER_HOST"`
	} `mapstructure:"SERVER"`

	JWT struct {
		Secret         string `mapstructure:"JWT_SECRET"`
		ExpirationHours int   `mapstructure:"JWT_EXPIRATION_HOURS"`
	} `mapstructure:"JWT"`

	RateLimit struct {
		RequestsPerMinute int  `mapstructure:"RATE_LIMIT_REQUESTS_PER_MINUTE"`
		Enabled          bool `mapstructure:"RATE_LIMIT_ENABLED"`
	} `mapstructure:"RATE_LIMIT"`

	Environment string `mapstructure:"ENVIRONMENT"`
}

// LoadConfig carrega as configurações do arquivo .env
func LoadConfig() (*Config, error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Vincular variáveis de ambiente explicitamente
	// DB
	viper.BindEnv("DB_HOST_DEV")
	viper.BindEnv("DB_PORT_DEV")
	viper.BindEnv("DB_USER_DEV")
	viper.BindEnv("DB_PASSWORD_DEV")
	viper.BindEnv("DB_NAME_DEV")

	// Redis
	viper.BindEnv("REDIS_HOST_DEV")
	viper.BindEnv("REDIS_PORT_DEV")
	viper.BindEnv("REDIS_DATABASE_DEV")

	// Server
	viper.BindEnv("SERVER_PORT")
	viper.BindEnv("SERVER_HOST")

	// JWT
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("JWT_EXPIRATION_HOURS")

	// Rate Limit
	viper.BindEnv("RATE_LIMIT_REQUESTS_PER_MINUTE")
	viper.BindEnv("RATE_LIMIT_ENABLED")

	// Environment
	viper.BindEnv("ENVIRONMENT")

	if err := viper.ReadInConfig(); err != nil {
		// Não retornar erro se o arquivo não existir, usar variáveis de ambiente
	}

	// Criar config manualmente a partir das variáveis de ambiente
	config := Config{}
	config.DB.Host = viper.GetString("DB_HOST_DEV")
	config.DB.Port = viper.GetString("DB_PORT_DEV")
	config.DB.User = viper.GetString("DB_USER_DEV")
	config.DB.Password = viper.GetString("DB_PASSWORD_DEV")
	config.DB.Name = viper.GetString("DB_NAME_DEV")

	config.Redis.Host = viper.GetString("REDIS_HOST_DEV")
	config.Redis.Port = viper.GetString("REDIS_PORT_DEV")
	config.Redis.Database = viper.GetInt("REDIS_DATABASE_DEV")

	config.Server.Port = viper.GetString("SERVER_PORT")
	config.Server.Host = viper.GetString("SERVER_HOST")

	config.JWT.Secret = viper.GetString("JWT_SECRET")
	config.JWT.ExpirationHours = viper.GetInt("JWT_EXPIRATION_HOURS")

	config.RateLimit.RequestsPerMinute = viper.GetInt("RATE_LIMIT_REQUESTS_PER_MINUTE")
	config.RateLimit.Enabled = viper.GetBool("RATE_LIMIT_ENABLED")

	config.Environment = viper.GetString("ENVIRONMENT")

	// Validação de variáveis obrigatórias
	if config.DB.Host == "" {
		return nil, errors.New("DB_HOST_DEV é obrigatório")
	}
	if config.DB.Name == "" {
		return nil, errors.New("DB_NAME_DEV é obrigatório")
	}
	if config.DB.User == "" {
		return nil, errors.New("DB_USER_DEV é obrigatório")
	}
	if config.DB.Password == "" {
		return nil, errors.New("DB_PASSWORD_DEV é obrigatório")
	}

	// Validação de JWT
	if config.JWT.Secret == "" {
		return nil, errors.New("JWT_SECRET é obrigatório")
	}
	if config.JWT.ExpirationHours == 0 {
		config.JWT.ExpirationHours = 24 // Padrão: 24 horas
	}

	// Valores padrão
	if config.DB.Port == "" {
		config.DB.Port = "5432"
	}
	if config.Server.Port == "" {
		config.Server.Port = "8080"
	}
	if config.Server.Host == "" {
		config.Server.Host = "localhost"
	}
	if config.Environment == "" {
		config.Environment = "development"
	}
	if config.RateLimit.RequestsPerMinute == 0 {
		config.RateLimit.RequestsPerMinute = 60 // Padrão: 60 requisições por minuto
	}
	// Rate limit habilitado por padrão se não especificado
	if !viper.IsSet("RATE_LIMIT_ENABLED") {
		config.RateLimit.Enabled = true
	}

	return &config, nil
}

// NewDB cria uma nova conexão com o banco de dados PostgreSQL
func NewDB(config *Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Password,
		config.DB.Name,
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	// Configuração de connection pooling
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	// Testar conexão
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("erro ao fazer ping no banco de dados: %w", err)
	}

	return db, nil
}

// NewRedis cria uma nova conexão com o Redis
func NewRedis(config *Config) (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s", config.Redis.Host, config.Redis.Port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       config.Redis.Database,
		Password: "", // Adicione se necessário
	})

	// Testar conexão
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("erro ao conectar ao Redis: %w", err)
	}

	return client, nil
}

// RunMigrations executa as migrações do banco de dados automaticamente
func RunMigrations(db *sqlx.DB, migrationsPath string) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("erro ao criar driver de migração: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("erro ao criar instância de migração: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("erro ao executar migrações: %w", err)
	}

	return nil
}