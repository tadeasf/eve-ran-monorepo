'use client'

import { useState, useEffect } from 'react'
import RegionPerformance from '../components/charts/RegionPerformance'
import ISKPerformance from '../components/charts/ISKPerformance'
import CharacterPerformance from '../components/charts/CharacterPerformance'
import { Kill, CharacterStats } from '@/lib/types'

interface RegionData {
  region: string
  kills: number
  isk: number
}

interface ISKData {
  date: string
  isk: number
}

interface CharacterData {
  name: string
  kills: number
  isk: number
}

export default function ChartsPage() {
  const [regionData, setRegionData] = useState<RegionData[]>([])
  const [iskData, setIskData] = useState<ISKData[]>([])
  const [characterData, setCharacterData] = useState<CharacterData[]>([])

  useEffect(() => {
    const fetchData = async () => {
      try {
        // Fetch kills data
        const killsResponse = await fetch('/api/kills')
        const killsData: Kill[] = await killsResponse.json()

        // Fetch characters data
        const charactersResponse = await fetch('/api/characters')
        const charactersData: CharacterStats[] = await charactersResponse.json()

        // Process kills data for region performance
        const regionMap = new Map<string, { kills: number; isk: number }>()
        killsData.forEach(kill => {
          const region = kill.ZkillData.LocationID.toString() // Assuming LocationID represents the region
          const currentData = regionMap.get(region) || { kills: 0, isk: 0 }
          regionMap.set(region, {
            kills: currentData.kills + 1,
            isk: currentData.isk + kill.ZkillData.TotalValue
          })
        })
        const processedRegionData: RegionData[] = Array.from(regionMap, ([region, data]) => ({
          region,
          kills: data.kills,
          isk: data.isk
        }))
        setRegionData(processedRegionData)

        // Process kills data for ISK performance
        const iskMap = new Map<string, number>()
        killsData.forEach(kill => {
          const date = kill.KillmailTime.split('T')[0]
          iskMap.set(date, (iskMap.get(date) || 0) + kill.ZkillData.TotalValue)
        })
        const processedIskData: ISKData[] = Array.from(iskMap, ([date, isk]) => ({ date, isk }))
        setIskData(processedIskData)

        // Process character data
        const processedCharacterData: CharacterData[] = charactersData.map(char => ({
          name: char.name,
          kills: char.kill_count,
          isk: char.total_isk
        }))
        setCharacterData(processedCharacterData)
      } catch (error) {
        console.error('Error fetching data:', error)
      }
    }

    fetchData()
  }, [])

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Performance Charts</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <RegionPerformance data={regionData} />
        <ISKPerformance data={iskData} />
        <CharacterPerformance data={characterData} />
      </div>
    </div>
  )
}