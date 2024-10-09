

import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import ClientLayout from "./ClientLayout"
import { SpeedInsights } from "@vercel/speed-insights/next"

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'EVE Ran',
  description: 'EVE Online Character Tracking',
  icons: {
    icon: '/public/icon.png',
  },
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <SpeedInsights />
        <ClientLayout>
          {children}
        </ClientLayout>
      </body>
    </html>
  )
}