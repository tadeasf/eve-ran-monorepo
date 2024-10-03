import { Area, AreaChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer } from "recharts"
import { TrendingUp } from "lucide-react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "../components/ui/chart"
import { formatISK } from '../../lib/utils'
import { ChartConfig } from '../../lib/types'

interface TotalIskChartProps {
  iskDestroyedOverTime: { date: string; isk: number }[]
  startDate: string
  endDate: string
  chartConfig: ChartConfig
}

export default function TotalIskChart({ iskDestroyedOverTime, startDate, endDate, chartConfig }: TotalIskChartProps) {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Total ISK Destroyed Over Time</CardTitle>
        <CardDescription>
          Showing total ISK destroyed for all characters in selected regions
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart
              data={iskDestroyedOverTime}
              margin={{
                top: 10,
                right: 30,
                left: 0,
                bottom: 0,
              }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis
                dataKey="date"
                tickFormatter={(value) => new Date(value).toLocaleDateString()}
              />
              <YAxis tickFormatter={(value) => formatISK(value)} />
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
              Showing ISK destroyed trend <TrendingUp className="h-4 w-4" />
            </div>
            <div className="flex items-center gap-2 leading-none text-muted-foreground">
              {startDate} - {endDate}
            </div>
          </div>
        </div>
      </CardFooter>
    </Card>
  )
}