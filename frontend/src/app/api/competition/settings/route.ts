import { NextResponse } from 'next/server'

const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(): Promise<NextResponse> {
  try {
    const response = await fetch(`${API_URL}/competition/settings`, {
      cache: 'no-store'
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching competition settings:', error)
    return NextResponse.json({ error: 'Failed to fetch competition settings' }, { status: 500 })
  }
}

export async function POST(request: Request): Promise<NextResponse> {
  try {
    const body = await request.json()
    const response = await fetch(`${API_URL}/competition/settings`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(body)
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error updating competition settings:', error)
    return NextResponse.json({ error: 'Failed to update competition settings' }, { status: 500 })
  }
}
