import { useState } from 'react'
import { CharacterStats, Kill } from '../../lib/types'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/app/components/ui/table"
import { formatISK } from '../../lib/utils'
import { ArrowUpDown } from 'lucide-react'
import { Button } from "@/app/components/ui/button"
import CharacterKills from './CharacterKills'

interface CharacterTableProps {
  characters: CharacterStats[]
  allKills: Kill[]
  startDate: string
  endDate: string
  selectedRegions: Array<{ id: number, name: string }>
}

export default function CharacterTable({ characters, allKills }: CharacterTableProps) {
  const [sortColumn, setSortColumn] = useState<keyof CharacterStats>('kill_count')
  const [sortDirection, setSortDirection] = useState<'asc' | 'desc'>('desc')
  const [currentPage, setCurrentPage] = useState(1)
  const [selectedCharacter, setSelectedCharacter] = useState<CharacterStats | null>(null)
  const itemsPerPage = 5

  const sortedCharacters = [...characters].sort((a, b) => {
    if (a[sortColumn] < b[sortColumn]) return sortDirection === 'asc' ? -1 : 1
    if (a[sortColumn] > b[sortColumn]) return sortDirection === 'asc' ? 1 : -1
    return 0
  })

  const paginatedCharacters = sortedCharacters.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  )

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

  const handleCharacterClick = (character: CharacterStats) => {
    setSelectedCharacter(character)
  }

  return (
    <>
      <Table>
        <TableHeader>
          <TableRow>
            <SortableHeader column="name">Name</SortableHeader>
            <SortableHeader column="kill_count">Kills</SortableHeader>
            <SortableHeader column="total_isk">ISK Destroyed</SortableHeader>
          </TableRow>
        </TableHeader>
        <TableBody>
          {paginatedCharacters.map((character) => (
            <TableRow key={character.character_id} onClick={() => handleCharacterClick(character)} className="cursor-pointer hover:bg-gray-100">
              <TableCell>{character.name}</TableCell>
              <TableCell>{character.kill_count}</TableCell>
              <TableCell>{formatISK(character.total_isk)}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      <div className="flex justify-between mt-4">
        <Button
          onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
          disabled={currentPage === 1}
        >
          Previous
        </Button>
        <span>Page {currentPage} of {Math.ceil(characters.length / itemsPerPage)}</span>
        <Button
          onClick={() => setCurrentPage(prev => Math.min(prev + 1, Math.ceil(characters.length / itemsPerPage)))}
          disabled={currentPage === Math.ceil(characters.length / itemsPerPage)}
        >
          Next
        </Button>
      </div>
      {selectedCharacter && (
        <CharacterKills
          character={selectedCharacter}
          kills={allKills.filter(kill => kill.CharacterID === selectedCharacter.character_id)}
          onClose={() => setSelectedCharacter(null)}
        />
      )}
    </>
  )
}