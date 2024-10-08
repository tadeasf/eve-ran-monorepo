"use client"

import { Bar, BarChart, XAxis, YAxis, ResponsiveContainer } from "recharts"
import { TrendingUp, TrendingDown } from "lucide-react"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/app/components/ui/card"
import { ChartContainer, ChartTooltip, ChartTooltipContent } from "@/app/components/ui/chart"
import { Kill, ChartConfig } from '../../lib/types'

// Add this new interface
interface SimpleCharacter {
  id: number;
  name: string;
}

interface Top10PointsProps {
  kills: Kill[]
  characters: SimpleCharacter[]  // Use the new SimpleCharacter interface here
  startDate: string
  endDate: string
  chartConfig: ChartConfig
}

export default function Top10Points({ kills, characters, startDate, endDate, chartConfig }: Top10PointsProps) {
  const characterPoints = kills.reduce((acc, kill) => {
    const characterId = kill.CharacterID
    const points = kill.ZkillData.Points
    acc[characterId] = (acc[characterId] || 0) + points
    return acc
  }, {} as Record<number, number>)

  const top10Points = Object.entries(characterPoints)
    .sort(([, a], [, b]) => b - a)
    .slice(0, 10)
    .map(([characterId, points]) => {
      const character = characters.find(char => char.id === parseInt(characterId))
      return { 
        name: character ? character.name : `Character ${characterId}`, 
        points 
      }
    })

  const calculateTrend = () => {
    if (top10Points.length === 0) return "0.0"
    
    const totalPoints = top10Points.reduce((sum, char) => sum + char.points, 0)
    const averagePoints = totalPoints / top10Points.length
    if (averagePoints === 0) return "0.0"
    
    const trend = ((top10Points[0].points - averagePoints) / averagePoints) * 100
    return trend.toFixed(1)
  }

  const trend = calculateTrend()
  const isTrendingUp = parseFloat(trend) > 0

  if (top10Points.length === 0) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Top 10 Points</CardTitle>
          <CardDescription>No data available for the selected period</CardDescription>
        </CardHeader>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Top 10 Points</CardTitle>
        <CardDescription>
          Showing top 10 characters with most points from {startDate} to {endDate}
        </CardDescription>
      </CardHeader>
      <CardContent>
        <ChartContainer config={chartConfig}>
          <ResponsiveContainer width="100%" height={300}>
            <BarChart data={top10Points} layout="vertical">
              <XAxis type="number" hide />
              <YAxis dataKey="name" type="category" width={100} />
              <ChartTooltip
                cursor={false}
                content={<ChartTooltipContent hideLabel />}
              />
              <Bar dataKey="points" fill={chartConfig.isk.color} radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartContainer>
      </CardContent>
      <CardFooter className="flex-col items-start gap-2 text-sm">
        <div className="flex gap-2 font-medium leading-none">
          {isTrendingUp ? (
            <>Top pointer trending up by {trend}% <TrendingUp className="h-4 w-4" /></>
          ) : (
            <>Top pointer trending down by {Math.abs(parseFloat(trend))}% <TrendingDown className="h-4 w-4" /></>
          )}
        </div>
        <div className="leading-none text-muted-foreground">
          Compared to the average of top 10 pointers
        </div>
      </CardFooter>
    </Card>
  )
}