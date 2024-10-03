import { NextResponse } from 'next/server'

async function fetchCharacterName(characterId: number): Promise<string> {
  const response = await fetch(`https://zkillboard.com/character/${characterId}/`)
  if (!response.ok) {
    throw new Error('Failed to fetch character name')
  }
  const html = await response.text()
  const nameMatch = html.match(/<meta name="description" content="([^:]+):/)
  return nameMatch ? nameMatch[1].trim() : 'Unknown'
}

export async function GET(
  request: Request,
  { params }: { params: { characterId: string } }
) {
  const characterId = parseInt(params.characterId, 10)
  
  try {
    const name = await fetchCharacterName(characterId)
    return NextResponse.json({ name })
  } catch (error) {
    console.error('Error fetching character name:', error)
    return NextResponse.json({ error: 'Failed to fetch character name' }, { status: 500 })
  }
}