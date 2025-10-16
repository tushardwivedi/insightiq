'use client'

import { useState } from 'react'
import { Brain, Mail, Lock, User, AlertCircle, Github } from 'lucide-react'
import Link from 'next/link'
import { useRouter } from 'next/navigation'
import { signUp, doesEmailExist } from 'supertokens-web-js/recipe/emailpassword'
import { getAuthorisationURLWithQueryParamsAndSetState } from 'supertokens-web-js/recipe/thirdparty'

export default function SignupPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [name, setName] = useState('')
  const [error, setError] = useState('')
  const [isLoading, setIsLoading] = useState(false)
  const router = useRouter()

  const handleEmailSignup = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setIsLoading(true)

    try {
      // Check if email already exists
      const emailExists = await doesEmailExist({ email })

      if (emailExists.doesExist) {
        setError('An account with this email already exists')
        setIsLoading(false)
        return
      }

      // Sign up with email/password using SuperTokens
      const response = await signUp({
        formFields: [
          { id: 'email', value: email },
          { id: 'password', value: password },
        ],
      })

      if (response.status === 'OK') {
        // Signup successful, redirect to app
        router.push('/app')
      } else if (response.status === 'FIELD_ERROR') {
        setError(response.formFields[0].error)
      } else {
        setError('Signup failed. Please try again.')
      }
    } catch (err: any) {
      setError(err.message || 'Signup failed. Please try again.')
    } finally {
      setIsLoading(false)
    }
  }

  const handleGitHubSignup = async () => {
    setError('')
    setIsLoading(true)

    try {
      const authUrl = await getAuthorisationURLWithQueryParamsAndSetState({
        thirdPartyId: 'github',
        frontendRedirectURI: `${window.location.origin}/auth/callback/github`,
      })

      // Redirect to GitHub OAuth
      window.location.href = authUrl
    } catch (err: any) {
      setError('GitHub signup failed. Please try again.')
      setIsLoading(false)
    }
  }

  const handleGoogleSignup = async () => {
    setError('')
    setIsLoading(true)

    try {
      const authUrl = await getAuthorisationURLWithQueryParamsAndSetState({
        thirdPartyId: 'google',
        frontendRedirectURI: `${window.location.origin}/auth/callback/google`,
      })

      // Redirect to Google OAuth
      window.location.href = authUrl
    } catch (err: any) {
      setError('Google signup failed. Please try again.')
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4">
      <div className="max-w-md w-full">
        <div className="text-center mb-8">
          <div className="inline-flex items-center gap-2 mb-4">
            <div className="w-12 h-12 rounded-xl flex items-center justify-center" style={{ background: 'var(--accent-color)' }}>
              <Brain className="w-7 h-7" style={{ color: 'var(--primary-background)' }} />
            </div>
          </div>
          <h1 className="text-3xl font-bold mb-2" style={{ color: 'var(--text-primary)' }}>
            Create your account
          </h1>
          <p style={{ color: 'var(--text-secondary)' }}>
            Get started with InsightIQ analytics
          </p>
        </div>

        <div className="card p-8">
          {error && (
            <div className="mb-6 p-4 rounded-lg flex items-start gap-3" style={{ background: 'rgba(252, 129, 129, 0.1)', border: '1px solid var(--error-color)' }}>
              <AlertCircle className="w-5 h-5 flex-shrink-0" style={{ color: 'var(--error-color)' }} />
              <span style={{ color: 'var(--error-color)' }}>{error}</span>
            </div>
          )}

          {/* Social Login Buttons */}
          <div className="space-y-3 mb-6">
            <button
              type="button"
              onClick={handleGitHubSignup}
              disabled={isLoading}
              className="w-full py-3 px-4 rounded-lg font-medium transition-all flex items-center justify-center gap-2"
              style={{
                background: '#24292e',
                color: 'white',
                border: '1px solid #444d56'
              }}
            >
              <Github className="w-5 h-5" />
              Continue with GitHub
            </button>

            <button
              type="button"
              onClick={handleGoogleSignup}
              disabled={isLoading}
              className="w-full py-3 px-4 rounded-lg font-medium transition-all flex items-center justify-center gap-2"
              style={{
                background: 'white',
                color: '#1f1f1f',
                border: '1px solid #dadce0'
              }}
            >
              <svg className="w-5 h-5" viewBox="0 0 24 24">
                <path fill="#4285F4" d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"/>
                <path fill="#34A853" d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"/>
                <path fill="#FBBC05" d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"/>
                <path fill="#EA4335" d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"/>
              </svg>
              Continue with Google
            </button>
          </div>

          <div className="relative mb-6">
            <div className="absolute inset-0 flex items-center">
              <div className="w-full border-t" style={{ borderColor: 'var(--border-color)' }}></div>
            </div>
            <div className="relative flex justify-center text-sm">
              <span className="px-4" style={{ background: 'var(--card-background)', color: 'var(--text-secondary)' }}>
                Or sign up with email
              </span>
            </div>
          </div>

          {/* Email/Password Signup Form */}
          <form onSubmit={handleEmailSignup} className="space-y-4">
            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>
                Full Name
              </label>
              <div className="relative">
                <User className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5" style={{ color: 'var(--text-secondary)' }} />
                <input
                  type="text"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  className="w-full pl-11 pr-4 py-3 rounded-lg"
                  style={{
                    background: 'var(--primary-background)',
                    color: 'var(--text-primary)',
                    border: '1px solid var(--border-color)'
                  }}
                  placeholder="John Doe"
                  required
                  disabled={isLoading}
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>
                Email Address
              </label>
              <div className="relative">
                <Mail className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5" style={{ color: 'var(--text-secondary)' }} />
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  className="w-full pl-11 pr-4 py-3 rounded-lg"
                  style={{
                    background: 'var(--primary-background)',
                    color: 'var(--text-primary)',
                    border: '1px solid var(--border-color)'
                  }}
                  placeholder="you@example.com"
                  required
                  disabled={isLoading}
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium mb-2" style={{ color: 'var(--text-primary)' }}>
                Password
              </label>
              <div className="relative">
                <Lock className="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5" style={{ color: 'var(--text-secondary)' }} />
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="w-full pl-11 pr-4 py-3 rounded-lg"
                  style={{
                    background: 'var(--primary-background)',
                    color: 'var(--text-primary)',
                    border: '1px solid var(--border-color)'
                  }}
                  placeholder="At least 8 characters"
                  required
                  minLength={8}
                  disabled={isLoading}
                />
              </div>
            </div>

            <button
              type="submit"
              disabled={isLoading || !email || !password || !name}
              className="btn-primary w-full py-3 rounded-lg font-semibold transition-all"
            >
              {isLoading ? (
                <div className="flex items-center justify-center gap-2">
                  <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-white"></div>
                  Creating account...
                </div>
              ) : (
                'Create Account'
              )}
            </button>
          </form>

          <div className="mt-6 text-center">
            <p className="text-sm" style={{ color: 'var(--text-secondary)' }}>
              Already have an account?{' '}
              <Link href="/login" className="font-medium hover:underline" style={{ color: 'var(--accent-color)' }}>
                Sign in
              </Link>
            </p>
          </div>
        </div>

        <div className="mt-6 text-center">
          <Link href="/" className="text-sm hover:underline" style={{ color: 'var(--text-secondary)' }}>
            ← Back to home
          </Link>
        </div>
      </div>
    </div>
  )
}
