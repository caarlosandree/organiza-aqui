# Backend - Organiza Aqui

API REST desenvolvida em Go usando Echo framework para gerenciamento financeiro, tarefas, hábitos e conhecimento pessoal.

## Stack Tecnológica

- **Go 1.25.5** - Linguagem de programação
- **Echo Framework v4** - Framework web HTTP
- **PostgreSQL 18** - Banco de dados relacional
- **Redis** - Cache e rate limiting
- **sqlx** - Extensão do database/sql
- **Viper** - Gerenciamento de configurações
- **Zap** - Logger estruturado de alta performance
- **golang-migrate** - Gerenciamento de migrations
- **Swagger** - Documentação interativa da API
- **JWT** - Autenticação baseada em tokens
- **Cron** - Agendamento de tarefas
- **BrasilAPI** - Integração com API de bancos brasileiros

## Estrutura do Projeto

```
backend/
  ├── cmd/server/              # Ponto de entrada da aplicação
  │   └── main.go              # Inicialização do servidor
  ├── internal/                # Código interno da aplicação
  │   ├── config/              # Configurações e conexões (DB, Redis)
  │   ├── handler/             # Handlers HTTP (controllers)
  │   │   ├── auth_handler.go
  │   │   ├── account_handler.go
  │   │   ├── transaction_handler.go
  │   │   ├── task_handler.go
  │   │   └── ...              # Outros handlers
  │   ├── service/             # Lógica de negócio
  │   │   ├── auth_service.go
  │   │   ├── transaction_service.go
  │   │   └── ...              # Outros serviços
  │   ├── repository/          # Acesso ao banco de dados
  │   │   ├── user_repository.go
  │   │   ├── transaction_repository.go
  │   │   └── ...              # Outros repositórios
  │   ├── dto/                 # Data Transfer Objects
  │   │   ├── auth.go
  │   │   ├── financial.go
  │   │   └── ...              # Outros DTOs
  │   ├── model/               # Modelos de dados
  │   │   ├── user.go
  │   │   ├── transaction.go
  │   │   └── ...              # Outros modelos
  │   ├── middleware/          # Middlewares HTTP
  │   │   ├── auth_middleware.go
  │   │   ├── rate_limit_middleware.go
  │   │   └── ...              # Outros middlewares
  │   ├── error/               # Erros customizados
  │   ├── validator/           # Validadores customizados
  │   └── util/                # Utilitários
  ├── pkg/                     # Código reutilizável
  │   ├── logger/              # Logger customizado (Zap)
  │   └── response/            # Helpers de resposta HTTP
  ├── migrations/              # Scripts de migration SQL
  ├── docs/                    # Documentação Swagger gerada
  ├── tests/                   # Testes de integração
  ├── Makefile                 # Comandos de build e desenvolvimento
  └── go.mod                   # Dependências do projeto
```

## Arquitetura

A aplicação segue o padrão de arquitetura em camadas:

1. **Handler Layer**: Recebe requisições HTTP, valida entrada e retorna respostas
2. **Service Layer**: Contém a lógica de negócio
3. **Repository Layer**: Abstrai o acesso ao banco de dados
4. **Model Layer**: Define as estruturas de dados

## Funcionalidades Principais

### Autenticação e Autorização
- Registro e login de usuários
- Autenticação JWT com expiração configurável
- Middleware de autenticação para rotas protegidas
- Gerenciamento de sessões

### Gestão Financeira
- **Contas**: Gerenciamento de contas bancárias e carteiras
- **Transações**: CRUD completo de transações financeiras
- **Categorias**: Organização de transações por categorias
- **Cartões de Crédito**: Gerenciamento de cartões, faturas e limites
- **Parcelas**: Sistema de parcelamento de transações
- **Recorrências**: Padrões de transações recorrentes
- **Períodos**: Agrupamento de transações por período
- **Importação**: Importação de extratos OFX e CSV
- **Analytics**: Relatórios e análises financeiras

### Gestão de Tarefas
- **Status de Tarefas**: Status customizáveis com ordenação (Lexorank)
- **Tarefas**: CRUD completo com suporte a reordenação
- **Timeline**: Linha do tempo de eventos

### Conhecimento Pessoal
- **Notas**: Sistema de anotações
- **Eventos de Calendário**: Gerenciamento de eventos
- **Hábitos**: Criação e acompanhamento de hábitos
- **Tracking de Hábitos**: Registro de execução de hábitos

### Integrações
- **BrasilAPI**: Sincronização automática de bancos brasileiros
- **Cron Jobs**: Sincronização semanal de bancos (domingos às 02:00)

## Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto com as seguintes variáveis:

```bash
# Banco de Dados
DB_HOST_DEV=localhost
DB_PORT_DEV=5432
DB_USER_DEV=postgres
DB_PASSWORD_DEV=sua_senha
DB_NAME_DEV=organiza_aqui

# Redis
REDIS_HOST_DEV=localhost
REDIS_PORT_DEV=6379
REDIS_DATABASE_DEV=0

# Servidor
SERVER_HOST=localhost
SERVER_PORT=8080

# JWT
JWT_SECRET=seu_secret_jwt_super_seguro
JWT_EXPIRATION_HOURS=24

# Rate Limiting
RATE_LIMIT_ENABLED=true
RATE_LIMIT_REQUESTS_PER_MINUTE=60

# Ambiente
ENVIRONMENT=development
```

### Setup Inicial

1. **Instalar dependências do Go:**
```bash
make deps
```

