import { NextResponse } from 'next/server'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function POST(request: Request) {
  try {
    const characterIds = await request.json()
    
    const response = await fetch(`${API_URL}/characters/batch`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(characterIds)
    })
    
    if (!response.ok) {
      throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
    }
    
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error batch adding characters:', error)
    return NextResponse.json({ error: 'Failed to batch add characters' }, { status: 500 })
  }
}