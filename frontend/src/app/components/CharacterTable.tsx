import { CharacterStats } from '../../lib/types'
import { formatISK } from '../../lib/utils'
import { Table, TableBody, TableCaption, TableCell, TableHead, TableHeader, TableRow } from "./ui/table"

interface CharacterTableProps {
  characters: CharacterStats[]
}

export default function CharacterTable({ characters }: CharacterTableProps) {
  return (
    <Table>
      <TableCaption>A list of your characters and their stats.</TableCaption>
      <TableHeader>
        <TableRow>
          <TableHead>Character Name</TableHead>
          <TableHead>Kill Count</TableHead>
          <TableHead>Total Value</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {characters.map((character) => (
          <TableRow key={character.character_id}>
            <TableCell>{character.name}</TableCell>
            <TableCell>{character.kill_count}</TableCell>
            <TableCell>{formatISK(character.total_value)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  )
}