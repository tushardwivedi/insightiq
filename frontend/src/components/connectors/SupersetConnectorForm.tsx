'use client'

import { useState } from 'react'
import { DataConnector, SupersetConnectorConfig } from '@/types'
import { apiClient } from '@/lib/api'

interface SupersetConnectorFormProps {
  onCancel: () => void
  onSuccess: (connector: DataConnector) => void
  connector?: DataConnector
}

export default function SupersetConnectorForm({ onCancel, onSuccess, connector }: SupersetConnectorFormProps) {
  const [formData, setFormData] = useState({
    name: connector?.name || '',
    url: connector?.config.url || '',
    username: (connector?.config as SupersetConnectorConfig)?.username || '',
    password: (connector?.config as SupersetConnectorConfig)?.password || '',
    bearer_token: (connector?.config as SupersetConnectorConfig)?.bearer_token || '',
  })
  const [loading, setLoading] = useState(false)
  const [testing, setTesting] = useState(false)
  const [testResult, setTestResult] = useState<{ success: boolean; message: string } | null>(null)
  const [authMethod, setAuthMethod] = useState<'credentials' | 'token'>('credentials')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setTestResult(null)

    try {
      const config: SupersetConnectorConfig = {
        url: formData.url.trim(),
        username: authMethod === 'credentials' ? formData.username.trim() : '',
        password: authMethod === 'credentials' ? formData.password : '',
        bearer_token: authMethod === 'token' ? formData.bearer_token.trim() : '',
      }

      const connectorData = {
        name: formData.name.trim(),
        type: 'superset' as const,
        config,
      }

      let result: DataConnector
      if (connector) {
        result = await apiClient.updateConnector(connector.id, connectorData)
      } else {
        result = await apiClient.createConnector(connectorData)
      }

      onSuccess(result)
    } catch (error: any) {
      console.error('Failed to save connector:', error)
      setTestResult({
        success: false,
        message: error.message || 'Failed to save connector. Please try again.'
      })
    } finally {
      setLoading(false)
    }
  }

  const handleTestConnection = async () => {
    setTesting(true)
    setTestResult(null)

    try {
      const config: SupersetConnectorConfig = {
        url: formData.url.trim(),
        username: authMethod === 'credentials' ? formData.username.trim() : '',
        password: authMethod === 'credentials' ? formData.password : '',
        bearer_token: authMethod === 'token' ? formData.bearer_token.trim() : '',
      }

      const result = await apiClient.testConnectorConfig({
        type: 'superset',
        config,
      })

      setTestResult({
        success: result.success,
        message: result.message || (result.success ? 'Connection successful!' : 'Connection failed')
      })
    } catch (error: any) {
      console.error('Connection test failed:', error)
      setTestResult({
        success: false,
        message: error.message || 'Connection test failed. Please check your settings.'
      })
    } finally {
      setTesting(false)
    }
  }

  const isFormValid = () => {
    const isValid = (
      formData.name.trim() &&
      formData.url.trim() &&
      (authMethod === 'credentials'
        ? formData.username.trim() && formData.password.trim()
        : formData.bearer_token.trim()
      )
    )
    console.log('Form validation:', {
      name: formData.name.trim(),
      url: formData.url.trim(),
      authMethod,
      username: formData.username.trim(),
      password: formData.password.trim(),
      bearer_token: formData.bearer_token.trim(),
      isValid
    })
    return isValid
  }

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-bold !text-black mb-1">
          {connector ? 'Edit' : 'Add'} Superset Connector
        </h3>
        <p className="text-sm font-medium text-gray-800">
          Connect to your Apache Superset instance for analytics data
        </p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-4">
        {/* Name */}
        <div>
          <label htmlFor="name" className="block text-sm font-bold !text-black mb-1">
            Connection Name *
          </label>
          <input
            type="text"
            id="name"
            value={formData.name}
            onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
            placeholder="My Superset Instance"
            className="w-full px-3 py-2 border-2 border-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white !text-black placeholder:text-gray-600 font-medium"
            required
          />
        </div>

        {/* URL */}
        <div>
          <label htmlFor="url" className="block text-sm font-bold !text-black mb-1">
            Superset URL *
          </label>
          <input
            type="url"
            id="url"
            value={formData.url}
            onChange={(e) => setFormData(prev => ({ ...prev, url: e.target.value }))}
            placeholder="https://superset.example.com"
            className="w-full px-3 py-2 border-2 border-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white !text-black placeholder:text-gray-600 font-medium"
            required
          />
        </div>

        {/* Authentication Method */}
        <div>
          <label className="block text-sm font-bold !text-black mb-2">
            Authentication Method
          </label>
          <div className="space-y-2">
            <label className="flex items-center">
              <input
                type="radio"
                value="credentials"
                checked={authMethod === 'credentials'}
                onChange={(e) => setAuthMethod(e.target.value as 'credentials')}
                className="mr-2"
              />
              <span className="text-sm font-medium !text-black">Username & Password</span>
            </label>
            <label className="flex items-center">
              <input
                type="radio"
                value="token"
                checked={authMethod === 'token'}
                onChange={(e) => setAuthMethod(e.target.value as 'token')}
                className="mr-2"
              />
              <span className="text-sm font-medium !text-black">Bearer Token</span>
            </label>
          </div>
        </div>

        {/* Authentication Fields */}
        {authMethod === 'credentials' ? (
          <>
            <div>
              <label htmlFor="username" className="block text-sm font-bold !text-black mb-1">
                Username *
              </label>
              <input
                type="text"
                id="username"
                value={formData.username}
                onChange={(e) => setFormData(prev => ({ ...prev, username: e.target.value }))}
                placeholder="admin"
                className="w-full px-3 py-2 border-2 border-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white !text-black placeholder:text-gray-600 font-medium"
                required={authMethod === 'credentials'}
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-bold !text-black mb-1">
                Password *
              </label>
              <input
                type="password"
                id="password"
                value={formData.password}
                onChange={(e) => setFormData(prev => ({ ...prev, password: e.target.value }))}
                placeholder="••••••••"
                className="w-full px-3 py-2 border-2 border-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white !text-black placeholder:text-gray-600 font-medium"
                required={authMethod === 'credentials'}
              />
            </div>
          </>
        ) : (
          <div>
            <label htmlFor="bearer_token" className="block text-sm font-bold !text-black mb-1">
              Bearer Token *
            </label>
            <textarea
              id="bearer_token"
              value={formData.bearer_token}
              onChange={(e) => setFormData(prev => ({ ...prev, bearer_token: e.target.value }))}
              placeholder="eyJhbGciOiJIUzI1NiIsInR5cCI6..."
              rows={3}
              className="w-full px-3 py-2 border-2 border-gray-400 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 resize-none bg-white !text-black placeholder:text-gray-600 font-medium"
              required={authMethod === 'token'}
            />
            <p className="text-xs font-medium text-gray-700 mt-1">
              Generate a bearer token from Superset's API settings
            </p>
          </div>
        )}

        {/* Test Result */}
        {testResult && (
          <div className={`p-3 rounded-md ${testResult.success ? 'bg-green-50 border border-green-200' : 'bg-red-50 border border-red-200'}`}>
            <div className="flex items-center gap-2">
              <span className={testResult.success ? 'text-green-600' : 'text-red-600'}>
                {testResult.success ? '✓' : '✗'}
              </span>
              <span className={`text-sm ${testResult.success ? 'text-green-800' : 'text-red-800'}`}>
                {testResult.message}
              </span>
            </div>
          </div>
        )}

        {/* Actions */}
        <div className="flex gap-3 pt-4">
          <button
            type="button"
            onClick={handleTestConnection}
            disabled={!isFormValid() || testing}
            className="flex-1 px-4 py-2 border border-blue-600 text-blue-600 rounded-md hover:bg-blue-50 disabled:bg-gray-100 disabled:text-gray-400 disabled:border-gray-300 disabled:cursor-not-allowed transition-colors"
          >
            {testing ? 'Testing...' : 'Test Connection'}
          </button>
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 border border-gray-300 text-gray-700 rounded-md hover:bg-gray-50 transition-colors"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={!isFormValid() || loading}
            className="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            {loading ? 'Saving...' : connector ? 'Update' : 'Save'}
          </button>
        </div>
      </form>
    </div>
  )
}