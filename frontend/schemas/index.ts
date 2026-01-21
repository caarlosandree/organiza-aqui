// Schemas Zod para validação
// Adicione seus schemas aqui conforme necessário

import { z } from 'zod'

// Exemplo de schema - ajuste conforme necessário
export const exampleSchema = z.object({
  name: z.string().min(3, 'Nome deve ter no mínimo 3 caracteres'),
  email: z.string().email('Email inválido'),
})

export type ExampleFormData = z.infer<typeof exampleSchema>
