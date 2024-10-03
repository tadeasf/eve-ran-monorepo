'use client'

import { useState, useCallback, useEffect } from 'react'
import { useQuery } from 'react-query'
import CharacterTable from '../components/CharacterTable'
import FilterControls from '../components/FilterControls'
import TotalKillsChart from '../components/TotalKillsChart'
import TotalIskChart from '../components/TotalIskChart'
import { Region, CharacterStats, Character, KillmailData, ChartConfig } from '../../lib/types'
import { formatISK } from '../../lib/utils'
import { Alert, AlertDescription, AlertTitle } from "../components/ui/alert"
import { Skeleton } from "../components/ui/skeleton"
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
  const [iskDestroyedOverTime, setIskDestroyedOverTime] = useState<{ date: string; isk: number }[]>([])

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
      const iskPerDay: { [date: string]: number } = {}
      combinedStats.forEach((character) => {
        character.kills.forEach((kill) => {
          const date = kill.killmail_time.split('T')[0]
          killsPerDay[date] = (killsPerDay[date] || 0) + 1
          iskPerDay[date] = (iskPerDay[date] || 0) + kill.total_value
        })
      })

      const sortedKillsOverTime = Object.entries(killsPerDay)
        .map(([date, kills]) => ({ date, kills }))
        .sort((a, b) => a.date.localeCompare(b.date))

      const sortedIskDestroyedOverTime = Object.entries(iskPerDay)
        .map(([date, isk]) => ({ date, isk }))
        .sort((a, b) => a.date.localeCompare(b.date))

      setKillsOverTime(sortedKillsOverTime)
      setIskDestroyedOverTime(sortedIskDestroyedOverTime)
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

  const chartConfig: ChartConfig = {
    kills: {
      label: "Total Kills",
      color: "hsl(var(--chart-1))",
    },
    isk: {
      label: "Total ISK Destroyed",
      color: "hsl(var(--chart-2))",
    },
  }

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
              <div className="grid grid-cols-1 md:grid-cols-2 gap-8 mb-8">
                <TotalKillsChart
                  killsOverTime={killsOverTime}
                  startDate={startDate}
                  endDate={endDate}
                  chartConfig={chartConfig}
                />
                <TotalIskChart
                  iskDestroyedOverTime={iskDestroyedOverTime}
                  startDate={startDate}
                  endDate={endDate}
                  chartConfig={chartConfig}
                />
              </div>
            </>
          )}
        </>
      )}
    </div>
  )
}