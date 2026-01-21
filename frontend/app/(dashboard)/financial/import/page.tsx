'use client'

import { useState } from 'react'
import { Upload, FileText } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { ImportDialog } from '@/components/financial/ImportDialog'
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert'
import { CheckCircle2, AlertCircle } from 'lucide-react'

export default function ImportPage() {
  const [csvDialogOpen, setCsvDialogOpen] = useState(false)
  const [ofxDialogOpen, setOfxDialogOpen] = useState(false)

  return (
    <div className="flex flex-col gap-6">
      <div>
        <h1 className="text-3xl font-bold">Importar Transações</h1>
        <p className="text-muted-foreground">
          Importe transações de arquivos CSV ou OFX
        </p>
      </div>

      <Alert>
        <AlertCircle className="h-4 w-4" />
        <AlertTitle>Formato dos Arquivos</AlertTitle>
        <AlertDescription>
          <p className="mt-2">
            <strong>CSV:</strong> O arquivo deve conter as colunas: Data, Descrição, Valor, Tipo
            (opcional: Categoria, Tags). Use vírgula ou ponto e vírgula como delimitador.
          </p>
          <p className="mt-2">
            <strong>OFX:</strong> Arquivo no formato OFX padrão exportado pelo seu banco.
          </p>
        </AlertDescription>
      </Alert>

      <div className="grid gap-4 md:grid-cols-2">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileText className="h-5 w-5" />
              Importar CSV
            </CardTitle>
            <CardDescription>
              Importe transações de um arquivo CSV
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => setCsvDialogOpen(true)} className="w-full">
              <Upload className="mr-2 h-4 w-4" />
              Selecionar Arquivo CSV
            </Button>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <FileText className="h-5 w-5" />
              Importar OFX
            </CardTitle>
            <CardDescription>
              Importe transações de um arquivo OFX
            </CardDescription>
          </CardHeader>
          <CardContent>
            <Button onClick={() => setOfxDialogOpen(true)} className="w-full">
              <Upload className="mr-2 h-4 w-4" />
              Selecionar Arquivo OFX
            </Button>
          </CardContent>
        </Card>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>Como Funciona</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-start gap-3">
            <CheckCircle2 className="h-5 w-5 text-green-600 mt-0.5" />
            <div>
              <h3 className="font-semibold">1. Seleção do Arquivo</h3>
              <p className="text-sm text-muted-foreground">
                Escolha o arquivo CSV ou OFX que deseja importar e selecione a conta de destino.
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <CheckCircle2 className="h-5 w-5 text-green-600 mt-0.5" />
            <div>
              <h3 className="font-semibold">2. Visualização</h3>
              <p className="text-sm text-muted-foreground">
                Visualize as transações que serão importadas e identifique duplicatas.
              </p>
            </div>
          </div>
          <div className="flex items-start gap-3">
            <CheckCircle2 className="h-5 w-5 text-green-600 mt-0.5" />
            <div>
              <h3 className="font-semibold">3. Importação</h3>
              <p className="text-sm text-muted-foreground">
                Confirme a importação. Transações duplicadas serão automaticamente ignoradas.
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <ImportDialog open={csvDialogOpen} onOpenChange={setCsvDialogOpen} type="csv" />
      <ImportDialog open={ofxDialogOpen} onOpenChange={setOfxDialogOpen} type="ofx" />
    </div>
  )
}
