// Copyright (C) 2024 Tadeáš Fořt
// 
// This file is part of EVE Ran Services.
// 
// EVE Ran Services is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// EVE Ran Services is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with EVE Ran Services.  If not, see <https://www.gnu.org/licenses/>.

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