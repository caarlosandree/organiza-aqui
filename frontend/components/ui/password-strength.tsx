'use client'

import {
  calculatePasswordStrength,
  getPasswordStrengthColor,
  getPasswordStrengthText,
  type PasswordStrength,
} from '@/utils/passwordStrength'

interface PasswordStrengthIndicatorProps {
  password: string
}

export function PasswordStrengthIndicator({
  password,
}: PasswordStrengthIndicatorProps) {
  if (!password) return null

  const { strength, score, feedback } = calculatePasswordStrength(password)
  const color = getPasswordStrengthColor(strength)
  const text = getPasswordStrengthText(strength)

  // Calcular largura da barra baseada no score (0-100)
  const widthPercentage = Math.min(score, 100)

  return (
    <div className="mt-2 space-y-2">
      {/* Barra de força */}
      <div className="flex items-center gap-2">
        <div className="flex-1 h-2 bg-gray-200 rounded-full overflow-hidden">
          <div
            className={`h-full transition-all duration-300 ${color}`}
            style={{ width: `${widthPercentage}%` }}
          />
        </div>
        <span className={`text-xs font-semibold ${getTextColor(strength)}`}>
          {text}
        </span>
      </div>

      {/* Feedback */}
      {feedback.length > 0 && (
        <ul className="text-xs text-[#8A8F8E] space-y-1">
          {feedback.map((item, index) => (
            <li key={index} className="flex items-center gap-1">
              <span className="text-red-500">•</span>
              {item}
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}

function getTextColor(strength: PasswordStrength): string {
  const colors = {
    'muito-fraco': 'text-red-500',
    'fraco': 'text-orange-500',
    'medio': 'text-yellow-500',
    'forte': 'text-green-500',
    'muito-forte': 'text-[#3FA7A0]',
  }
  return colors[strength]
}
