export type PasswordStrength = 'muito-fraco' | 'fraco' | 'medio' | 'forte' | 'muito-forte'

export interface PasswordStrengthResult {
  strength: PasswordStrength
  score: number // 0-100
  feedback: string[]
}

export function calculatePasswordStrength(password: string): PasswordStrengthResult {
  if (!password) {
    return {
      strength: 'muito-fraco',
      score: 0,
      feedback: [],
    }
  }

  let score = 0
  const feedback: string[] = []

  // Critérios de força
  const hasMinLength = password.length >= 8
  const hasMaxLength = password.length >= 12
  const hasLowercase = /[a-z]/.test(password)
  const hasUppercase = /[A-Z]/.test(password)
  const hasNumbers = /[0-9]/.test(password)
  const hasSpecialChars = /[^a-zA-Z0-9]/.test(password)
  const hasNoRepeating = !/(.)\1{2,}/.test(password)
  const hasNoCommonPatterns = !/(123|abc|qwe|password|admin)/i.test(password)

  // Pontuação
  if (hasMinLength) score += 15
  if (hasMaxLength) score += 10
  if (hasLowercase) score += 10
  if (hasUppercase) score += 15
  if (hasNumbers) score += 15
  if (hasSpecialChars) score += 20
  if (hasNoRepeating) score += 10
  if (hasNoCommonPatterns) score += 5

  // Feedback
  if (!hasMinLength) feedback.push('Use pelo menos 8 caracteres')
  if (!hasLowercase) feedback.push('Adicione letras minúsculas')
  if (!hasUppercase) feedback.push('Adicione letras maiúsculas')
  if (!hasNumbers) feedback.push('Adicione números')
  if (!hasSpecialChars) feedback.push('Adicione caracteres especiais (!@#$%...)')

  // Determinar força
  let strength: PasswordStrength
  if (score < 30) {
    strength = 'muito-fraco'
  } else if (score < 50) {
    strength = 'fraco'
  } else if (score < 70) {
    strength = 'medio'
  } else if (score < 85) {
    strength = 'forte'
  } else {
    strength = 'muito-forte'
  }

  return { strength, score, feedback }
}

export function getPasswordStrengthColor(strength: PasswordStrength): string {
  const colors = {
    'muito-fraco': 'bg-red-500',
    'fraco': 'bg-orange-500',
    'medio': 'bg-yellow-500',
    'forte': 'bg-green-500',
    'muito-forte': 'bg-[#3FA7A0]', // Cor primária do projeto
  }
  return colors[strength]
}

export function getPasswordStrengthText(strength: PasswordStrength): string {
  const texts = {
    'muito-fraco': 'Muito Fraco',
    'fraco': 'Fraco',
    'medio': 'Médio',
    'forte': 'Forte',
    'muito-forte': 'Muito Forte',
  }
  return texts[strength]
}
