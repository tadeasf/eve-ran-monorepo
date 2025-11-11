import { NextResponse } from 'next/server'

const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(
  request: Request,
  { params }: { params: Promise<{ month: string; year: string }> }
): Promise<NextResponse> {
  try {
    const { month, year } = await params
    const response = await fetch(`${API_URL}/competition/history/${month}/${year}`, {
      cache: 'no-store'
    })
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching competition history:', error)
    return NextResponse.json({ error: 'Failed to fetch competition history' }, { status: 500 })
  }
}
