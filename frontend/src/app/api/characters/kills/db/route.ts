import { NextResponse } from 'next/server'

export async function GET(request: Request, { params }: { params: { id: string } }) {
  const { id } = params
  const { searchParams } = new URL(request.url)
  const url = new URL(`https://ran.api.tadeasfort.com/characters/${id}/kills/db`)
  
  searchParams.forEach((value, key) => {
    url.searchParams.append(key, value)
  })

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}