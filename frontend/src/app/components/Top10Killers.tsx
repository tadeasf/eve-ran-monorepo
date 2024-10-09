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

"use client"

import { Bar, BarChart, XAxis, YAxis, ResponsiveContainer } from "recharts"
import { TrendingUp, TrendingDown } from "lucide-react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/app/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/app/components/ui/chart"
import { CharacterStats } from '../../lib/types'
import { ChartConfig } from '../../lib/types'

interface Top10KillersProps {
  characters: CharacterStats[]
  startDate: string
  endDate: string
  chartConfig: ChartConfig
}

export default function Top10Killers({ characters, startDate, endDate, chartConfig }: Top10KillersProps) {
  const top10Killers = characters
    .sort((a, b) => b.kill_count - a.kill_count)
    .slice(0, 10)
    .map(char => ({ name: char.name, kills: char.kill_count }))

  const calculateTrend = () => {
    if (top10Killers.length === 0) return "0.0"
    
    const totalKills = top10Killers.reduce((sum, char) => sum + char.kills, 0)
    const averageKills = totalKills / top10Killers.length
    if (averageKills === 0) return "0.0"
    
    const trend = ((top10Killers[0].kills - averageKills) / averageKills) * 100
    return trend.toFixed(1)
  }

  const trend = calculateTrend()
  const isTrendingUp = parseFloat(trend) > 0

  if (top10Killers.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Top 10 Killers</CardTitle>
          <CardDescription>No data available for the selected period</CardDescription>
        </CardHeader>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Top 10 Killers</CardTitle>
        <CardDescription>
          Showing top 10 characters with most kills from {startDate} to {endDate}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={top10Killers} layout="vertical">
              <XAxis type="number" hide />
              <YAxis dataKey="name" type="category" width={100} />
              <ChartTooltip
                cursor={false}
                content={<ChartTooltipContent hideLabel />}
              />
              <Bar dataKey="kills" fill={chartConfig.kills.color} radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 font-medium leading-none">
          {isTrendingUp ? (
            <>Top killer trending up by {trend}% <TrendingUp className="h-4 w-4" /></>
          ) : (
            <>Top killer trending down by {Math.abs(parseFloat(trend))}% <TrendingDown className="h-4 w-4" /></>
          )}
        </div>
        <div className="leading-none text-muted-foreground">
          Compared to the average of top 10 killers
        </div>
      </CardFooter>
    </Card>
  )
}