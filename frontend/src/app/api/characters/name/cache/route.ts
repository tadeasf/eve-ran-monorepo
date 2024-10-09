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

export async function POST(request: Request) {
  const { id, name } = await request.json()

  try {
    const response = await fetch(`${API_URL}/characters/name/cache`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id, name }),
    })

    if (!response.ok) {
      throw new Error('Failed to cache character name')
    }

    return NextResponse.json({ message: 'Character name cached successfully' })
  } catch (error) {
    console.error('Error caching character name:', error)
    return NextResponse.json({ error: 'Failed to cache character name' }, { status: 500 })
  }
}