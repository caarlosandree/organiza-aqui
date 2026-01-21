# ERD - Organiza Aqui

## Diagrama de Entidade-Relacionamento

### Core (Autenticação e Usuários)

```
users
├── id (UUID, PK)
├── email (VARCHAR(255), UNIQUE)
├── password_hash (VARCHAR(255))
├── name (VARCHAR(255))
├── created_at (TIMESTAMP)
└── updated_at (TIMESTAMP)

auth_sessions
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── token_hash (VARCHAR(255))
├── expires_at (TIMESTAMP)
└── created_at (TIMESTAMP)
```

### Financeiro (Vault)

```
accounts
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── name (VARCHAR(255))
├── type (VARCHAR(50)) -- 'checking', 'savings', 'credit', 'investment'
├── balance (BIGINT) -- em centavos
├── currency (VARCHAR(3)) -- 'BRL', 'USD', etc.
└── created_at (TIMESTAMP)

categories
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── name (VARCHAR(255))
├── parent_id (UUID, FK -> categories.id, NULL)
├── path (VARCHAR(255)) -- Materialized Path
├── type (VARCHAR(50)) -- 'income', 'expense'
├── color (VARCHAR(7)) -- hex color
└── created_at (TIMESTAMP)

transactions
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── account_id (UUID, FK -> accounts.id)
├── category_id (UUID, FK -> categories.id, NULL)
├── type (VARCHAR(50)) -- 'income', 'expense', 'transfer'
├── amount (BIGINT) -- em centavos
├── description (TEXT)
├── date (DATE)
└── created_at (TIMESTAMP)

recurrence_patterns
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── transaction_id (UUID, FK -> transactions.id)
├── frequency (VARCHAR(50)) -- 'daily', 'weekly', 'monthly', 'yearly'
├── interval (INTEGER) -- a cada X dias/semanas/meses
├── end_date (DATE, NULL) -- NULL = infinito
└── created_at (TIMESTAMP)
```

### Tarefas e Agenda (Planner)

```
task_statuses
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── name (VARCHAR(255))
├── order (INTEGER) -- ordem de exibição
├── color (VARCHAR(7)) -- hex color
└── created_at (TIMESTAMP)

tasks
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── status_id (UUID, FK -> task_statuses.id)
├── title (VARCHAR(255))
├── description (TEXT)
├── lexorank (VARCHAR(50)) -- para ordenação drag-and-drop
├── due_date (DATE, NULL)
├── completed_at (TIMESTAMP, NULL)
├── amount (BIGINT, NULL) -- valor monetário opcional
├── category_id (UUID, FK -> categories.id, NULL) -- categoria financeira opcional
└── created_at (TIMESTAMP)

calendar_events
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── title (VARCHAR(255))
├── description (TEXT)
├── start_time (TIMESTAMP)
├── end_time (TIMESTAMP)
├── all_day (BOOLEAN)
└── created_at (TIMESTAMP)
```

### Anotações (Knowledge)

```
notes
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── title (VARCHAR(255))
├── content (TEXT)
├── tags (TEXT[]) -- array de tags
├── created_at (TIMESTAMP)
└── updated_at (TIMESTAMP)

habits
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── name (VARCHAR(255))
├── frequency (VARCHAR(50)) -- 'daily', 'weekly', etc.
├── streak_count (INTEGER) -- contador de dias consecutivos
└── created_at (TIMESTAMP)
```

### Timeline (Agregação)

```
timeline_items
├── id (UUID, PK)
├── user_id (UUID, FK -> users.id)
├── item_type (VARCHAR(50)) -- 'transaction', 'task', 'event', 'note'
├── item_id (UUID) -- ID do item no módulo específico
├── occurred_at (TIMESTAMP) -- quando o evento ocorreu
├── title (VARCHAR(255)) -- título para exibição
├── metadata (JSONB) -- dados agregados para performance
└── created_at (TIMESTAMP)
```

## Relacionamentos

### Hierárquicos
- `categories.parent_id` -> `categories.id` (auto-referência)
- `tasks.category_id` -> `categories.id` (opcional, para integração financeiro)

### Transacionais
- `transactions.account_id` -> `accounts.id`
- `transactions.category_id` -> `categories.id`
- `tasks.status_id` -> `task_statuses.id`

### Timeline (Polimórfico)
- `timeline_items.item_id` referencia diferentes tabelas baseado em `item_type`
- Não há FK direta, mas há integridade lógica via service

## Índices Recomendados

### Core
- `idx_users_email` em `users(email)`
- `idx_auth_sessions_user` em `auth_sessions(user_id, expires_at)`

### Financeiro
- `idx_accounts_user` em `accounts(user_id)`
- `idx_categories_user_parent` em `categories(user_id, parent_id)`
- `idx_categories_path` em `categories(path)` -- para queries hierárquicas
- `idx_transactions_user_date` em `transactions(user_id, date DESC)`
- `idx_transactions_account` em `transactions(account_id)`
- `idx_transactions_category` em `transactions(category_id)`

### Tarefas
- `idx_task_statuses_user_order` em `task_statuses(user_id, order)`
- `idx_tasks_user_status` em `tasks(user_id, status_id)`
- `idx_tasks_lexorank` em `tasks(user_id, lexorank)`
- `idx_tasks_due_date` em `tasks(due_date)` -- para filtros de vencimento

### Timeline
- `idx_timeline_user_date` em `timeline_items(user_id, occurred_at DESC)`
- `idx_timeline_type` em `timeline_items(item_type, item_id)`

## Constraints

### Foreign Keys
- Todas as tabelas com `user_id` têm FK para `users.id`
- `ON DELETE CASCADE` para dados dependentes do usuário
- `ON DELETE RESTRICT` para categorias com transações

### Unique
- `users(email)` -- email único
- `accounts(user_id, name)` -- nome de conta único por usuário
- `task_statuses(user_id, name)` -- nome de status único por usuário

### Check
- `transactions.amount > 0` -- valor deve ser positivo
- `accounts.balance` pode ser negativo (débito)
- `task_statuses.order >= 0` -- ordem não negativa
