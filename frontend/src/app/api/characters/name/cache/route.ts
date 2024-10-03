import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL

export async function POST(request: Request) {
  const { id, name } = await request.json()

  try {
    const response = await fetch(`${API_URL}/characters/name/cache`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, name }),
    })

    if (!response.ok) {
      throw new Error('Failed to cache character name')
    }

    return NextResponse.json({ message: 'Character name cached successfully' })
  } catch (error) {
    console.error('Error caching character name:', error)
    return NextResponse.json({ error: 'Failed to cache character name' }, { status: 500 })
  }
}