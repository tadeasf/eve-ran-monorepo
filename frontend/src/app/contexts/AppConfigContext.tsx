'use client'

import React, { createContext, useContext } from 'react'

export interface AppConfig {
  corporation_name: string
}

interface AppConfigContextType {
  config: AppConfig
  loading: boolean
}

const config: AppConfig = {
  corporation_name: process.env.NEXT_PUBLIC_CORPORATION_NAME || 'Tundragon Corp'
}

const AppConfigContext = createContext<AppConfigContextType>({
  config,
  loading: false
})

export function AppConfigProvider({ children }: { children: React.ReactNode }) {
  return (
    <AppConfigContext.Provider value={{ config, loading: false }}>
      {children}
    </AppConfigContext.Provider>
  )
}

export function useAppConfig() {
  return useContext(AppConfigContext)
}
