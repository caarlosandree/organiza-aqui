package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

// Repositories contém todos os repositórios
type Repositories struct {
	User                UserRepository
	Auth                AuthRepository
	Account             AccountRepository
	Bank                BankRepository
	Category            CategoryRepository
	Transaction         TransactionRepository
	TransactionPeriod   TransactionPeriodRepository
	Recurrence          RecurrenceRepository
	CreditCard          CreditCardRepository
	CreditCardBill      CreditCardBillRepository
	TaskStatus          TaskStatusRepository
	Task                TaskRepository
	Timeline            TimelineRepository
	CalendarEvent       CalendarEventRepository
	Note                NoteRepository
	Habit               HabitRepository
	HabitTracking       HabitTrackingRepository
}

// NewRepositories cria uma nova instância de Repositories
func NewRepositories(db *sqlx.DB, redisClient *redis.Client) *Repositories {
	return &Repositories{
		User:              NewUserRepository(db, redisClient),
		Auth:              NewAuthRepository(db, redisClient),
		Account:           NewAccountRepository(db),
		Bank:              NewBankRepository(db),
		Category:          NewCategoryRepository(db),
		Transaction:       NewTransactionRepository(db),
		TransactionPeriod: NewTransactionPeriodRepository(db),
		Recurrence:        NewRecurrenceRepository(db),
		CreditCard:        NewCreditCardRepository(db),
		CreditCardBill:    NewCreditCardBillRepository(db),
		TaskStatus:        NewTaskStatusRepository(db),
		Task:              NewTaskRepository(db),
		Timeline:          NewTimelineRepository(db),
		CalendarEvent:     NewCalendarEventRepository(db),
		Note:              NewNoteRepository(db),
		Habit:             NewHabitRepository(db),
		HabitTracking:     NewHabitTrackingRepository(db),
	}
}