2. **Instalar ferramentas de desenvolvimento:**
```bash
make install-tools
```

3. **Configurar variáveis de ambiente:**
```bash
# Criar arquivo .env com as variáveis acima
```

4. **Executar migrations:**
```bash
make migrate-up
```

## Executando a Aplicação

### Desenvolvimento
```bash
make run
```

O servidor será iniciado em `http://localhost:8080` (ou na porta configurada).

### Build de Produção
```bash
make build
./bin/server
```

### Health Check

Verifique se o servidor está rodando:
```bash
curl http://localhost:8080/health
```

## Documentação da API

Após iniciar o servidor, a documentação Swagger interativa estará disponível em:
```
http://localhost:8080/swagger/index.html
```

### Gerar Documentação

Para atualizar a documentação Swagger após alterações:
```bash
make swagger
```

## Endpoints Principais

### Autenticação
- `POST /api/v1/auth/register` - Registrar novo usuário
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `GET /api/v1/auth/me` - Obter dados do usuário autenticado

### Contas
- `GET /api/v1/accounts` - Listar contas
- `POST /api/v1/accounts` - Criar conta
- `GET /api/v1/accounts/:id` - Obter conta
- `PUT /api/v1/accounts/:id` - Atualizar conta
- `PUT /api/v1/accounts/:id/initial-balance` - Atualizar saldo inicial
- `POST /api/v1/accounts/:id/recalculate-balance` - Recalcular saldo
- `DELETE /api/v1/accounts/:id` - Deletar conta

### Transações
- `GET /api/v1/transactions` - Listar transações
- `GET /api/v1/transactions/statement` - Obter extrato
- `POST /api/v1/transactions` - Criar transação
- `GET /api/v1/transactions/:id` - Obter transação
- `PUT /api/v1/transactions/:id` - Atualizar transação
- `DELETE /api/v1/transactions/:id` - Deletar transação
- `PATCH /api/v1/transactions/:id/status` - Atualizar status

### Tarefas
- `GET /api/v1/tasks` - Listar tarefas
- `POST /api/v1/tasks` - Criar tarefa
- `GET /api/v1/tasks/:id` - Obter tarefa
- `PUT /api/v1/tasks/:id` - Atualizar tarefa
- `DELETE /api/v1/tasks/:id` - Deletar tarefa
- `POST /api/v1/tasks/reorder` - Reordenar tarefas
- `POST /api/v1/tasks/:id/complete` - Completar tarefa
- `POST /api/v1/tasks/:id/uncomplete` - Desmarcar tarefa como completa

### Analytics
- `GET /api/v1/analytics/income-expense` - Receitas e despesas por período
- `GET /api/v1/analytics/category-breakdown` - Breakdown por categoria
- `GET /api/v1/analytics/monthly-trend` - Tendência mensal
- `GET /api/v1/analytics/patrimonio-liquido` - Patrimônio líquido
- `GET /api/v1/analytics/calendario-vencimentos` - Calendário de vencimentos
- `GET /api/v1/analytics/gastos-por-tag` - Gastos por tag

*Para ver todos os endpoints, consulte a documentação Swagger.*

## Middlewares

A aplicação utiliza os seguintes middlewares (na ordem de execução):

1. **RequestIDMiddleware**: Adiciona ID único a cada requisição
2. **StructuredLoggerMiddleware**: Logging estruturado de requisições
3. **RecoverMiddleware**: Recuperação de panics
4. **CORS**: Configuração de CORS
5. **RateLimitMiddleware**: Limitação de taxa de requisições (usando Redis)
6. **AuthMiddleware**: Validação de tokens JWT (apenas em rotas protegidas)

## Comandos Úteis

### Desenvolvimento
```bash
# Executar aplicação
make run

# Build da aplicação
make build

# Instalar dependências
make deps

# Instalar ferramentas de desenvolvimento
make install-tools
```

### Testes
```bash
# Executar todos os testes
make test

# Executar testes com cobertura
make test-coverage
```

### Migrations
```bash
# Criar nova migration
make migrate-create

# Executar migrations
make migrate-up

# Reverter migrations
make migrate-down
```

### Qualidade de Código
```bash
# Executar linter
make lint
```

### Documentação
```bash
# Gerar documentação Swagger
make swagger
```

### Limpeza
```bash
# Remover arquivos gerados
make clean
```

## Características Técnicas

### Banco de Dados
- Connection pooling configurado (máx. 100 conexões, 10 idle)
- Migrations automáticas na inicialização
- Suporte a transações

### Cache e Performance
- Redis para rate limiting
- Connection pooling otimizado

### Segurança
- Autenticação JWT
- Rate limiting configurável
- Validação de entrada
- CORS configurado

### Logging
- Logging estruturado com Zap
- Diferentes níveis por ambiente
- Request ID para rastreamento

### Agendamento
- Cron job para sincronização semanal de bancos
- Execução automática na inicialização

## Desenvolvimento

### Estrutura de Código
- Separação clara de responsabilidades
- Interfaces para desacoplamento
- DTOs para transferência de dados
- Validação de entrada

### Convenções
- Handlers recebem services via injeção de dependência
- Services contêm lógica de negócio
- Repositories abstraem acesso a dados
- Erros customizados para melhor tratamento

## Contribuindo

1. Crie uma branch para sua feature
2. Faça suas alterações
3. Execute os testes: `make test`
4. Execute o linter: `make lint`
5. Atualize a documentação Swagger se necessário: `make swagger`
6. Faça commit e push

## Licença

Este projeto é parte do sistema Organiza Aqui.