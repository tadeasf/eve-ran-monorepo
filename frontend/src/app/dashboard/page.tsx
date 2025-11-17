'use client'

import { useState, useCallback, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import Image from 'next/image'
import FilterControls from '../components/FilterControls'
import TotalKillsChart from '../components/TotalKillsChart'
import TotalIskChart from '../components/TotalIskChart'
import { Region, CharacterStats, Character, ChartConfig, Kill } from '../../lib/types'
import { Skeleton } from "@/components/ui/skeleton"
import { Progress } from "@/components/ui/progress"
import Top10Killers from '../components/Top10Killers'
import Top10Points from '../components/Top10Points'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import {
    Table,
    TableBody,
    TableCaption,
    TableCell,
    TableHead,
    TableHeader,
    TableRow,
} from '@/components/ui/table'
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from '@/components/ui/popover'
import { Badge } from '@/components/ui/badge'
import { BarChart3, ArrowUpDown } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { formatISK } from '../../lib/utils'
import { apiFetchJson } from '@/lib/api'

const fetchRegions = async (): Promise<Region[]> => {
  return apiFetchJson<Region[]>('/regions')
}

const getDaysAgoDate = (daysAgo: number): string => {
  const date = new Date()
  date.setDate(date.getDate() - daysAgo)
  return date.toISOString().split('T')[0]
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

// Load dashboard settings from localStorage
const loadDashboardSettings = () => {
  if (typeof window !== 'undefined') {
    try {
      const saved = localStorage.getItem('dashboard_settings')
      if (saved) {
        return JSON.parse(saved)
      }
    } catch (e) {
      console.error('Failed to parse dashboard settings:', e)
    }
  }
  return null
}

export default function Dashboard() {
  const [characters, setCharacters] = useState<CharacterStats[]>([])
  const [selectedRegions, setSelectedRegions] = useState<Array<{ id: number, name: string }>>([])
  const [settingsLoaded, setSettingsLoaded] = useState(false)

  // Initialize dates with settings or defaults
  const [startDate, setStartDate] = useState<string>(() => {
    const settings = loadDashboardSettings()
    if (settings && settings.defaultStartDate) {
      return getDaysAgoDate(parseInt(settings.defaultStartDate))
    }
    return getTwoWeeksAgoMonday()
  })

  const [endDate, setEndDate] = useState<string>(() => {
    const settings = loadDashboardSettings()
    if (settings && settings.defaultEndDate) {
      return getDaysAgoDate(parseInt(settings.defaultEndDate))
    }
    return getTodayDate()
  })

  const [isLoading, setIsLoading] = useState(false)
  const [killsOverTime, setKillsOverTime] = useState<{ date: string; kills: number }[]>([])
  const [iskDestroyedOverTime, setIskDestroyedOverTime] = useState<{ date: string; isk: number }[]>([])
  const [allKills, setAllKills] = useState<Kill[]>([])
  const [sortColumn, setSortColumn] = useState<keyof CharacterStats>('kill_count')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')
  const [systems, setSystems] = useState<Map<number, string>>(new Map())

  const { data: regions, isLoading: isRegionsLoading, error: regionsError } = useQuery({
    queryKey: ['regions'],
    queryFn: fetchRegions,
  })

  useEffect(() => {
    if (regions && regions.length > 0 && !settingsLoaded) {
      // Load saved settings or use defaults
      const settings = loadDashboardSettings()

      if (settings && settings.defaultRegions && settings.defaultRegions.length > 0) {
        // Use saved region preferences
        const savedRegions = regions
          .filter(region => settings.defaultRegions.includes(region.region_id))
          .map(region => ({
            id: region.region_id,
            name: region.name
          }))
        setSelectedRegions(savedRegions)
      } else {
        // Select all regions by default
        setSelectedRegions(regions.map(region => ({
          id: region.region_id,
          name: region.name
        })))
      }

      setSettingsLoaded(true)
    }
  }, [regions, settingsLoaded])

  // Fetch systems data
  useEffect(() => {
    const fetchSystems = async () => {
      try {
        const response = await fetch('/api/systems')
        if (!response.ok) throw new Error('Failed to fetch systems')
        const systemsData = await response.json()
        const systemsMap = new Map<number, string>()
        systemsData.forEach((system: { system_id: number; name: string }) => {
          systemsMap.set(system.system_id, system.name)
        })
        setSystems(systemsMap)
      } catch (error) {
        console.error('Error fetching systems:', error)
      }
    }
    fetchSystems()
  }, [])

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

  const handleSort = (column: keyof CharacterStats) => {
    if (column === sortColumn) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortColumn(column)
      setSortDirection('desc')
    }
  }

  const sortedCharacters = [...characters].sort((a, b) => {
    if (a[sortColumn] < b[sortColumn]) return sortDirection === 'asc' ? -1 : 1
    if (a[sortColumn] > b[sortColumn]) return sortDirection === 'asc' ? 1 : -1
    return 0
  })

  const SortableHeader = ({ column, children }: { column: keyof CharacterStats, children: React.ReactNode }) => (
    <TableHead onClick={() => handleSort(column)} className="cursor-pointer">
      <div className="flex items-center">
        {children}
        <ArrowUpDown className="ml-2 h-4 w-4" />
      </div>
    </TableHead>
  )

  return (
    <div className="container mx-auto px-4 py-8 space-y-8">
      {/* Header */}
      <div className="space-y-2">
        <h1 className="text-4xl font-bold flex items-center gap-3">
          <BarChart3 className="size-10 text-primary" />
          PVP Performance Dashboard
        </h1>
        <p className="text-muted-foreground text-lg">
          Track killmail statistics and combat performance for all characters
        </p>
      </div>

      {isRegionsLoading ? (
        <Progress value={33} className="w-full" />
      ) : (
        <>
          {/* Filters Section */}
          <Card>
            <CardHeader>
              <CardTitle>Filters</CardTitle>
              <CardDescription>
                Select regions and date range to view killmail data
              </CardDescription>
            </CardHeader>
            <CardContent>
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
            </CardContent>
          </Card>

          {isLoading ? (
            <Skeleton className="w-full h-[400px]" />
          ) : (
            <>
              {/* Charts Section */}
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
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

              {/* Character Table Section */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-2xl">Character Statistics</CardTitle>
                  <CardDescription>
                    Detailed performance metrics for all tracked characters
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <Table>
                    <TableCaption>
                      {characters.length > 0 ? `Showing ${characters.length} character${characters.length > 1 ? 's' : ''}` : 'No characters found'}
                    </TableCaption>
                    <TableHeader>
                      <TableRow>
                        <SortableHeader column="name">Character</SortableHeader>
                        <SortableHeader column="kill_count">Kills</SortableHeader>
                        <SortableHeader column="total_isk">ISK Destroyed</SortableHeader>
                        <TableHead className="text-right">Actions</TableHead>
                      </TableRow>
                    </TableHeader>
                    <TableBody>
                      {sortedCharacters.length > 0 ? (
                        sortedCharacters.map((character) => {
                          const characterKills = allKills.filter(kill => kill.CharacterID === character.character_id)
                          return (
                            <TableRow key={character.character_id}>
                              <TableCell className="font-medium">
                                <div className="flex items-center gap-2">
                                  <Image
                                    src={`https://images.evetech.net/characters/${character.character_id}/portrait?size=32`}
                                    alt={character.name}
                                    width={32}
                                    height={32}
                                    className="size-8 rounded-full"
                                    unoptimized
                                  />
                                  {character.name}
                                </div>
                              </TableCell>
                              <TableCell>{character.kill_count.toLocaleString()}</TableCell>
                              <TableCell className="font-mono">{formatISK(character.total_isk)}</TableCell>
                              <TableCell className="text-right">
                                <Popover>
                                  <PopoverTrigger asChild>
                                    <Button variant="outline" size="sm">
                                      View Kills
                                    </Button>
                                  </PopoverTrigger>
                                  <PopoverContent className="w-[700px] max-h-[600px] overflow-auto" align="end">
                                    <div className="space-y-4">
                                      <div>
                                        <h3 className="font-semibold text-lg">{character.name} - Top Kills</h3>
                                        <p className="text-sm text-muted-foreground">
                                          {characterKills.length > 10 
                                            ? `Top 10 most valuable kills (${characterKills.length} total)`
                                            : `${characterKills.length} kill${characterKills.length !== 1 ? 's' : ''} in selected period`
                                          }
                                        </p>
                                      </div>
                                      {characterKills.length > 0 ? (
                                        <Table>
                                          <TableHeader>
                                            <TableRow>
                                              <TableHead>Date</TableHead>
                                              <TableHead>System</TableHead>
                                              <TableHead className="text-right">ISK Value</TableHead>
                                              <TableHead className="text-right">Points</TableHead>
                                            </TableRow>
                                          </TableHeader>
                                          <TableBody>
                                            {characterKills
                                              .sort((a, b) => b.ZkillData.TotalValue - a.ZkillData.TotalValue)
                                              .slice(0, 10)
                                              .map((kill) => (
                                                <TableRow key={kill.KillmailID}>
                                                  <TableCell className="text-sm">
                                                    <a
                                                      href={`https://zkillboard.com/kill/${kill.KillmailID}/`}
                                                      target="_blank"
                                                      rel="noopener noreferrer"
                                                      className="text-primary hover:underline cursor-pointer"
                                                    >
                                                      {new Date(kill.KillmailTime).toLocaleDateString()}
                                                    </a>
                                                  </TableCell>
                                                  <TableCell className="text-sm">
                                                    {systems.get(kill.SolarSystemID) || `System ${kill.SolarSystemID}`}
                                                  </TableCell>
                                                  <TableCell className="text-right font-mono text-sm">
                                                    {formatISK(kill.ZkillData.TotalValue)}
                                                  </TableCell>
                                                  <TableCell className="text-right text-sm">
                                                    <Badge variant="secondary">
                                                      {kill.ZkillData.Points || 0}
                                                    </Badge>
                                                  </TableCell>
                                                </TableRow>
                                              ))}
                                          </TableBody>
                                        </Table>
                                      ) : (
                                        <p className="text-center text-muted-foreground py-4">
                                          No kills found for this character in the selected period.
                                        </p>
                                      )}
                                    </div>
                                  </PopoverContent>
                                </Popover>
                              </TableCell>
                            </TableRow>
                          )
                        })
                      ) : (
                        <TableRow>
                          <TableCell colSpan={4} className="h-24 text-center">
                            No character data available.
                          </TableCell>
                        </TableRow>
                      )}
                    </TableBody>
                  </Table>
                </CardContent>
              </Card>
            </>
          )}
        </>
      )}
    </div>
  )
}