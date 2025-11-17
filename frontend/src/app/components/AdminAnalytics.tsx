'use client'

import { useState, useEffect, useMemo } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select"
import { Skeleton } from "@/components/ui/skeleton"
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, PieChart, Pie, Cell, LineChart, Line } from 'recharts'
import { Character, Kill, System } from '@/lib/types'
import { Trophy, Target, TrendingUp } from 'lucide-react'
import { apiFetchJson } from '@/lib/api'

const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8', '#82CA9D', '#FFC658', '#FF6B9D']

type DateRange = 'today' | '7days' | 'month' | 'year' | 'all'

const fetchCharacters = async (): Promise<Character[]> => {
  return apiFetchJson<Character[]>('/characters')
}

const fetchAllKills = async (): Promise<Kill[]> => {
  return apiFetchJson<Kill[]>('/kills')
}

const fetchSystems = async (): Promise<System[]> => {
  return apiFetchJson<System[]>('/systems')
}

const getDateRangeFilter = (range: DateRange): { start: Date; end: Date } => {
  const end = new Date()
  const start = new Date()

  switch (range) {
    case 'today':
      start.setHours(0, 0, 0, 0)
      break
    case '7days':
      start.setDate(start.getDate() - 7)
      break
    case 'month':
      start.setMonth(start.getMonth() - 1)
      break
    case 'year':
      start.setFullYear(start.getFullYear() - 1)
      break
    case 'all':
      start.setFullYear(2000) // Far enough back
      break
  }

  return { start, end }
}

