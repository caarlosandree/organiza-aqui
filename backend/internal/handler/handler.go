package handler

import (
	"github.com/labstack/echo/v4"
	appMiddleware "github.com/organiza-aqui/backend/internal/middleware"
	"github.com/organiza-aqui/backend/internal/service"
)

// Handlers contém todos os handlers
type Handlers struct {
	Auth              *AuthHandler
	Account           *AccountHandler
	Bank              *BankHandler
	Category          *CategoryHandler
	Transaction       *TransactionHandler
	TransactionPeriod *TransactionPeriodHandler
	Recurrence        *RecurrenceHandler
	Installment       *InstallmentHandler
	CreditCard        *CreditCardHandler
	Import            *ImportHandler
	Analytics         *AnalyticsHandler
	TaskStatus        *TaskStatusHandler
	Task              *TaskHandler
	Timeline          *TimelineHandler
	CalendarEvent     *CalendarEventHandler
	Note              *NoteHandler
	Habit             *HabitHandler
	HabitTracking     *HabitTrackingHandler
}

// NewHandlers cria uma nova instância de Handlers
func NewHandlers(services *service.Services) *Handlers {
	return &Handlers{
		Auth:              NewAuthHandler(services.Auth),
		Account:           NewAccountHandler(services.Account),
		Bank:              NewBankHandler(services.Bank),
		Category:          NewCategoryHandler(services.Category),
		Transaction:       NewTransactionHandler(services.Transaction),
		TransactionPeriod: NewTransactionPeriodHandler(services.TransactionPeriod),
		Recurrence:        NewRecurrenceHandler(services.Recurrence),
		Installment:       NewInstallmentHandler(services.Installment),
		CreditCard:        NewCreditCardHandler(services.CreditCard),
		Import:            NewImportHandler(services.Import),
		Analytics:         NewAnalyticsHandler(services.Analytics),
		TaskStatus:        NewTaskStatusHandler(services.TaskStatus),
		Task:              NewTaskHandler(services.Task),
		Timeline:          NewTimelineHandler(services.Timeline),
		CalendarEvent:     NewCalendarEventHandler(services.CalendarEvent),
		Note:              NewNoteHandler(services.Note),
		Habit:             NewHabitHandler(services.Habit),
		HabitTracking:     NewHabitTrackingHandler(services.HabitTracking),
	}
}

