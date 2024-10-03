'use client'

import { useQuery, useQueryClient, useMutation } from 'react-query'
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../components/ui/table"
import { Button } from "../components/ui/button"
import { Trash2 } from "lucide-react"
import AddCharacterForm from '../components/AddCharacterForm'
import { Character } from '../../lib/types'
import { Skeleton } from "../components/ui/skeleton"
import { Progress } from "../components/ui/progress"

const fetchCharacters = async (): Promise<Character[]> => {
  const response = await fetch('/api/characters')
  if (!response.ok) {
    throw new Error('Failed to fetch characters')
  }
  return response.json()
}

export default function Characters() {
  const queryClient = useQueryClient()
  const { data: characters = [], isLoading, error } = useQuery<Character[]>('characters', fetchCharacters)

  const deleteMutation = useMutation(
    (characterId: number) => fetch(`/api/characters/${characterId}`, { method: 'DELETE' }),
    {
      onSuccess: () => queryClient.invalidateQueries('characters'),
    }
  )

  const addMutation = useMutation(
    async (characterId: number) => {
      const response = await fetch('/api/characters', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: characterId }),
      })
      if (!response.ok) {
        const errorText = await response.text();
        console.error('Error response:', errorText);
        throw new Error(`Failed to add character: ${response.status} ${response.statusText}`);
      }
      return response.json()
    },
    {
      onSuccess: () => queryClient.invalidateQueries('characters'),
      onError: (error) => {
        console.error('Mutation error:', error);
      }
    }
  )

  if (error) return <div>Error: {(error as Error).message}</div>

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">Character Management</h1>
      {isLoading ? (
        <Progress value={75} className="w-full mb-4" />
      ) : (
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Character ID</TableHead>
              <TableHead>Character Name</TableHead>
              <TableHead>zKillboard Link</TableHead>
              <TableHead>Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {characters.map((character) => (
              <TableRow key={character.id}>
                <TableCell>{character.id}</TableCell>
                <TableCell>{character.name}</TableCell>
                <TableCell>
                  <a
                    href={`https://zkillboard.com/character/${character.id}/`}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="text-blue-500 hover:underline"
                  >
                    View on zKillboard
                  </a>
                </TableCell>
                <TableCell>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => deleteMutation.mutate(character.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      )}
      <div className="mt-8">
        <h2 className="text-2xl font-bold mb-4">Add New Character</h2>
        {isLoading ? (
          <Skeleton className="h-[100px] w-full" />
        ) : (
          <AddCharacterForm onAddCharacter={(id) => addMutation.mutate(id)} />
        )}
      </div>
    </div>
  )
}