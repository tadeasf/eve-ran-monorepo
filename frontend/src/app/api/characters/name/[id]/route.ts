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
import fetch from 'node-fetch'

const API_URL = process.env.NEXT_PUBLIC_API_URL

async function fetchCharacterName(characterId: number): Promise<string> {
  const response = await fetch(`https://zkillboard.com/character/${characterId}/`, {
    headers: {
      'User-Agent': 'EVE RAN Application (https://github.com/tadeasfort/eve-ran-monorepo)'
    }
  })
  if (!response.ok) {
    throw new Error(`Failed to fetch character name: ${response.status} ${response.statusText}`)
  }
  const html = await response.text()
  const nameMatch = html.match(/<meta name="description" content="([^:]+):/)
  return nameMatch ? nameMatch[1].trim() : 'Unknown'
}

export async function GET(
  request: Request,
  { params }: { params: { id: string } }
) {
  const characterId = params.id
  
  if (!characterId) {
    return NextResponse.json({ error: 'Character ID is required' }, { status: 400 })
  }

  try {
    const name = await fetchCharacterName(parseInt(characterId, 10))
    
    await fetch(`${API_URL}/characters/name/cache`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: characterId, name }),
    })

    return NextResponse.json({ name })
  } catch (error) {
    console.error('Error fetching character name:', error)
    return NextResponse.json({ error: 'Failed to fetch character name' }, { status: 500 })
  }
}