// SetupRoutes configura todas as rotas da API
func SetupRoutes(e *echo.Echo, handlers *Handlers) {
	// Versionamento da API
	v1 := e.Group("/api/v1")

	// Rotas públicas (sem autenticação)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", handlers.Auth.Register)
		auth.POST("/login", handlers.Auth.Login)
	}

	// Rotas protegidas (com autenticação)
	// O middleware precisa do authService, vamos passar através do handler
	protected := v1.Group("", appMiddleware.AuthMiddleware(handlers.Auth.GetAuthService()))
	{
		// Auth
		protected.POST("/auth/logout", handlers.Auth.Logout)
		protected.GET("/auth/me", handlers.Auth.Me)

		// Accounts
		accounts := protected.Group("/accounts")
		{
			accounts.GET("", handlers.Account.GetAccounts)
			accounts.POST("", handlers.Account.CreateAccount)
			accounts.GET("/:id", handlers.Account.GetAccount)
			accounts.PUT("/:id", handlers.Account.UpdateAccount)
			accounts.PUT("/:id/initial-balance", handlers.Account.UpdateInitialBalance)
			accounts.POST("/:id/recalculate-balance", handlers.Account.RecalculateBalance)
			accounts.DELETE("/:id", handlers.Account.DeleteAccount)
		}

		// Banks
		banks := protected.Group("/banks")
		{
			banks.GET("", handlers.Bank.GetAllBanks)
			banks.GET("/:id", handlers.Bank.GetBank)
			banks.POST("/sync", handlers.Bank.SyncBanks)
		}

		// Categories
		categories := protected.Group("/categories")
		{
			categories.GET("", handlers.Category.GetCategories)
			categories.POST("", handlers.Category.CreateCategory)
			categories.GET("/:id", handlers.Category.GetCategory)
			categories.PUT("/:id", handlers.Category.UpdateCategory)
			categories.DELETE("/:id", handlers.Category.DeleteCategory)
		}

		// Transactions
		transactions := protected.Group("/transactions")
		{
			transactions.GET("", handlers.Transaction.GetTransactions)
			transactions.GET("/statement", handlers.Transaction.GetStatement)
			transactions.POST("", handlers.Transaction.CreateTransaction)
			transactions.GET("/:id", handlers.Transaction.GetTransaction)
			transactions.PUT("/:id", handlers.Transaction.UpdateTransaction)
			transactions.DELETE("/:id", handlers.Transaction.DeleteTransaction)
			transactions.PATCH("/:id/status", handlers.Transaction.UpdateTransactionStatus)
		}

		// Transaction Periods
		transactionPeriods := protected.Group("/transaction-periods")
		{
			transactionPeriods.GET("", handlers.TransactionPeriod.GetTransactionPeriods)
			transactionPeriods.GET("/:id", handlers.TransactionPeriod.GetTransactionPeriod)
		}

		// Installments
		installments := protected.Group("/installments")
		{
			installments.POST("", handlers.Installment.CreateInstallments)
			installments.GET("/:parent_id", handlers.Installment.GetInstallments)
			installments.PUT("/:id", handlers.Installment.UpdateInstallment)
			installments.DELETE("/:id/cancel-future", handlers.Installment.CancelFutureInstallments)
		}

		// Credit Cards
		creditCards := protected.Group("/credit-cards")
		{
			creditCards.GET("", handlers.CreditCard.GetCreditCards)
			creditCards.POST("", handlers.CreditCard.CreateCreditCard)
			creditCards.GET("/:id", handlers.CreditCard.GetCreditCard)
			creditCards.PUT("/:id", handlers.CreditCard.UpdateCreditCard)
			creditCards.DELETE("/:id", handlers.CreditCard.DeleteCreditCard)
			creditCards.GET("/:id/available-limit", handlers.CreditCard.GetAvailableLimit)
			creditCards.GET("/:id/bill-projection", handlers.CreditCard.GetBillProjection)
			creditCards.GET("/:id/bills", handlers.CreditCard.GetBills)
			creditCards.POST("/:id/bills/:bill_id/close", handlers.CreditCard.CloseBill)
			creditCards.POST("/:id/bills/:bill_id/pay", handlers.CreditCard.PayBill)
		}

		// Import
		importGroup := protected.Group("/import")
		{
			importGroup.POST("/ofx", handlers.Import.ImportOFX)
			importGroup.POST("/csv", handlers.Import.ImportCSV)
			importGroup.POST("/ofx/preview", handlers.Import.PreviewOFX)
			importGroup.POST("/csv/preview", handlers.Import.PreviewCSV)
		}

		// Recurrence
		recurrence := protected.Group("/recurrence")
		{
			recurrence.GET("", handlers.Recurrence.GetPatterns)
			recurrence.POST("", handlers.Recurrence.CreatePattern)
			recurrence.GET("/:id", handlers.Recurrence.GetPattern)
			recurrence.PUT("/:id", handlers.Recurrence.UpdatePattern)
			recurrence.DELETE("/:id", handlers.Recurrence.DeletePattern)
			recurrence.POST("/generate", handlers.Recurrence.GenerateTransactions)
		}

		// Analytics
		analytics := protected.Group("/analytics")
		{
			analytics.GET("/income-expense", handlers.Analytics.GetIncomeExpenseByPeriod)
			analytics.GET("/category-breakdown", handlers.Analytics.GetCategoryBreakdown)
			analytics.GET("/monthly-trend", handlers.Analytics.GetMonthlyTrend)
			analytics.GET("/patrimonio-liquido", handlers.Analytics.GetPatrimonioLiquido)
			analytics.GET("/calendario-vencimentos", handlers.Analytics.GetCalendarioVencimentos)
			analytics.GET("/gastos-por-tag", handlers.Analytics.GetGastosPorTag)
		}

		// Task Statuses
		taskStatuses := protected.Group("/task-statuses")
		{
			taskStatuses.GET("", handlers.TaskStatus.GetStatuses)
			taskStatuses.POST("", handlers.TaskStatus.CreateStatus)
			taskStatuses.GET("/:id", handlers.TaskStatus.GetStatus)
			taskStatuses.PUT("/:id", handlers.TaskStatus.UpdateStatus)
			taskStatuses.DELETE("/:id", handlers.TaskStatus.DeleteStatus)
			taskStatuses.POST("/reorder", handlers.TaskStatus.ReorderStatuses)
		}

		// Tasks
		tasks := protected.Group("/tasks")
		{
			tasks.GET("", handlers.Task.GetTasks)
			tasks.POST("", handlers.Task.CreateTask)
			tasks.GET("/:id", handlers.Task.GetTask)
			tasks.PUT("/:id", handlers.Task.UpdateTask)
			tasks.DELETE("/:id", handlers.Task.DeleteTask)
			tasks.POST("/reorder", handlers.Task.ReorderTask)
			tasks.POST("/:id/complete", handlers.Task.CompleteTask)
			tasks.POST("/:id/uncomplete", handlers.Task.UncompleteTask)
		}

		// Timeline
		timeline := protected.Group("/timeline")
		{
			timeline.GET("", handlers.Timeline.GetTimelineEvents)
			timeline.GET("/summary", handlers.Timeline.GetTimelineSummary)
		}

		// Calendar Events
		calendarEvents := protected.Group("/calendar-events")
		{
			calendarEvents.GET("", handlers.CalendarEvent.GetEvents)
			calendarEvents.POST("", handlers.CalendarEvent.CreateEvent)
			calendarEvents.GET("/:id", handlers.CalendarEvent.GetEvent)
			calendarEvents.PUT("/:id", handlers.CalendarEvent.UpdateEvent)
			calendarEvents.DELETE("/:id", handlers.CalendarEvent.DeleteEvent)
		}

		// Notes
		notes := protected.Group("/notes")
		{
			notes.GET("", handlers.Note.GetNotes)
			notes.POST("", handlers.Note.CreateNote)
			notes.GET("/:id", handlers.Note.GetNote)
			notes.PUT("/:id", handlers.Note.UpdateNote)
			notes.DELETE("/:id", handlers.Note.DeleteNote)
		}

		// Habits
		habits := protected.Group("/habits")
		{
			habits.GET("", handlers.Habit.GetHabits)
			habits.POST("", handlers.Habit.CreateHabit)
			habits.GET("/:id", handlers.Habit.GetHabit)
			habits.PUT("/:id", handlers.Habit.UpdateHabit)
			habits.DELETE("/:id", handlers.Habit.DeleteHabit)
			habits.GET("/:id/stats", handlers.Habit.GetHabitStats)
			habits.GET("/:habit_id/tracking", handlers.HabitTracking.GetTrackingByHabit)
		}

		// Habit Tracking
		habitTracking := protected.Group("/habit-tracking")
		{
			habitTracking.POST("", handlers.HabitTracking.CreateTracking)
			habitTracking.GET("/:id", handlers.HabitTracking.GetTracking)
			habitTracking.PUT("/:id", handlers.HabitTracking.UpdateTracking)
			habitTracking.DELETE("/:id", handlers.HabitTracking.DeleteTracking)
		}
	}
}