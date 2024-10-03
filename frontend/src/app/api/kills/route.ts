import { NextResponse } from 'next/server'

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url)
  const url = new URL('https://ran.api.next.tadeasfort.com/kills')
  
  searchParams.forEach((value, key) => {
    url.searchParams.append(key, value)
  })

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}