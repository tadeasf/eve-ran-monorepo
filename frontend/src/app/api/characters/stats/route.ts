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