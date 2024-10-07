import { Area, AreaChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer } from "recharts"
import { TrendingUp } from "lucide-react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "../components/ui/chart"
import { ChartConfig } from '../../lib/types'

interface TotalKillsChartProps {
  killsOverTime: { date: string; kills: number }[]
  startDate: string
  endDate: string
  chartConfig: ChartConfig
}

export default function TotalKillsChart({ killsOverTime, startDate, endDate, chartConfig }: TotalKillsChartProps) {
  if (killsOverTime.length === 0) {
    return <p>No kill data available for the selected period.</p>
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Total Kills</CardTitle>
        <CardDescription>
          Showing total kills for all characters in selected regions
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <ResponsiveContainer width="100%" height={300}>
            <AreaChart
              data={killsOverTime}
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
              <YAxis />
              <ChartTooltip
                content={<ChartTooltipContent indicator="line" />}
              />
              <Area
                type="monotone"
                dataKey="kills"
                stroke={chartConfig.kills.color}
                fill={chartConfig.kills.color}
              />
            </AreaChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
      <CardFooter>
        <div className="flex w-full items-start gap-2 text-sm">
          <div className="grid gap-2">
            <div className="flex items-center gap-2 font-medium leading-none">
              Showing kills trend <TrendingUp className="h-4 w-4" />
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