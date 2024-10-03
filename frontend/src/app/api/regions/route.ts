import { NextResponse } from 'next/server'

export async function GET() {
  const response = await fetch('https://ran.api.tadeasfort.com/regions')
  const data = await response.json()
  return NextResponse.json(data)
}