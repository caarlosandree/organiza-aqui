'use client'

import { useState } from 'react'
import { useForm, Controller } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { X } from 'lucide-react'
import {
  createNoteSchema,
  updateNoteSchema,
  type CreateNoteFormData,
  type UpdateNoteFormData,
} from '@/schemas/knowledgeSchema'
import type { Note } from '@/types/knowledge'

interface NoteFormProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  note?: Note
  onSubmit: (data: CreateNoteFormData | UpdateNoteFormData) => void
  isSubmitting?: boolean
}

export const NoteForm = ({
  open,
  onOpenChange,
  note,
  onSubmit,
  isSubmitting = false,
}: NoteFormProps) => {
  const isEditing = !!note
  const [tagInput, setTagInput] = useState('')

  const {
    register,
    handleSubmit,
    control,
    formState: { errors },
    reset,
    watch,
    setValue,
  } = useForm({
    resolver: zodResolver(isEditing ? updateNoteSchema : createNoteSchema),
    defaultValues: note
      ? {
          title: note.title,
          content: note.content,
          tags: note.tags || [],
          is_pinned: note.is_pinned,
        }
      : {
          tags: [],
          is_pinned: false,
        },
  })

  const tags = watch('tags') || []

  const handleAddTag = () => {
    if (tagInput.trim() && !tags.includes(tagInput.trim())) {
      setValue('tags', [...tags, tagInput.trim()])
      setTagInput('')
    }
  }

  const handleRemoveTag = (tagToRemove: string) => {
    setValue(
      'tags',
      tags.filter((tag) => tag !== tagToRemove)
    )
  }

  const handleFormSubmit = (data: CreateNoteFormData | UpdateNoteFormData) => {
    onSubmit(data)
    if (!isEditing) {
      reset()
      setTagInput('')
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{isEditing ? 'Editar Anotação' : 'Nova Anotação'}</DialogTitle>
          <DialogDescription>
            {isEditing
              ? 'Atualize as informações da anotação'
              : 'Preencha os dados para criar uma nova anotação'}
          </DialogDescription>
        </DialogHeader>
        <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="title">Título</Label>
            <Input
              id="title"
              {...register('title')}
              className={errors.title ? 'border-destructive' : ''}
              placeholder="Digite o título da anotação"
            />
            {errors.title && (
              <p className="text-sm text-destructive">{errors.title.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="content">Conteúdo</Label>
            <Textarea
              id="content"
              {...register('content')}
              className={errors.content ? 'border-destructive' : ''}
              placeholder="Digite o conteúdo da anotação"
              rows={10}
            />
            {errors.content && (
              <p className="text-sm text-destructive">{errors.content.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="tags">Tags</Label>
            <div className="flex gap-2">
              <Input
                id="tags"
                value={tagInput}
                onChange={(e) => setTagInput(e.target.value)}
                onKeyDown={(e) => {
                  if (e.key === 'Enter') {
                    e.preventDefault()
                    handleAddTag()
                  }
                }}
                placeholder="Digite uma tag e pressione Enter"
              />
              <Button type="button" onClick={handleAddTag} variant="outline">
                Adicionar
              </Button>
            </div>
            {tags.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-2">
                {tags.map((tag) => (
                  <Badge key={tag} variant="secondary" className="gap-1">
                    {tag}
                    <button
                      type="button"
                      onClick={() => handleRemoveTag(tag)}
                      className="ml-1 hover:bg-destructive/20 rounded-full p-0.5"
                    >
                      <X className="h-3 w-3" />
                    </button>
                  </Badge>
                ))}
              </div>
            )}
          </div>

          <div className="space-y-2">
            <Controller
              name="is_pinned"
              control={control}
              render={({ field }) => (
                <div className="flex items-center space-x-2">
                  <Checkbox
                    id="is_pinned"
                    checked={field.value}
                    onCheckedChange={field.onChange}
                  />
                  <Label htmlFor="is_pinned" className="cursor-pointer">
                    Fixar anotação
                  </Label>
                </div>
              )}
            />
          </div>

          <DialogFooter>
            <Button
              type="button"
              variant="outline"
              onClick={() => onOpenChange(false)}
              disabled={isSubmitting}
            >
              Cancelar
            </Button>
            <Button type="submit" disabled={isSubmitting}>
              {isSubmitting ? 'Salvando...' : isEditing ? 'Atualizar' : 'Criar'}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
