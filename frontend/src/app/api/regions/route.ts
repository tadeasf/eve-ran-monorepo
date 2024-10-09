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
import { Region } from '@/lib/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL

export async function GET(): Promise<NextResponse> {
  try {
    const response = await fetch(`${API_URL}/regions`)
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    const data: Region[] = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error fetching regions:', error)
    return NextResponse.json({ error: 'Failed to fetch regions' }, { status: 500 })
  }
}