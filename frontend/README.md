# Frontend - Organiza Aqui

Interface web desenvolvida em Next.js para gerenciamento financeiro, tarefas, hábitos e conhecimento pessoal.

## Stack Tecnológica

### Core
- **Next.js 16.1.3** - Framework React com App Router
- **React 19.2.3** - Biblioteca de interface
- **TypeScript 5.9.3** - Tipagem estática

### Gerenciamento de Estado
- **Zustand 5.0.10** - Gerenciamento de estado global
- **TanStack Query 5.90.19** - Gerenciamento de estado do servidor e cache

### UI e Estilização
- **Tailwind CSS 4.1.18** - Framework CSS utilitário
- **shadcn/ui** - Componentes UI baseados em Radix UI
- **Radix UI** - Componentes acessíveis e sem estilo
- **Lucide React** - Biblioteca de ícones
- **next-themes** - Suporte a tema claro/escuro

### Formulários e Validação
- **React Hook Form 7.71.1** - Gerenciamento de formulários
- **Zod 4.3.5** - Validação de schemas
- **@hookform/resolvers** - Resolvers para integração

### Utilitários
- **Axios 1.13.2** - Cliente HTTP
- **date-fns 4.1.0** - Manipulação de datas
- **Recharts 3.6.0** - Gráficos e visualizações
- **@dnd-kit** - Drag and drop para ordenação
- **nuqs** - Gerenciamento de query strings
- **Sonner** - Sistema de notificações toast

## Estrutura do Projeto

```
frontend/
  ├── app/                      # App Router do Next.js
  │   ├── (dashboard)/         # Grupo de rotas do dashboard
  │   │   ├── calendar/        # Página de calendário
  │   │   ├── financial/      # Páginas financeiras
  │   │   ├── habits/          # Página de hábitos
  │   │   ├── notes/           # Página de notas
  │   │   ├── tasks/           # Página de tarefas
  │   │   ├── layout.tsx       # Layout do dashboard
  │   │   └── page.tsx         # Dashboard principal
  │   ├── login/               # Página de login
  │   ├── register/            # Página de registro
  │   ├── layout.tsx           # Layout raiz
  │   ├── page.tsx             # Página inicial
  │   └── globals.css          # Estilos globais
  ├── components/              # Componentes React
  │   ├── dashboard/           # Widgets do dashboard
  │   │   ├── FinancialWidget.tsx
  │   │   ├── TasksWidget.tsx
  │   │   └── TimelineWidget.tsx
  │   ├── financial/           # Componentes financeiros
  │   │   ├── AccountCard.tsx
  │   │   ├── TransactionForm.tsx
  │   │   ├── TransactionTable.tsx
  │   │   └── ...
  │   ├── knowledge/           # Componentes de conhecimento
  │   │   ├── CalendarEventForm.tsx
  │   │   ├── HabitForm.tsx
  │   │   └── NoteForm.tsx
  │   ├── layout/              # Componentes de layout
  │   │   ├── AuthGuard.tsx
  │   │   ├── Header.tsx
  │   │   ├── Sidebar.tsx
  │   │   └── ...
  │   ├── tasks/               # Componentes de tarefas
  │   │   ├── TaskCard.tsx
  │   │   ├── TaskForm.tsx
  │   │   └── TaskKanban.tsx
  │   └── ui/                  # Componentes UI base (shadcn)
  │       ├── button.tsx
  │       ├── card.tsx
  │       ├── dialog.tsx
  │       └── ...
  ├── hooks/                   # Custom hooks
  │   ├── mutations/           # Hooks de mutação (React Query)
  │   │   ├── useAuth.ts
  │   │   ├── useAccountMutations.ts
  │   │   └── ...
  │   ├── queries/             # Hooks de query (React Query)
  │   │   ├── useAccounts.ts
  │   │   ├── useTransactions.ts
  │   │   └── ...
  │   └── useToast.ts          # Hook de toast
  ├── services/                # Serviços de API
  │   ├── authService.ts
  │   ├── financialService.ts
  │   ├── taskService.ts
  │   └── ...
  ├── stores/                  # Stores Zustand
  │   ├── authStore.ts         # Estado de autenticação
  │   ├── privacyStore.ts      # Estado de privacidade
  │   └── toastStore.ts        # Estado de toasts
  ├── schemas/                 # Schemas Zod
  │   ├── authSchema.ts
  │   ├── financialSchema.ts
  │   ├── taskSchema.ts
  │   └── ...
  ├── types/                   # Tipos TypeScript
  │   ├── financial.ts
  │   ├── task.ts
  │   └── ...
  ├── lib/                     # Utilitários e configurações
  │   ├── axios.ts             # Cliente Axios configurado
  │   ├── utils.ts             # Funções utilitárias
  │   └── theme.ts             # Configuração de tema
  ├── utils/                   # Utilitários adicionais
  │   ├── currencies.ts
  │   └── passwordStrength.ts
  ├── providers/               # Providers React
  │   └── QueryProvider.tsx    # Provider do React Query
  ├── public/                  # Arquivos estáticos
  │   ├── logo-dark.png
  │   └── logo-light.png
  ├── components.json          # Configuração do shadcn/ui
  ├── next.config.ts           # Configuração do Next.js
  ├── tsconfig.json            # Configuração do TypeScript
  └── package.json             # Dependências do projeto
```

