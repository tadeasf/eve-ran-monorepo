'use client'

import { useState, useEffect } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"
import { Label } from "@/components/ui/label"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Checkbox } from "@/components/ui/checkbox"
import { useToast } from "@/hooks/use-toast"
import { Save, RotateCcw } from 'lucide-react'
import { Region } from '@/lib/types'
import { apiFetchJson } from '@/lib/api'

interface DashboardSettings {
  defaultRegions: number[]
  defaultStartDate: string // e.g., "14" for 14 days ago
  defaultEndDate: string // e.g., "0" for today
}

const DEFAULT_SETTINGS: DashboardSettings = {
  defaultRegions: [], // Empty means all regions
  defaultStartDate: '14',
  defaultEndDate: '0',
}

const fetchRegions = async (): Promise<Region[]> => {
  return apiFetchJson<Region[]>('/regions')
}

export function AdminSettings() {
  const { toast } = useToast()
  const [settings, setSettings] = useState<DashboardSettings>(DEFAULT_SETTINGS)
  const [selectedRegions, setSelectedRegions] = useState<number[]>([])

  const { data: regions, isLoading } = useQuery({
    queryKey: ['regions'],
    queryFn: fetchRegions,
  })

  // Load settings from localStorage on mount
  useEffect(() => {
    const saved = localStorage.getItem('dashboard_settings')
    if (saved) {
      try {
        const parsed = JSON.parse(saved)
        setSettings(parsed)
        setSelectedRegions(parsed.defaultRegions)
      } catch (e) {
        console.error('Failed to parse saved settings:', e)
      }
    }
  }, [])

  const handleRegionToggle = (regionId: number) => {
    setSelectedRegions(prev => {
      if (prev.includes(regionId)) {
        return prev.filter(id => id !== regionId)
      } else {
        return [...prev, regionId]
      }
    })
  }

  const handleSelectAllRegions = () => {
    if (regions) {
      if (selectedRegions.length === regions.length) {
        setSelectedRegions([])
      } else {
        setSelectedRegions(regions.map(r => r.region_id))
      }
    }
  }

  const handleSave = () => {
    const newSettings: DashboardSettings = {
      defaultRegions: selectedRegions,
      defaultStartDate: settings.defaultStartDate,
      defaultEndDate: settings.defaultEndDate,
    }

    localStorage.setItem('dashboard_settings', JSON.stringify(newSettings))
    setSettings(newSettings)

    toast({
      title: "Settings saved successfully!",
      description: "Your dashboard preferences have been updated and will be applied on your next dashboard visit.",
    })
  }

  const handleReset = () => {
    localStorage.removeItem('dashboard_settings')
    setSettings(DEFAULT_SETTINGS)
    setSelectedRegions([])

    toast({
      title: "Settings reset",
      description: "Dashboard preferences have been reset to defaults.",
    })
  }

  const getDaysAgoDate = (daysAgo: number): string => {
    const date = new Date()
    date.setDate(date.getDate() - daysAgo)
    return date.toISOString().split('T')[0]
  }

  if (isLoading) {
    return (
      <div className="space-y-4">
        <div className="animate-pulse space-y-4">
          <div className="h-8 w-64 bg-gray-200 rounded"></div>
          <div className="h-64 bg-gray-200 rounded"></div>
        </div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader>
          <CardTitle>Default Date Range</CardTitle>
          <CardDescription>
            Set the default date range for dashboard views. Enter the number of days to look back from today.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="startDate">Start Date (days ago)</Label>
              <Input
                id="startDate"
                type="number"
                min="0"
                value={settings.defaultStartDate}
                onChange={(e) => setSettings({ ...settings, defaultStartDate: e.target.value })}
                placeholder="14"
              />
              <p className="text-xs text-muted-foreground">
                Will show data from: {getDaysAgoDate(parseInt(settings.defaultStartDate) || 0)}
              </p>
            </div>

            <div className="space-y-2">
              <Label htmlFor="endDate">End Date (days ago)</Label>
              <Input
                id="endDate"
                type="number"
                min="0"
                value={settings.defaultEndDate}
                onChange={(e) => setSettings({ ...settings, defaultEndDate: e.target.value })}
                placeholder="0"
              />
              <p className="text-xs text-muted-foreground">
                Will show data to: {getDaysAgoDate(parseInt(settings.defaultEndDate) || 0)}
              </p>
            </div>
          </div>

          <div className="flex items-center gap-2 p-4 bg-muted rounded-lg">
            <div className="text-sm">
              <strong>Preview:</strong> Dashboard will show kills from the last{' '}
              {parseInt(settings.defaultStartDate) || 0} days
              {parseInt(settings.defaultEndDate) > 0 && ` until ${parseInt(settings.defaultEndDate)} days ago`}
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Default Regions</CardTitle>
          <CardDescription>
            Select which regions should be displayed by default. Leave all unchecked to show all regions.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-between pb-4 border-b">
            <p className="text-sm text-muted-foreground">
              {selectedRegions.length === 0
                ? 'All regions selected (default)'
                : `${selectedRegions.length} region(s) selected`}
            </p>
            <Button variant="outline" size="sm" onClick={handleSelectAllRegions}>
              {selectedRegions.length === regions?.length ? 'Deselect All' : 'Select All'}
            </Button>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 max-h-[400px] overflow-y-auto p-2">
            {regions?.map((region) => (
              <div key={region.region_id} className="flex items-center space-x-2">
                <Checkbox
                  id={`region-${region.region_id}`}
                  checked={selectedRegions.includes(region.region_id)}
                  onCheckedChange={() => handleRegionToggle(region.region_id)}
                />
                <Label
                  htmlFor={`region-${region.region_id}`}
                  className="text-sm font-normal cursor-pointer"
                >
                  {region.name}
                </Label>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      <div className="flex justify-end gap-4">
        <Button variant="outline" onClick={handleReset}>
          <RotateCcw className="mr-2 h-4 w-4" />
          Reset to Defaults
        </Button>
        <Button onClick={handleSave}>
          <Save className="mr-2 h-4 w-4" />
          Save Settings
        </Button>
      </div>
    </div>
  )
}
