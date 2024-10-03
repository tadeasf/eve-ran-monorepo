'use client'

import { QueryClient, QueryClientProvider } from 'react-query'
import { ReactQueryDevtools } from 'react-query/devtools'
import { ThemeProvider } from "./components/ThemeProvider"
import { MainNav } from './components/NavigationMenu'

const queryClient = new QueryClient()

export default function ClientLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <QueryClientProvider client={queryClient}>
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
      <ReactQueryDevtools initialIsOpen={false} />
    </QueryClientProvider>
  )
}