import { NextResponse } from 'next/server'

const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(): Promise<NextResponse> {
  try {
    const response = await fetch(`${API_URL}/competition/ytd-winners`, {
      cache: 'no-store'
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching year-to-date winners:', error)
    return NextResponse.json({ error: 'Failed to fetch year-to-date winners' }, { status: 500 })
  }
}
