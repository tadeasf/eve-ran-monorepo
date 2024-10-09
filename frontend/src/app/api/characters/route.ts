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
import { Character } from '../../../lib/types'

const API_URL = process.env.NEXT_PUBLIC_API_URL

export async function GET() {
  try {
    const response = await fetch(`${API_URL}/characters`)
    if (!response.ok) {
      throw new Error('Failed to fetch characters')
    }
    const characters: Character[] = await response.json()
    return NextResponse.json(characters)
  } catch (error) {
    console.error('Error fetching characters:', error)
    return NextResponse.json({ error: 'Failed to fetch characters' }, { status: 500 })
  }
}

export async function POST(request: Request) {
  try {
    const body = await request.json()
    const response = await fetch(`${API_URL}/characters`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: body.id }),
    })

    if (!response.ok) {
      if (response.status === 409) {
        return NextResponse.json({ error: 'Character already exists' }, { status: 409 })
      } else if (response.status === 500) {
        return NextResponse.json({ error: 'Server error occurred' }, { status: 500 })
      } else {
        return NextResponse.json({ error: `HTTP error! status: ${response.status}` }, { status: response.status })
      }
    }

    const data = await response.json()
    return NextResponse.json(data)
  } catch (error) {
    console.error('Error adding character:', error)
    return NextResponse.json({ error: 'Failed to add character' }, { status: 500 })
  }
}