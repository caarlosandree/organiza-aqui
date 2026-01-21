package error

import "errors"

// Erros genéricos da aplicação
var (
	// ErrNotFound é retornado quando um recurso não é encontrado
	ErrNotFound = errors.New("recurso não encontrado")

	// ErrInvalidInput é retornado quando dados de entrada são inválidos
	ErrInvalidInput = errors.New("dados de entrada inválidos")

	// ErrUnauthorized é retornado quando o usuário não está autorizado
	ErrUnauthorized = errors.New("não autorizado")

	// ErrForbidden é retornado quando o acesso é negado
	ErrForbidden = errors.New("acesso negado")

	// ErrInternalServer é retornado para erros internos do servidor
	ErrInternalServer = errors.New("erro interno do servidor")

	// ErrConflict é retornado quando há conflito (ex: recurso já existe)
	ErrConflict = errors.New("conflito")
)

// Erros específicos de entidades
var (
	// ErrTaskNotFound é retornado quando uma tarefa não é encontrada
	ErrTaskNotFound = errors.New("tarefa não encontrada")

	// ErrAccountNotFound é retornado quando uma conta não é encontrada
	ErrAccountNotFound = errors.New("conta não encontrada")

	// ErrTransactionNotFound é retornado quando uma transação não é encontrada
	ErrTransactionNotFound = errors.New("transação não encontrada")

	// ErrCategoryNotFound é retornado quando uma categoria não é encontrada
	ErrCategoryNotFound = errors.New("categoria não encontrada")

	// ErrNoteNotFound é retornado quando uma anotação não é encontrada
	ErrNoteNotFound = errors.New("anotação não encontrada")

	// ErrHabitNotFound é retornado quando um hábito não é encontrado
	ErrHabitNotFound = errors.New("hábito não encontrado")

	// ErrCalendarEventNotFound é retornado quando um evento de calendário não é encontrado
	ErrCalendarEventNotFound = errors.New("evento não encontrado")

	// ErrStatusNotFound é retornado quando um status não é encontrado
	ErrStatusNotFound = errors.New("status não encontrado")

	// ErrRecurrencePatternNotFound é retornado quando um padrão de recorrência não é encontrado
	ErrRecurrencePatternNotFound = errors.New("padrão não encontrado")

	// ErrCreditCardNotFound é retornado quando um cartão de crédito não é encontrado
	ErrCreditCardNotFound = errors.New("cartão não encontrado")

	// ErrBillNotFound é retornado quando uma fatura não é encontrada
	ErrBillNotFound = errors.New("fatura não encontrada")

	// ErrTrackingNotFound é retornado quando um tracking não é encontrado
	ErrTrackingNotFound = errors.New("tracking não encontrado")

	// ErrUnauthorizedAccess é retornado quando um recurso não pertence ao usuário
	ErrUnauthorizedAccess = errors.New("recurso não pertence ao usuário")
)

// Erros de autenticação e autorização
var (
	// ErrEmailAlreadyExists é retornado quando um email já está cadastrado
	ErrEmailAlreadyExists = errors.New("email já cadastrado")

	// ErrUsernameAlreadyExists é retornado quando um username já está cadastrado
	ErrUsernameAlreadyExists = errors.New("username já cadastrado")

	// ErrInvalidCredentials é retornado quando as credenciais são inválidas
	ErrInvalidCredentials = errors.New("credenciais inválidas")
)

// Erros de negócio
var (
	// ErrAlreadyClosed é retornado quando uma fatura já está fechada
	ErrAlreadyClosed = errors.New("fatura já está fechada ou paga")

	// ErrAlreadyPaid é retornado quando uma fatura já está paga
	ErrAlreadyPaid = errors.New("fatura já está paga")

	// ErrNotClosed é retornado quando uma operação requer que a fatura esteja fechada
	ErrNotClosed = errors.New("fatura deve estar fechada para ser paga")

	// ErrHasSubcategories é retornado quando uma categoria tem subcategorias e não pode ser deletada
	ErrHasSubcategories = errors.New("não é possível deletar categoria com subcategorias")

	// ErrReferenceDateFuture é retornado quando a data de referência é futura
	ErrReferenceDateFuture = errors.New("data de referência não pode ser futura")

	// ErrParentTransactionNotFound é retornado quando uma transação mãe não é encontrada
	ErrParentTransactionNotFound = errors.New("transação mãe não encontrada")
)