export function AdminAnalytics() {
  const [selectedCharacterId, setSelectedCharacterId] = useState<number | null>(null)
  const [dateRange, setDateRange] = useState<DateRange>('7days')

  const { data: characters, isLoading: isLoadingChars } = useQuery({
    queryKey: ['characters'],
    queryFn: fetchCharacters,
  })
  const { data: allKills, isLoading: isLoadingKills } = useQuery({
    queryKey: ['kills'],
    queryFn: fetchAllKills,
  })
  const { data: systems, isLoading: isLoadingSystems } = useQuery({
    queryKey: ['systems'],
    queryFn: fetchSystems,
  })

  // Auto-select first character
  useEffect(() => {
    if (characters && characters.length > 0 && !selectedCharacterId) {
      setSelectedCharacterId(characters[0].id)
    }
  }, [characters, selectedCharacterId])

  // Create system lookup map
  const systemsMap = useMemo(() => {
    if (!systems) return new Map<number, string>()
    return new Map(systems.map(s => [s.system_id, s.name]))
  }, [systems])

  // Filter kills by date range and character
  const filteredKills = useMemo(() => {
    if (!allKills || !selectedCharacterId) return []

    const { start, end } = getDateRangeFilter(dateRange)

    return allKills.filter(kill => {
      const killDate = new Date(kill.KillmailTime)
      return kill.CharacterID === selectedCharacterId &&
             killDate >= start &&
             killDate <= end
    })
  }, [allKills, selectedCharacterId, dateRange])

  const analytics = useMemo(() => {
    if (filteredKills.length === 0) {
      return {
        totalKills: 0,
        totalISK: 0,
        totalPoints: 0,
        regionActivity: [],
        hourlyActivity: [],
        dailyActivity: [],
        topSystems: []
      }
    }

    const systemActivityMap = new Map<number, { kills: number; isk: number; name: string }>()
    const hourMap = new Map<number, number>()
    const dailyMap = new Map<string, { kills: number; isk: number }>()

    let totalISK = 0
    let totalPoints = 0

    filteredKills.forEach(kill => {
      totalISK += kill.ZkillData.TotalValue
      totalPoints += kill.ZkillData.Points

      // System activity with proper names
      const systemName = systemsMap.get(kill.SolarSystemID) || `Unknown System (${kill.SolarSystemID})`
      if (!systemActivityMap.has(kill.SolarSystemID)) {
        systemActivityMap.set(kill.SolarSystemID, { kills: 0, isk: 0, name: systemName })
      }
      const systemData = systemActivityMap.get(kill.SolarSystemID)!
      systemData.kills++
      systemData.isk += kill.ZkillData.TotalValue

      // Hourly activity
      const hour = new Date(kill.KillmailTime).getHours()
      hourMap.set(hour, (hourMap.get(hour) || 0) + 1)

      // Daily activity
      const date = kill.KillmailTime.split('T')[0]
      if (!dailyMap.has(date)) {
        dailyMap.set(date, { kills: 0, isk: 0 })
      }
      const dailyData = dailyMap.get(date)!
      dailyData.kills++
      dailyData.isk += kill.ZkillData.TotalValue
    })

    // Convert to arrays and sort
    const regionActivity = Array.from(systemActivityMap.entries())
      .map(([, data]) => ({ region: data.name, kills: data.kills, isk: data.isk }))
      .sort((a, b) => b.kills - a.kills)
      .slice(0, 10)

    const topSystems = Array.from(systemActivityMap.entries())
      .map(([, data]) => ({ system: data.name, kills: data.kills }))
      .sort((a, b) => b.kills - a.kills)
      .slice(0, 10)

    const hourlyActivity = Array.from(hourMap.entries())
      .map(([hour, kills]) => ({ hour: `${hour}:00`, kills }))
      .sort((a, b) => parseInt(a.hour) - parseInt(b.hour))

    const dailyActivity = Array.from(dailyMap.entries())
      .map(([date, data]) => ({ date, ...data }))
      .sort((a, b) => a.date.localeCompare(b.date))
      .slice(-30) // Last 30 days

    return {
      totalKills: filteredKills.length,
      totalISK,
      totalPoints,
      regionActivity,
      hourlyActivity,
      dailyActivity,
      topSystems
    }
  }, [filteredKills, systemsMap])

  const selectedCharacter = characters?.find(c => c.id === selectedCharacterId)

  if (isLoadingChars || isLoadingKills || isLoadingSystems) {
    return (
      <div className="space-y-6">
        <Skeleton className="w-full h-[200px]" />
        <Skeleton className="w-full h-[400px]" />
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Filters */}
      <div className="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
        <div className="space-y-1">
          <h2 className="text-2xl font-bold">Character Analytics</h2>
          <p className="text-sm text-muted-foreground">Deep dive into individual character performance</p>
        </div>
        <div className="flex flex-col sm:flex-row gap-3 w-full sm:w-auto">
          <Select
            value={dateRange}
            onValueChange={(value) => setDateRange(value as DateRange)}
          >
            <SelectTrigger className="w-full sm:w-[180px]">
              <SelectValue placeholder="Select date range" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="today">Today</SelectItem>
              <SelectItem value="7days">Last 7 Days</SelectItem>
              <SelectItem value="month">Last Month</SelectItem>
              <SelectItem value="year">Last Year</SelectItem>
              <SelectItem value="all">All Time</SelectItem>
            </SelectContent>
          </Select>

          <Select
            value={selectedCharacterId?.toString() || ''}
            onValueChange={(value) => setSelectedCharacterId(parseInt(value))}
          >
            <SelectTrigger className="w-full sm:w-[250px]">
              <SelectValue placeholder="Select a character" />
            </SelectTrigger>
            <SelectContent>
              {characters?.map((character) => (
                <SelectItem key={character.id} value={character.id.toString()}>
                  {character.name}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {selectedCharacter && (
        <>
          {/* Stats Overview */}
          <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Kills</CardTitle>
                <Target className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{analytics.totalKills.toLocaleString()}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total ISK Destroyed</CardTitle>
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {(analytics.totalISK / 1000000000).toFixed(2)}B ISK
                </div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Total Points</CardTitle>
                <Trophy className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">{analytics.totalPoints.toLocaleString()}</div>
              </CardContent>
            </Card>
            <Card>
              <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                <CardTitle className="text-sm font-medium">Avg ISK/Kill</CardTitle>
                <TrendingUp className="h-4 w-4 text-muted-foreground" />
              </CardHeader>
              <CardContent>
                <div className="text-2xl font-bold">
                  {analytics.totalKills > 0
                    ? (analytics.totalISK / analytics.totalKills / 1000000).toFixed(2)
                    : 0}M
                </div>
              </CardContent>
            </Card>
          </div>

          {/* Charts */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Daily Activity */}
            <Card>
              <CardHeader>
                <CardTitle>Kill Activity</CardTitle>
                <CardDescription>Daily kills and ISK destroyed</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <LineChart data={analytics.dailyActivity}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="date" tick={{ fontSize: 12 }} />
                    <YAxis yAxisId="left" />
                    <YAxis yAxisId="right" orientation="right" />
                    <Tooltip />
                    <Legend />
                    <Line yAxisId="left" type="monotone" dataKey="kills" stroke="#8884d8" name="Kills" />
                    <Line yAxisId="right" type="monotone" dataKey="isk" stroke="#82ca9d" name="ISK (B)" />
                  </LineChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            {/* Hourly Activity */}
            <Card>
              <CardHeader>
                <CardTitle>Activity by Hour</CardTitle>
                <CardDescription>When is this character most active?</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={analytics.hourlyActivity}>
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis dataKey="hour" />
                    <YAxis />
                    <Tooltip />
                    <Bar dataKey="kills" fill="#8884d8" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            {/* Top Systems */}
            <Card>
              <CardHeader>
                <CardTitle>Most Active Systems</CardTitle>
                <CardDescription>Top 10 systems by kill count</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <BarChart data={analytics.topSystems} layout="vertical">
                    <CartesianGrid strokeDasharray="3 3" />
                    <XAxis type="number" />
                    <YAxis type="category" dataKey="system" width={120} tick={{ fontSize: 10 }} />
                    <Tooltip />
                    <Bar dataKey="kills" fill="#82ca9d" />
                  </BarChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>

            {/* System Distribution */}
            <Card>
              <CardHeader>
                <CardTitle>System Distribution</CardTitle>
                <CardDescription>Top systems by activity</CardDescription>
              </CardHeader>
              <CardContent>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={analytics.regionActivity}
                      dataKey="kills"
                      nameKey="region"
                      cx="50%"
                      cy="50%"
                      outerRadius={80}
                      label={(entry) => entry.region.length > 15 ? entry.region.substring(0, 15) + '...' : entry.region}
                    >
                      {analytics.regionActivity.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                      ))}
                    </Pie>
                    <Tooltip />
                  </PieChart>
                </ResponsiveContainer>
              </CardContent>
            </Card>
          </div>
        </>
      )}
    </div>
  )
}
