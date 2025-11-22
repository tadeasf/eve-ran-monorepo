import { NextResponse } from 'next/server'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production'
  ? 'http://api:8080'
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const { searchParams } = new URL(request.url)
  const startDate = searchParams.get('startDate')
  const endDate = searchParams.get('endDate')

  const url = new URL(`${API_URL}/characters/${params.id}/analytics`)

  if (startDate) url.searchParams.append('startDate', startDate)
  if (endDate) url.searchParams.append('endDate', endDate)

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}
