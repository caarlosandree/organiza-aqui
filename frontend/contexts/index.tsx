// React Contexts
// Adicione seus contexts aqui conforme necessário

'use client'

import { createContext, useContext, ReactNode } from 'react'

// Exemplo de context - ajuste conforme necessário
type ExampleContextType = Record<string, never>

const ExampleContext = createContext<ExampleContextType | undefined>(undefined)

export const ExampleProvider = ({ children }: { children: ReactNode }) => {
  // Lógica do provider
  return (
    <ExampleContext.Provider value={{}}>
      {children}
    </ExampleContext.Provider>
  )
}

export const useExampleContext = () => {
  const context = useContext(ExampleContext)
  if (context === undefined) {
    throw new Error('useExampleContext must be used within ExampleProvider')
  }
  return context
}
