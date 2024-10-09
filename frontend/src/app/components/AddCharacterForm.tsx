

import { useState } from 'react'
import { Button } from "./ui/button"
import { Input } from "./ui/input"

interface AddCharacterFormProps {
  onAddCharacter: (characterId: number) => void
}

export default function AddCharacterForm({ onAddCharacter }: AddCharacterFormProps) {
  const [characterId, setCharacterId] = useState('')

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (characterId) {
      const id = Number(characterId)
      onAddCharacter(id)
      setCharacterId('')
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
    </form>
  )
}