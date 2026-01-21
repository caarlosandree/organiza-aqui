/**
 * Tema Centralizado "Calma Moderna"
 * 
 * Este arquivo contém as constantes e utilitários para o tema do sistema,
 * que transmite calma, confiança, clareza e equilíbrio emocional.
 */

/**
 * Paleta de cores em formato hexadecimal
 */
export const themeColors = {
  // Cores principais
  primary: '#3FA7A0', // Verde Azulado Suave
  primaryHover: '#2E8F89', // Hover do primário
  secondary: '#1F4E5F', // Azul Profundo Calmante
  accent: '#A8DADC', // Verde Claro Orgânico
  
  // Backgrounds e superfícies
  background: '#F4F6F5', // Cinza Quente Claro
  card: '#FFFFFF', // Branco para cards
  
  // Textos
  textPrimary: '#1F2933', // Texto principal
  textSecondary: '#8A8F8E', // Cinza Médio Neutro
  
  // Estados
  success: '#5FB3A2', // Sucesso calmo
  warning: '#E9C46A', // Alerta suave
  error: '#E76F51', // Erro controlado
} as const

/**
 * Paleta de cores em formato OKLCH (para uso em CSS)
 */
export const themeColorsOKLCH = {
  primary: 'oklch(0.67 0.08 186)', // #3FA7A0
  primaryHover: 'oklch(0.62 0.08 186)', // #2E8F89
  secondary: 'oklch(0.35 0.06 225)', // #1F4E5F
  accent: 'oklch(0.85 0.05 186)', // #A8DADC
  background: 'oklch(0.965 0.003 160)', // #F4F6F5
  card: 'oklch(1 0 0)', // #FFFFFF
  textPrimary: 'oklch(0.18 0.005 260)', // #1F2933
  textSecondary: 'oklch(0.55 0.003 160)', // #8A8F8E
  success: 'oklch(0.71 0.09 182)', // #5FB3A2
  warning: 'oklch(0.80 0.12 75)', // #E9C46A
  error: 'oklch(0.66 0.18 30)', // #E76F51
} as const

/**
 * Mapeamento de uso das cores por elemento da interface
 */
export const themeUsage = {
  background: themeColors.background,
  header: themeColors.secondary,
  sidebar: themeColors.secondary,
  card: themeColors.card,
  buttonPrimary: themeColors.primary,
  buttonPrimaryHover: themeColors.primaryHover,
  textPrimary: themeColors.textPrimary,
  textSecondary: themeColors.textSecondary,
  success: themeColors.success,
  warning: themeColors.warning,
  error: themeColors.error,
} as const

/**
 * Tipo TypeScript para as cores do tema
 */
export type ThemeColor = typeof themeColors[keyof typeof themeColors]

/**
 * Tipo TypeScript para uso das cores
 */
export type ThemeUsage = keyof typeof themeUsage

/**
 * Helper para obter uma cor do tema por nome
 */
export function getThemeColor(colorName: keyof typeof themeColors): string {
  return themeColors[colorName]
}

/**
 * Helper para obter uma cor do tema em formato OKLCH
 */
export function getThemeColorOKLCH(colorName: keyof typeof themeColorsOKLCH): string {
  return themeColorsOKLCH[colorName]
}

/**
 * Helper para obter cor de uso específico
 */
export function getThemeUsage(usage: ThemeUsage): string {
  return themeUsage[usage]
}

/**
 * Mapeamento de cores de prioridade (para tarefas)
 * Usa variações da paleta Calma Moderna
 */
export const priorityColors = {
  low: {
    background: 'bg-accent/20', // Verde Claro Orgânico suave
    text: 'text-accent-foreground',
    border: 'border-accent/30',
  },
  medium: {
    background: 'bg-warning/20', // Alerta suave
    text: 'text-warning-foreground',
    border: 'border-warning/30',
  },
  high: {
    background: 'bg-primary/20', // Verde Azulado suave
    text: 'text-primary-foreground',
    border: 'border-primary/30',
  },
  urgent: {
    background: 'bg-destructive/20', // Erro controlado
    text: 'text-destructive',
    border: 'border-destructive/30',
  },
} as const

/**
 * Tipo para prioridades
 */
export type Priority = keyof typeof priorityColors

/**
 * Helper para obter classes CSS de prioridade
 */
export function getPriorityClasses(priority: Priority) {
  return priorityColors[priority]
}