/** @type {import('next').NextConfig} */
const nextConfig = {
  // Increase timeout for LLM processing (default is 30s)
  experimental: {
    proxyTimeout: 120000, // 120 seconds in milliseconds
  },
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://agent:8080/api/:path*',
      },
    ]
  },
  async headers() {
    return [
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'X-XSS-Protection',
            value: '1; mode=block',
          },
          {
            key: 'Referrer-Policy',
            value: 'strict-origin-when-cross-origin',
          },
          {
            key: 'Content-Security-Policy',
            value: "default-src 'self'; script-src 'self' 'unsafe-eval' 'unsafe-inline' https://cdnjs.cloudflare.com; style-src 'self' 'unsafe-inline' https://cdnjs.cloudflare.com; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' http://localhost:8080 http://127.0.0.1:8080 /api/*; media-src 'self' blob:;",
          },
        ],
      },
    ];
  },
}

export default nextConfig