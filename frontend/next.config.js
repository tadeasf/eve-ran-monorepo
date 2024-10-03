/** @type {import('next').NextConfig} */
const nextConfig = {
    async rewrites() {
      return [
        {
          source: '/api/:path*',
          destination: 'https://ran.backend.tadeasfort.com/:path*',
        },
      ]
    },
  }
  
  module.exports = nextConfig