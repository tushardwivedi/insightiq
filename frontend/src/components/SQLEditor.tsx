'use client'

import { useState, useEffect, useRef } from 'react'
import { Copy, Check } from 'lucide-react'

interface Props {
  value: string
  onChange: (value: string) => void
  disabled?: boolean
}

export default function SQLEditor({ value, onChange, disabled = false }: Props) {
  const [isCopied, setIsCopied] = useState(false)
  const codeRef = useRef<HTMLElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)

  useEffect(() => {
    // Load Prism.js dynamically for syntax highlighting
    const loadPrism = async () => {
      if (typeof window !== 'undefined' && !window.Prism) {
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
    // Re-highlight when value changes
    if (window.Prism && codeRef.current) {
      window.Prism.highlightElement(codeRef.current)
    }
  }, [value])

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(value)
      setIsCopied(true)
      setTimeout(() => setIsCopied(false), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  const handleScroll = () => {
    if (textareaRef.current && codeRef.current) {
      const pre = codeRef.current.parentElement
      if (pre) {
        pre.scrollTop = textareaRef.current.scrollTop
        pre.scrollLeft = textareaRef.current.scrollLeft
      }
    }
  }

  return (
    <div className="sql-editor-container">
      <div className="sql-editor-header">
        <span className="sql-editor-label">SQL Query</span>
        <button
          className="sql-copy-button"
          onClick={handleCopy}
          aria-label="Copy SQL to clipboard"
          title="Copy to clipboard"
          type="button"
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
      </div>

      <div className="sql-editor-wrapper">
        {/* Syntax highlighted display (background) */}
        <pre className="sql-editor-highlight" aria-hidden="true">
          <code ref={codeRef} className="language-sql">
            {value + '\n'}
          </code>
        </pre>

        {/* Actual textarea (foreground, transparent text) */}
        <textarea
          ref={textareaRef}
          value={value}
          onChange={(e) => onChange(e.target.value)}
          onScroll={handleScroll}
          className="sql-editor-textarea"
          disabled={disabled}
          spellCheck="false"
          autoCapitalize="off"
          autoComplete="off"
          autoCorrect="off"
        />
      </div>
    </div>
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
