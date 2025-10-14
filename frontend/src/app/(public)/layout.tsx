import Link from 'next/link'

export default function PublicLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <div className="min-h-screen flex flex-col">
      {/* Navigation Bar */}
      <nav className="landing-nav sticky top-0 z-50">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2">
            <div className="w-8 h-8 rounded-lg flex items-center justify-center" style={{ background: 'var(--accent-color)' }}>
              <span className="text-xl font-bold" style={{ color: 'var(--primary-background)' }}>IQ</span>
            </div>
            <span className="text-xl font-bold" style={{ color: 'var(--text-primary)' }}>InsightIQ</span>
          </Link>

          <div className="flex items-center gap-4">
            <Link
              href="/login"
              className="px-6 py-2 rounded-lg font-medium transition-all"
              style={{
                background: 'var(--accent-color)',
                color: 'var(--primary-background)'
              }}
            >
              Login
            </Link>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="flex-1">
        {children}
      </main>

      {/* Footer */}
      <footer className="border-t" style={{ borderColor: 'var(--border-color)', background: 'var(--surface-color)' }}>
        <div className="container mx-auto px-4 py-8">
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div>
              <h3 className="font-bold mb-3" style={{ color: 'var(--text-primary)' }}>InsightIQ</h3>
              <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
                Self-hosted AI analytics platform for your data
              </p>
            </div>
            <div>
              <h4 className="font-semibold mb-3" style={{ color: 'var(--text-primary)' }}>Resources</h4>
              <ul className="space-y-2 text-sm" style={{ color: 'var(--text-secondary)' }}>
                <li><a href="#" className="hover:underline">Documentation</a></li>
                <li><a href="#" className="hover:underline">Installation Guide</a></li>
                <li><a href="#" className="hover:underline">API Reference</a></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold mb-3" style={{ color: 'var(--text-primary)' }}>Connect</h4>
              <ul className="space-y-2 text-sm" style={{ color: 'var(--text-secondary)' }}>
                <li><a href="#" className="hover:underline">GitHub</a></li>
                <li><a href="#" className="hover:underline">Issues</a></li>
                <li><a href="#" className="hover:underline">License</a></li>
              </ul>
            </div>
          </div>
          <div className="mt-8 pt-8 border-t text-center text-sm" style={{ borderColor: 'var(--border-color)', color: 'var(--text-secondary)' }}>
            © 2024 InsightIQ. Open Source Software.
          </div>
        </div>
      </footer>
    </div>
  )
}
