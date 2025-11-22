"use client"

import { useState, useEffect } from 'react'
import { useAuth } from '@/app/contexts/AuthContext'
import { useAppConfig } from '@/app/contexts/AppConfigContext'
import { Button } from '@/app/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/app/components/ui/card'
import { Sidebar, SidebarContent, SidebarFooter, SidebarHeader, SidebarProvider, SidebarTrigger } from '@/components/ui/sidebar'
import { Badge } from '@/components/ui/badge'
import { Separator } from '@/components/ui/separator'

import {
    Shield,
    Users,
    Settings,
    BarChart3,
    Database,
    LogOut,
    Plus,
    Upload,
    Trophy
} from 'lucide-react'
import { BatchCharacterManager } from './BatchCharacterManager'
import { AdminCharacterManager } from './AdminCharacterManager'
import { AdminAnalytics } from './AdminAnalytics'
import { AdminSettings } from './AdminSettings'
import { AdminCompetitionSettings } from './AdminCompetitionSettings'
import { getDashboardStats } from '@/app/lib/eveApi'

type AdminView = 'overview' | 'characters' | 'batch-add' | 'analytics' | 'settings' | 'competition-settings'

export function AdminDashboard() {
    const [activeView, setActiveView] = useState<AdminView>('overview')
    const [dashboardStats, setDashboardStats] = useState<{
        total_characters: number;
        total_kills: number;
        total_regions: number;
        recent_activity: number;
    } | null>(null)
    const [isLoadingStats, setIsLoadingStats] = useState(true)
    const { logout } = useAuth()
    const { config } = useAppConfig()

    // Load dashboard stats on component mount
    useEffect(() => {
        const loadStats = async () => {
            try {
                setIsLoadingStats(true)
                const stats = await getDashboardStats()
                setDashboardStats(stats)
            } catch (error) {
                console.error('Failed to load dashboard stats:', error)
            } finally {
                setIsLoadingStats(false)
            }
        }

        loadStats()
    }, [])

    const sidebarItems = [
        { id: 'overview' as AdminView, label: 'Overview', icon: BarChart3 },
        { id: 'characters' as AdminView, label: 'Characters', icon: Users },
        { id: 'batch-add' as AdminView, label: 'Batch Add', icon: Plus },
        { id: 'analytics' as AdminView, label: 'Analytics', icon: Database },
        { id: 'competition-settings' as AdminView, label: 'Competition Settings', icon: Trophy },
        { id: 'settings' as AdminView, label: 'Settings', icon: Settings },
    ]

    const renderContent = () => {
        switch (activeView) {
            case 'overview':
                return (
                    <div className="space-y-6">
                        <div>
                            <h1 className="text-3xl font-bold">Admin Overview</h1>
                            <p className="text-muted-foreground">
                                Welcome to the {config.corporation_name} admin panel
                            </p>
                        </div>

                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
                            <Card>
                                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                    <CardTitle className="text-sm font-medium">Total Characters</CardTitle>
                                    <Users className="size-4 text-muted-foreground" />
                                </CardHeader>
                                <CardContent>
                                    <div className="text-2xl font-bold">
                                        {isLoadingStats ? '...' : (dashboardStats?.total_characters || 0)}
                                    </div>
                                    <p className="text-xs text-muted-foreground">
                                        Characters in system
                                    </p>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                    <CardTitle className="text-sm font-medium">Total Kills</CardTitle>
                                    <BarChart3 className="size-4 text-muted-foreground" />
                                </CardHeader>
                                <CardContent>
                                    <div className="text-2xl font-bold">
                                        {isLoadingStats ? '...' : (dashboardStats?.total_kills || 0).toLocaleString()}
                                    </div>
                                    <p className="text-xs text-muted-foreground">
                                        Tracked kills
                                    </p>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                    <CardTitle className="text-sm font-medium">Active Regions</CardTitle>
                                    <Database className="size-4 text-muted-foreground" />
                                </CardHeader>
                                <CardContent>
                                    <div className="text-2xl font-bold">
                                        {isLoadingStats ? '...' : (dashboardStats?.total_regions || 0)}
                                    </div>
                                    <p className="text-xs text-muted-foreground">
                                        EVE Online regions
                                    </p>
                                </CardContent>
                            </Card>

                            <Card>
                                <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                                    <CardTitle className="text-sm font-medium">Recent Activity</CardTitle>
                                    <Users className="size-4 text-muted-foreground" />
                                </CardHeader>
                                <CardContent>
                                    <div className="text-2xl font-bold">
                                        {isLoadingStats ? '...' : (dashboardStats?.recent_activity || 0).toLocaleString()}
                                    </div>
                                    <p className="text-xs text-muted-foreground">
                                        Kills last 7 days
                                    </p>
                                </CardContent>
                            </Card>
                        </div>

                        <Card>
                            <CardHeader>
                                <CardTitle>Quick Actions</CardTitle>
                                <CardDescription>
                                    Commonly used admin functions
                                </CardDescription>
                            </CardHeader>
                            <CardContent className="flex flex-wrap gap-4">
                                <Button
                                    onClick={() => setActiveView('batch-add')}
                                    variant="outline"
                                    className="flex items-center gap-2"
                                >
                                    <Plus className="size-4" />
                                    Batch Add Characters
                                </Button>
                                <Button
                                    onClick={() => setActiveView('analytics')}
                                    variant="outline"
                                    className="flex items-center gap-2"
                                >
                                    <BarChart3 className="size-4" />
                                    View Analytics
                                </Button>
                                <Button
                                    variant="outline"
                                    className="flex items-center gap-2"
                                >
                                    <Upload className="size-4" />
                                    Import Data
                                </Button>
                            </CardContent>
                        </Card>
                    </div>
                )

            case 'batch-add':
                return <BatchCharacterManager />

            case 'characters':
                return (
                    <div className="space-y-6">
                        <div>
                            <h1 className="text-3xl font-bold">Character Management</h1>
                            <p className="text-muted-foreground">
                                Search and manage characters using the EVE ESI API
                            </p>
                        </div>
                        <AdminCharacterManager />
                    </div>
                )

            case 'analytics':
                return <AdminAnalytics />

            case 'competition-settings':
                return (
                    <div className="space-y-6">
                        <div>
                            <h1 className="text-3xl font-bold">Competition Settings</h1>
                            <p className="text-muted-foreground">
                                Configure competition parameters and tracking metrics
                            </p>
                        </div>
                        <AdminCompetitionSettings />
                    </div>
                )

            case 'settings':
                return (
                    <div className="space-y-6">
                        <div>
                            <h1 className="text-3xl font-bold">Dashboard Settings</h1>
                            <p className="text-muted-foreground">
                                Configure your default dashboard filters and preferences
                            </p>
                        </div>
                        <AdminSettings />
                    </div>
                )

            default:
                return null
        }
    }

    return (
        <SidebarProvider>
            <div className="flex h-screen w-full">
                <Sidebar>
                    <SidebarHeader className="p-4">
                        <div className="flex items-center space-x-2">
                            <Shield className="size-8 text-primary" />
                            <div>
                                <h2 className="text-lg font-semibold">Admin Panel</h2>
                                <p className="text-sm text-muted-foreground">{config.corporation_name}</p>
                            </div>
                        </div>
                    </SidebarHeader>

                    <SidebarContent className="p-4">
                        <nav className="space-y-2">
                            {sidebarItems.map((item) => {
                                const Icon = item.icon
                                return (
                                    <Button
                                        key={item.id}
                                        variant={activeView === item.id ? "default" : "ghost"}
                                        className="w-full justify-start"
                                        onClick={() => setActiveView(item.id)}
                                    >
                                        <Icon className="mr-2 size-4" />
                                        {item.label}
                                    </Button>
                                )
                            })}
                        </nav>
                    </SidebarContent>

                    <SidebarFooter className="p-4">
                        <Separator className="mb-4" />
                        <div className="flex items-center space-x-2 mb-4">
                            <div className="size-8 rounded-full bg-primary flex items-center justify-center">
                                <Shield className="size-4 text-primary-foreground" />
                            </div>
                            <div className="flex-1">
                                <p className="text-sm font-medium">Admin User</p>
                                <p className="text-xs text-muted-foreground">admin@{config.corporation_name.toLowerCase().replace(/\s+/g, '.')}</p>
                            </div>
                        </div>
                        <Button
                            variant="outline"
                            size="sm"
                            onClick={logout}
                            className="w-full"
                        >
                            <LogOut className="mr-2 size-4" />
                            Sign Out
                        </Button>
                    </SidebarFooter>
                </Sidebar>

                <main className="flex-1 overflow-hidden">
                    <div className="flex items-center justify-between p-4 border-b">
                        <SidebarTrigger />
                        <Badge variant="secondary">Admin Mode</Badge>
                    </div>
                    <div className="p-6 overflow-auto h-full">
                        {renderContent()}
                    </div>
                </main>
            </div>
        </SidebarProvider>
    )
}
