# Migration Guide - EVE Ran v2.1.0

This guide helps you migrate from the previous version to v2.1.0 with API key authentication and TanStack Query v5.

## Prerequisites

1. **Install Dependencies:**
   ```bash
   cd frontend
   npm install
   ```

2. **Generate API Key:**
   ```bash
   openssl rand -base64 32
   ```

## Backend Migration

### 1. Update Environment Variables

Add to `backend/.env`:
```env
API_KEY=your_generated_api_key_here
REDIS_HOST=redis
REDIS_PORT=6379
```

### 2. Regenerate Swagger Documentation

```bash
cd backend
swag init -g src/main.go
```

### 3. Test Health Endpoints

```bash
# Basic health check (no auth required)
curl http://localhost:12922/health

# Readiness check (no auth required)
curl http://localhost:12922/ready

# Test protected endpoint without key (should fail)
curl http://localhost:12922/characters

# Test protected endpoint with key (should succeed)
curl -H "X-API-Key: your_api_key" http://localhost:12922/characters
```

## Frontend Migration

### 1. Update Environment Variables

Create/update `frontend/.env.local`:
```env
NEXT_PUBLIC_API_URL=https://api.tundragon.space
NEXT_PUBLIC_API_KEY=your_generated_api_key_here
```

### 2. Update Layout to Include QueryProvider

**File:** `frontend/src/app/layout.tsx`

```tsx
import { QueryProvider } from './components/QueryProvider';

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>
        <QueryProvider>
          <ThemeProvider>
            {children}
          </ThemeProvider>
        </QueryProvider>
      </body>
    </html>
  );
}
```

### 3. Migrate API Calls

**Before (direct fetch):**
```tsx
const response = await fetch(`${API_URL}/characters`);
const data = await response.json();
```

**After (with API key):**
```tsx
import { apiFetchJson } from '@/lib/api';

const data = await apiFetchJson('/characters');
```

### 4. Migrate React Query Hooks

**Before (react-query v3):**
```tsx
import { useQuery } from 'react-query';

const { data, isLoading } = useQuery('characters', async () => {
  const res = await fetch(`${API_URL}/characters`);
  return res.json();
}, {
  cacheTime: 300000,
});
```

**After (TanStack Query v5):**
```tsx
import { useQuery } from '@tanstack/react-query';
import { apiFetchJson } from '@/lib/api';

const { data, isLoading } = useQuery({
  queryKey: ['characters'],
  queryFn: () => apiFetchJson<Character[]>('/characters'),
  gcTime: 300000, // Previously cacheTime
});
```

### 5. Update Mutations

**Before:**
```tsx
import { useMutation, useQueryClient } from 'react-query';

const mutation = useMutation(
  (newCharacter) => fetch('/api/characters', {
    method: 'POST',
    body: JSON.stringify(newCharacter),
  }),
  {
    onSuccess: () => {
      queryClient.invalidateQueries('characters');
    },
  }
);
```

**After:**
```tsx
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { apiFetch } from '@/lib/api';

const mutation = useMutation({
  mutationFn: (newCharacter) => apiFetch('/characters', {
    method: 'POST',
    body: JSON.stringify(newCharacter),
  }),
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['characters'] });
  },
});
```

### 6. Add Loading States with Skeletons

**Before:**
```tsx
if (isLoading) return <div>Loading...</div>;
```

**After:**
```tsx
import { CharacterTableSkeleton } from '@/components/ui/skeleton';

if (isLoading) return <CharacterTableSkeleton />;
```

## Key Changes Summary

### Breaking Changes

1. **All API endpoints require X-API-Key header** (except health checks and Swagger)
2. **React Query v3 → TanStack Query v5** syntax changes
3. **Import paths changed:** `react-query` → `@tanstack/react-query`

### New Features

1. **Health check endpoints:** `/health`, `/ready`, `/live`
2. **Data validation:** All ESI/zKillboard responses validated
3. **Connection pooling:** Optimized database performance
4. **Enhanced indexing:** Faster queries
5. **Skeleton loaders:** Better UX during loading
6. **Activity feed:** New page at `/activity`

## Common Migration Issues

### Issue: "Invalid or missing API key"

**Solution:**
- Ensure `API_KEY` in `backend/.env` matches `NEXT_PUBLIC_API_KEY` in `frontend/.env.local`
- Restart both backend and frontend after changing env files
- Check that the key is being sent in the header

### Issue: "Module '@tanstack/react-query' not found"

**Solution:**
```bash
cd frontend
npm install @tanstack/react-query @tanstack/react-query-devtools
```

### Issue: "cacheTime is not a valid option"

**Solution:**
Replace `cacheTime` with `gcTime` in all query configurations.

### Issue: "CORS error when calling API"

**Solution:**
- Check that your origin is in the allowed origins list in `backend/src/main.go`
- Ensure `X-API-Key` is in the allowed headers
- Verify the API URL in `frontend/.env.local` is correct

## Testing Checklist

- [ ] Backend starts without errors
- [ ] Health endpoints return 200 OK
- [ ] API key authentication works (test with/without key)
- [ ] Frontend starts without errors
- [ ] Characters page loads
- [ ] Kills page loads
- [ ] Competition page loads
- [ ] Activity feed page loads
- [ ] Swagger documentation loads
- [ ] All API calls include X-API-Key header

## Rollback Instructions

If you need to rollback:

1. **Remove API key middleware:**
   ```bash
   git checkout main -- backend/src/main.go
   git checkout main -- backend/src/middleware/apikey.go
   ```

2. **Revert React Query:**
   ```bash
   cd frontend
   npm install react-query@^3.39.3
   npm uninstall @tanstack/react-query @tanstack/react-query-devtools
   ```

3. **Restart services:**
   ```bash
   docker-compose restart
   ```

## Support

For issues or questions:
- Check the [Implementation Updates](./IMPLEMENTATION_UPDATES.md) documentation
- Review the [Deployment Guide](./DEPLOYMENT.md)
- Open an issue on GitHub

---

**Migration Guide Version:** 1.0  
**Target Version:** EVE Ran v2.1.0  
**Last Updated:** November 17, 2025
