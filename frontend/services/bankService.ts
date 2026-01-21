import api from '@/lib/axios'
import type { Bank } from '@/types/bank'

export const bankService = {
  getBanks: async (): Promise<Bank[]> => {
    const response = await api.get<Bank[]>('/banks')
    return response.data
  },

  getBank: async (id: string): Promise<Bank> => {
    const response = await api.get<Bank>(`/banks/${id}`)
    return response.data
  },

  syncBanks: async (): Promise<void> => {
    await api.post('/banks/sync')
  },
}
