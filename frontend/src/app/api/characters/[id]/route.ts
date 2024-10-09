

import { NextResponse } from 'next/server'

const API_URL = process.env.NEXT_PUBLIC_API_URL

export async function DELETE(request: Request, { params }: { params: { id: string } }) {
  const { id } = params
  const response = await fetch(`${API_URL}/characters/${id}`, {
    method: 'DELETE',
  })
  if (response.ok) {
    return NextResponse.json({ message: 'Character deleted successfully' })
  } else {
    return NextResponse.json({ error: 'Failed to delete character' }, { status: response.status })
  }
}