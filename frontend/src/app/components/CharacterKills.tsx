

import { useState, useMemo } from 'react'
import { CharacterStats, Kill } from '../../lib/types'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/app/components/ui/table"
import { formatISK, formatDate } from '../../lib/utils'
import { ArrowUpDown, X } from 'lucide-react'
import { Button } from "@/app/components/ui/button"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/app/components/ui/select"

interface CharacterKillsProps {
  character: CharacterStats
  kills: Kill[]
  onClose: () => void
}

type SortableColumn = 'KillmailTime' | 'FittedValue' | 'DroppedValue' | 'DestroyedValue' | 'TotalValue' | 'Points'
type SoloFilter = 'all' | 'solo' | 'not-solo'

export default function CharacterKills({ character, kills, onClose }: CharacterKillsProps) {
  const [sortColumn, setSortColumn] = useState<SortableColumn>('KillmailTime')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')
  const [soloFilter, setSoloFilter] = useState<SoloFilter>('all')

  const filteredAndSortedKills = useMemo(() => {
    let filteredKills = kills;
    if (soloFilter !== 'all') {
      filteredKills = kills.filter(kill => 
        soloFilter === 'solo' ? kill.ZkillData.Solo : !kill.ZkillData.Solo
      );
    }

    return filteredKills.sort((a, b) => {
      let aValue, bValue;
      if (sortColumn === 'KillmailTime') {
        aValue = a[sortColumn];
        bValue = b[sortColumn];
      } else {
        aValue = a.ZkillData[sortColumn];
        bValue = b.ZkillData[sortColumn];
      }
      if (aValue < bValue) return sortDirection === 'asc' ? -1 : 1
      if (aValue > bValue) return sortDirection === 'asc' ? 1 : -1
      return 0
    })
  }, [kills, soloFilter, sortColumn, sortDirection])

  const handleSort = (column: SortableColumn) => {
    if (column === sortColumn) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortColumn(column)
      setSortDirection('desc')
    }
  }

  const SortableHeader = ({ column, children }: { column: SortableColumn, children: React.ReactNode }) => (
    <TableHead onClick={() => handleSort(column)} className="cursor-pointer">
      <div className="flex items-center">
        {children}
        <ArrowUpDown className="ml-2 h-4 w-4" />
      </div>
    </TableHead>
  )

  return (
    <div className="mt-8">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-2xl font-bold">Kills for {character.name}</h2>
        <div className="flex items-center space-x-4">
          <Select value={soloFilter} onValueChange={(value: SoloFilter) => setSoloFilter(value)}>
            <SelectTrigger className="w-[180px]">
              <SelectValue placeholder="Filter by solo" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">All Kills</SelectItem>
              <SelectItem value="solo">Solo Kills</SelectItem>
              <SelectItem value="not-solo">Not Solo Kills</SelectItem>
            </SelectContent>
          </Select>
          <Button onClick={onClose}><X className="h-4 w-4" /></Button>
        </div>
      </div>
      <Table>
        <TableHeader>
          <TableRow>
            <SortableHeader column="KillmailTime">Time</SortableHeader>
            <TableHead>Solo</TableHead>
            <SortableHeader column="FittedValue">Fitted Value</SortableHeader>
            <SortableHeader column="DroppedValue">Dropped Value</SortableHeader>
            <SortableHeader column="DestroyedValue">Destroyed Value</SortableHeader>
            <SortableHeader column="TotalValue">Total Value</SortableHeader>
            <SortableHeader column="Points">Points</SortableHeader>
            <TableHead>zKillboard</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {filteredAndSortedKills.map((kill) => (
            <TableRow key={kill.KillmailID}>
              <TableCell>{formatDate(kill.KillmailTime)}</TableCell>
              <TableCell>{kill.ZkillData.Solo ? 'Solo' : 'Not Solo'}</TableCell>
              <TableCell>{formatISK(kill.ZkillData.FittedValue)}</TableCell>
              <TableCell>{formatISK(kill.ZkillData.DroppedValue)}</TableCell>
              <TableCell>{formatISK(kill.ZkillData.DestroyedValue)}</TableCell>
              <TableCell>{formatISK(kill.ZkillData.TotalValue)}</TableCell>
              <TableCell>{kill.ZkillData.Points}</TableCell>
              <TableCell>
                <a href={`https://zkillboard.com/kill/${kill.ZkillData.KillmailID}/`} target="_blank" rel="noopener noreferrer" className="text-blue-500 hover:underline">
                  View
                </a>
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  )
}