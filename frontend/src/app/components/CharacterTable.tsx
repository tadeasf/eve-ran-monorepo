import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/app/components/ui/table"
import { CharacterStats } from '@/lib/types'

interface CharacterTableProps {
  characters: CharacterStats[]
}

export default function CharacterTable({ characters }: CharacterTableProps) {
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Character Name</TableHead>
          <TableHead>Kill Count</TableHead>
          <TableHead>Total ISK Destroyed</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {characters.map((character) => (
          <TableRow key={character.character_id}>
            <TableCell>{character.name}</TableCell>
            <TableCell>{character.kill_count}</TableCell>
            <TableCell>{character.total_isk.toLocaleString()} ISK</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}