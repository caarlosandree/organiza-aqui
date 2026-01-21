'use client'

import { CreditCardCard } from './CreditCardCard'
import type { CreditCard } from '@/types/financial'

interface CreditCardListProps {
  creditCards: CreditCard[]
  onEdit?: (card: CreditCard) => void
  onDelete?: (id: string) => void
}

export function CreditCardList({ creditCards, onEdit }: CreditCardListProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {creditCards.map((creditCard) => (
        <CreditCardCard
          key={creditCard.id}
          creditCard={creditCard}
          onClick={() => onEdit?.(creditCard)}
        />
      ))}
    </div>
  )
}
