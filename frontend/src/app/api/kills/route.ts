import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url)
  const url = new URL(`${API_URL}/kills`)
  
  searchParams.forEach((value, key) => {
    url.searchParams.append(key, value)
  })

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}