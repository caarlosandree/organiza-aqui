export interface Account {
  id: string
  user_id: string
  name: string
  type: 'checking' | 'savings' | 'credit' | 'investment'
  balance: number // em centavos
  currency: string
  bank_id: string
  initial_balance?: number // saldo inicial em centavos
  initial_balance_date?: string // data de referência do saldo inicial
  created_at: string
}

export interface Category {
  id: string
  user_id: string
  name: string
  parent_id?: string
  path: string
  type: 'income' | 'expense'
  color: string
  created_at: string
  children?: Category[]
}

export interface Transaction {
  id: string
  user_id: string
  account_id: string
  category_id?: string
  type: 'income' | 'expense' | 'transfer'
  amount: number // em centavos
  description: string
  date: string
  reference_month?: string // formato YYYY-MM
  status: 'pending' | 'paid' | 'cancelled'
  tags?: string[]
  to_account_id?: string
  parent_transaction_id?: string
  installment_number?: number
  total_installments?: number
  period_id?: string
  created_at: string
}

export interface CreateAccountRequest {
  name: string
  type: 'checking' | 'savings' | 'credit' | 'investment'
  currency: string
  bank_id: string
}

export interface UpdateAccountRequest {
  name: string
  type: 'checking' | 'savings' | 'credit' | 'investment'
  currency: string
  bank_id: string
}

export interface CreateCategoryRequest {
  name: string
  parent_id?: string
  type: 'income' | 'expense'
  color: string
}

export interface UpdateCategoryRequest {
  name: string
  parent_id?: string
  type: 'income' | 'expense'
  color: string
}

export interface CreateTransactionRequest {
  account_id: string
  category_id?: string | null
  type: 'income' | 'expense' | 'transfer'
  amount: number
  description?: string
  date: string
  reference_month?: string // formato YYYY-MM
  status?: 'pending' | 'paid' | 'cancelled'
  tags?: string[]
  to_account_id?: string
  total_installments?: number
}

export interface UpdateTransactionRequest {
  account_id: string
  category_id?: string | null
  type: 'income' | 'expense' | 'transfer'
  amount: number
  description?: string
  date: string
  reference_month?: string // formato YYYY-MM
  status?: 'pending' | 'paid' | 'cancelled'
  tags?: string[]
  to_account_id?: string
}

export interface TransactionFilters {
  account_id?: string
  category_id?: string
  type?: 'income' | 'expense' | 'transfer'
  status?: 'pending' | 'paid' | 'cancelled'
  tags?: string[]
  parent_transaction_id?: string
  start_date?: string
  end_date?: string
  limit?: number
  offset?: number
}

export interface TransactionsResponse {
  data: Transaction[]
  total: number
}

// Credit Card Types
export interface CreditCard {
  id: string
  user_id: string
  name: string
  account_id: string
  limit_amount: number // em centavos
  closing_day: number
  due_day: number
  color: string
  created_at: string
  updated_at: string
}

export interface CreateCreditCardRequest {
  name: string
  account_id: string
  limit_amount: number
  closing_day: number
  due_day: number
  color: string
}

export interface UpdateCreditCardRequest {
  name: string
  account_id: string
  limit_amount: number
  closing_day: number
  due_day: number
  color: string
}

export interface CreditCardBill {
  id: string
  credit_card_id: string
  month: number
  year: number
  status: 'open' | 'closed' | 'paid'
  closing_date: string
  due_date: string
  total_amount: number // em centavos, calculado dinamicamente
  payment_transaction_id?: string
  created_at: string
  updated_at: string
}

export interface PayBillRequest {
  account_id: string
  date: string
}

export interface ProjecaoFaturaResponse {
  credit_card_id: string
  month: number
  year: number
  projected_amount: number // em centavos
  closing_date: string
  due_date: string
}

// Installment Types
export interface UpdateInstallmentRequest {
  amount: number
  description: string
  date: string
  status?: 'pending' | 'paid' | 'cancelled'
  scope: 'this' | 'this_and_future' | 'all'
}

// Import Types
export interface ImportOFXRequest {
  account_id: string
  file: File | Blob
}

export interface ImportCSVRequest {
  account_id: string
  file: File | Blob
  delimiter?: string
}

export interface ImportPreviewResponse {
  total_transactions: number
  duplicates: number
  new_transactions: number
  transactions: Transaction[]
}

export interface ImportResponse {
  total_processed: number
  duplicates: number
  created: number
  errors: number
}

// Analytics Types
export interface PatrimonioLiquidoResponse {
  total_patrimonio: number // em centavos
  total_contas: number // em centavos
  total_faturas: number // em centavos
  contas?: Account[]
}

export interface VencimentoItem {
  id: string
  type: 'transaction' | 'bill'
  description: string
  amount: number // em centavos
  date: string
  status: string
}

export interface CalendarioVencimentosResponse {
  items: VencimentoItem[]
}

export interface TagBreakdown {
  tag: string
  total_amount: number // em centavos
  transaction_count: number
  percentage?: number
}

export interface GastosPorTagResponse {
  tags: TagBreakdown[]
}

// Transaction Period Types
export interface TransactionPeriod {
  id: string
  user_id: string
  account_id: string
  account_name?: string
  period_type: 'bank' | 'credit_card'
  year: number
  month: number
  status: 'open' | 'closed' | 'archived'
  created_at: string
  updated_at: string
}

export interface PeriodWithTransactions {
  period: TransactionPeriod
  transactions: Transaction[]
  total_income: number
  total_bank_expense: number
  total_credit_card_expense: number
  balance: number // receitas - despesas bancárias (sem cartão)
}

export interface TransactionPeriodFilters {
  account_id?: string
  period_type?: 'bank' | 'credit_card'
  year?: number
  month?: number
  status?: 'open' | 'closed' | 'archived'
  reference_month?: string // formato YYYY-MM
}
