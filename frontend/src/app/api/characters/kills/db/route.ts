import { NextResponse } from 'next/server'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production' 
  ? 'http://api:8080' 
  : process.env.NEXT_PUBLIC_API_URL

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url)
  const id = searchParams.get('id')
  
  if (!id) {
    return NextResponse.json({ error: 'Character ID is required' }, { status: 400 })
  }

  const url = new URL(`${API_URL}/characters/${id}/kills/db`)
  
  searchParams.forEach((value, key) => {
    if (key !== 'id') {
      url.searchParams.append(key, value)
    }
  })

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}