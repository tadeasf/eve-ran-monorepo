import { NextResponse } from 'next/server'

export async function GET(
  request: Request,
  { params }: { params: { regionID: string } }
) {
  const { regionID } = params
  const { searchParams } = new URL(request.url)
  const startDate = searchParams.get('startDate')
  const endDate = searchParams.get('endDate')
  
  const url = new URL(`http://localhost:8080/kills/region/${regionID}`)
  if (startDate) url.searchParams.append('startDate', startDate)
  if (endDate) url.searchParams.append('endDate', endDate)

  try {
    const response = await fetch(url.toString())
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching data:', error)
    return NextResponse.json({ error: 'Failed to fetch data from backend' }, { status: 500 })
  }
}