import { NextResponse } from 'next/server'

export async function DELETE(request: Request, { params }: { params: { id: string } }) {
  const { id } = params
  const response = await fetch(`http://localhost:8080/characters/${id}`, {
    method: 'DELETE',
  })
  if (response.ok) {
    return NextResponse.json({ message: 'Character deleted successfully' })
  } else {
    return NextResponse.json({ error: 'Failed to delete character' }, { status: response.status })
  }
}