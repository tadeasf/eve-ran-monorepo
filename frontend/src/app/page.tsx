"use client"

import { useState, useEffect } from 'react'
import Image from 'next/image'
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
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import {
    Dialog,
    DialogContent,
    DialogDescription,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog'
import { Trophy, Target, TrendingUp, Calendar, Eye } from 'lucide-react'
import {
    getCurrentCompetitionStandings,
    getCompetitionSettings,
    getYearToDateWinners,
    getAllRegions,
    type CompetitionStanding,
    type CompetitionSettings,
    type CompetitionWinner
} from '@/app/lib/eveApi'

export default function Home() {
    const [standings, setStandings] = useState<CompetitionStanding[]>([])
    const [settings, setSettings] = useState<CompetitionSettings | null>(null)
    const [ytdWinners, setYtdWinners] = useState<CompetitionWinner[]>([])
    const [isLoading, setIsLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)
    const [allRegions, setAllRegions] = useState<{ region_id: number; name: string }[]>([])
    const [isLoadingRegions, setIsLoadingRegions] = useState(false)

    useEffect(() => {
        const loadData = async () => {
            try {
                setIsLoading(true)
                setError(null)
                
                const [standingsData, settingsData, winnersData] = await Promise.all([
                    getCurrentCompetitionStandings(),
                    getCompetitionSettings(),
                    getYearToDateWinners()
                ])
                
                setStandings(standingsData)
                setSettings(settingsData)
                setYtdWinners(winnersData)
            } catch (err) {
                console.error('Error loading competition data:', err)
                setError('Failed to load competition data. Please try again later.')
            } finally {
                setIsLoading(false)
            }
        }

        loadData()
    }, [])

    const formatValue = (value: number, metric: string) => {
        if (metric === 'isk_destroyed') {
            return `${(value / 1_000_000_000).toFixed(2)}B ISK`
        }
        return value.toLocaleString()
    }

    const getMetricLabel = (metric: string) => {
        return metric === 'isk_destroyed' ? 'ISK Destroyed' : 'Kill Count'
    }

    const getRankBadge = (rank: number) => {
        if (rank === 1) {
            return <Badge className="bg-yellow-500 hover:bg-yellow-600">🥇 1st</Badge>
        } else if (rank === 2) {
            return <Badge className="bg-gray-400 hover:bg-gray-500">🥈 2nd</Badge>
        } else if (rank === 3) {
            return <Badge className="bg-amber-600 hover:bg-amber-700">🥉 3rd</Badge>
        }
        return <Badge variant="outline">{rank}</Badge>
    }

    const getMonthName = (month: number) => {
        const months = ['January', 'February', 'March', 'April', 'May', 'June',
                       'July', 'August', 'September', 'October', 'November', 'December']
        return months[month - 1]
    }

    if (isLoading) {
        return (
            <div className="container mx-auto px-4 py-8">
                <div className="flex items-center justify-center min-h-[400px]">
                    <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary"></div>
                </div>
            </div>
        )
    }

    if (error) {
        return (
            <div className="container mx-auto px-4 py-8">
                <Card className="border-destructive">
                    <CardHeader>
                        <CardTitle className="text-destructive">Error</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p>{error}</p>
                    </CardContent>
                </Card>
            </div>
        )
    }

    const currentMonth = new Date().toLocaleString('default', { month: 'long', year: 'numeric' })

    return (
        <div className="container mx-auto px-4 py-8 space-y-8">
            {/* Header */}
            <div className="space-y-2">
                <h1 className="text-4xl font-bold flex items-center gap-3">
                    <Trophy className="size-10 text-yellow-500" />
                    Competition Leaderboard
                </h1>
                <p className="text-muted-foreground text-lg">
                    Current standings for {currentMonth}
                </p>
            </div>

            {/* Competition Info Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Competition Metric</CardTitle>
                        <Target className="size-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">
                            {settings ? getMetricLabel(settings.metric) : 'Loading...'}
                        </div>
                        <p className="text-xs text-muted-foreground">
                            Current tracking metric
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Participants</CardTitle>
                        <TrendingUp className="size-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">
                            {standings.length}
                        </div>
                        <p className="text-xs text-muted-foreground">
                            Active competitors
                        </p>
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                        <CardTitle className="text-sm font-medium">Regions</CardTitle>
                        <Calendar className="size-4 text-muted-foreground" />
                    </CardHeader>
                    <CardContent>
                        <div className="text-2xl font-bold">
                            {!settings?.regions || settings.regions.length === 0 ? 'All' : settings.regions.length}
                        </div>
                        <p className="text-xs text-muted-foreground">
                            Tracked regions
                        </p>
                        {settings?.regions && settings.regions.length > 0 && (
                            <Dialog>
                                <DialogTrigger asChild>
                                    <Button
                                        variant="outline"
                                        size="sm"
                                        className="w-full mt-2"
                                        onClick={async () => {
                                            if (allRegions.length === 0) {
                                                setIsLoadingRegions(true)
                                                try {
                                                    const regions = await getAllRegions()
                                                    setAllRegions(regions)
                                                } catch (err) {
                                                    console.error('Failed to fetch regions:', err)
                                                } finally {
                                                    setIsLoadingRegions(false)
                                                }
                                            }
                                        }}
                                    >
                                        <Eye className="size-4 mr-2" />
                                        Show Regions
                                    </Button>
                                </DialogTrigger>
                                <DialogContent className="max-w-md">
                                    <DialogHeader>
                                        <DialogTitle>Tracked Regions</DialogTitle>
                                        <DialogDescription>
                                            The following regions are being tracked for this competition
                                        </DialogDescription>
                                    </DialogHeader>
                                    <div className="max-h-96 overflow-y-auto">
                                        {isLoadingRegions ? (
                                            <div className="flex items-center justify-center py-8">
                                                <div className="size-8 animate-spin rounded-full border-4 border-primary border-t-transparent" />
                                            </div>
                                        ) : (
                                            <div className="space-y-2">
                                                {allRegions
                                                    .filter(region => settings?.regions.includes(region.region_id))
                                                    .map(region => (
                                                        <div
                                                            key={region.region_id}
                                                            className="flex items-center justify-between rounded-lg border p-3"
                                                        >
                                                            <span className="font-medium">{region.name}</span>
                                                            <Badge variant="secondary" className="ml-2">
                                                                {region.region_id}
                                                            </Badge>
                                                        </div>
                                                    ))}
                                            </div>
                                        )}
                                    </div>
                                </DialogContent>
                            </Dialog>
                        )}
                    </CardContent>
                </Card>
            </div>

            {/* Current Standings Table */}
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl">
                        Current Standings
                    </CardTitle>
                    <CardDescription>
                        Live leaderboard for {currentMonth}
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    <Table>
                        <TableCaption>
                            Competition standings are updated in real-time
                        </TableCaption>
                        <TableHeader>
                            <TableRow>
                                <TableHead className="w-[100px]">Rank</TableHead>
                                <TableHead>Character</TableHead>
                                <TableHead className="text-right">
                                    {settings ? getMetricLabel(settings.metric) : 'Value'}
                                </TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {standings.length > 0 ? (
                                standings.map((standing) => (
                                    <TableRow key={standing.character_id}>
                                        <TableCell>{getRankBadge(standing.rank)}</TableCell>
                                        <TableCell className="font-medium">
                                            <div className="flex items-center gap-2">
                                                <Image
                                                    src={`https://images.evetech.net/characters/${standing.character_id}/portrait?size=32`}
                                                    alt={standing.character_name}
                                                    width={32}
                                                    height={32}
                                                    className="size-8 rounded-full"
                                                    unoptimized
                                                />
                                                {standing.character_name}
                                            </div>
                                        </TableCell>
                                        <TableCell className="text-right font-mono">
                                            {settings && formatValue(standing.value, settings.metric)}
                                        </TableCell>
                                    </TableRow>
                                ))
                            ) : (
                                <TableRow>
                                    <TableCell colSpan={3} className="h-24 text-center">
                                        No competition data available yet.
                                    </TableCell>
                                </TableRow>
                            )}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>

            {/* Year-to-Date Winners Section */}
            <Separator />
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl">
                        {new Date().getFullYear()} Winners
                    </CardTitle>
                    <CardDescription>
                        Monthly champions for the current year
                    </CardDescription>
                </CardHeader>
                <CardContent>
                    {ytdWinners.length > 0 ? (
                        <Table>
                            <TableCaption>
                                First place winners for each month
                            </TableCaption>
                            <TableHeader>
                                <TableRow>
                                    <TableHead>Month</TableHead>
                                    <TableHead>Champion</TableHead>
                                    <TableHead className="text-right">
                                        {settings && getMetricLabel(settings.metric)}
                                    </TableHead>
                                </TableRow>
                            </TableHeader>
                            <TableBody>
                                {ytdWinners.map((winner) => (
                                    <TableRow key={`${winner.month}-${winner.year}`}>
                                        <TableCell className="font-medium">
                                            {getMonthName(winner.month)}
                                        </TableCell>
                                        <TableCell>
                                            <div className="flex items-center gap-2">
                                                <Image
                                                    src={`https://images.evetech.net/characters/${winner.character_id}/portrait?size=32`}
                                                    alt={winner.character_name}
                                                    width={32}
                                                    height={32}
                                                    className="size-8 rounded-full"
                                                    unoptimized
                                                />
                                                {winner.character_name}
                                            </div>
                                        </TableCell>
                                        <TableCell className="text-right font-mono">
                                            {formatValue(winner.value, winner.metric)}
                                        </TableCell>
                                    </TableRow>
                                ))}
                            </TableBody>
                        </Table>
                    ) : (
                        <div className="text-center py-8 text-muted-foreground">
                            <p>No previous monthly winners yet.</p>
                            <p className="text-sm mt-2">Winners will appear here after the first month is completed and results are saved.</p>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}