import { z } from 'zod'
import { calculatePasswordStrength } from '@/utils/passwordStrength'

// Validação customizada de força de senha
const passwordStrengthRefine = z.string().refine(
  (password) => {
    const { strength } = calculatePasswordStrength(password)
    // Aceitar apenas senhas médias ou mais fortes
    return strength !== 'muito-fraco' && strength !== 'fraco'
  },
  {
    message:
      'A senha deve ser pelo menos média (use letras maiúsculas, minúsculas, números e caracteres especiais)',
  }
)

export const loginSchema = z.object({
  identifier: z.string().min(1, 'Email ou usuário é obrigatório'),
  password: z.string().min(6, 'Senha deve ter no mínimo 6 caracteres'),
})

export const registerSchema = z
  .object({
    email: z.string().email('Email inválido'),
    username: z
      .string()
      .min(3, 'Username deve ter no mínimo 3 caracteres')
      .max(50, 'Username deve ter no máximo 50 caracteres')
      .regex(/^[a-zA-Z0-9]+$/, 'Username deve conter apenas letras e números'),
    password: passwordStrengthRefine.min(8, 'Senha deve ter no mínimo 8 caracteres'),
    confirmPassword: z.string().min(8, 'Confirmação de senha é obrigatória'),
    name: z
      .string()
      .min(3, 'Nome deve ter no mínimo 3 caracteres')
      .max(255, 'Nome deve ter no máximo 255 caracteres'),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'As senhas não coincidem',
    path: ['confirmPassword'],
  })

export type LoginFormData = z.infer<typeof loginSchema>
export type RegisterFormData = z.infer<typeof registerSchema>
