export interface Currency {
  code: string
  name: string
  symbol: string
}

export const currencies: Currency[] = [
  { code: 'BRL', name: 'Real Brasileiro', symbol: 'R$' },
  { code: 'USD', name: 'Dólar Americano', symbol: '$' },
  { code: 'EUR', name: 'Euro', symbol: '€' },
  { code: 'GBP', name: 'Libra Esterlina', symbol: '£' },
  { code: 'JPY', name: 'Iene Japonês', symbol: '¥' },
  { code: 'CAD', name: 'Dólar Canadense', symbol: 'C$' },
  { code: 'AUD', name: 'Dólar Australiano', symbol: 'A$' },
  { code: 'CHF', name: 'Franco Suíço', symbol: 'CHF' },
  { code: 'CNY', name: 'Yuan Chinês', symbol: '¥' },
  { code: 'ARS', name: 'Peso Argentino', symbol: '$' },
]

export const getCurrencyByCode = (code: string): Currency | undefined => {
  return currencies.find((currency) => currency.code === code)
}

export const getDefaultCurrency = (): Currency => {
  return currencies[0] // BRL
}
