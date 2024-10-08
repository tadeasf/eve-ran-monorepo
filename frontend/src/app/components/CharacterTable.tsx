import { useState, useMemo } from 'react'
import { CharacterStats } from '../../lib/types'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/app/components/ui/table"
import { formatISK } from '../../lib/utils'
import { ArrowUpDown } from 'lucide-react'

interface CharacterTableProps {
  characters: CharacterStats[]
}

export default function CharacterTable({ characters }: CharacterTableProps) {
  const [sortColumn, setSortColumn] = useState<keyof CharacterStats>('kill_count')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')

  const sortedCharacters = useMemo(() => {
    return [...characters].sort((a, b) => {
      if (a[sortColumn] < b[sortColumn]) return sortDirection === 'asc' ? -1 : 1
      if (a[sortColumn] > b[sortColumn]) return sortDirection === 'asc' ? 1 : -1
      return 0
    })
  }, [characters, sortColumn, sortDirection])

  const handleSort = (column: keyof CharacterStats) => {
    if (column === sortColumn) {
      setSortDirection(sortDirection === 'asc' ? 'desc' : 'asc')
    } else {
      setSortColumn(column)
      setSortDirection('desc')
    }
  }

  const SortableHeader = ({ column, children }: { column: keyof CharacterStats, children: React.ReactNode }) => (
    <TableHead onClick={() => handleSort(column)} className="cursor-pointer">
      <div className="flex items-center">
        {children}
        <ArrowUpDown className="ml-2 h-4 w-4" />
      </div>
    </TableHead>
  )

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <SortableHeader column="name">Name</SortableHeader>
          <SortableHeader column="kill_count">Kills</SortableHeader>
          <SortableHeader column="total_isk">ISK Destroyed</SortableHeader>
        </TableRow>
      </TableHeader>
      <TableBody>
        {sortedCharacters.map((character) => (
          <TableRow key={character.character_id}>
            <TableCell>{character.name}</TableCell>
            <TableCell>{character.kill_count}</TableCell>
            <TableCell>{formatISK(character.total_isk)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}