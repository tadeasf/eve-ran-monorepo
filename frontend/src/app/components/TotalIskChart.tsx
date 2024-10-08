import { Area, AreaChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer } from "recharts"
import { TrendingUp, TrendingDown } from "lucide-react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "../components/ui/chart"
import { ChartConfig } from '../../lib/types'

interface TotalIskChartProps {
  iskDestroyedOverTime: { date: string; isk: number }[]
  startDate: string
  endDate: string
  chartConfig: ChartConfig
}

export default function TotalIskChart({ iskDestroyedOverTime, chartConfig }: TotalIskChartProps) {
  if (!iskDestroyedOverTime || iskDestroyedOverTime.length === 0) {
    return <p>No ISK data available for the selected period.</p>
  }

  const formatYAxis = (value: number) => {
    if (value >= 1e9) return `${(value / 1e9).toFixed(1)}B`
    if (value >= 1e6) return `${(value / 1e6).toFixed(1)}M`
    if (value >= 1e3) return `${(value / 1e3).toFixed(1)}K`
    return value.toString()
  }

  const calculateWeeklyTrend = () => {
    const twoWeeksAgo = new Date()
    twoWeeksAgo.setDate(twoWeeksAgo.getDate() - 14)
    
    const lastTwoWeeks = iskDestroyedOverTime.filter(day => new Date(day.date) >= twoWeeksAgo)
    const lastWeek = lastTwoWeeks.slice(-7)
    const previousWeek = lastTwoWeeks.slice(0, 7)

    const lastWeekTotal = lastWeek.reduce((sum, day) => sum + day.isk, 0)
    const previousWeekTotal = previousWeek.reduce((sum, day) => sum + day.isk, 0)

    const percentageChange = ((lastWeekTotal - previousWeekTotal) / previousWeekTotal) * 100
    return percentageChange.toFixed(1)
  }

  const trend = calculateWeeklyTrend()
  const isTrendingUp = parseFloat(trend) > 0

  const lastTwoWeeks = iskDestroyedOverTime.slice(-14)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Total ISK Destroyed</CardTitle>
        <CardDescription>
          Showing total ISK destroyed for all characters in selected regions (Last 2 weeks)
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart
              data={lastTwoWeeks}
              margin={{
                top: 10,
                right: 30,
                left: 0,
                bottom: 10,
              }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="date"
                tickFormatter={(value) => new Date(value).toLocaleDateString()}
              />
              <YAxis tickFormatter={formatYAxis} />
              <ChartTooltip
                content={<ChartTooltipContent indicator="line" />}
              />
              <Area
                type="monotone"
                dataKey="isk"
                stroke={chartConfig.isk.color}
                fill={chartConfig.isk.color}
              />
            </AreaChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
      <CardFooter>
        <div className="flex w-full items-start gap-2 text-sm">
          <div className="grid gap-2">
            <div className="flex items-center gap-2 font-medium leading-none">
              {isTrendingUp ? (
                <>Trending up by {trend}% this week <TrendingUp className="h-4 w-4" /></>
              ) : (
                <>Trending down by {Math.abs(parseFloat(trend))}% this week <TrendingDown className="h-4 w-4" /></>
              )}
            </div>
            <div className="flex items-center gap-2 leading-none text-muted-foreground">
              Last 2 weeks
            </div>
          </div>
        </div>
      </CardFooter>
    </Card>
  )
}