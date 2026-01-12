import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { AuthProvider } from '@/contexts/AuthContext'
import { ThemeProvider } from '@/contexts/ThemeContext'
import SuperTokensProvider from '@/components/SuperTokensProvider'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'InsightIQ - Analytics Voice Agent',
  description: 'AI-powered analytics with voice interaction',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={inter.className}>
        <SuperTokensProvider>
          <AuthProvider>
            <ThemeProvider>
              <div className="min-h-screen" style={{ background: 'var(--primary-background)', color: 'var(--text-primary)' }}>
                {children}
              </div>
            </ThemeProvider>
          </AuthProvider>
        </SuperTokensProvider>
      </body>
    </html>
  )
}