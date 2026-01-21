<div align="center">

# 🗂️ Organiza Aqui

**Sistema completo de organização pessoal que integra gestão financeira, tarefas, hábitos e conhecimento em uma única plataforma moderna e intuitiva.**

[![Status](https://img.shields.io/badge/status-em%20desenvolvimento-yellow)](https://github.com)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg)](CONTRIBUTING.md)
[![Issues](https://img.shields.io/badge/issues-welcome-brightgreen)](../../issues)
[![Pull Requests](https://img.shields.io/badge/PRs-welcome-brightgreen)](../../pulls)

### 🛠️ Stack Principal

[![Go Version](https://img.shields.io/badge/go-1.25+-00ADD8?logo=go&logoColor=white)](https://golang.org)
[![Node Version](https://img.shields.io/badge/node-18+-339933?logo=node.js&logoColor=white)](https://nodejs.org)
[![Next.js](https://img.shields.io/badge/Next.js-16.1-black?logo=next.js&logoColor=white)](https://nextjs.org)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react&logoColor=black)](https://react.dev)
[![TypeScript](https://img.shields.io/badge/TypeScript-5.9-3178C6?logo=typescript&logoColor=white)](https://www.typescriptlang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-18-336791?logo=postgresql&logoColor=white)](https://www.postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-DC382D?logo=redis&logoColor=white)](https://redis.io)

### 📦 Tecnologias

![Echo](https://img.shields.io/badge/Echo-4.15-00ADD8?logo=go&logoColor=white)
![Tailwind CSS](https://img.shields.io/badge/Tailwind-4.1-38B2AC?logo=tailwind-css&logoColor=white)
![shadcn/ui](https://img.shields.io/badge/shadcn%2Fui-latest-000000?logo=react&logoColor=white)
![TanStack Query](https://img.shields.io/badge/TanStack%20Query-5.90-FF4154?logo=react-query&logoColor=white)
![Zustand](https://img.shields.io/badge/Zustand-5.0-443F48?logo=react&logoColor=white)

[Funcionalidades](#-funcionalidades-atuais) • [Roadmap](#-funcionalidades-futuras) • [Documentação](#-documentação) • [Desenvolvimento](#️-desenvolvimento)

</div>

---

## 📑 Índice

- [📌 Sobre o Projeto](#-sobre-o-projeto)
- [🎯 Funcionalidades Atuais](#-funcionalidades-atuais)
  - [💰 Gestão Financeira](#-gestão-financeira-vault)
  - [✅ Gestão de Tarefas](#-gestão-de-tarefas-planner)
  - [📚 Conhecimento Pessoal](#-conhecimento-pessoal-knowledge-base)
  - [📊 Dashboard](#-dashboard)
  - [🔐 Autenticação e Segurança](#-autenticação-e-segurança)
- [🚀 Funcionalidades Futuras](#-funcionalidades-futuras)
- [📈 Plano de Implantação](#-plano-de-implantação)
- [⚡ Desempenho](#-desempenho)
- [🏗️ Arquitetura](#️-arquitetura)
- [📚 Documentação](#-documentação)
- [🛠️ Desenvolvimento](#️-desenvolvimento)
- [📊 Status do Projeto](#-status-do-projeto)
- [🤝 Contribuindo](#-contribuindo)
- [📝 Licença](#-licença)

---

## 📌 Sobre o Projeto

O **Organiza Aqui** é uma aplicação web desenvolvida para ajudar pessoas a organizarem suas vidas de forma eficiente, oferecendo ferramentas integradas para controle financeiro, gerenciamento de tarefas, acompanhamento de hábitos e armazenamento de conhecimento pessoal.

O sistema foi projetado com foco em **usabilidade**, **performance** e **escalabilidade**, proporcionando uma experiência fluida e responsiva em diferentes dispositivos.

### 🎯 Objetivos

- ✅ **Centralizar** todas as informações pessoais em um único lugar
- ✅ **Simplificar** o controle financeiro e planejamento
- ✅ **Facilitar** o gerenciamento de tarefas e projetos
- ✅ **Acompanhar** hábitos e metas pessoais
- ✅ **Armazenar** e organizar conhecimento pessoal

### ✨ Destaques

- 🚀 **Performance**: Otimizado para respostas rápidas e experiência fluida
- 🔒 **Segurança**: Autenticação JWT e validação de dados
- 📱 **Responsivo**: Interface adaptável a diferentes dispositivos
- 🎨 **Moderno**: UI/UX moderna e intuitiva
- 📊 **Analytics**: Relatórios e análises detalhadas
- 🔄 **Integrações**: Sincronização com APIs externas


## 🎯 Funcionalidades Atuais

### 💰 Gestão Financeira (Vault)

#### Contas e Transações
- **Gerenciamento de Contas**: Criação e controle de múltiplas contas bancárias, carteiras e investimentos
- **Transações**: Registro completo de receitas, despesas e transferências
- **Categorização**: Sistema hierárquico de categorias para organização financeira
- **Saldo Automático**: Cálculo automático de saldos com suporte a saldo inicial
- **Histórico Completo**: Visualização de todas as transações com filtros avançados

#### Cartões de Crédito
- **Gerenciamento de Cartões**: Cadastro e controle de múltiplos cartões de crédito
- **Faturas**: Controle de faturas com fechamento e pagamento
- **Limites**: Acompanhamento de limite disponível e utilizado
- **Projeções**: Projeção de faturas futuras baseadas em transações

#### Parcelamento e Recorrências
- **Parcelas**: Sistema completo de parcelamento de transações
- **Recorrências**: Padrões de transações recorrentes (diárias, semanais, mensais, anuais)
- **Geração Automática**: Geração automática de transações baseadas em padrões
- **Cancelamento Inteligente**: Cancelamento de parcelas futuras sem afetar histórico

#### Importação de Dados
- **OFX**: Importação de extratos bancários no formato OFX
- **CSV**: Importação de transações via arquivo CSV
- **Preview**: Visualização prévia antes de confirmar importação
- **Validação**: Validação automática de dados importados

#### Analytics e Relatórios
- **Receitas e Despesas**: Análise de receitas e despesas por período
- **Breakdown por Categoria**: Visualização de gastos por categoria
- **Tendências Mensais**: Gráficos de tendência ao longo do tempo
- **Patrimônio Líquido**: Cálculo e visualização do patrimônio líquido
- **Calendário de Vencimentos**: Visualização de contas a pagar e receber
- **Gastos por Tag**: Análise de gastos agrupados por tags

#### Integrações
- **BrasilAPI**: Sincronização automática de bancos brasileiros
- **Sincronização Automática**: Atualização semanal de dados bancários via cron job

### ✅ Gestão de Tarefas (Planner)

#### Sistema Kanban
- **Status Customizáveis**: Criação e personalização de colunas de status
- **Drag and Drop**: Reordenação intuitiva de tarefas entre status
- **Ordenação Flexível**: Sistema Lexorank para ordenação eficiente
- **Visualização Clara**: Interface Kanban moderna e responsiva

#### Funcionalidades de Tarefas
- **CRUD Completo**: Criação, edição, visualização e exclusão de tarefas
- **Reordenação**: Reordenação de tarefas dentro do mesmo status
- **Completar/Desmarcar**: Marcação de tarefas como concluídas
- **Filtros e Busca**: Localização rápida de tarefas

### 📚 Conhecimento Pessoal (Knowledge Base)

#### Notas
- **Sistema de Anotações**: Criação e organização de notas pessoais
- **CRUD Completo**: Gerenciamento completo de notas
- **Busca**: Localização rápida de conteúdo

#### Eventos de Calendário
- **Agenda**: Criação e gerenciamento de eventos no calendário
- **Visualização**: Visualização de eventos por data
- **CRUD Completo**: Gerenciamento completo de eventos

#### Hábitos
- **Criação de Hábitos**: Definição de hábitos a serem acompanhados
- **Tracking**: Registro diário de execução de hábitos
- **Estatísticas**: Visualização de estatísticas de acompanhamento
- **Histórico**: Histórico completo de execuções

### 📊 Dashboard

- **Visão Geral**: Dashboard centralizado com widgets informativos
- **Widget Financeiro**: Resumo financeiro rápido
- **Widget de Tarefas**: Tarefas pendentes e em andamento
- **Timeline**: Linha do tempo de eventos importantes
- **Gráficos**: Visualizações gráficas de dados financeiros

### 🔐 Autenticação e Segurança

- **Registro e Login**: Sistema completo de autenticação
- **JWT**: Autenticação baseada em tokens JWT
- **Proteção de Rotas**: Middleware de autenticação para rotas protegidas
- **Sessões**: Gerenciamento de sessões de usuário
- **Segurança**: Validação e sanitização de dados

## 🚀 Funcionalidades Futuras

### Planejadas para Próximas Versões

#### Gestão Financeira
- [ ] **Orçamentos**: Sistema de orçamentos mensais e anuais
- [ ] **Metas Financeiras**: Definição e acompanhamento de metas
- [ ] **Relatórios Avançados**: Relatórios personalizáveis e exportáveis
- [ ] **Multi-moeda**: Suporte a múltiplas moedas
- [ ] **Investimentos**: Rastreamento de investimentos e carteiras
- [ ] **Débitos Automáticos**: Controle de débitos automáticos
- [ ] **Conciliação Bancária**: Ferramenta de conciliação automática

#### Gestão de Tarefas
- [ ] **Subtarefas**: Sistema de subtarefas e dependências
- [ ] **Etiquetas**: Sistema de etiquetas para organização
- [ ] **Prioridades**: Sistema de priorização de tarefas
- [ ] **Prazos e Lembretes**: Notificações de prazos
- [ ] **Templates**: Templates de tarefas reutilizáveis
- [ ] **Colaboração**: Compartilhamento de tarefas entre usuários

#### Conhecimento Pessoal
- [ ] **Wikis Pessoais**: Sistema de wikis para conhecimento estruturado
- [ ] **Tags e Categorias**: Sistema avançado de organização
- [ ] **Busca Full-Text**: Busca avançada em todo o conteúdo
- [ ] **Exportação**: Exportação de notas em diferentes formatos
- [ ] **Templates de Notas**: Templates para diferentes tipos de notas

#### Integrações
- [ ] **APIs de Bancos**: Integração direta com APIs bancárias
- [ ] **Notificações Push**: Notificações em tempo real
- [ ] **Email**: Envio de relatórios e lembretes por email
- [ ] **Calendários Externos**: Sincronização com Google Calendar, Outlook
- [ ] **Backup Automático**: Sistema de backup automático na nuvem

#### Mobile
- [ ] **Aplicativo Mobile**: Aplicativo nativo iOS e Android
- [ ] **Notificações Mobile**: Notificações push no mobile
- [ ] **Modo Offline**: Funcionalidade offline para uso sem internet

#### Colaboração
- [ ] **Compartilhamento**: Compartilhamento de dados entre usuários
- [ ] **Família/Grupos**: Contas compartilhadas para famílias ou grupos
- [ ] **Permissões**: Sistema de permissões granulares

#### Analytics Avançados
- [ ] **IA e Machine Learning**: Previsões e insights inteligentes
- [ ] **Análise de Padrões**: Identificação automática de padrões de gastos
- [ ] **Recomendações**: Recomendações personalizadas baseadas em dados

## 📈 Plano de Implantação

### Fase 1: MVP (Concluída) ✅
- Sistema básico de autenticação
- Gestão financeira essencial (contas, transações, categorias)
- Gestão básica de tarefas
- Dashboard inicial
- Interface responsiva

### Fase 2: Funcionalidades Core (Concluída) ✅
- Sistema completo de cartões de crédito
- Parcelamento e recorrências
- Importação de dados (OFX, CSV)
- Analytics básicos
- Sistema de hábitos e notas
- Eventos de calendário

### Fase 3: Melhorias e Otimizações (Em Andamento) 🔄
- Otimizações de performance
- Melhorias na UX/UI
- Expansão de analytics
- Refinamento de funcionalidades existentes

### Fase 4: Funcionalidades Avançadas (Planejada) 📅
- Orçamentos e metas
- Investimentos
- Multi-moeda
- Relatórios avançados
- Templates e automações

### Fase 5: Mobile e Integrações (Planejada) 📅
- Aplicativo mobile nativo
- Integrações com APIs bancárias
- Sincronização com calendários externos
- Notificações push

### Fase 6: Colaboração e IA (Futuro) 🔮
- Compartilhamento e colaboração
- IA para insights e previsões
- Análise de padrões avançada
- Recomendações personalizadas

## ⚡ Desempenho

### Otimizações Implementadas

#### Backend
- **Connection Pooling**: Pool de conexões otimizado (máx. 100 conexões)
- **Query Optimization**: Queries otimizadas com índices estratégicos
- **Caching**: Cache com Redis para rate limiting e dados frequentes
- **Migrations Automáticas**: Migrations executadas automaticamente na inicialização
- **Graceful Shutdown**: Encerramento gracioso do servidor

#### Frontend
- **React Query**: Cache inteligente de dados do servidor
- **Code Splitting**: Divisão automática de código por rota
- **Lazy Loading**: Carregamento sob demanda de componentes
- **Otimização de Imagens**: Otimização automática de imagens via Next.js
- **Server Components**: Uso de Server Components quando possível

### Métricas de Performance

#### Tempo de Resposta
- **API**: Respostas médias < 200ms para operações CRUD
- **Frontend**: First Contentful Paint < 1.5s
- **Time to Interactive**: < 3s

#### Escalabilidade
- **Concorrência**: Suporte a múltiplas requisições simultâneas
- **Banco de Dados**: Índices otimizados para queries frequentes
- **Cache**: Estratégia de cache para reduzir carga no banco

#### Monitoramento
- **Logging Estruturado**: Logs estruturados para análise
- **Request ID**: Rastreamento de requisições via Request ID
- **Error Tracking**: Tratamento e logging de erros

### Melhorias Contínuas

- Monitoramento de performance em produção
- Otimização contínua de queries
- Análise de gargalos
- Implementação de melhorias baseadas em métricas reais

## 🏗️ Arquitetura

O projeto segue uma arquitetura moderna e escalável:

- **Backend**: API REST com arquitetura em camadas (Handler → Service → Repository)
- **Frontend**: Aplicação Next.js com App Router e Server Components
- **Banco de Dados**: PostgreSQL com migrations versionadas
- **Cache**: Redis para rate limiting e cache
- **Autenticação**: JWT com tokens seguros

### 🛠️ Stack Tecnológica

<details>
<summary><b>Clique para ver a stack completa</b></summary>

#### Backend
- **Go 1.25+** - Linguagem de programação
- **Echo Framework** - Framework web HTTP
- **PostgreSQL 18** - Banco de dados relacional
- **Redis** - Cache e rate limiting
- **JWT** - Autenticação baseada em tokens
- **Swagger** - Documentação interativa da API

#### Frontend
- **Next.js 16** - Framework React com App Router
- **React 19** - Biblioteca de interface
- **TypeScript** - Tipagem estática
- **Tailwind CSS** - Framework CSS utilitário
- **shadcn/ui** - Componentes UI
- **TanStack Query** - Gerenciamento de estado do servidor
- **Zustand** - Gerenciamento de estado global

Para mais detalhes, consulte os [READMEs individuais](#-documentação).

</details>

## 📚 Documentação

Para informações detalhadas sobre cada parte do projeto, consulte:

| Documentação | Descrição |
|-------------|-----------|
| [📘 Backend README](./backend/README.md) | Documentação completa do backend, incluindo stack tecnológica, estrutura, endpoints e configuração |
| [📗 Frontend README](./frontend/README.md) | Documentação completa do frontend, incluindo componentes, hooks, services e configuração |
| [📊 Schema Design](./docs/schema-design.md) | Decisões de design do banco de dados e estrutura |
| [🗄️ ERD](./docs/erd.md) | Diagrama de Entidade-Relacionamento |

## 🛠️ Desenvolvimento

### Pré-requisitos

Antes de começar, certifique-se de ter instalado:

- [Node.js](https://nodejs.org/) 18+ (para frontend)
- [Go](https://golang.org/) 1.25+ (para backend)
- [PostgreSQL](https://www.postgresql.org/) 18+
- [Redis](https://redis.io/)

### 🚀 Início Rápido

1. **Clone o repositório**
```bash
git clone <repository-url>
cd organiza-aqui
```

2. **Configure o Backend**
```bash
cd backend
# Configure o arquivo .env com as variáveis necessárias
make deps          # Instalar dependências
make migrate-up    # Executar migrations
make run           # Iniciar servidor
```

3. **Configure o Frontend**
```bash
cd frontend
# Configure o arquivo .env.local com NEXT_PUBLIC_API_URL
npm install        # Instalar dependências
npm run dev        # Iniciar servidor de desenvolvimento
```

4. **Acesse a aplicação**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- Swagger Docs: http://localhost:8080/swagger/index.html

> 📖 Para instruções detalhadas, consulte os [READMEs individuais](#-documentação) de cada parte do projeto.

## 📊 Status do Projeto

| Fase | Status | Progresso |
|------|--------|-----------|
| MVP | ✅ Concluída | 100% |
| Funcionalidades Core | ✅ Concluída | 100% |
| Melhorias e Otimizações | 🔄 Em Andamento | ~60% |
| Funcionalidades Avançadas | 📅 Planejada | 0% |
| Mobile e Integrações | 📅 Planejada | 0% |
| Colaboração e IA | 🔮 Futuro | 0% |

### 📈 Estatísticas

![GitHub repo size](https://img.shields.io/github/repo-size/organiza-aqui/organiza-aqui)
![GitHub language count](https://img.shields.io/github/languages/count/organiza-aqui/organiza-aqui)
![GitHub top language](https://img.shields.io/github/languages/top/organiza-aqui/organiza-aqui)
![GitHub last commit](https://img.shields.io/github/last-commit/organiza-aqui/organiza-aqui)

## 🤝 Contribuindo

Contribuições são muito bem-vindas! Este projeto segue algumas diretrizes:

1. **Fork o projeto**
2. **Crie uma branch** para sua feature (`git checkout -b feature/AmazingFeature`)
3. **Commit suas mudanças** (`git commit -m 'Add some AmazingFeature'`)
4. **Push para a branch** (`git push origin feature/AmazingFeature`)
5. **Abra um Pull Request**

### 📋 Antes de Contribuir

- Leia os [READMEs individuais](#-documentação) para entender a estrutura
- Siga os padrões de código estabelecidos
- Adicione testes quando apropriado
- Atualize a documentação se necessário

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 📧 Contato

Para dúvidas, sugestões ou problemas:

- 📧 Abra uma [Issue](../../issues) no GitHub
- 💬 Entre em contato através do repositório

---

<div align="center">

**Organiza Aqui** - Organize sua vida de forma eficiente e inteligente.

Feito com ❤️ para ajudar você a se organizar melhor

[⬆ Voltar ao topo](#-organiza-aqui)

</div>

---

## 🏷️ Tags e Tópicos

Este projeto utiliza as seguintes tags e tópicos:

**Categoria**: `organização-pessoal` `produtividade` `gestão-pessoal` `life-management`

**Funcionalidades**: `gestão-financeira` `tarefas` `kanban` `hábitos` `notas` `calendário` `dashboard` `analytics`

**Tecnologias**: `nextjs` `react` `typescript` `go` `golang` `postgresql` `redis` `api-rest` `jwt` `swagger` `tailwindcss` `shadcn-ui`

**Padrões**: `arquitetura-em-camadas` `rest-api` `jwt-auth` `migrations` `docker-ready`

---

<div align="center">

### ⭐ Se este projeto foi útil para você, considere dar uma estrela!

[![GitHub stars](https://img.shields.io/github/stars/organiza-aqui/organiza-aqui.svg?style=social&label=Star)](../../stargazers)
[![GitHub forks](https://img.shields.io/github/forks/organiza-aqui/organiza-aqui.svg?style=social&label=Fork)](../../network/members)
[![GitHub watchers](https://img.shields.io/github/watchers/organiza-aqui/organiza-aqui.svg?style=social&label=Watch)](../../watchers)

</div>
