import { NextResponse } from 'next/server'
import { Character } from '../../../lib/types'

export async function GET() {
  try {
    const response = await fetch('https://ran.api.next.tadeasfort.com/characters')
    if (!response.ok) {
      throw new Error('Failed to fetch characters')
    }
    const characters: Character[] = await response.json()
    return NextResponse.json(characters)
  } catch (error) {
    console.error('Error fetching characters:', error)
    return NextResponse.json({ error: 'Failed to fetch characters' }, { status: 500 })
  }
}

export async function POST(request: Request) {
  try {
    const body = await request.json()
    const response = await fetch('https://ran.api.next.tadeasfort.com/characters', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    })

    if (!response.ok) {
      const errorText = await response.text();
      console.error('Backend API error:', errorText);
      return NextResponse.json({ error: `Backend API error: ${response.status} ${response.statusText}` }, { status: response.status });
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error adding character:', error);
    return NextResponse.json({ error: 'Failed to add character' }, { status: 500 })
  }
}