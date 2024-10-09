// Copyright (C) 2024 Tadeáš Fořt
// 
// This file is part of EVE Ran Services.
// 
// EVE Ran Services is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// EVE Ran Services is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with EVE Ran Services.  If not, see <https://www.gnu.org/licenses/>.

import { Area, AreaChart, CartesianGrid, XAxis, YAxis, ResponsiveContainer } from "recharts"
import { TrendingUp, TrendingDown } from "lucide-react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "../components/ui/chart"
import { ChartConfig } from '../../lib/types'

interface TotalKillsChartProps {
  killsOverTime: { date: string; kills: number }[]
  startDate: string
  endDate: string
  chartConfig: ChartConfig
}

export default function TotalKillsChart({ killsOverTime, chartConfig }: TotalKillsChartProps) {
  if (killsOverTime.length === 0) {
    return <p>No kill data available for the selected period.</p>
  }

  const calculateWeeklyTrend = () => {
    const twoWeeksAgo = new Date()
    twoWeeksAgo.setDate(twoWeeksAgo.getDate() - 14)
    
    const lastTwoWeeks = killsOverTime.filter(day => new Date(day.date) >= twoWeeksAgo)
    const lastWeek = lastTwoWeeks.slice(-7)
    const previousWeek = lastTwoWeeks.slice(0, 7)

    const lastWeekTotal = lastWeek.reduce((sum, day) => sum + day.kills, 0)
    const previousWeekTotal = previousWeek.reduce((sum, day) => sum + day.kills, 0)

    const percentageChange = ((lastWeekTotal - previousWeekTotal) / previousWeekTotal) * 100
    return percentageChange.toFixed(1)
  }

  const trend = calculateWeeklyTrend()
  const isTrendingUp = parseFloat(trend) > 0

  const lastTwoWeeks = killsOverTime.slice(-14)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Total Kills</CardTitle>
        <CardDescription>
          Showing total kills for all characters in selected regions (Last 2 weeks)
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