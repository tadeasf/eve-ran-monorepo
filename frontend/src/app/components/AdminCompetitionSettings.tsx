"use client"

import { useState, useEffect } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/app/components/ui/card'
import { Button } from '@/app/components/ui/button'
import { Label } from '@/components/ui/label'
import { RadioGroup, RadioGroupItem } from '@/components/ui/radio-group'
import { Checkbox } from '@/components/ui/checkbox'
import { Badge } from '@/components/ui/badge'
import { Save, Trophy, Loader2 } from 'lucide-react'
import { useToast } from '@/hooks/use-toast'
import {
    getCompetitionSettings,
    updateCompetitionSettings,
    type CompetitionSettings
} from '@/app/lib/eveApi'

interface Region {
    region_id: number
    name: string
}

export function AdminCompetitionSettings() {
    const [settings, setSettings] = useState<CompetitionSettings>({
        metric: 'kill_count',
        regions: [],
        active: true
    })
    const [availableRegions, setAvailableRegions] = useState<Region[]>([])
    const [isLoading, setIsLoading] = useState(true)
    const [isSaving, setIsSaving] = useState(false)
    const { toast } = useToast()

    useEffect(() => {
        const loadData = async () => {
            try {
                setIsLoading(true)
                
                // Load current settings
                try {
                    const currentSettings = await getCompetitionSettings()
                    setSettings(currentSettings)
                } catch {
                    // If no settings exist yet, use defaults
                    console.log('No existing settings found, using defaults')
                }
                
                // Load available regions
                const response = await fetch('/api/regions')
                if (response.ok) {
                    const regions = await response.json()
                    setAvailableRegions(regions)
                }
            } catch (err) {
                console.error('Error loading data:', err)
                toast({
                    title: 'Error',
                    description: 'Failed to load competition settings',
                    variant: 'destructive'
                })
            } finally {
                setIsLoading(false)
            }
        }

        loadData()
    }, [toast])

    const handleMetricChange = (value: string) => {
        setSettings(prev => ({
            ...prev,
            metric: value as 'isk_destroyed' | 'kill_count'
        }))
    }

    const handleRegionToggle = (regionId: number) => {
        setSettings(prev => {
            const currentRegions = prev.regions || []
            const isSelected = currentRegions.includes(regionId)
            
            return {
                ...prev,
                regions: isSelected
                    ? currentRegions.filter(id => id !== regionId)
                    : [...currentRegions, regionId]
            }
        })
    }

    const handleSelectAllRegions = () => {
        setSettings(prev => ({
            ...prev,
            regions: []
        }))
    }

    const handleSaveSettings = async () => {
        try {
            setIsSaving(true)
            await updateCompetitionSettings(settings)
            
            toast({
                title: 'Success',
                description: 'Competition settings saved successfully',
            })
        } catch (err) {
            console.error('Error saving settings:', err)
            toast({
                title: 'Error',
                description: 'Failed to save competition settings',
                variant: 'destructive'
            })
        } finally {
            setIsSaving(false)
        }
    }

    if (isLoading) {
        return (
            <div className="flex items-center justify-center min-h-[400px]">
                <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary"></div>
            </div>
        )
    }

    return (
        <div className="space-y-6">
            <Card>
                <CardHeader>
                    <CardTitle className="text-2xl flex items-center gap-2">
                        <Trophy className="size-6" />
                        Competition Configuration
                    </CardTitle>
                    <CardDescription>
                        Set the metric and regions for the monthly competition
                    </CardDescription>
                </CardHeader>
                <CardContent className="space-y-8">
                    {/* Metric Selection */}
                    <div className="space-y-4">
                        <Label className="text-lg font-semibold">Competition Metric</Label>
                        <RadioGroup
                            value={settings.metric}
                            onValueChange={handleMetricChange}
                            className="space-y-3"
                        >
                            <div className="flex items-center space-x-3 border rounded-lg p-4 hover:bg-accent cursor-pointer">
                                <RadioGroupItem value="kill_count" id="kill_count" />
                                <Label htmlFor="kill_count" className="cursor-pointer flex-1">
                                    <div className="font-medium">Kill Count</div>
                                    <div className="text-sm text-muted-foreground">
                                        Track total number of kills per character
                                    </div>
                                </Label>
                            </div>
                            <div className="flex items-center space-x-3 border rounded-lg p-4 hover:bg-accent cursor-pointer">
                                <RadioGroupItem value="isk_destroyed" id="isk_destroyed" />
                                <Label htmlFor="isk_destroyed" className="cursor-pointer flex-1">
                                    <div className="font-medium">ISK Destroyed</div>
                                    <div className="text-sm text-muted-foreground">
                                        Track total ISK value destroyed per character
                                    </div>
                                </Label>
                            </div>
                        </RadioGroup>
                    </div>

                    {/* Region Selection */}
                    <div className="space-y-4">
                        <div className="flex items-center justify-between">
                            <Label className="text-lg font-semibold">Tracked Regions</Label>
                            <Button
                                variant="outline"
                                size="sm"
                                onClick={handleSelectAllRegions}
                            >
                                Select All Regions
                            </Button>
                        </div>
                        <p className="text-sm text-muted-foreground">
                            Select specific regions to track, or leave empty to track all regions
                        </p>
                        
                        {settings.regions.length === 0 && (
                            <Badge variant="secondary" className="mb-2">
                                All Regions Selected
                            </Badge>
                        )}
                        
                        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3 max-h-96 overflow-y-auto border rounded-lg p-4">
                            {availableRegions.map((region) => (
                                <div
                                    key={region.region_id}
                                    className="flex items-center space-x-2"
                                >
                                    <Checkbox
                                        id={`region-${region.region_id}`}
                                        checked={settings.regions.includes(region.region_id)}
                                        onCheckedChange={() => handleRegionToggle(region.region_id)}
                                    />
                                    <Label
                                        htmlFor={`region-${region.region_id}`}
                                        className="text-sm cursor-pointer"
                                    >
                                        {region.name}
                                    </Label>
                                </div>
                            ))}
                        </div>
                    </div>

                    {/* Save Button */}
                    <div className="flex justify-end pt-4">
                        <Button
                            onClick={handleSaveSettings}
                            disabled={isSaving}
                            size="lg"
                        >
                            {isSaving ? (
                                <>
                                    <Loader2 className="mr-2 size-4 animate-spin" />
                                    Saving...
                                </>
                            ) : (
                                <>
                                    <Save className="mr-2 size-4" />
                                    Save Settings
                                </>
                            )}
                        </Button>
                    </div>
                </CardContent>
            </Card>

            {/* Info Card */}
            <Card className="border-blue-500/50 bg-blue-500/5">
                <CardHeader>
                    <CardTitle className="text-lg">ℹ️ How Competition Works</CardTitle>
                </CardHeader>
                <CardContent className="space-y-2 text-sm text-muted-foreground">
                    <p>
                        • Competitions run on a monthly basis, from the 1st to the last day of each month.
                    </p>
                    <p>
                        • The leaderboard is updated in real-time based on the selected metric.
                    </p>
                    <p>
                        • At the end of each month, the top performers are recorded in the competition history.
                    </p>
                    <p>
                        • You can change settings at any time, but changes will only affect the current month&apos;s competition.
                    </p>
                </CardContent>
            </Card>
        </div>
    )
}
