'use client'

import { useState } from 'react'
import Link from 'next/link'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { Eye, EyeOff, User, Lock, LayoutDashboard } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Checkbox } from '@/components/ui/checkbox'
import { useLogin } from '@/hooks/mutations/useAuth'
import { loginSchema, type LoginFormData } from '@/schemas/authSchema'
import { useToast } from '@/hooks/useToast'
import Image from 'next/image'

export default function LoginPage() {
  const loginMutation = useLogin()
  const toast = useToast()
  const [showPassword, setShowPassword] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  })

  const onSubmit = async (data: LoginFormData) => {
    try {
      await loginMutation.mutateAsync(data)
      toast.success('Bem-vindo de volta! Redirecionando...')
      // O redirecionamento já é feito pelo hook useLogin no onSuccess
    } catch (error: unknown) {
      const err = error as { response?: { data?: { error?: string } }; message?: string }
      const errorMessage =
        err?.response?.data?.error || err?.message || 'Erro ao fazer login'
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
            Sua jornada para uma vida mais equilibrada e produtiva começa com
            um clique. Domine sua rotina, alcance seus objetivos.
          </p>
        </div>
      </div>

      {/* Lado Direito: Formulário de Login (30%) */}
      <div className="w-full md:w-[30%] bg-[#F4F6F5] flex flex-col justify-center items-center p-8 sm:p-12 z-30 shadow-2xl">
        <div className="w-full max-w-md">
          {/* Logo/Branding */}
          <div className="mb-10 text-center md:text-left">
            <div className="inline-flex items-center justify-center p-3 rounded-xl bg-[#1F4E5F] mb-4">
              <LayoutDashboard className="text-white w-8 h-8" />
            </div>
            <h2 className="text-3xl font-bold text-[#1F2933]">
              Bem-vindo de volta
            </h2>
            <p className="text-[#8A8F8E] mt-2">
              Acesse sua conta para gerenciar sua vida.
            </p>
          </div>

          {/* Form */}
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
            {/* Email ou Username */}
            <div>
              <Label
                htmlFor="identifier"
                className="block text-sm font-semibold text-[#1F2933] mb-2"
              >
                E-mail ou Usuário
              </Label>
              <div className="relative">
                <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-[#8A8F8E]">
                  <User className="w-5 h-5" />
                </span>
                <Input
                  id="identifier"
                  type="text"
                  placeholder="ex: joao@email.com ou joaosilva"
                  className={`w-full pl-10 pr-4 py-3 rounded-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-[#3FA7A0] focus:border-transparent transition bg-white ${
                    errors.identifier ? 'border-red-500' : ''
                  }`}
                  {...register('identifier')}
                />
              </div>
              {errors.identifier && (
                <p className="text-sm text-red-500 mt-1">
                  {errors.identifier.message}
                </p>
              )}
            </div>

            {/* Password */}
            <div>
              <div className="flex justify-between items-center mb-2">
                <Label
                  htmlFor="password"
                  className="block text-sm font-semibold text-[#1F2933]"
                >
                  Senha
                </Label>
                <Link
                  href="#"
                  className="text-sm font-medium text-[#3FA7A0] hover:underline"
                >
                  Esqueceu a senha?
                </Link>
              </div>
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
              {errors.password && (
                <p className="text-sm text-red-500 mt-1">
                  {errors.password.message}
                </p>
              )}
            </div>

            {/* Remember Me */}
            <div className="flex items-center space-x-2">
              <Checkbox id="remember" />
              <Label
                htmlFor="remember"
                className="text-sm text-[#1F2933] cursor-pointer"
              >
                Lembrar de mim
              </Label>
            </div>

            {/* Login Button */}
            <Button
              type="submit"
              disabled={loginMutation.isPending}
              className="w-full bg-[#3FA7A0] hover:bg-[#2E8F89] text-white font-bold py-3 px-4 rounded-lg shadow-lg transform active:scale-[0.98] transition-all"
            >
              {loginMutation.isPending ? 'Entrando...' : 'Entrar no Sistema'}
            </Button>
          </form>

          {/* Create Account */}
          <div className="mt-8 text-center">
            <p className="text-[#8A8F8E]">Ainda não tem uma conta?</p>
            <Link href="/register">
              <Button
                variant="outline"
                className="mt-2 w-full py-3 px-4 rounded-lg border-2 border-[#1F4E5F] text-[#1F4E5F] font-semibold hover:bg-[#1F4E5F] hover:text-white transition-colors duration-300"
              >
                Criar Nova Conta
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
