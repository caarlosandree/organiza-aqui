package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/model"
	"github.com/redis/go-redis/v9"
)

// AuthRepository define a interface para operações de autenticação
type AuthRepository interface {
	CreateSession(ctx context.Context, session *model.AuthSession) error
	FindSessionByTokenHash(ctx context.Context, tokenHash string) (*model.AuthSession, error)
	DeleteSession(ctx context.Context, id uuid.UUID) error
	DeleteExpiredSessions(ctx context.Context) error
	DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error
}

type authRepository struct {
	db    *sqlx.DB
	redis *redis.Client
}

// NewAuthRepository cria uma nova instância de AuthRepository
func NewAuthRepository(db *sqlx.DB, redisClient *redis.Client) AuthRepository {
	return &authRepository{
		db:    db,
		redis: redisClient,
	}
}

func (r *authRepository) CreateSession(ctx context.Context, session *model.AuthSession) error {
	// Serializar sessão para JSON
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("erro ao serializar sessão: %w", err)
	}

	// Calcular TTL baseado em expiresAt
	ttl := time.Until(session.ExpiresAt)
	if ttl <= 0 {
		return fmt.Errorf("sessão já expirada")
	}

	// Chave Redis: session:{tokenHash}
	key := fmt.Sprintf("session:%s", session.TokenHash)

	// Armazenar no Redis com TTL
	if err := r.redis.Set(ctx, key, sessionJSON, ttl).Err(); err != nil {
		return fmt.Errorf("erro ao criar sessão no Redis: %w", err)
	}

	return nil
}

func (r *authRepository) FindSessionByTokenHash(ctx context.Context, tokenHash string) (*model.AuthSession, error) {
	// Chave Redis: session:{tokenHash}
	key := fmt.Sprintf("session:%s", tokenHash)

	// Buscar no Redis
	sessionJSON, err := r.redis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			// Sessão não encontrada no Redis
			return nil, nil
		}
		return nil, fmt.Errorf("erro ao buscar sessão no Redis: %w", err)
	}

	// Deserializar sessão
	var session model.AuthSession
	if err := json.Unmarshal([]byte(sessionJSON), &session); err != nil {
		return nil, fmt.Errorf("erro ao deserializar sessão: %w", err)
	}

	// Verificar se ainda não expirou (TTL do Redis já faz isso, mas verificamos por segurança)
	if time.Now().After(session.ExpiresAt) {
		// Sessão expirada, remover do Redis
		r.redis.Del(ctx, key)
		return nil, nil
	}

	return &session, nil
}

func (r *authRepository) DeleteSession(ctx context.Context, id uuid.UUID) error {
	// Para deletar por ID, precisamos buscar todas as sessões do usuário
	// ou usar um padrão de busca. Como não temos tokenHash aqui, vamos usar
	// um padrão de busca por ID. Mas isso é ineficiente.
	// Alternativa: usar um índice adicional no Redis: session:id:{id} -> tokenHash
	// Por enquanto, vamos buscar no banco para encontrar o tokenHash e deletar no Redis
	
	// Buscar sessão no banco para obter tokenHash
	var tokenHash string
	query := `SELECT token_hash FROM auth_sessions WHERE id = $1`
	err := r.db.GetContext(ctx, &tokenHash, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			// Sessão não existe, considerar sucesso
			return nil
		}
		return fmt.Errorf("erro ao buscar sessão no banco: %w", err)
	}

	// Deletar do Redis
	key := fmt.Sprintf("session:%s", tokenHash)
	if err := r.redis.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("erro ao deletar sessão no Redis: %w", err)
	}

	// Também deletar do banco (para manter consistência durante transição)
	_, err = r.db.ExecContext(ctx, `DELETE FROM auth_sessions WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar sessão no banco: %w", err)
	}

	return nil
}

func (r *authRepository) DeleteExpiredSessions(ctx context.Context) error {
	// No Redis, as sessões expiradas são removidas automaticamente pelo TTL
	// Não precisamos fazer nada aqui, mas podemos limpar o banco se necessário
	query := `DELETE FROM auth_sessions WHERE expires_at < $1`
	_, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		return fmt.Errorf("erro ao deletar sessões expiradas no banco: %w", err)
	}
	return nil
}

func (r *authRepository) DeleteSessionsByUserID(ctx context.Context, userID uuid.UUID) error {
	// Buscar todas as sessões do usuário no banco para obter os tokenHash
	var sessions []struct {
		TokenHash string `db:"token_hash"`
	}
	query := `SELECT token_hash FROM auth_sessions WHERE user_id = $1`
	err := r.db.SelectContext(ctx, &sessions, query, userID)
	if err != nil {
		return fmt.Errorf("erro ao buscar sessões do usuário: %w", err)
	}

	// Deletar cada sessão do Redis
	for _, session := range sessions {
		key := fmt.Sprintf("session:%s", session.TokenHash)
		if err := r.redis.Del(ctx, key).Err(); err != nil {
			// Log erro mas continua com as outras
			continue
		}
	}

	// Deletar do banco também
	_, err = r.db.ExecContext(ctx, `DELETE FROM auth_sessions WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("erro ao deletar sessões do usuário no banco: %w", err)
	}

	return nil
}
