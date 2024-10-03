import { NextResponse } from 'next/server'

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url)
  const regionIDs = searchParams.getAll('regionID')
  const startDate = searchParams.get('startDate')
  const endDate = searchParams.get('endDate')

  const url = new URL('http://localhost:8080/characters/stats')
  
  regionIDs.forEach(regionID => url.searchParams.append('regionID', regionID))
  if (startDate) url.searchParams.append('startDate', startDate)
  if (endDate) url.searchParams.append('endDate', endDate)

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}