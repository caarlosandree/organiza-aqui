'use client'

import { useState, useRef } from 'react'
import { Upload, FileText, AlertCircle, Landmark, CreditCard, CloudUpload } from 'lucide-react'
import { Input } from '@/components/ui/input'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { useAccounts } from '@/hooks/queries/useAccounts'
import { useCreditCards } from '@/hooks/queries/useCreditCards'
import { usePreviewCSV, useImportCSV, usePreviewOFX, useImportOFX } from '@/hooks/mutations/useImportMutations'
import { ImportPreview } from './ImportPreview'
import { cn } from '@/lib/utils'

interface ImportDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  type: 'csv' | 'ofx'
}

export function ImportDialog({ open, onOpenChange, type }: ImportDialogProps) {
  const [context, setContext] = useState<'bank' | 'credit_card'>('bank')
  const [selectedAccount, setSelectedAccount] = useState<string>('')
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [showPreview, setShowPreview] = useState(false)
  const [isDragging, setIsDragging] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const { data: accounts } = useAccounts()
  const { data: creditCards } = useCreditCards()
  const previewCSVMutation = usePreviewCSV()
  const importCSVMutation = useImportCSV()
  const previewOFXMutation = usePreviewOFX()
  const importOFXMutation = useImportOFX()

  // Usar os hooks corretos baseado no tipo
  const previewMutation = type === 'csv' ? previewCSVMutation : previewOFXMutation
  const importMutation = type === 'csv' ? importCSVMutation : importOFXMutation

  // Filtrar contas baseado no contexto
  const availableAccounts = context === 'credit_card'
    ? accounts?.filter(acc => acc.type === 'credit' || creditCards?.some(cc => cc.account_id === acc.id))
    : accounts?.filter(acc => acc.type !== 'credit')

  // Quando o contexto muda, limpar conta selecionada
  const handleContextChange = (newContext: 'bank' | 'credit_card') => {
    setContext(newContext)
    setSelectedAccount('')
  }

  const validateFile = (file: File): boolean => {
    const extension = file.name.split('.').pop()?.toLowerCase()
    if (type === 'csv') {
      return extension === 'csv'
    } else {
      return extension === 'ofx'
    }
  }

  const handleFileSelect = (file: File | null) => {
    if (!file) return
    
    if (!validateFile(file)) {
      // Poderia mostrar um toast aqui, mas por enquanto apenas retorna
      return
    }
    
    setSelectedFile(file)
    setShowPreview(false)
  }

  const handleFileInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    handleFileSelect(file || null)
  }

  const handleDragEnter = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(true)
  }

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    // Só remove o estado de dragging se sair da área de drop
    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
      setIsDragging(false)
    }
  }

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
  }

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault()
    e.stopPropagation()
    setIsDragging(false)

    const files = e.dataTransfer.files
    if (files && files.length > 0) {
      handleFileSelect(files[0])
    }
  }

  const handleClickUpload = () => {
    fileInputRef.current?.click()
  }

  const handlePreview = async () => {
    if (!selectedFile || !selectedAccount) return

    try {
      if (type === 'csv') {
        await previewCSVMutation.mutateAsync({ accountId: selectedAccount, file: selectedFile })
      } else {
        await previewOFXMutation.mutateAsync({ accountId: selectedAccount, file: selectedFile })
      }
      setShowPreview(true)
    } catch (error) {
      // Erro já é tratado pelo hook de mutation
      console.error('Erro ao fazer preview:', error)
    }
  }

  // Resetar estado quando o dialog fechar
  const handleOpenChange = (open: boolean) => {
    if (!open) {
      setSelectedFile(null)
      setSelectedAccount('')
      setShowPreview(false)
      setContext('bank')
    }
    onOpenChange(open)
  }

  const handleImport = async (externalIds: string[]) => {
    if (!selectedFile || !selectedAccount) return

    try {
      if (type === 'csv') {
        await importCSVMutation.mutateAsync({ accountId: selectedAccount, file: selectedFile })
      } else {
        await importOFXMutation.mutateAsync({ accountId: selectedAccount, file: selectedFile, externalIds })
      }
      onOpenChange(false)
      setSelectedFile(null)
      setSelectedAccount('')
      setShowPreview(false)
      setContext('bank') // Resetar para o padrão
    } catch (error) {
      // Erro já é tratado pelo hook de mutation
      console.error('Erro ao importar:', error)
    }
  }

  const fileTypeLabel = type.toUpperCase()

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Importar {fileTypeLabel}</DialogTitle>
          <DialogDescription>
            Selecione a conta e o arquivo {fileTypeLabel} para importar
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Seleção de Contexto: Banco ou Cartão */}
          <div className="grid grid-cols-2 gap-3 p-1 bg-muted rounded-lg">
            <button
              type="button"
              className={cn(
                'py-2 px-3 rounded-md text-sm font-medium transition-all flex items-center justify-center gap-2',
                context === 'bank'
                  ? 'bg-background shadow text-foreground'
                  : 'text-muted-foreground hover:text-foreground'
              )}
              onClick={() => handleContextChange('bank')}
            >
              <Landmark size={16} /> Conta Banco
            </button>
            <button
              type="button"
              className={cn(
                'py-2 px-3 rounded-md text-sm font-medium transition-all flex items-center justify-center gap-2',
                context === 'credit_card'
                  ? 'bg-background shadow text-foreground'
                  : 'text-muted-foreground hover:text-foreground'
              )}
              onClick={() => handleContextChange('credit_card')}
            >
              <CreditCard size={16} /> Cartão Crédito
            </button>
          </div>

          <div className="space-y-2">
            <Label htmlFor="account">
              {context === 'credit_card' ? 'Cartão vinculado à conta' : 'Conta Bancária'}
            </Label>
            <Select value={selectedAccount} onValueChange={setSelectedAccount}>
              <SelectTrigger id="account">
                <SelectValue placeholder="Selecione a conta" />
              </SelectTrigger>
              <SelectContent>
                {availableAccounts?.map((account) => (
                  <SelectItem key={account.id} value={account.id}>
                    {account.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="file">Arquivo {fileTypeLabel}</Label>
            
            {/* Input file oculto */}
            <Input
              ref={fileInputRef}
              id="file"
              type="file"
              accept={type === 'csv' ? '.csv' : '.ofx'}
              onChange={handleFileInputChange}
              className="hidden"
            />

            {/* Área de Drag and Drop */}
            <div
              onDragEnter={handleDragEnter}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
              onClick={handleClickUpload}
              className={cn(
                'relative border-2 border-dashed rounded-lg p-8 transition-all cursor-pointer',
                'hover:border-primary hover:bg-primary/5',
                isDragging
                  ? 'border-primary bg-primary/10 scale-[1.02]'
                  : 'border-muted-foreground/25',
                selectedFile && 'border-primary/50 bg-primary/5'
              )}
            >
              <div className="flex flex-col items-center justify-center gap-4 text-center">
                {isDragging ? (
                  <>
                    <CloudUpload className="h-12 w-12 text-primary animate-bounce" />
                    <div className="space-y-1">
                      <p className="text-sm font-medium text-primary">Solte o arquivo aqui</p>
                      <p className="text-xs text-muted-foreground">
                        Arquivo {fileTypeLabel} será importado
                      </p>
                    </div>
                  </>
                ) : selectedFile ? (
                  <>
                    <FileText className="h-12 w-12 text-primary" />
                    <div className="space-y-1">
                      <p className="text-sm font-medium">{selectedFile.name}</p>
                      <p className="text-xs text-muted-foreground">
                        Clique para selecionar outro arquivo
                      </p>
                    </div>
                  </>
                ) : (
                  <>
                    <Upload className="h-12 w-12 text-muted-foreground" />
                    <div className="space-y-1">
                      <p className="text-sm font-medium">
                        Arraste e solte o arquivo {fileTypeLabel} aqui
                      </p>
                      <p className="text-xs text-muted-foreground">
                        ou clique para selecionar um arquivo
                      </p>
                    </div>
                  </>
                )}
              </div>
            </div>
          </div>

          {previewMutation.isError && (
            <Alert variant="destructive">
              <AlertCircle className="h-4 w-4" />
              <AlertTitle>Erro</AlertTitle>
              <AlertDescription>
                {previewMutation.error instanceof Error
                  ? previewMutation.error.message
                  : 'Erro ao fazer preview do arquivo'}
              </AlertDescription>
            </Alert>
          )}

          {showPreview && previewMutation.data && (
            <ImportPreview
              preview={previewMutation.data}
              onImport={handleImport}
              isLoading={importMutation.isPending}
            />
          )}

          {!showPreview && (
            <div className="flex gap-2">
              <Button
                type="button"
                onClick={handlePreview}
                disabled={!selectedFile || !selectedAccount || previewMutation.isPending}
                className="flex-1"
              >
                <Upload className="h-4 w-4 mr-2" />
                {previewMutation.isPending ? 'Processando...' : 'Visualizar'}
              </Button>
              <Button type="button" variant="outline" onClick={() => handleOpenChange(false)}>
                Cancelar
              </Button>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
