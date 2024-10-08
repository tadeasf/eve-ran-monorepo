'use client'

import { useState, useCallback, useEffect } from 'react'
import { useQuery } from 'react-query'
import CharacterTable from '../components/CharacterTable'
import FilterControls from '../components/FilterControls'
import TotalKillsChart from '../components/TotalKillsChart'
import TotalIskChart from '../components/TotalIskChart'
import { Region, CharacterStats, Character, ChartConfig, Kill } from '../../lib/types'
import { Skeleton } from "../components/ui/skeleton"
import { Progress } from "../components/ui/progress"
import Top10Killers from '../components/Top10Killers'
import Top10Points from '../components/Top10Points'

const fetchRegions = async (): Promise<Region[]> => {
  const response = await fetch('/api/regions')
  if (!response.ok) {
    throw new Error('Failed to fetch regions')
  }
  return response.json()
}

const getTwoWeeksAgoMonday = () => {
  const today = new Date()
  const twoWeeksAgo = new Date(today.setDate(today.getDate() - 14))
  const monday = new Date(twoWeeksAgo.setDate(twoWeeksAgo.getDate() - (twoWeeksAgo.getDay() + 6) % 7))
  return monday.toISOString().split('T')[0]
}

const getTodayDate = () => {
  return new Date().toISOString().split('T')[0]
}

export default function Dashboard() {
  const [characters, setCharacters] = useState<CharacterStats[]>([])
  const [selectedRegions, setSelectedRegions] = useState<Array<{ id: number, name: string }>>([])
  const [startDate, setStartDate] = useState<string>(getTwoWeeksAgoMonday())
  const [endDate, setEndDate] = useState<string>(getTodayDate())
  const [isLoading, setIsLoading] = useState(false)
  const [killsOverTime, setKillsOverTime] = useState<{ date: string; kills: number }[]>([])
  const [iskDestroyedOverTime, setIskDestroyedOverTime] = useState<{ date: string; isk: number }[]>([])
  const [allKills, setAllKills] = useState<Kill[]>([])

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

      const killsPromises = selectedRegions.map(async (region) => {
        const response = await fetch(`/api/kills/region/${region.id}?startDate=${startDate}&endDate=${endDate}`)
        if (!response.ok) {
          throw new Error(`Failed to fetch data for region ${region.id}`)
        }
        return response.json()
      })

      const regionKills = await Promise.all(killsPromises)
      const allKills = regionKills.flat().filter((kill: Kill) => {
        const killDate = new Date(kill.KillmailTime).toISOString().split('T')[0]
        return killDate >= startDate && killDate <= endDate
      })
      setAllKills(allKills)

      const characterStats = characterData.map((character) => {
        const characterKills = allKills.filter((kill: Kill) => kill.CharacterID === character.id)
        const killCount = characterKills.length
        const totalIsk = characterKills.reduce((sum, kill) => sum + kill.ZkillData.TotalValue, 0)
        const totalValue = characterKills.reduce((sum, kill) => sum + kill.ZkillData.TotalValue, 0)

        return {
          character_id: character.id,
          name: character.name,
          kill_count: killCount,
          total_isk: totalIsk,
          total_value: totalValue
        }
      })

      setCharacters(characterStats)

      const killsOverTime = allKills.reduce((acc, kill) => {
        const date = kill.KillmailTime.split('T')[0]
        acc[date] = (acc[date] || 0) + 1
        return acc
      }, {} as Record<string, number>)

      const iskDestroyedOverTime = allKills.reduce((acc, kill) => {
        const date = kill.KillmailTime.split('T')[0]
        acc[date] = (acc[date] || 0) + (kill.ZkillData.TotalValue || 0)
        return acc
      }, {} as Record<string, number>)

      setKillsOverTime(Object.entries(killsOverTime).map(([date, kills]) => ({ date, kills: typeof kills === 'number' ? kills : 0 })))
      setIskDestroyedOverTime(Object.entries(iskDestroyedOverTime).map(([date, isk]) => ({ date, isk: typeof isk === 'number' ? isk : 0 })))

      setIsLoading(false)
    } catch (error) {
      console.error('Failed to fetch character stats:', error)
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
                <CharacterTable 
                  characters={characters}
                  allKills={allKills}
                  startDate={startDate}
                  endDate={endDate}
                  selectedRegions={selectedRegions}
                />
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
                <Top10Killers
                  characters={characters}
                  startDate={startDate}
                  endDate={endDate}
                  chartConfig={chartConfig}
                />
                <Top10Points
                  kills={allKills}
                  characters={characters.map(char => ({ id: char.character_id, name: char.name }))}
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