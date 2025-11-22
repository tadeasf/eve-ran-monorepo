import Script from 'next/script'
import './globals.css'
import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import ClientLayout from "./ClientLayout"
import { SpeedInsights } from "@vercel/speed-insights/next"

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: process.env.NEXT_PUBLIC_CORPORATION_NAME || 'Tundragon Corp',
  description: 'EVE Online Intelligence Hub - Character Management & Analytics',
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
      <Script
        defer
        src="https://analytics.zizcon.cz/script.js"
        data-website-id="740cd735-4fc5-4568-b96f-cc9f2cadde5b"
      />
    </html>
  )
}
