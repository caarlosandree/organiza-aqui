/**
 * Converte centavos para formato de moeda
 * @param cents Valor em centavos
 * @param currency Código da moeda (padrão: BRL)
 * @returns String formatada (ex: "R$ 100,50")
 */
export function formatCurrency(cents: number, currency: string = 'BRL'): string {
  const value = cents / 100
  return new Intl.NumberFormat('pt-BR', {
    style: 'currency',
    currency: currency,
  }).format(value)
}

/**
 * Converte valor de moeda para centavos
 * @param value Valor em formato de moeda (ex: "100.50")
 * @returns Valor em centavos
 */
export function parseCurrencyToCents(value: string): number {
  // Remove tudo exceto números e vírgula/ponto
  const cleaned = value.replace(/[^\d,.-]/g, '')
  // Substitui vírgula por ponto
  const normalized = cleaned.replace(',', '.')
  const floatValue = parseFloat(normalized) || 0
  return Math.round(floatValue * 100)
}
