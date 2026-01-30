import Link from 'next/link'
import { Github } from 'lucide-react'

export function Footer() {
  return (
    <footer className="relative bg-[#0F172A] border-t border-[#4FD1C5]/10">
      <div className="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-[#4FD1C5]/50 to-transparent" />

      <div className="container mx-auto px-4 py-12">
        <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
          <div className="md:col-span-1">
            <Link href="/" className="flex items-center gap-3 mb-4">
              <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-[#4FD1C5] to-[#38B2AC] flex items-center justify-center">
                <span className="text-lg font-bold text-[#0F172A]">IQ</span>
              </div>
              <span className="text-xl font-bold text-white">InsightIQ</span>
            </Link>
            <p className="text-sm text-[#A0AEC0] leading-relaxed">
              Self-hosted AI analytics platform. Chat with your databases using natural language, voice, and SQL.
            </p>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-4">Product</h4>
            <ul className="space-y-3 text-sm">
              <li>
                <Link href="/#features" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                  Features
                </Link>
              </li>
              <li>
                <Link href="/#installation" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                  Installation
                </Link>
              </li>
              <li>
                <Link href="/docs" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                  Documentation
                </Link>
              </li>
              <li>
                <Link href="/docs/api" className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors">
                  API Reference
                </Link>
              </li>
            </ul>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-4">Resources</h4>
            <ul className="space-y-3 text-sm">
              <li>
                <a
                  href="https://github.com/tushardwivedi/insightiq"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors flex items-center gap-2"
                >
                  <Github className="w-4 h-4" />
                  GitHub
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/tushardwivedi/insightiq/issues"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors"
                >
                  Report Issues
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/tushardwivedi/insightiq/blob/main/LICENSE"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors"
                >
                  MIT License
                </a>
              </li>
            </ul>
          </div>

          <div>
            <h4 className="font-semibold text-white mb-4">Powered By</h4>
            <div className="flex flex-wrap gap-2">
              {['Next.js', 'Go', 'PostgreSQL', 'Ollama', 'Whisper', 'Qdrant'].map((tech) => (
                <span
                  key={tech}
                  className="px-3 py-1 text-xs rounded-full bg-[#2D3748]/50 text-[#A0AEC0] border border-[#4A5568]/30"
                >
                  {tech}
                </span>
              ))}
            </div>
          </div>
        </div>

        <div className="mt-12 pt-8 border-t border-[#4A5568]/30 flex flex-col md:flex-row items-center justify-between gap-4">
          <p className="text-sm text-[#A0AEC0]">
            &copy; {new Date().getFullYear()} InsightIQ. Open Source Software under MIT License.
          </p>
          <div className="flex items-center gap-6">
            <a
              href="https://github.com/tushardwivedi/insightiq"
              target="_blank"
              rel="noopener noreferrer"
              className="text-[#A0AEC0] hover:text-[#4FD1C5] transition-colors"
            >
              <Github className="w-5 h-5" />
            </a>
          </div>
        </div>
      </div>
    </footer>
  )
}
