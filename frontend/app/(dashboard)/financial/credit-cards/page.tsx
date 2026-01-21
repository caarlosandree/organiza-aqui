'use client'

import { useState } from 'react'
import { Plus } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog'
import { CreditCardList } from '@/components/financial/CreditCardList'
import { CreditCardForm } from '@/components/financial/CreditCardForm'
import { useCreditCards } from '@/hooks/queries/useCreditCards'
import { useCreateCreditCard, useUpdateCreditCard, useDeleteCreditCard } from '@/hooks/mutations/useCreditCardMutations'
import { type CreditCardFormData } from '@/schemas/financialSchema'
import type { CreditCard } from '@/types/financial'
import { Loader2 } from 'lucide-react'

export default function CreditCardsPage() {
  const [isDialogOpen, setIsDialogOpen] = useState(false)
  const [editingCard, setEditingCard] = useState<CreditCard | null>(null)

  const { data: creditCards, isLoading } = useCreditCards()
  const createMutation = useCreateCreditCard()
  const updateMutation = useUpdateCreditCard()
  const deleteMutation = useDeleteCreditCard()

  const handleSubmit = async (data: CreditCardFormData) => {
    try {
      if (editingCard) {
        await updateMutation.mutateAsync({ id: editingCard.id, data })
      } else {
        await createMutation.mutateAsync(data)
      }
      setIsDialogOpen(false)
      setEditingCard(null)
    } catch (error) {
      console.error('Erro ao salvar cartão:', error)
    }
  }

  const handleEdit = (card: CreditCard) => {
    setEditingCard(card)
    setIsDialogOpen(true)
  }

  const handleDelete = async (id: string) => {
    try {
      await deleteMutation.mutateAsync(id)
      setIsDialogOpen(false)
      setEditingCard(null)
    } catch (error) {
      console.error('Erro ao deletar cartão:', error)
    }
  }

  return (
    <div className="flex flex-col gap-4">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold">Cartões de Crédito</h1>
          <p className="text-muted-foreground">
            Gerencie seus cartões de crédito e faturas
          </p>
        </div>
        <Dialog open={isDialogOpen} onOpenChange={setIsDialogOpen}>
          <DialogTrigger asChild>
            <Button onClick={() => setEditingCard(null)}>
              <Plus className="mr-2 h-4 w-4" />
              Novo Cartão
            </Button>
          </DialogTrigger>
          <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
            <DialogHeader>
              <DialogTitle>
                {editingCard ? 'Editar Cartão' : 'Novo Cartão de Crédito'}
              </DialogTitle>
              <DialogDescription>
                {editingCard
                  ? 'Atualize os dados do cartão de crédito'
                  : 'Cadastre um novo cartão de crédito'}
              </DialogDescription>
            </DialogHeader>
            <CreditCardForm
              creditCard={editingCard || undefined}
              onSubmit={handleSubmit}
              onCancel={() => {
                setIsDialogOpen(false)
                setEditingCard(null)
              }}
              onDelete={editingCard ? handleDelete : undefined}
              isLoading={createMutation.isPending || updateMutation.isPending || deleteMutation.isPending}
            />
          </DialogContent>
        </Dialog>
      </div>

      {isLoading ? (
        <div className="flex justify-center items-center p-8">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      ) : creditCards && creditCards.length > 0 ? (
        <CreditCardList
          creditCards={creditCards}
          onEdit={handleEdit}
          onDelete={handleDelete}
        />
      ) : (
        <div className="text-center p-8 text-muted-foreground">
          <p>Nenhum cartão de crédito cadastrado.</p>
          <p className="text-sm">Clique em &quot;Novo Cartão&quot; para começar.</p>
        </div>
      )}
    </div>
  )
}
