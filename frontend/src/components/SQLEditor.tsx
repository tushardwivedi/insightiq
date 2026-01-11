'use client'

import { useState, useEffect, useRef, useCallback } from 'react'
import { Copy, Check, Database, Zap, AlertTriangle, Maximize2, Minimize2 } from 'lucide-react'
import { motion, AnimatePresence } from 'framer-motion'

interface Props {
  value: string
  onChange: (value: string) => void
  disabled?: boolean
}

interface SQLSuggestion {
  text: string
  type: 'keyword' | 'function' | 'table' | 'column'
  description?: string
}

export default function SQLEditor({ value, onChange, disabled = false }: Props) {
  const [isCopied, setIsCopied] = useState(false)
  const [isFullscreen, setIsFullscreen] = useState(false)
  const [showLineNumbers, setShowLineNumbers] = useState(true)
  const [validationErrors, setValidationErrors] = useState<string[]>([])
  const [suggestions, setSuggestions] = useState<SQLSuggestion[]>([])
  const [showSuggestions, setShowSuggestions] = useState(false)
  const codeRef = useRef<HTMLElement>(null)
  const textareaRef = useRef<HTMLTextAreaElement>(null)
  const wrapperRef = useRef<HTMLDivElement>(null)

  const sqlKeywords = [
    'SELECT', 'FROM', 'WHERE', 'AND', 'OR', 'NOT', 'IN', 'LIKE', 'BETWEEN',
    'JOIN', 'LEFT', 'RIGHT', 'INNER', 'OUTER', 'ON', 'AS', 'ORDER', 'BY',
    'GROUP', 'HAVING', 'LIMIT', 'OFFSET', 'UNION', 'ALL', 'DISTINCT',
    'COUNT', 'SUM', 'AVG', 'MIN', 'MAX', 'CASE', 'WHEN', 'THEN', 'ELSE',
    'END', 'NULL', 'IS', 'NULL', 'CAST', 'CONVERT', 'COALESCE', 'IFNULL',
    'ASC', 'DESC', 'TOP', 'WITH', 'TIES', 'OVER', 'PARTITION', 'ROW_NUMBER',
    'RANK', 'DENSE_RANK', 'NTILE', 'LEAD', 'LAG', 'FIRST_VALUE', 'LAST_VALUE'
  ]

  const commonQueries: { name: string; template: string }[] = [
    { name: 'Select All', template: 'SELECT * FROM {table} LIMIT 100' },
    { name: 'Count Records', template: 'SELECT COUNT(*) as total FROM {table}' },
    { name: 'Sum Column', template: 'SELECT SUM({column}) as total FROM {table}' },
    { name: 'Average', template: 'SELECT AVG({column}) as average FROM {table}' },
    { name: 'Group By', template: 'SELECT {column}, COUNT(*) as count FROM {table} GROUP BY {column}' },
    { name: 'Top Records', template: 'SELECT * FROM {table} ORDER BY {column} DESC LIMIT 10' },
    { name: 'Distinct', template: 'SELECT DISTINCT {column} FROM {table}' },
    { name: 'Joins', template: 'SELECT * FROM {table1} t1 JOIN {table2} t2 ON t1.id = t2.{table1}_id' },
  ]

  useEffect(() => {
    const loadPrism = async () => {
      if (typeof window !== 'undefined' && !window.Prism) {
        const link = document.createElement('link')
        link.rel = 'stylesheet'
        link.href = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/themes/prism-tomorrow.min.css'
        document.head.appendChild(link)

        const script = document.createElement('script')
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/prism.min.js'
        script.async = true
        document.body.appendChild(script)

        const sqlScript = document.createElement('script')
        sqlScript.src = 'https://cdnjs.cloudflare.com/ajax/libs/prism/1.29.0/components/prism-sql.min.js'
        sqlScript.async = true
        document.body.appendChild(sqlScript)

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
    if (window.Prism && codeRef.current) {
      window.Prism.highlightElement(codeRef.current)
    }
  }, [value])

  const validateSQL = useCallback((sql: string) => {
    const errors: string[] = []
    const lowerSql = sql.toLowerCase()

    if (sql.trim()) {
      if (!sql.toUpperCase().startsWith('SELECT')) {
        errors.push('Only SELECT queries are allowed for security reasons')
      }

      const dangerousKeywords = ['drop', 'delete', 'truncate', 'alter', 'create', 'insert', 'update', 'grant', 'revoke']
      const foundDangerous = dangerousKeywords.filter(kw => lowerSql.includes(kw))
      if (foundDangerous.length > 0) {
        errors.push(`Restricted keywords found: ${foundDangerous.join(', ')}`)
      }

      const unclosedQuotes = (sql.match(/'/g) || []).length % 2 !== 0
      if (unclosedQuotes) {
        errors.push('Unclosed single quote detected')
      }

      const unclosedParens = (sql.match(/\(/g) || []).length !== (sql.match(/\)/g) || []).length
      if (unclosedParens) {
        errors.push('Unbalanced parentheses detected')
      }
    }

    setValidationErrors(errors)
    return errors
  }, [])

  useEffect(() => {
    validateSQL(value)
  }, [value, validateSQL])

  const getCurrentWord = () => {
    const cursorPosition = textareaRef.current?.selectionStart || 0
    const textBeforeCursor = value.substring(0, cursorPosition)
    const words = textBeforeCursor.split(/[\s,()]+/)
    return words[words.length - 1] || ''
  }

  const updateSuggestions = useCallback(() => {
    const word = getCurrentWord()
    if (word.length < 1) {
      setSuggestions([])
      setShowSuggestions(false)
      return
    }

    const matches = sqlKeywords
      .filter(kw => kw.toLowerCase().startsWith(word.toLowerCase()))
      .map(kw => ({ text: kw, type: 'keyword' as const }))

    setSuggestions(matches.slice(0, 5))
    setShowSuggestions(matches.length > 0)
  }, [value])

  useEffect(() => {
    updateSuggestions()
  }, [updateSuggestions])

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

  const getLineCount = () => {
    return value.split('\n').length
  }

  const insertTemplate = (template: string) => {
    onChange(template)
    textareaRef.current?.focus()
  }

  const lineCount = getLineCount()

  const EditorContainer = isFullscreen ? 'fixed inset-0 z-50' : 'relative'
  const EditorMaxHeight = isFullscreen ? 'h-screen' : 'max-h-96'

  return (
    <motion.div
      layout
      className={`${EditorContainer} sql-editor-container`}
    >
      {isFullscreen && (
        <div
          className="absolute inset-0 bg-black/50 backdrop-blur-sm"
          onClick={() => setIsFullscreen(false)}
        />
      )}

      <motion.div
        layout
        className={`relative bg-[#1e1e1e] rounded-lg overflow-hidden ${isFullscreen ? 'h-full' : ''}`}
        style={{ height: isFullscreen ? '100%' : 'auto' }}
      >
        <div className="flex items-center justify-between px-4 py-2 bg-[#2d2d2d] border-b border-[#404040]">
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2">
              <div className="flex gap-1.5">
                <div className="w-3 h-3 rounded-full bg-[#ff5f56]" />
                <div className="w-3 h-3 rounded-full bg-[#ffbd2e]" />
                <div className="w-3 h-3 rounded-full bg-[#27c93f]" />
              </div>
            </div>
            <div className="flex items-center gap-2 text-sm text-gray-400 ml-2">
              <Database className="w-4 h-4" />
              <span>SQL Query</span>
            </div>
          </div>

          <div className="flex items-center gap-2">
            <AnimatePresence>
              {validationErrors.length > 0 && (
                <motion.div
                  initial={{ opacity: 0, scale: 0.9 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.9 }}
                  className="flex items-center gap-1 px-2 py-1 rounded bg-red-500/20 text-red-400 text-xs"
                >
                  <AlertTriangle className="w-3 h-3" />
                  {validationErrors.length} issue{validationErrors.length > 1 ? 's' : ''}
                </motion.div>
              )}
            </AnimatePresence>

            <div className="relative group">
              <button
                className="p-1.5 rounded hover:bg-[#404040] transition-colors text-gray-400 hover:text-white"
                title="Query Templates"
              >
                <Zap className="w-4 h-4" />
              </button>
              <div className="absolute right-0 top-full mt-1 w-48 bg-[#2d2d2d] border border-[#404040] rounded-lg shadow-xl opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-10">
                <div className="p-2 text-xs text-gray-400 border-b border-[#404040]">
                  Quick Templates
                </div>
                {commonQueries.map((q, i) => (
                  <button
                    key={i}
                    onClick={() => insertTemplate(q.template)}
                    className="w-full px-3 py-2 text-left text-sm text-gray-300 hover:bg-[#404040] hover:text-white transition-colors flex items-center gap-2"
                  >
                    <span className="text-primary text-xs">{q.name}</span>
                  </button>
                ))}
              </div>
            </div>

            <button
              onClick={() => setShowLineNumbers(!showLineNumbers)}
              className={`p-1.5 rounded transition-colors ${showLineNumbers ? 'text-primary bg-primary/20' : 'text-gray-400 hover:text-white hover:bg-[#404040]'}`}
              title="Toggle line numbers"
            >
              <span className="text-xs font-mono">{lineCount}</span>
            </button>

            <button
              onClick={() => setIsFullscreen(!isFullscreen)}
              className="p-1.5 rounded hover:bg-[#404040] transition-colors text-gray-400 hover:text-white"
              title={isFullscreen ? "Exit fullscreen" : "Fullscreen"}
            >
              {isFullscreen ? <Minimize2 className="w-4 h-4" /> : <Maximize2 className="w-4 h-4" />}
            </button>

            <button
              onClick={handleCopy}
              className="flex items-center gap-1.5 px-3 py-1.5 rounded bg-[#404040] hover:bg-[#505050] transition-colors text-gray-300 hover:text-white text-sm"
            >
              {isCopied ? (
                <>
                  <Check className="w-4 h-4 text-green-400" />
                  <span className="text-green-400">Copied!</span>
                </>
              ) : (
                <>
                  <Copy className="w-4 h-4" />
                  <span>Copy</span>
                </>
              )}
            </button>
          </div>
        </div>

        <div className={`flex ${EditorMaxHeight} overflow-auto`}>
          {showLineNumbers && (
            <div className="flex-shrink-0 w-10 py-4 bg-[#252525] text-center border-r border-[#404040] select-none">
              {Array.from({ length: lineCount }, (_, i) => (
                <div
                  key={i}
                  className="text-xs text-gray-500 font-mono leading-6 hover:text-gray-300"
                >
                  {i + 1}
                </div>
              ))}
            </div>
          )}

          <div className="flex-1 relative min-w-0">
            <pre className="absolute inset-0 m-0 p-4 bg-[#1e1e1e] font-mono text-sm leading-6 overflow-auto" aria-hidden="true">
              <code ref={codeRef} className="language-sql">
                {value || '\n'}
              </code>
            </pre>

            <textarea
              ref={textareaRef}
              value={value}
              onChange={(e) => onChange(e.target.value)}
              onScroll={handleScroll}
              className="absolute inset-0 w-full h-full p-4 bg-transparent border-none outline-none resize-none font-mono text-sm leading-6 text-transparent caret-primary"
              disabled={disabled}
              spellCheck="false"
              autoCapitalize="off"
              autoComplete="off"
              autoCorrect="off"
              style={{
                color: 'transparent',
                caretColor: '#4fd1c5',
                whiteSpace: 'pre-wrap',
                wordWrap: 'break-word',
              }}
            />
          </div>
        </div>

        <AnimatePresence>
          {validationErrors.length > 0 && (
            <motion.div
              initial={{ opacity: 0, height: 0 }}
              animate={{ opacity: 1, height: 'auto' }}
              exit={{ opacity: 0, height: 0 }}
              className="border-t border-[#404040]"
            >
              {validationErrors.map((error, i) => (
                <div
                  key={i}
                  className="flex items-center gap-2 px-4 py-2 text-sm text-red-400 bg-red-500/10 border-l-2 border-red-500"
                >
                  <AlertTriangle className="w-4 h-4 flex-shrink-0" />
                  <span>{error}</span>
                </div>
              ))}
            </motion.div>
          )}
        </AnimatePresence>

        <div className="flex items-center justify-between px-4 py-2 bg-[#252525] border-t border-[#404040] text-xs text-gray-500">
          <div className="flex items-center gap-4">
            <span>{value.length} characters</span>
            <span>{lineCount} lines</span>
          </div>
          <div className="flex items-center gap-2">
            <span>SELECT only for security</span>
          </div>
        </div>
      </motion.div>
    </motion.div>
  )
}

declare global {
  interface Window {
    Prism?: {
      highlightElement: (element: Element) => void
      highlightAll: () => void
    }
  }
}
