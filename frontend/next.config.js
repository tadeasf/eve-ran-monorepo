/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    domains: ['images.evetech.net'],
  },
  experimental: {
    optimizePackageImports: ['@radix-ui/react-icons', 'lucide-react']
  },
  webpack: (config, { isServer }) => {
    // Fix for node-fetch warning
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        'node-fetch': false,
      }
    }
    return config
  },
}

module.exports = nextConfig