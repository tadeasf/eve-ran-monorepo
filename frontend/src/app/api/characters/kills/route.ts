import { NextResponse } from 'next/server'

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url)
  const id = searchParams.get('id')

  if (!id) {
    return NextResponse.json({ error: 'Character ID is required' }, { status: 400 })
  }

  const url = new URL(`https://ran.api.tadeasfort.com/characters/${id}/kills`)
  
  searchParams.forEach((value, key) => {
    if (key !== 'id') {
      url.searchParams.append(key, value)
    }
  })

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}