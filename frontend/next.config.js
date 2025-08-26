/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['images.evetech.net'],
  },
  experimental: {
    optimizePackageImports: ['@radix-ui/react-icons', 'lucide-react']
  },
}

module.exports = nextConfig