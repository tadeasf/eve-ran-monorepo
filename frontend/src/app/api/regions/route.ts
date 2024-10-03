import { NextResponse } from 'next/server'

export async function GET() {
  try {
    const response = await fetch('https://ran.backend.tadeasfort.com/regions')
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching regions:', error)
    return NextResponse.json({ error: 'Failed to fetch regions' }, { status: 500 })
  }
}