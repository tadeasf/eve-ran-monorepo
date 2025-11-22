'use client'

import { QueryClient, QueryClientProvider } from 'react-query'
import { ReactQueryDevtools } from 'react-query/devtools'
import { ThemeProvider } from "./components/ThemeProvider"
import { AuthProvider } from "./contexts/AuthContext"
import { AppConfigProvider } from "./contexts/AppConfigContext"
import { MainNav } from './components/NavigationMenu'

const queryClient = new QueryClient()

export default function ClientLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <QueryClientProvider client={queryClient}>
      <AppConfigProvider>
        <AuthProvider>
          <ThemeProvider
            attribute="class"
            defaultTheme="dark"
            enableSystem
            disableTransitionOnChange
          >
            <header className="bg-background border-b">
              <div className="container mx-auto">
                <MainNav />
              </div>
            </header>
            <main className="container mx-auto py-6">
              {children}
            </main>
          </ThemeProvider>
        </AuthProvider>
      </AppConfigProvider>
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}