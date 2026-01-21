'use client'

import { useState } from 'react'
import { Upload, FileText, AlertCircle } from 'lucide-react'
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
import { usePreviewCSV, useImportCSV } from '@/hooks/mutations/useImportMutations'
import { ImportPreview } from './ImportPreview'

interface ImportDialogProps {
  open: boolean
  onOpenChange: (open: boolean) => void
  type: 'csv' | 'ofx'
}

export function ImportDialog({ open, onOpenChange, type }: ImportDialogProps) {
  const [selectedAccount, setSelectedAccount] = useState<string>('')
  const [selectedFile, setSelectedFile] = useState<File | null>(null)
  const [showPreview, setShowPreview] = useState(false)
  const { data: accounts } = useAccounts()
  const previewMutation = usePreviewCSV()
  const importMutation = useImportCSV()

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      setSelectedFile(file)
      setShowPreview(false)
    }
  }

  const handlePreview = async () => {
    if (!selectedFile || !selectedAccount) return

    previewMutation.mutate(
      { accountId: selectedAccount, file: selectedFile },
      {
        onSuccess: () => {
          setShowPreview(true)
        },
      }
    )
  }

  const handleImport = async () => {
    if (!selectedFile || !selectedAccount) return

    importMutation.mutate(
      { accountId: selectedAccount, file: selectedFile },
      {
        onSuccess: () => {
          onOpenChange(false)
          setSelectedFile(null)
          setSelectedAccount('')
          setShowPreview(false)
        },
      }
    )
  }

  const fileTypeLabel = type.toUpperCase()

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Importar {fileTypeLabel}</DialogTitle>
          <DialogDescription>
            Selecione a conta e o arquivo {fileTypeLabel} para importar
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="account">Conta</Label>
            <Select value={selectedAccount} onValueChange={setSelectedAccount}>
              <SelectTrigger id="account">
                <SelectValue placeholder="Selecione a conta" />
              </SelectTrigger>
              <SelectContent>
                {accounts?.map((account) => (
                  <SelectItem key={account.id} value={account.id}>
                    {account.name}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <Label htmlFor="file">Arquivo {fileTypeLabel}</Label>
            <div className="flex items-center gap-2">
              <Input
                id="file"
                type="file"
                accept={type === 'csv' ? '.csv' : '.ofx'}
                onChange={handleFileSelect}
                className="flex-1"
              />
            </div>
            {selectedFile && (
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <FileText className="h-4 w-4" />
                <span>{selectedFile.name}</span>
              </div>
            )}
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
              <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
                Cancelar
              </Button>
            </div>
          )}
        </div>
      </DialogContent>
    </Dialog>
  )
}
