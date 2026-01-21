'use client'

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts'
import { formatCurrency } from '@/utils/currency'
import { useSpendingByTag } from '@/hooks/queries/useAnalytics'
import { usePrivacyStore } from '@/stores/privacyStore'
import { startOfMonth, endOfMonth, format } from 'date-fns'

interface ChartDataItem {
  tag: string
  total: number
  count: number
  percentage: number
}

interface CustomTooltipProps {
  active?: boolean
  payload?: Array<{ payload: ChartDataItem }>
  isPrivacyMode: boolean
}

function CustomTooltip({ active, payload, isPrivacyMode }: CustomTooltipProps) {
  if (active && payload && payload.length) {
    const data = payload[0].payload as ChartDataItem
    return (
      <div className="bg-background border rounded-lg p-3 shadow-lg">
        <p className="font-medium">{data.tag}</p>
        <p className="text-sm text-muted-foreground">
          {isPrivacyMode ? '••••' : formatCurrency(data.total * 100, 'BRL')}
        </p>
        <p className="text-xs text-muted-foreground">
          {data.count} {data.count === 1 ? 'transação' : 'transações'}
        </p>
        {data.percentage > 0 && (
          <p className="text-xs text-muted-foreground">
            {data.percentage.toFixed(1)}% do total
          </p>
        )}
      </div>
    )
  }
  return null
}

interface SpendingByTagChartProps {
  startDate?: string
  endDate?: string
}

export function SpendingByTagChart({ startDate, endDate }: SpendingByTagChartProps) {
  const { isPrivacyMode } = usePrivacyStore()
  const defaultStartDate = startDate || format(startOfMonth(new Date()), 'yyyy-MM-dd')
  const defaultEndDate = endDate || format(endOfMonth(new Date()), 'yyyy-MM-dd')

  const { data: spendingByTag, isLoading } = useSpendingByTag(defaultStartDate, defaultEndDate)

  if (isLoading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Gastos por Tag</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center p-4">Carregando...</div>
        </CardContent>
      </Card>
    )
  }

  if (!spendingByTag || spendingByTag.tags.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Gastos por Tag</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-center p-4 text-muted-foreground">
            Nenhum gasto encontrado para o período selecionado
          </div>
        </CardContent>
      </Card>
    )
  }

  const chartData: ChartDataItem[] = spendingByTag.tags.map((tag) => ({
    tag: tag.tag || 'Sem tag',
    total: isPrivacyMode ? 0 : tag.total_amount / 100, // Converter centavos para reais
    count: tag.transaction_count,
    percentage: tag.percentage || 0,
  }))

  return (
    <Card>
      <CardHeader>
        <CardTitle>Gastos por Tag</CardTitle>
      </CardHeader>
      <CardContent>
        <ResponsiveContainer width="100%" height={350}>
          <BarChart data={chartData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="tag"
              angle={-45}
              textAnchor="end"
              height={100}
              interval={0}
              tick={{ fontSize: 12 }}
            />
            <YAxis
              tickFormatter={(value) => (isPrivacyMode ? '•••' : formatCurrency(value * 100, 'BRL'))}
              tick={{ fontSize: 12 }}
            />
            <Tooltip content={<CustomTooltip isPrivacyMode={isPrivacyMode} />} />
            <Legend />
            <Bar
              dataKey="total"
              fill="#3b82f6"
              name="Total Gasto"
              radius={[4, 4, 0, 0]}
            />
          </BarChart>
        </ResponsiveContainer>
      </CardContent>
    </Card>
  )
}
