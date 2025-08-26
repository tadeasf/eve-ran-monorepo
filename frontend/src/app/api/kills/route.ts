import { NextResponse } from 'next/server'
import { Kill } from '@/lib/types'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(request: Request): Promise<NextResponse> {
  const { searchParams } = new URL(request.url)
  const url = new URL(`${API_URL}/kills`)
  
  searchParams.forEach((value, key) => {
    url.searchParams.append(key, value)
  })

  try {
    const response = await fetch(url.toString())
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data: Kill[] = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching kills:', error)
    return NextResponse.json({ error: 'Failed to fetch kills' }, { status: 500 })
  }
}