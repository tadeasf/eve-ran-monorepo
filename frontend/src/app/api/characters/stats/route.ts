

import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url)
  const regionIDs = searchParams.getAll('regionID')
  const startDate = searchParams.get('startDate')
  const endDate = searchParams.get('endDate')

  const url = new URL(`${API_URL}/characters/stats`)
  
  regionIDs.forEach(regionID => url.searchParams.append('regionID', regionID))
  if (startDate) url.searchParams.append('startDate', startDate)
  if (endDate) url.searchParams.append('endDate', endDate)

  const response = await fetch(url.toString())
  const data = await response.json()
  return NextResponse.json(data)
}