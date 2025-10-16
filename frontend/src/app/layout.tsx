import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { AuthProvider } from '@/contexts/AuthContext'
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
    <html lang="en">
      <body className={inter.className}>
        <SuperTokensProvider>
          <AuthProvider>
            <div className="min-h-screen" style={{ background: 'var(--primary-background)', color: 'var(--text-primary)' }}>
              {children}
            </div>
          </AuthProvider>
        </SuperTokensProvider>
      </body>
    </html>
  )
}