package service

import (
	"github.com/jmoiron/sqlx"
	"github.com/organiza-aqui/backend/internal/config"
	"github.com/organiza-aqui/backend/internal/repository"
	"github.com/redis/go-redis/v9"
)

// Services contém todos os services
type Services struct {
	Auth              AuthService
	Account           AccountService
	Bank              BankService
	Category          CategoryService
	Transaction       TransactionService
	TransactionPeriod TransactionPeriodService
	Recurrence        RecurrenceService
	Installment       InstallmentService
	CreditCard        CreditCardService
	Import            ImportService
	Analytics         AnalyticsService
	TaskStatus        TaskStatusService
	Task              TaskService
	Timeline          TimelineService
	CalendarEvent     CalendarEventService
	Note              NoteService
	Habit             HabitService
	HabitTracking     HabitTrackingService
}

// NewServices cria uma nova instância de Services
func NewServices(db *sqlx.DB, repositories *repository.Repositories, cfg *config.Config, redisClient *redis.Client) *Services {
	transactionPeriodService := NewTransactionPeriodService(
		repositories.TransactionPeriod,
		repositories.Transaction,
		repositories.Account,
	)
	installmentService := NewInstallmentService(db, repositories.Transaction, repositories.Account, repositories.Category)
	transactionService := NewTransactionService(
		db,
		repositories.Transaction,
		repositories.Account,
		repositories.Category,
		transactionPeriodService,
		installmentService,
	)
	
	return &Services{
		Auth: NewAuthService(
			repositories.User,
			repositories.Auth,
			redisClient,
			cfg.JWT.Secret,
			cfg.JWT.ExpirationHours,
		),
		Account:            NewAccountService(db, repositories.Account, repositories.Transaction),
		Bank:               NewBankService(repositories.Bank),
		Category:           NewCategoryService(repositories.Category),
		Transaction:        transactionService,
		TransactionPeriod:  transactionPeriodService,
		Recurrence:         NewRecurrenceService(repositories.Recurrence, repositories.Transaction, transactionService),
		Installment:        installmentService,
		CreditCard:         NewCreditCardService(repositories.CreditCard, repositories.CreditCardBill, repositories.Account, repositories.Transaction),
		Import:             NewImportService(db, repositories.Transaction, repositories.Account, repositories.CreditCard, transactionPeriodService),
		Analytics:          NewAnalyticsService(repositories.Transaction, repositories.Account, repositories.CreditCard, repositories.CreditCardBill),
		TaskStatus:         NewTaskStatusService(repositories.TaskStatus, repositories.Task),
		Task:               NewTaskService(db, repositories.Task, repositories.TaskStatus, repositories.Transaction, repositories.Account),
		Timeline:           NewTimelineService(repositories.Timeline, repositories.Transaction, repositories.Task),
		CalendarEvent:      NewCalendarEventService(repositories.CalendarEvent),
		Note:               NewNoteService(repositories.Note),
		Habit:              NewHabitService(repositories.Habit, repositories.HabitTracking),
		HabitTracking:      NewHabitTrackingService(repositories.HabitTracking, repositories.Habit),
	}
}
