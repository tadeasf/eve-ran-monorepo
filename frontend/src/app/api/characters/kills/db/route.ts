

import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL

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