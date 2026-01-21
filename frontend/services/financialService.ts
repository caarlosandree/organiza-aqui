import api from '@/lib/axios'
import type {
  Account,
  Category,
  Transaction,
  CreateAccountRequest,
  UpdateAccountRequest,
  CreateCategoryRequest,
  UpdateCategoryRequest,
  CreateTransactionRequest,
  UpdateTransactionRequest,
  TransactionFilters,
  TransactionsResponse,
  CreditCard,
  CreateCreditCardRequest,
  UpdateCreditCardRequest,
  CreditCardBill,
  PayBillRequest,
  ProjecaoFaturaResponse,
  UpdateInstallmentRequest,
  ImportPreviewResponse,
  ImportResponse,
  PatrimonioLiquidoResponse,
  CalendarioVencimentosResponse,
  GastosPorTagResponse,
  TransactionPeriod,
  TransactionPeriodFilters,
  PeriodWithTransactions,
} from '@/types/financial'

export const financialService = {
  // Accounts
  getAccounts: async (): Promise<Account[]> => {
    const response = await api.get<Account[]>('/accounts')
    return response.data
  },

  getAccount: async (id: string): Promise<Account> => {
    const response = await api.get<Account>(`/accounts/${id}`)
    return response.data
  },

  createAccount: async (data: CreateAccountRequest): Promise<Account> => {
    const response = await api.post<Account>('/accounts', data)
    return response.data
  },

  updateAccount: async (id: string, data: UpdateAccountRequest): Promise<Account> => {
    const response = await api.put<Account>(`/accounts/${id}`, data)
    return response.data
  },

  deleteAccount: async (id: string): Promise<void> => {
    await api.delete(`/accounts/${id}`)
  },

  updateInitialBalance: async (id: string, data: { balance: number; date: string }): Promise<Account> => {
    const response = await api.put<Account>(`/accounts/${id}/initial-balance`, data)
    return response.data
  },

  recalculateBalance: async (id: string): Promise<Account> => {
    const response = await api.post<Account>(`/accounts/${id}/recalculate-balance`)
    return response.data
  },

  // Categories
  getCategories: async (params?: { type?: string; tree?: boolean }): Promise<Category[]> => {
    const response = await api.get<Category[]>('/categories', { params })
    return response.data
  },

  getCategory: async (id: string): Promise<Category> => {
    const response = await api.get<Category>(`/categories/${id}`)
    return response.data
  },

  createCategory: async (data: CreateCategoryRequest): Promise<Category> => {
    const response = await api.post<Category>('/categories', data)
    return response.data
  },

  updateCategory: async (id: string, data: UpdateCategoryRequest): Promise<Category> => {
    const response = await api.put<Category>(`/categories/${id}`, data)
    return response.data
  },

  deleteCategory: async (id: string): Promise<void> => {
    await api.delete(`/categories/${id}`)
  },

  // Transactions
  getTransactions: async (filters?: TransactionFilters): Promise<TransactionsResponse> => {
    const response = await api.get<TransactionsResponse>('/transactions', { params: filters })
    return response.data
  },

  getTransaction: async (id: string): Promise<Transaction> => {
    const response = await api.get<Transaction>(`/transactions/${id}`)
    return response.data
  },

  createTransaction: async (data: CreateTransactionRequest): Promise<Transaction> => {
    const response = await api.post<Transaction>('/transactions', data)
    return response.data
  },

  updateTransaction: async (id: string, data: UpdateTransactionRequest): Promise<Transaction> => {
    const response = await api.put<Transaction>(`/transactions/${id}`, data)
    return response.data
  },

  deleteTransaction: async (id: string): Promise<void> => {
    await api.delete(`/transactions/${id}`)
  },

  updateTransactionStatus: async (id: string, status: 'pending' | 'paid' | 'cancelled'): Promise<Transaction> => {
    const response = await api.patch<Transaction>(`/transactions/${id}/status`, { status })
    return response.data
  },

  // Credit Cards
  getCreditCards: async (): Promise<CreditCard[]> => {
    const response = await api.get<CreditCard[]>('/credit-cards')
    return response.data
  },

  getCreditCard: async (id: string): Promise<CreditCard> => {
    const response = await api.get<CreditCard>(`/credit-cards/${id}`)
    return response.data
  },

  createCreditCard: async (data: CreateCreditCardRequest): Promise<CreditCard> => {
    const response = await api.post<CreditCard>('/credit-cards', data)
    return response.data
  },

  updateCreditCard: async (id: string, data: UpdateCreditCardRequest): Promise<CreditCard> => {
    const response = await api.put<CreditCard>(`/credit-cards/${id}`, data)
    return response.data
  },

  deleteCreditCard: async (id: string): Promise<void> => {
    await api.delete(`/credit-cards/${id}`)
  },

  getAvailableLimit: async (id: string): Promise<{ available_limit: number }> => {
    const response = await api.get<{ available_limit: number }>(`/credit-cards/${id}/available-limit`)
    return response.data
  },

  getBillProjection: async (id: string): Promise<ProjecaoFaturaResponse> => {
    const response = await api.get<ProjecaoFaturaResponse>(`/credit-cards/${id}/bill-projection`)
    return response.data
  },

  getBills: async (id: string): Promise<CreditCardBill[]> => {
    const response = await api.get<CreditCardBill[]>(`/credit-cards/${id}/bills`)
    return response.data
  },

  closeBill: async (creditCardId: string, billId: string): Promise<void> => {
    await api.post(`/credit-cards/${creditCardId}/bills/${billId}/close`)
  },

  payBill: async (creditCardId: string, billId: string, data: PayBillRequest): Promise<CreditCardBill> => {
    const response = await api.post<CreditCardBill>(`/credit-cards/${creditCardId}/bills/${billId}/pay`, data)
    return response.data
  },

  // Installments
  createInstallments: async (data: CreateTransactionRequest): Promise<Transaction> => {
    const response = await api.post<Transaction>('/installments', data)
    return response.data
  },

  getInstallments: async (parentId: string): Promise<Transaction[]> => {
    const response = await api.get<Transaction[]>(`/installments/${parentId}`)
    return response.data
  },

  updateInstallment: async (id: string, data: UpdateInstallmentRequest): Promise<Transaction> => {
    const response = await api.put<Transaction>(`/installments/${id}`, data)
    return response.data
  },

  cancelFutureInstallments: async (id: string): Promise<void> => {
    await api.delete(`/installments/${id}/cancel-future`)
  },

  // Import
  importOFX: async (accountId: string, file: File, externalIds?: string[]): Promise<ImportResponse> => {
    const formData = new FormData()
    formData.append('account_id', accountId)
    formData.append('file', file)
    if (externalIds && externalIds.length > 0) {
      formData.append('external_ids', JSON.stringify(externalIds))
    }
    const response = await api.post<ImportResponse>('/import/ofx', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  importCSV: async (accountId: string, file: File, delimiter?: string): Promise<ImportResponse> => {
    const formData = new FormData()
    formData.append('account_id', accountId)
    formData.append('file', file)
    if (delimiter) {
      formData.append('delimiter', delimiter)
    }
    const response = await api.post<ImportResponse>('/import/csv', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  previewOFX: async (accountId: string, file: File): Promise<ImportPreviewResponse> => {
    const formData = new FormData()
    formData.append('account_id', accountId)
    formData.append('file', file)
    const response = await api.post<ImportPreviewResponse>('/import/ofx/preview', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  previewCSV: async (accountId: string, file: File, delimiter?: string): Promise<ImportPreviewResponse> => {
    const formData = new FormData()
    formData.append('account_id', accountId)
    formData.append('file', file)
    if (delimiter) {
      formData.append('delimiter', delimiter)
    }
    const response = await api.post<ImportPreviewResponse>('/import/csv/preview', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    })
    return response.data
  },

  // Analytics
  getPatrimonioLiquido: async (): Promise<PatrimonioLiquidoResponse> => {
    const response = await api.get<PatrimonioLiquidoResponse>('/analytics/patrimonio-liquido')
    return response.data
  },

  getCalendarioVencimentos: async (days?: number): Promise<CalendarioVencimentosResponse> => {
    const response = await api.get<CalendarioVencimentosResponse>('/analytics/calendario-vencimentos', {
      params: days ? { days } : undefined,
    })
    return response.data
  },

  getGastosPorTag: async (startDate: string, endDate: string): Promise<GastosPorTagResponse> => {
    const response = await api.get<GastosPorTagResponse>('/analytics/gastos-por-tag', {
      params: { start_date: startDate, end_date: endDate },
    })
    return response.data
  },

  // Transaction Periods
  getTransactionPeriods: async (filters?: TransactionPeriodFilters): Promise<TransactionPeriod[]> => {
    const response = await api.get<TransactionPeriod[]>('/transaction-periods', {
      params: filters,
    })
    return response.data
  },

  getTransactionPeriod: async (id: string): Promise<PeriodWithTransactions> => {
    const response = await api.get<PeriodWithTransactions>(`/transaction-periods/${id}`)
    return response.data
  },
}
