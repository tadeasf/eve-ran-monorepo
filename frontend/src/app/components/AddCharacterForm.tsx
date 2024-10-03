import { useState } from 'react'
import { Button } from "./ui/button"
import { Input } from "./ui/input"

interface AddCharacterFormProps {
  onAddCharacter: (characterId: number) => void
  getCharacterName: (characterId: number) => Promise<string>
}

export default function AddCharacterForm({ onAddCharacter, getCharacterName }: AddCharacterFormProps) {
  const [characterId, setCharacterId] = useState('')
  const [characterName, setCharacterName] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    if (characterId) {
      try {
        const id = Number(characterId)
        const name = await getCharacterName(id)
        setCharacterName(name)
        onAddCharacter(id)
        setCharacterId('')
      } catch (error) {
        console.error('Error fetching character name:', error)
        setCharacterName('Error: Unable to fetch character name')
      }
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col space-y-2 mb-6">
      <div className="flex items-center space-x-2">
        <Input
          type="number"
          value={characterId}
          onChange={(e) => setCharacterId(e.target.value)}
          placeholder="Enter Character ID"
        />
        <Button type="submit">Add Character</Button>
      </div>
      {characterName && (
        <p className="text-sm text-gray-500">Character Name: {characterName}</p>
      )}
    </form>
  )
}