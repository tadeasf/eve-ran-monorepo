#!/bin/bash

# Fix all incorrect UI component import paths in the frontend

cd /home/tadeas/eve-ran-monorepo/frontend

echo "Fixing UI component import paths..."

# Find all TypeScript/TSX files and replace the paths
find src -type f \( -name "*.ts" -o -name "*.tsx" \) -exec sed -i \
  -e 's|@/app/components/ui/|@/components/ui/|g' \
  -e 's|"\.\./components/ui/|"@/components/ui/|g' \
  -e "s|'\.\./components/ui/|'@/components/ui/|g" \
  {} +

echo "✓ Fixed all UI component import paths"
echo "Paths changed:"
echo "  @/app/components/ui/* → @/components/ui/*"
echo "  ../components/ui/* → @/components/ui/*"
