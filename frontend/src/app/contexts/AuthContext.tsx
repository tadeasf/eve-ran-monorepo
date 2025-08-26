"use client"

import React, { createContext, useContext, useState, useEffect } from 'react'

interface AuthContextType {
    isAuthenticated: boolean
    login: (username: string, password: string) => boolean
    logout: () => void
    loading: boolean
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function AuthProvider({ children }: { children: React.ReactNode }) {
    const [isAuthenticated, setIsAuthenticated] = useState(false)
    const [loading, setLoading] = useState(true)

    useEffect(() => {
        // Check if user is already authenticated on mount
        const authStatus = localStorage.getItem('tundragon_auth')
        if (authStatus === 'true') {
            setIsAuthenticated(true)
        }
        setLoading(false)
    }, [])

    const login = (username: string, password: string): boolean => {
        // Simple authentication against environment variables
        const adminUsername = process.env.NEXT_PUBLIC_ADMIN_USERNAME || 'admin'
        const adminPassword = process.env.NEXT_PUBLIC_ADMIN_PASSWORD || 'TnDrg2024SecCorp'

        if (username === adminUsername && password === adminPassword) {
            setIsAuthenticated(true)
            localStorage.setItem('tundragon_auth', 'true')
            return true
        }
        return false
    }

    const logout = () => {
        setIsAuthenticated(false)
        localStorage.removeItem('tundragon_auth')
    }

    return (
        <AuthContext.Provider value={{ isAuthenticated, login, logout, loading }}>
            {children}
        </AuthContext.Provider>
    )
}

export function useAuth() {
    const context = useContext(AuthContext)
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider')
    }
    return context
}
