import { NextResponse } from 'next/server'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

async function fetchCharacterName(characterId: number): Promise<string> {
  try {
    // First try EVE ESI API
    const response = await fetch(`https://esi.evetech.net/latest/characters/${characterId}/`, {
      headers: {
        'User-Agent': 'EVE RAN Application (https://github.com/tadeasfort/eve-ran-monorepo)'
      }
    })
    
    if (response.ok) {
      const data = await response.json() as { name?: string }
      return data.name || 'Unknown'
    }
    
    // Fallback to zkillboard scraping if ESI fails
    const zkbResponse = await fetch(`https://zkillboard.com/character/${characterId}/`, {
      headers: {
        'User-Agent': 'EVE RAN Application (https://github.com/tadeasfort/eve-ran-monorepo)'
      }
    })
    
    if (zkbResponse.ok) {
      const html = await zkbResponse.text()
      const nameMatch = html.match(/<meta name="description" content="([^:]+):/)
      return nameMatch ? nameMatch[1].trim() : 'Unknown'
    }
    
    return 'Unknown'
  } catch (error) {
    console.error(`Failed to fetch character name for ${characterId}:`, error)
    return 'Unknown'
  }
}

export async function GET(
  _request: Request,
  { params }: { params: { id: string } }
) {
  const characterId = params.id
  
  if (!characterId) {
    return NextResponse.json({ error: 'Character ID is required' }, { status: 400 })
  }

  try {
    const name = await fetchCharacterName(parseInt(characterId, 10))
    
    await fetch(`${API_URL}/characters/name/cache`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: characterId, name }),
    })

    return NextResponse.json({ name })
  } catch (error) {
    console.error('Error fetching character name:', error)
    return NextResponse.json({ error: 'Failed to fetch character name' }, { status: 500 })
  }
}