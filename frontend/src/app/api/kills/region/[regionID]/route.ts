// Copyright (C) 2024 Tadeáš Fořt
// 
// This file is part of EVE Ran Services.
// 
// EVE Ran Services is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
// 
// EVE Ran Services is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with EVE Ran Services.  If not, see <https://www.gnu.org/licenses/>.

import { NextResponse } from 'next/server'
import { Kill } from '@/lib/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL
const FETCH_TIMEOUT = 120000 // 120 seconds timeout

export async function GET(
  request: Request,
  { params }: { params: { regionID: string } }
): Promise<NextResponse> {
  const { regionID } = params
  const { searchParams } = new URL(request.url)
  const startDate = searchParams.get('startDate')
  const endDate = searchParams.get('endDate')
  
  const url = new URL(`${API_URL}/kills/region/${regionID}`)
  if (startDate) url.searchParams.append('startDate', startDate)
  if (endDate) url.searchParams.append('endDate', endDate)

  console.log(`Fetching data from: ${url.toString()}`)

  try {
    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), FETCH_TIMEOUT)

    const response = await fetch(url.toString(), {
      signal: controller.signal,
    })

    clearTimeout(timeoutId)

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }

    const data: Kill[] = await response.json()
    console.log(`Received data for region ${regionID}:`, data.length ? `${data.length} kills` : 'No kills')

    return NextResponse.json(data)
  } catch (error: unknown) {
    if (error instanceof Error && error.name === 'AbortError') {
      console.error(`Request timed out after ${FETCH_TIMEOUT / 1000} seconds`)
      return NextResponse.json({ error: 'Request timed out' }, { status: 504 })
    }
    console.error('Error fetching data:', error)
    return NextResponse.json({ error: 'Failed to fetch data from backend' }, { status: 500 })
  }
}