"use client"

import { useAuth } from '@/app/contexts/AuthContext'
import { LoginForm } from '@/app/components/LoginForm'
import { AdminDashboard } from '@/app/components/AdminDashboard'

export default function AdminPage() {
  const { isAuthenticated, loading } = useAuth()

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-primary"></div>
      </div>
    )
  }

  if (!isAuthenticated) {
    return <LoginForm />
  }

  return <AdminDashboard />
}