'use client'

import { useState, useEffect, useRef } from 'react'
import { Copy, Check, ChevronDown, ChevronUp } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'

interface Props {
  query: string
  defaultCollapsed?: boolean
}

export default function GeneratedQueryViewer({ query, defaultCollapsed = true }: Props) {
  const [isCollapsed, setIsCollapsed] = useState(defaultCollapsed)
  const [isCopied, setIsCopied] = useState(false)
  const codeRef = useRef<HTMLElement>(null)

  useEffect(() => {
    // Load Prism.js dynamically for syntax highlighting
    const loadPrism = async () => {
      if (typeof window !== 'undefined') {
        // Load Prism CSS
        const link = document.createElement('link')
        link.rel = 'stylesheet'
        link.href = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css'
        document.head.appendChild(link)

        // Load Prism JS
        const script = document.createElement('script')
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js'
        script.async = true
        document.body.appendChild(script)

        // Load SQL language support
        const sqlScript = document.createElement('script')
        sqlScript.src = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-sql.min.js'
        sqlScript.async = true
        document.body.appendChild(sqlScript)

        // Highlight code after scripts load
        sqlScript.onload = () => {
          if (window.Prism && codeRef.current) {
            window.Prism.highlightElement(codeRef.current)
          }
        }
      }
    }

    loadPrism()
  }, [])

  useEffect(() => {
    // Re-highlight when query changes
    if (window.Prism && codeRef.current) {
      window.Prism.highlightElement(codeRef.current)
    }
  }, [query])

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(query)
      setIsCopied(true)
      setTimeout(() => setIsCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  return (
    <motion.div
      initial={{ opacity: 0, x: -20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: 0.2 }}
      className="sql-viewer-container"
      data-collapsed={isCollapsed}
    >
      <div
        className="viewer-header"
        onClick={() => setIsCollapsed(!isCollapsed)}
      >
        <h3 className="viewer-title">
          Generated SQL Query
          {isCollapsed && (
            <span className="text-xs ml-2" style={{ color: 'var(--text-secondary)', fontWeight: 'normal' }}>
              (Click to expand)
            </span>
          )}
        </h3>
        <div className="viewer-actions" onClick={(e) => e.stopPropagation()}>
          <button
            className="copy-button"
            onClick={handleCopy}
            aria-label="Copy SQL to clipboard"
            title="Copy to clipboard"
          >
            {isCopied ? (
              <>
                <Check className="w-4 h-4" />
                <span>Copied!</span>
              </>
            ) : (
              <>
                <Copy className="w-4 h-4" />
                <span>Copy</span>
              </>
            )}
          </button>
          <button
            className="collapse-button"
            onClick={() => setIsCollapsed(!isCollapsed)}
            aria-label={isCollapsed ? "Expand query" : "Collapse query"}
          >
            {isCollapsed ? <ChevronDown className="w-5 h-5" /> : <ChevronUp className="w-5 h-5" />}
          </button>
        </div>
      </div>

      <AnimatePresence>
        {!isCollapsed && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.3, ease: 'easeInOut' }}
            className="code-block-wrapper"
          >
            <pre className="sql-pre">
              <code ref={codeRef} className="language-sql">
                {query}
              </code>
            </pre>
          </motion.div>
        )}
      </AnimatePresence>
    </motion.div>
  )
}

// Extend Window interface for Prism
declare global {
  interface Window {
    Prism?: {
      highlightElement: (element: Element) => void
      highlightAll: () => void
    }
  }
}
