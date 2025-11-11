import { NextResponse } from 'next/server'

// Use internal service name in production, external URL in development
const API_URL = process.env.NODE_ENV === 'production'
  ? 'http://api:8080'
  : process.env.NEXT_PUBLIC_API_URL

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