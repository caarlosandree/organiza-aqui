import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface PrivacyState {
  isPrivacyMode: boolean
  togglePrivacy: () => void
  setPrivacyMode: (enabled: boolean) => void
}

export const usePrivacyStore = create<PrivacyState>()(
  persist(
    (set) => ({
      isPrivacyMode: false,
      togglePrivacy: () => set((state) => ({ isPrivacyMode: !state.isPrivacyMode })),
      setPrivacyMode: (enabled: boolean) => set({ isPrivacyMode: enabled }),
    }),
    {
      name: 'privacy-storage',
    }
  )
)
