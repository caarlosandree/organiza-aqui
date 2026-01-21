-- Create timeline_events table (view unificada de eventos)
-- Esta tabela não armazena dados, mas serve como índice para agregação
-- Os dados reais vêm de transactions, tasks, calendar_events, etc.

CREATE TABLE IF NOT EXISTS timeline_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    entity_type VARCHAR(50) NOT NULL, -- 'transaction', 'task', 'calendar_event', 'note'
    entity_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    event_date TIMESTAMP NOT NULL,
    metadata JSONB, -- Dados adicionais específicos do tipo de entidade
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_timeline_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT check_entity_type CHECK (entity_type IN ('transaction', 'task', 'calendar_event', 'note'))
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_timeline_user ON timeline_events(user_id);
CREATE INDEX IF NOT EXISTS idx_timeline_date ON timeline_events(user_id, event_date DESC);
CREATE INDEX IF NOT EXISTS idx_timeline_entity ON timeline_events(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_timeline_type_date ON timeline_events(user_id, entity_type, event_date DESC);
