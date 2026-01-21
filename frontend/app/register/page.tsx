'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useForm, useWatch } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Eye, EyeOff, User, Lock, Mail, LayoutDashboard, AtSign } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { useRegister } from '@/hooks/mutations/useAuth'
import { registerSchema, type RegisterFormData } from '@/schemas/authSchema'
import { useToast } from '@/hooks/useToast'
import { PasswordStrengthIndicator } from '@/components/ui/password-strength'
import Image from 'next/image'

export default function RegisterPage() {
  const registerMutation = useRegister()
  const toast = useToast()
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)
  const [usernameSuggestions, setUsernameSuggestions] = useState<string[]>([])

  const {
    register,
    handleSubmit,
    formState: { errors },
    setValue,
    control,
  } = useForm<RegisterFormData>({
    resolver: zodResolver(registerSchema),
  })

  // Monitorar o campo password para exibir indicador de força
  const password = useWatch({ control, name: 'password' })

  const onSubmit = async (data: RegisterFormData) => {
    try {
      setUsernameSuggestions([])
      await registerMutation.mutateAsync({
        email: data.email,
        username: data.username,
        password: data.password,
        name: data.name,
      })
      toast.success('Conta criada com sucesso! Redirecionando...')
      // O redirecionamento já é feito pelo hook useRegister no onSuccess
    } catch (error: unknown) {
      const err = error as { response?: { data?: { suggestions?: string[]; message?: string; error?: string } }; message?: string }
      // Verificar se há sugestões de username na resposta
      if (err?.response?.data?.suggestions) {
        setUsernameSuggestions(err.response.data.suggestions)
      }
      
      const errorMessage =
        err?.response?.data?.message ||
        err?.response?.data?.error ||
        err?.message ||
        'Erro ao criar conta'
      toast.error(errorMessage)
    }
  }

  return (
    <div className="flex w-full min-h-screen overflow-hidden">
      {/* Lado Esquerdo: Imagem (70%) */}
      <div className="hidden md:block md:w-[70%] relative overflow-hidden">
        <div className="absolute inset-0 bg-[#1F4E5F] opacity-20 z-10" />
        <Image
          src="https://images.unsplash.com/photo-1484480974693-6ca0a78fb36b?ixlib=rb-4.0.3&auto=format&fit=crop&w=2072&q=80"
          alt="Organização e Planejamento de Vida"
          fill
          className="object-cover scale-105 transition-transform duration-[10s] hover:scale-100"
          priority
        />
        {/* Overlay Content */}
        <div className="absolute bottom-16 left-16 z-20 text-white max-w-xl">
          <h1 className="text-5xl font-bold mb-4 drop-shadow-lg">
            Organiza Aqui
          </h1>
          <p className="text-xl font-light opacity-90 drop-shadow-md">
            Comece sua jornada para uma vida mais equilibrada e produtiva.
            Crie sua conta e transforme sua rotina.
          </p>
        </div>
      </div>

      {/* Lado Direito: Formulário de Registro (30%) */}
      <div className="w-full md:w-[30%] bg-[#F4F6F5] flex flex-col justify-center items-center p-8 sm:p-12 z-30 shadow-2xl overflow-y-auto">
        <div className="w-full max-w-md">
          {/* Logo/Branding */}
          <div className="mb-10 text-center md:text-left">
            <div className="inline-flex items-center justify-center p-3 rounded-xl bg-[#1F4E5F] mb-4">
              <LayoutDashboard className="text-white w-8 h-8" />
            </div>
            <h2 className="text-3xl font-bold text-[#1F2933]">
              Criar Nova Conta
            </h2>
            <p className="text-[#8A8F8E] mt-2">
              Preencha os dados para começar.
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Nome */}
            <div>
              <Label
                htmlFor="name"
                className="block text-sm font-semibold text-[#1F2933] mb-2"
              >
                Nome Completo
              </Label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-[#8A8F8E]">
                  <User className="w-5 h-5" />
                </span>
                <Input
                  id="name"
                  type="text"
                  placeholder="ex: João Silva"
                  className={`w-full pl-10 pr-4 py-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-[#3FA7A0] focus:border-transparent transition bg-white ${
                    errors.name ? 'border-red-500' : ''
                  }`}
                  {...register('name')}
                />
              </div>
              {errors.name && (
                <p className="text-sm text-red-500 mt-1">
                  {errors.name.message}
                </p>
              )}
            </div>

            {/* Email */}
            <div>
              <Label
                htmlFor="email"
                className="block text-sm font-semibold text-[#1F2933] mb-2"
              >
                E-mail
              </Label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-[#8A8F8E]">
                  <Mail className="w-5 h-5" />
                </span>
                <Input
                  id="email"
                  type="email"
                  placeholder="ex: joao@email.com"
                  className={`w-full pl-10 pr-4 py-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-[#3FA7A0] focus:border-transparent transition bg-white ${
                    errors.email ? 'border-red-500' : ''
                  }`}
                  {...register('email')}
                />
              </div>
              {errors.email && (
                <p className="text-sm text-red-500 mt-1">
                  {errors.email.message}
                </p>
              )}
            </div>

            {/* Username */}
            <div>
              <Label
                htmlFor="username"
                className="block text-sm font-semibold text-[#1F2933] mb-2"
              >
                Username
              </Label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-[#8A8F8E]">
                  <AtSign className="w-5 h-5" />
                </span>
                <Input
                  id="username"
                  type="text"
                  placeholder="ex: joaosilva"
                  className={`w-full pl-10 pr-4 py-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-[#3FA7A0] focus:border-transparent transition bg-white ${
                    errors.username ? 'border-red-500' : ''
                  }`}
                  {...register('username')}
                />
              </div>
              {errors.username && (
                <p className="text-sm text-red-500 mt-1">
                  {errors.username.message}
                </p>
              )}
              {usernameSuggestions.length > 0 && (
                <div className="mt-2">
                  <p className="text-sm text-[#8A8F8E] mb-2">
                    Este username já está em uso. Sugestões:
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {usernameSuggestions.map((suggestion) => (
                      <button
                        key={suggestion}
                        type="button"
                        onClick={() => {
                          setValue('username', suggestion)
                          setUsernameSuggestions([])
                        }}
                        className="px-3 py-1 text-sm bg-[#3FA7A0] text-white rounded-lg hover:bg-[#2E8F89] transition-colors"
                      >
                        {suggestion}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* Password */}
            <div>
              <Label
                htmlFor="password"
                className="block text-sm font-semibold text-[#1F2933] mb-2"
              >
                Senha
              </Label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-[#8A8F8E]">
                  <Lock className="w-5 h-5" />
                </span>
                <Input
                  id="password"
                  type={showPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  className={`w-full pl-10 pr-12 py-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-[#3FA7A0] focus:border-transparent transition bg-white ${
                    errors.password ? 'border-red-500' : ''
                  }`}
                  {...register('password')}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center text-[#8A8F8E] hover:text-[#1F4E5F] transition-colors"
                  aria-label={
                    showPassword ? 'Ocultar senha' : 'Mostrar senha'
                  }
                >
                  {showPassword ? (
                    <EyeOff className="w-5 h-5" />
                  ) : (
                    <Eye className="w-5 h-5" />
                  )}
                </button>
              </div>
              {password && <PasswordStrengthIndicator password={password} />}
              {errors.password && (
                <p className="text-sm text-red-500 mt-1">
                  {errors.password.message}
                </p>
              )}
            </div>

            {/* Confirm Password */}
            <div>
              <Label
                htmlFor="confirmPassword"
                className="block text-sm font-semibold text-[#1F2933] mb-2"
              >
                Confirmar Senha
              </Label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-[#8A8F8E]">
                  <Lock className="w-5 h-5" />
                </span>
                <Input
                  id="confirmPassword"
                  type={showConfirmPassword ? 'text' : 'password'}
                  placeholder="••••••••"
                  className={`w-full pl-10 pr-12 py-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-[#3FA7A0] focus:border-transparent transition bg-white ${
                    errors.confirmPassword ? 'border-red-500' : ''
                  }`}
                  {...register('confirmPassword')}
                />
                <button
                  type="button"
                  onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                  className="absolute inset-y-0 right-0 pr-3 flex items-center text-[#8A8F8E] hover:text-[#1F4E5F] transition-colors"
                  aria-label={
                    showConfirmPassword
                      ? 'Ocultar confirmação de senha'
                      : 'Mostrar confirmação de senha'
                  }
                >
                  {showConfirmPassword ? (
                    <EyeOff className="w-5 h-5" />
                  ) : (
                    <Eye className="w-5 h-5" />
                  )}
                </button>
              </div>
              {errors.confirmPassword && (
                <p className="text-sm text-red-500 mt-1">
                  {errors.confirmPassword.message}
                </p>
              )}
            </div>

            {/* Register Button */}
            <Button
              type="submit"
              disabled={registerMutation.isPending}
              className="w-full bg-[#3FA7A0] hover:bg-[#2E8F89] text-white font-bold py-3 px-4 rounded-lg shadow-lg transform active:scale-[0.98] transition-all"
            >
              {registerMutation.isPending
                ? 'Criando conta...'
                : 'Criar Conta'}
            </Button>
          </form>

          {/* Back to Login */}
          <div className="mt-8 text-center">
            <p className="text-[#8A8F8E]">Já tem uma conta?</p>
            <Link href="/login">
              <Button
                variant="outline"
                className="mt-2 w-full py-3 px-4 rounded-lg border-2 border-[#1F4E5F] text-[#1F4E5F] font-semibold hover:bg-[#1F4E5F] hover:text-white transition-colors duration-300"
              >
                Fazer Login
              </Button>
            </Link>
          </div>

          {/* Footer / Support */}
          <div className="mt-12 text-center text-xs text-[#8A8F8E]">
            &copy; 2024 Organiza Aqui. Todos os direitos reservados.
          </div>
        </div>
      </div>
    </div>
  )
}
