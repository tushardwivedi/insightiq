'use client'

import { useState, useEffect } from 'react'
import { DataConnector, ConnectorTestResult } from '@/types'
import { apiClient } from '@/lib/api'
import ConnectorCard from './ConnectorCard'
import SupersetConnectorForm from './connectors/SupersetConnectorForm'

interface ConnectorSidebarProps {
  isOpen: boolean
  onClose: () => void
}

export default function ConnectorSidebar({ isOpen, onClose }: ConnectorSidebarProps) {
  const [connectors, setConnectors] = useState<DataConnector[]>([])
  const [loading, setLoading] = useState(false)
  const [showAddForm, setShowAddForm] = useState(false)
  const [selectedConnectorType, setSelectedConnectorType] = useState<string>('')
  const [isHovered, setIsHovered] = useState(false)

  useEffect(() => {
    if (isOpen) {
      loadConnectors()
    }
  }, [isOpen])

  const loadConnectors = async () => {
    try {
      setLoading(true)
      const response = await apiClient.getConnectors()
      setConnectors(response.data || [])
    } catch (error) {
      console.error('Failed to load connectors:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleConnectorAdded = (connector: DataConnector) => {
    setConnectors(prev => [...prev, connector])
    setShowAddForm(false)
    setSelectedConnectorType('')
  }

  const handleConnectorUpdated = (updatedConnector: DataConnector) => {
    setConnectors(prev =>
      prev.map(c => c.id === updatedConnector.id ? updatedConnector : c)
    )
  }

  const handleConnectorDeleted = (connectorId: string) => {
    setConnectors(prev => prev.filter(c => c.id !== connectorId))
  }

  const handleTestConnection = async (connector: DataConnector) => {
    try {
      setConnectors(prev =>
        prev.map(c =>
          c.id === connector.id ? { ...c, status: 'testing' } : c
        )
      )

      const result = await apiClient.testConnector(connector.id)

      setConnectors(prev =>
        prev.map(c =>
          c.id === connector.id
            ? {
                ...c,
                status: result.success ? 'connected' : 'error',
                last_tested: new Date().toISOString()
              }
            : c
        )
      )
    } catch (error) {
      console.error('Connection test failed:', error)
      setConnectors(prev =>
        prev.map(c =>
          c.id === connector.id ? { ...c, status: 'error' } : c
        )
      )
    }
  }

  const connectorTypes = [
    {
      type: 'superset',
      name: 'Apache Superset',
      description: 'Connect to Apache Superset for analytics dashboards',
      icon: '📊',
    },
    {
      type: 'postgres',
      name: 'PostgreSQL',
      description: 'Connect to PostgreSQL database',
      icon: '🐘',
    },
    {
      type: 'mysql',
      name: 'MySQL',
      description: 'Connect to MySQL database',
      icon: '🐬',
    },
    {
      type: 'api',
      name: 'REST API',
      description: 'Connect to external REST API',
      icon: '🌐',
    },
  ]

  return (
    <>
      {/* Mobile Backdrop - only shown when explicitly opened on mobile */}
      {isOpen && (
        <div
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={onClose}
        />
      )}

      {/* Collapsible Sidebar with Hover Effect */}
      <div
        className={`
          fixed top-0 left-0 h-full shadow-xl z-50 transform transition-all duration-300 ease-in-out
          ${isOpen ? 'w-80 translate-x-0' : isHovered ? 'w-80 translate-x-0' : 'w-16 translate-x-0'}
          lg:relative lg:z-auto
          hover:shadow-2xl
        `}
        style={{ background: 'var(--surface-color)' }}
        onMouseEnter={() => setIsHovered(true)}
        onMouseLeave={() => setIsHovered(false)}
      >
        <div className="flex flex-col h-full">
          {/* Header */}
          <div className="flex items-center justify-between p-4 border-b" style={{ borderColor: 'var(--border-color)' }}>
            {/* Collapsed State - Show only icon */}
            {!isOpen && !isHovered ? (
              <div className="flex items-center justify-center w-full">
                <svg className="w-6 h-6" style={{ color: 'var(--accent-color)' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 1.79 4 4 4h8c2.21 0 4-1.79 4-4V7c0-2.21-1.79-4-4-4H8c-2.21 0-4 1.79-4 4z" />
                </svg>
              </div>
            ) : (
              /* Expanded State - Show full header */
              <>
                <h2 className="text-lg font-semibold" style={{ color: 'var(--text-primary)' }}>Data Connectors</h2>
                <button
                  onClick={onClose}
                  className="p-2 rounded-md transition-colors lg:hidden"
                  style={{ color: 'var(--text-primary)' }}
                  onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'var(--hover-surface)'}
                  onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'transparent'}
                >
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </>
            )}
          </div>

          {/* Content */}
          <div className="flex-1 overflow-y-auto p-4">
            {/* Collapsed State - Show minimal content */}
            {!isOpen && !isHovered ? (
              <div className="flex flex-col items-center space-y-4">
                {/* Collapsed Add Button - Icon only */}
                <button
                  onClick={() => setShowAddForm(true)}
                  className="w-8 h-8 rounded-lg transition-colors flex items-center justify-center"
                  style={{ background: 'var(--accent-color)', color: 'var(--primary-background)' }}
                  title="Add Data Source"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                  </svg>
                </button>

                {/* Collapsed Connector Status Indicators */}
                {connectors.slice(0, 3).map((connector) => (
                  <div
                    key={connector.id}
                    className="w-8 h-8 rounded-lg flex items-center justify-center text-xs font-medium"
                    style={{
                      backgroundColor: connector.status === 'connected' ? '#10b981' :
                                     connector.status === 'testing' ? '#f59e0b' : '#ef4444',
                      color: 'white'
                    }}
                    title={`${connector.name} - ${connector.status}`}
                  >
                    {connector.name.charAt(0).toUpperCase()}
                  </div>
                ))}

                {connectors.length > 3 && (
                  <div className="text-xs text-center" style={{ color: 'var(--text-secondary)' }}>
                    +{connectors.length - 3}
                  </div>
                )}
              </div>
            ) : (
              /* Expanded State - Show full content */
              !showAddForm ? (
                <>
                  {/* Add Connector Button */}
                  <button
                    onClick={() => setShowAddForm(true)}
                    className="btn-primary w-full mb-4 px-4 py-3 rounded-lg transition-colors flex items-center justify-center gap-2"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    Add Data Source
                  </button>

                  {/* Connectors List */}
                  {loading ? (
                    <div className="flex items-center justify-center py-8">
                      <div className="animate-spin rounded-full h-8 w-8 border-b-2" style={{ borderColor: 'var(--accent-color)' }}></div>
                    </div>
                  ) : connectors.length === 0 ? (
                    <div className="text-center py-8" style={{ color: 'var(--text-secondary)' }}>
                      <svg className="w-12 h-12 mx-auto mb-4" style={{ color: 'var(--border-color)' }} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 7v10c0 2.21 1.79 4 4 4h8c2.21 0 4-1.79 4-4V7c0-2.21-1.79-4-4-4H8c-2.21 0-4 1.79-4 4z" />
                      </svg>
                      <p className="text-sm">No data sources connected</p>
                      <p className="text-xs mt-1" style={{ color: 'var(--text-secondary)' }}>Add your first connector to get started</p>
                    </div>
                  ) : (
                    <div className="space-y-3">
                      {connectors.map((connector) => (
                        <ConnectorCard
                          key={connector.id}
                          connector={connector}
                          onTest={() => handleTestConnection(connector)}
                          onEdit={(connector) => handleConnectorUpdated(connector)}
                          onDelete={() => handleConnectorDeleted(connector.id)}
                        />
                      ))}
                    </div>
                  )}
                </>
              ) : (
                /* Add Form - only show when expanded */
                <>
                  {/* Back Button */}
                  <button
                    onClick={() => {
                      setShowAddForm(false)
                      setSelectedConnectorType('')
                    }}
                    className="mb-4 flex items-center gap-2 transition-colors"
                    style={{ color: 'var(--text-secondary)' }}
                    onMouseEnter={(e) => e.currentTarget.style.color = 'var(--text-primary)'}
                    onMouseLeave={(e) => e.currentTarget.style.color = 'var(--text-secondary)'}
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                    </svg>
                    Back to connectors
                  </button>

                  {!selectedConnectorType ? (
                    <>
                      <h3 className="text-md font-medium mb-4" style={{ color: 'var(--text-primary)' }}>Choose Connector Type</h3>
                      <div className="space-y-3">
                        {connectorTypes.map((type) => (
                          <button
                            key={type.type}
                            onClick={() => setSelectedConnectorType(type.type)}
                            className="w-full p-4 text-left border rounded-lg transition-colors"
                            style={{ borderColor: 'var(--border-color)' }}
                            onMouseEnter={(e) => {
                              e.currentTarget.style.borderColor = 'var(--accent-color)'
                              e.currentTarget.style.background = 'var(--hover-surface)'
                            }}
                            onMouseLeave={(e) => {
                              e.currentTarget.style.borderColor = 'var(--border-color)'
                              e.currentTarget.style.background = 'transparent'
                            }}
                          >
                            <div className="flex items-start gap-3">
                              <span className="text-2xl">{type.icon}</span>
                              <div className="flex-1">
                                <h4 className="font-medium" style={{ color: 'var(--text-primary)' }}>{type.name}</h4>
                                <p className="text-sm mt-1" style={{ color: 'var(--text-secondary)' }}>{type.description}</p>
                              </div>
                            </div>
                          </button>
                        ))}
                      </div>
                    </>
                  ) : selectedConnectorType === 'superset' ? (
                    <SupersetConnectorForm
                      onCancel={() => setSelectedConnectorType('')}
                      onSuccess={handleConnectorAdded}
                    />
                  ) : (
                    <div className="text-center py-8" style={{ color: 'var(--text-secondary)' }}>
                      <p>Connector form for {selectedConnectorType} coming soon!</p>
                      <button
                        onClick={() => setSelectedConnectorType('')}
                        className="mt-4 px-4 py-2 rounded-md transition-colors"
                        style={{
                          background: 'var(--hover-surface)',
                          color: 'var(--text-primary)',
                          border: '1px solid var(--border-color)'
                        }}
                      >
                        Go Back
                      </button>
                    </div>
                  )}
                </>
              )
            )}
          </div>

          {/* Footer */}
          <div className="p-4 border-t" style={{ borderColor: 'var(--border-color)' }}>
            {!isOpen && !isHovered ? (
              /* Collapsed Footer - Show count only */
              <div className="text-xs text-center font-medium" style={{ color: 'var(--text-secondary)' }}>
                {connectors.length}
              </div>
            ) : (
              /* Expanded Footer - Show full text */
              <div className="text-xs text-center" style={{ color: 'var(--text-secondary)' }}>
                {connectors.length} connector{connectors.length !== 1 ? 's' : ''} configured
              </div>
            )}
          </div>
        </div>
      </div>
    </>
  )
}