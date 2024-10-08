import { useState, useMemo } from 'react'
import { CharacterStats } from '../../lib/types'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/app/components/ui/table"
import { formatISK } from '../../lib/utils'

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

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead onClick={() => handleSort('name')} className="cursor-pointer">Name</TableHead>
          <TableHead onClick={() => handleSort('kill_count')} className="cursor-pointer">Kills</TableHead>
          <TableHead onClick={() => handleSort('total_isk')} className="cursor-pointer">ISK Destroyed</TableHead>
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