## Arquitetura

A aplicação segue uma arquitetura moderna baseada em:

### Camadas

1. **Presentation Layer (App Router)**
   - Rotas e layouts do Next.js
   - Páginas e grupos de rotas

2. **Component Layer**
   - Componentes reutilizáveis organizados por domínio
   - Componentes UI base (shadcn/ui)
   - Componentes de layout

3. **Business Logic Layer**
   - Custom hooks (queries e mutations)
   - Services (comunicação com API)
   - Stores (estado global)

4. **Data Layer**
   - React Query para cache e sincronização
   - Zustand para estado global
   - Axios para requisições HTTP

## Funcionalidades Principais

### Autenticação
- Login e registro de usuários
- Proteção de rotas com `AuthGuard`
- Gerenciamento de token JWT
- Interceptores Axios para autenticação automática

### Dashboard
- Visão geral financeira
- Widgets de tarefas
- Timeline de eventos
- Gráficos e análises

### Gestão Financeira
- **Contas**: Gerenciamento de contas bancárias
- **Transações**: CRUD completo de transações
- **Categorias**: Organização por categorias
- **Cartões de Crédito**: Gerenciamento de cartões e faturas
- **Parcelas**: Sistema de parcelamento
- **Importação**: Importação de extratos OFX e CSV
- **Analytics**: Relatórios e gráficos financeiros
- **Calendário de Vencimentos**: Visualização de contas a pagar

### Gestão de Tarefas
- **Kanban Board**: Visualização em colunas
- **Status Customizáveis**: Criação e ordenação de status
- **Drag and Drop**: Reordenação de tarefas
- **Filtros e Busca**: Localização rápida de tarefas

### Conhecimento Pessoal
- **Notas**: Sistema de anotações
- **Eventos de Calendário**: Gerenciamento de eventos
- **Hábitos**: Criação e acompanhamento
- **Tracking**: Registro de execução de hábitos

### UI/UX
- **Tema Claro/Escuro**: Suporte a múltiplos temas
- **Responsivo**: Design adaptável a diferentes telas
- **Acessibilidade**: Componentes baseados em Radix UI
- **Notificações**: Sistema de toast para feedback
- **Loading States**: Estados de carregamento
- **Error Handling**: Tratamento de erros

## Configuração

### Variáveis de Ambiente

Crie um arquivo `.env.local` na raiz do projeto:

```bash
# URL da API Backend
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Instalação

1. **Instalar dependências:**
```bash
npm install
# ou
yarn install
# ou
pnpm install
```

2. **Configurar variáveis de ambiente:**
```bash
# Criar arquivo .env.local com as variáveis acima
```

3. **Executar em desenvolvimento:**
```bash
npm run dev
# ou
yarn dev
# ou
pnpm dev
```

A aplicação estará disponível em `http://localhost:3000`.

## Scripts Disponíveis

### Desenvolvimento
```bash
# Iniciar servidor de desenvolvimento
npm run dev

# Build de produção
npm run build

# Iniciar servidor de produção
npm start
```

### Qualidade de Código
```bash
# Executar linter
npm run lint

# Verificar tipos TypeScript
npm run typecheck

# Formatar código
npm run format
```

## Estrutura de Rotas

### Rotas Públicas
- `/` - Página inicial
- `/login` - Login de usuário
- `/register` - Registro de usuário

### Rotas Protegidas (Dashboard)
- `/dashboard` - Dashboard principal
- `/financial` - Gestão financeira
- `/tasks` - Gestão de tarefas
- `/habits` - Gestão de hábitos
- `/notes` - Sistema de notas
- `/calendar` - Calendário de eventos

## Componentes Principais

### Layout
- **AuthGuard**: Proteção de rotas autenticadas
- **Sidebar**: Navegação lateral
- **Header**: Cabeçalho da aplicação
- **ThemeProvider**: Gerenciamento de tema

