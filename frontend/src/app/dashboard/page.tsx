'use client'

import { useState, useCallback, useEffect } from 'react'
import { useQuery } from 'react-query'
import { TrendingUp } from "lucide-react"
import { Area, AreaChart, CartesianGrid, XAxis, ResponsiveContainer } from "recharts"
import CharacterTable from '../components/CharacterTable'
import FilterControls from '../components/FilterControls'
import { Region, CharacterStats, Character, KillmailData } from '../../lib/types'
import { formatISK } from '../../lib/utils'
import { Alert, AlertDescription, AlertTitle } from "../components/ui/alert"
import { Skeleton } from "../components/ui/skeleton"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../components/ui/card"
import { ChartConfig, ChartContainer, ChartTooltip, ChartTooltipContent } from "../components/ui/chart"
import { Progress } from "../components/ui/progress"

const fetchRegions = async (): Promise<Region[]> => {
  const response = await fetch('/api/regions')
  if (!response.ok) {
    throw new Error('Failed to fetch regions')
  }
  return response.json()
}

const getLastMonday = () => {
  const today = new Date()
  const lastMonday = new Date(today.setDate(today.getDate() - (today.getDay() + 6) % 7))
  return lastMonday.toISOString().split('T')[0]
}

const getTodayDate = () => {
  return new Date().toISOString().split('T')[0]
}

export default function Dashboard() {
  const [characters, setCharacters] = useState<CharacterStats[]>([])
  const [selectedRegions, setSelectedRegions] = useState<Array<{ id: number, name: string }>>([])
  const [startDate, setStartDate] = useState<string>(getLastMonday())
  const [endDate, setEndDate] = useState<string>(getTodayDate())
  const [isLoading, setIsLoading] = useState(false)
  const [killsOverTime, setKillsOverTime] = useState<{ date: string; kills: number }[]>([])

  const { data: regions, isLoading: isRegionsLoading, error: regionsError } = useQuery<Region[]>('regions', fetchRegions)

  useEffect(() => {
    if (regions) {
      const placid = regions.find(r => r.name === 'Placid')
      const syndicate = regions.find(r => r.name === 'Syndicate')
      if (placid && syndicate) {
        setSelectedRegions([
          { id: placid.region_id, name: placid.name },
          { id: syndicate.region_id, name: syndicate.name }
        ])
      }
    }
  }, [regions])

  const fetchCharacterStats = useCallback(async () => {
    setIsLoading(true)
    try {
      const characterResponse = await fetch('/api/characters')
      const characterData: Character[] = await characterResponse.json()

      const statsPromises = selectedRegions.map(async (region) => {
        const response = await fetch(`/api/kills/region/${region.id}?startDate=${startDate}&endDate=${endDate}`)
        if (!response.ok) {
          throw new Error(`Failed to fetch data for region ${region.id}`)
        }
        return response.json()
      })

      const regionStats = await Promise.all(statsPromises)
      const startDateTime = new Date(startDate).getTime()
      const endDateTime = new Date(endDate).getTime()

      const combinedStats = characterData.map((character: Character) => {
        const stats = regionStats.flatMap((regionStat: { data: KillmailData[] }) => 
          (regionStat.data || []).filter((stat: KillmailData) => 
            stat.character_id === character.id &&
            new Date(stat.killmail_time).getTime() >= startDateTime &&
            new Date(stat.killmail_time).getTime() <= endDateTime
          )
        )
        const totalValue = stats.reduce((sum: number, stat: KillmailData) => sum + stat.total_value, 0)
        return {
          character_id: character.id,
          name: character.name,
          kill_count: stats.length,
          total_value: totalValue,
          formatted_total_value: formatISK(totalValue),
          kills: stats,
        }
      })

      const killsPerDay: { [date: string]: number } = {}
      combinedStats.forEach((character) => {
        character.kills.forEach((kill) => {
          const date = kill.killmail_time.split('T')[0]
          killsPerDay[date] = (killsPerDay[date] || 0) + 1
        })
      })

      const sortedKillsOverTime = Object.entries(killsPerDay)
        .map(([date, kills]) => ({ date, kills }))
        .sort((a, b) => a.date.localeCompare(b.date))

      setKillsOverTime(sortedKillsOverTime)
      setCharacters(combinedStats)
    } catch (error) {
      console.error('Failed to fetch character stats:', error)
    } finally {
      setIsLoading(false)
    }
  }, [selectedRegions, startDate, endDate])

  useEffect(() => {
    if (selectedRegions.length > 0) {
      fetchCharacterStats()
    }
  }, [selectedRegions, fetchCharacterStats])

  if (regionsError) {
    return <div>Error loading regions: {(regionsError as Error).message}</div>
  }

  const chartConfig = {
    kills: {
      label: "Total Kills",
      color: "hsl(var(--chart-1))",
    },
  } satisfies ChartConfig

  return (
    <div className="container mx-auto px-4 py-8">
      <Alert className="mb-8">
        <AlertTitle>Welcome to the EVE Ran Dashboard</AlertTitle>
        <AlertDescription>
          Use the filters below to narrow down your character statistics.
        </AlertDescription>
      </Alert>
      {isRegionsLoading ? (
        <Progress value={33} className="w-full mb-8" />
      ) : (
        <>
          <FilterControls
            regions={regions || []}
            selectedRegions={selectedRegions}
            setSelectedRegions={setSelectedRegions}
            startDate={startDate}
            setStartDate={setStartDate}
            endDate={endDate}
            setEndDate={setEndDate}
            onApplyFilters={fetchCharacterStats}
            isLoading={isLoading}
          />
          {isLoading ? (
            <Skeleton className="w-full h-[400px] mb-8" />
          ) : (
            <>
              <div className="mb-8">
                <CharacterTable characters={characters} />
              </div>
              <Card className="mb-8">
                <CardHeader>
                  <CardTitle>Total Kills Over Time</CardTitle>
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
                          bottom: 0,
                        }}
                      >
                        <CartesianGrid strokeDasharray="3 3" />
                        <XAxis
                          dataKey="date"
                          tickFormatter={(value) => new Date(value).toLocaleDateString()}
                        />
                        <ChartTooltip
                          content={<ChartTooltipContent indicator="line" />}
                        />
                        <Area
                          type="monotone"
                          dataKey="kills"
                          stroke="#8884d8"
                          fill="#8884d8"
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
            </>
          )}
        </>
      )}
    </div>
  )
}