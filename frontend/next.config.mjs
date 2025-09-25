/** @type {import('next').NextConfig} */
const nextConfig = {
  async rewrites() {
    return [
      {
        source: '/api/:path*',
        destination: 'http://agent:8080/api/:path*',
      },
    ]
  },
}

export default nextConfig