### Financial
- **TransactionTable**: Tabela de transações
- **TransactionForm**: Formulário de transações
- **AccountCard**: Card de conta
- **CreditCardList**: Lista de cartões
- **NetWorthCard**: Card de patrimônio líquido
- **SpendingByTagChart**: Gráfico de gastos por tag
- **UpcomingBillsCalendar**: Calendário de vencimentos

### Tasks
- **TaskKanban**: Board Kanban de tarefas
- **TaskCard**: Card de tarefa
- **TaskForm**: Formulário de tarefa

### Dashboard
- **FinancialWidget**: Widget financeiro
- **TasksWidget**: Widget de tarefas
- **TimelineWidget**: Widget de timeline

## Hooks Customizados

### Queries (React Query)
- `useAccounts` - Listar contas
- `useTransactions` - Listar transações
- `useTasks` - Listar tarefas
- `useCategories` - Listar categorias
- `useAnalytics` - Dados de analytics
- E muitos outros...

### Mutations (React Query)
- `useAuth` - Autenticação (login, register, logout)
- `useAccountMutations` - Mutations de contas
- `useTransactionMutations` - Mutations de transações
- `useTaskMutations` - Mutations de tarefas
- E muitos outros...

### Utilitários
- `useToast` - Hook para exibir notificações
- `useMobile` - Detecção de dispositivo móvel

## Stores (Zustand)

### AuthStore
Gerencia estado de autenticação:
- Usuário atual
- Token JWT
- Estado de autenticação
- Métodos: login, logout, setUser, setToken

### PrivacyStore
Gerencia preferências de privacidade

### ToastStore
Gerencia notificações toast

## Services

Serviços organizados por domínio que encapsulam chamadas à API:

- **authService**: Autenticação
- **financialService**: Operações financeiras
- **taskService**: Operações de tarefas
- **knowledgeService**: Notas, eventos, hábitos
- **bankService**: Operações com bancos
- **timelineService**: Timeline de eventos

## Validação de Formulários

A aplicação utiliza **Zod** para validação de schemas e **React Hook Form** para gerenciamento de formulários:

- Schemas definidos em `schemas/`
- Validação client-side
- Mensagens de erro personalizadas
- Integração com React Hook Form via `@hookform/resolvers`

## Estilização

### Tailwind CSS
- Framework CSS utilitário
- Configuração customizada
- Variáveis CSS para temas

### shadcn/ui
- Componentes base acessíveis
- Customizáveis via Tailwind
- Baseados em Radix UI

### Tema
- Suporte a tema claro/escuro
- Gerenciado via `next-themes`
- Persistência de preferência

## Integração com Backend

### Cliente Axios
- Configurado em `lib/axios.ts`
- Base URL configurável via variável de ambiente
- Interceptor para adicionar token JWT automaticamente
- Interceptor para tratamento de erros 401 (redirecionamento)

### Autenticação
- Token armazenado no localStorage
- Adicionado automaticamente em todas as requisições
- Remoção automática em caso de erro 401

## Performance

### React Query
- Cache automático de queries
- Refetch inteligente
- Stale time configurado (1 minuto)
- Retry configurado (1 tentativa)

### Next.js
- App Router para melhor performance
- Server Components quando possível
- Code splitting automático
- Otimização de imagens

## Desenvolvimento

### Convenções
- Componentes em PascalCase
- Hooks com prefixo `use`
- Services em camelCase
- Types/interfaces em PascalCase
- Arquivos de componentes: `ComponentName.tsx`

### Estrutura de Componentes
```typescript
// Exemplo de estrutura de componente
'use client' // Se necessário

import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'

export function ComponentName() {
  // Lógica do componente
  return <div>...</div>
}
```

### Estrutura de Hooks
```typescript
// Exemplo de hook de query
import { useQuery } from '@tanstack/react-query'
import { service } from '@/services'

export function useResource() {
  return useQuery({
    queryKey: ['resource'],
    queryFn: () => service.getResource(),
  })
}
```

## Deploy

### Build de Produção
```bash
npm run build
npm start
```

### Variáveis de Ambiente em Produção
Certifique-se de configurar `NEXT_PUBLIC_API_URL` com a URL do backend em produção.

### Vercel (Recomendado)
A aplicação pode ser facilmente deployada na Vercel:
1. Conecte seu repositório
2. Configure as variáveis de ambiente
3. Deploy automático a cada push

## Contribuindo

1. Crie uma branch para sua feature
2. Faça suas alterações
3. Execute o linter: `npm run lint`
4. Verifique tipos: `npm run typecheck`
5. Formate o código: `npm run format`
6. Faça commit e push

## Licença

Este projeto é parte do sistema Organiza Aqui.
