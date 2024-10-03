import { NextResponse } from 'next/server'
import { Character } from '../../../lib/types'

async function fetchCharacterName(characterId: number): Promise<string> {
  const response = await fetch(`https://zkillboard.com/character/${characterId}/`)
  if (!response.ok) {
    throw new Error('Failed to fetch character name')
  }
  const html = await response.text()
  const nameMatch = html.match(/<meta name="description" content="([^:]+):/)
  return nameMatch ? nameMatch[1].trim() : 'Unknown'
}

export async function GET() {
  try {
    const response = await fetch('https://ran.api.next.tadeasfort.com/characters')
    const characters: Character[] = await response.json()

    const charactersWithNames = await Promise.all(
      characters.map(async (character) => {
        if (!character.name) {
          const name = await fetchCharacterName(character.id)
          // Update the character with the new name in your backend
          await fetch(`https://ran.api.next.tadeasfort.com/characters/${character.id}`, {
            method: 'PATCH',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ name }),
          })
          return { ...character, name }
        }
        return character
      })
    )

    return NextResponse.json(charactersWithNames)
  } catch (error) {
    console.error('Error fetching characters:', error)
    return NextResponse.json({ error: 'Failed to fetch characters' }, { status: 500 })
  }
}

export async function POST(request: Request) {
  try {
    const body = await request.json()
    const { id } = body

    // Fetch the character name from zKillboard
    const name = await fetchCharacterName(id)

    // Add the character with the fetched name to your backend
    const response = await fetch('https://ran.api.next.tadeasfort.com/characters', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ...body, name }),
    })

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error adding character:', error)
    return NextResponse.json({ error: 'Failed to add character' }, { status: 500 })
  }
}