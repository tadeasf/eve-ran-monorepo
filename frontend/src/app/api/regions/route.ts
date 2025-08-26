import { NextResponse } from 'next/server'
import { Region } from '@/lib/types'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(): Promise<NextResponse> {
  try {
    const response = await fetch(`${API_URL}/regions`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data: Region[] = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching regions:', error)
    return NextResponse.json({ error: 'Failed to fetch regions' }, { status: 500 })
  }
}