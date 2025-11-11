import { NextResponse } from 'next/server'

const API_URL = process.env.NODE_ENV === 'production'
  ? 'http://api:8080'
  : process.env.NEXT_PUBLIC_API_URL

export async function GET() {
  try {
    const response = await fetch(`${API_URL}/systems`)
    if (!response.ok) {
      throw new Error('Failed to fetch systems')
    }
    const systems = await response.json()
    return NextResponse.json(systems)
  } catch (error) {
    console.error('Error fetching systems:', error)
    return NextResponse.json({ error: 'Failed to fetch systems' }, { status: 500 })
  }
}
