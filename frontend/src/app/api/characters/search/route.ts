import { NextResponse } from 'next/server'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(request: Request) {
  try {
    const { searchParams } = new URL(request.url)
    const search = searchParams.get('search')
    
    if (!search) {
      return NextResponse.json({ error: 'Search term is required' }, { status: 400 })
    }

    const response = await fetch(`${API_URL}/characters/search?search=${encodeURIComponent(search)}`)
    
    if (!response.ok) {
      throw new Error(`Backend API error: ${response.status} ${response.statusText}`)
    }
    
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error searching characters:', error)
    return NextResponse.json({ error: 'Failed to search characters' }, { status: 500 })
  }
}