# Schema Design - Organiza Aqui

## Decisões de Design

### Tipos de Dados

#### IDs
- **UUID** como PRIMARY KEY em todas as tabelas
- Usar `github.com/google/uuid` para geração
- Facilita distribuição e evita problemas de sequência

#### Valores Monetários
- **BIGINT** para armazenar valores em centavos
- **NUNCA** usar FLOAT ou DECIMAL para valores monetários
- Exemplo: R$ 100,50 = 10050 centavos
- No Go: usar `int64` para representar centavos

#### Datas e Timestamps
- **TIMESTAMP** para created_at, updated_at, expires_at
- **DATE** para datas sem hora (due_date, date em transactions)
- **TIMESTAMP NULL** para campos opcionais (completed_at)

#### Strings
- **VARCHAR(255)** para nomes, títulos curtos
- **TEXT** para descrições e conteúdo longo
- **VARCHAR(50)** para códigos, tipos, status
- **VARCHAR(7)** para cores hexadecimais (#RRGGBB)

### Estruturas Hierárquicas

#### Categorias (Materialized Path)
- Campo `path` armazena o caminho hierárquico (ex: "1.2.5")
- Campo `parent_id` para referência direta ao pai
- Facilita queries de árvore sem recursão complexa
- Exemplo: "Alimentação > Restaurante > Fast Food" = path "1.5.12"

#### Lexorank para Tarefas
- Campo `lexorank VARCHAR(50)` para ordenação
- Permite inserção entre itens sem reordenar tudo
- Algoritmo: string ordenável que permite posicionamento entre dois valores
- Exemplo: entre "a" e "b" pode inserir "a5", entre "a5" e "b" pode inserir "a7"

### Índices

#### Índices Essenciais
- `user_id` em todas as tabelas (filtro mais comum)
- `(user_id, date)` em transactions (filtros por período)
- `(user_id, occurred_at)` em timeline_items (dashboard)
- `(user_id, status_id)` em tasks (kanban)
- `(user_id, parent_id)` em categories (hierarquia)

#### Índices Compostos
- Criar índices compostos para queries frequentes
- Exemplo: `idx_transactions_user_date` em `(user_id, date DESC)`

### Constraints

#### Foreign Keys
- Todas as referências devem ter FOREIGN KEY constraints
- ON DELETE CASCADE apenas onde faz sentido (ex: deletar user deleta tudo)
- ON DELETE RESTRICT para dados críticos (ex: não deletar categoria com transações)

#### Unique Constraints
- `(user_id, email)` em users (se multi-user)
- `(user_id, name)` em accounts (nomes únicos por usuário)
- `(user_id, lexorank)` em tasks (garantir unicidade do lexorank)

### Normalização

#### Nível 3NF
- Evitar redundância de dados
- Exemplo: não armazenar saldo calculado, calcular sob demanda ou via trigger

#### Denormalização Estratégica
- Campo `balance` em accounts (atualizado via trigger ou service)
- Campo `path` em categories (materialized path)
- Campo `metadata JSONB` em timeline_items (dados agregados para performance)

### Performance

#### Queries Otimizadas
- Usar JOINs ao invés de N+1 queries
- Usar projeções (SELECT específico) ao invés de SELECT *
- Paginação para listas grandes (LIMIT/OFFSET ou cursor-based)

#### Connection Pooling
- Configurar pool adequado no sqlx
- MaxIdleConns: 10
- MaxOpenConns: 100
- ConnMaxLifetime: 1 hora

### Segurança

#### Senhas
- Hash com bcrypt (cost 10-12)
- Nunca armazenar senha em texto plano
- Campo `password_hash` ao invés de `password`

#### Tokens JWT
- Armazenar hash do token em `auth_sessions`
- Campo `expires_at` para expiração
- Limpar sessões expiradas periodicamente

### Migrations

#### Versionamento
- Usar `golang-migrate` via CLI
- Nomenclatura: `{version}_{description}.{direction}.sql`
- Versões sequenciais: 000001, 000002, etc.
- Sempre criar `.up.sql` e `.down.sql`

#### Boas Práticas
- Migrations atômicas (uma mudança por migration)
- Testar migrations em desenvolvimento antes de produção
- Manter histórico de migrations no repositório
