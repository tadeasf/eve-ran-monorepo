# EVE Ran v2.1.0 - Deployment Checklist

## Pre-Deployment Checklist

### 1. Environment Configuration ⚙️

- [ ] Generate secure API key: `openssl rand -base64 32`
- [ ] Update `backend/.env` with:
  - [ ] Database credentials
  - [ ] Redis configuration
  - [ ] Generated API_KEY
- [ ] Update `frontend/.env.local` with:
  - [ ] NEXT_PUBLIC_API_URL
  - [ ] NEXT_PUBLIC_API_KEY (same as backend)
- [ ] Verify environment files are in `.gitignore`

### 2. Backend Setup 🔧

- [ ] Install Go dependencies: `cd backend && go mod download`
- [ ] Regenerate Swagger documentation: `swag init -g src/main.go`
- [ ] Test compilation: `go build -o main src/main.go`
- [ ] Run database migrations (automatic on startup)
- [ ] Verify health endpoints:
  - [ ] `curl http://localhost:12922/health`
  - [ ] `curl http://localhost:12922/ready`
- [ ] Test API key authentication:
  - [ ] Without key (should fail): `curl http://localhost:12922/characters`
  - [ ] With key (should succeed): `curl -H "X-API-Key: your_key" http://localhost:12922/characters`

### 3. Frontend Setup 💻

- [ ] Install Node dependencies: `cd frontend && npm install`
- [ ] Verify TanStack Query packages:
  - [ ] `@tanstack/react-query`
  - [ ] `@tanstack/react-query-devtools`
- [ ] Build frontend: `npm run build`
- [ ] Test development server: `npm run dev`
- [ ] Verify API calls include X-API-Key header

### 4. Database 🗄️

- [ ] PostgreSQL is running
- [ ] Database credentials are correct
- [ ] Connection pool limits are appropriate for your server
- [ ] Verify new indexes are created (check logs on startup)
- [ ] Test database connectivity via `/ready` endpoint
- [ ] Backup database before deployment

### 5. Redis 📦

- [ ] Redis is running
- [ ] Redis host/port configured correctly
- [ ] Test Redis connectivity via `/ready` endpoint
- [ ] Verify cache TTL settings (default: 5 minutes)

### 6. Docker Deployment 🐳

- [ ] Update `docker-compose.yml` if needed
- [ ] Build images: `docker-compose build`
- [ ] Start services: `docker-compose up -d`
- [ ] Check logs: `docker-compose logs -f`
- [ ] Verify all containers are running: `docker-compose ps`

### 7. API Documentation 📚

- [ ] Swagger UI accessible at `/swagger/index.html`
- [ ] Test API endpoints via Swagger
- [ ] Verify API key field shows in Swagger UI
- [ ] All new endpoints documented:
  - [ ] GET /kills (with pagination)
  - [ ] DELETE /characters/batch
  - [ ] GET /kills/{killmail_id}/comments
  - [ ] POST /kills/{killmail_id}/comments
  - [ ] PUT /comments/{id}
  - [ ] DELETE /comments/{id}

### 8. Testing 🧪

**Backend Tests:**
- [ ] Health endpoints respond correctly
- [ ] Database queries return paginated results
- [ ] Ship type filtering works
- [ ] Bulk character delete works
- [ ] Kill comments CRUD operations work
- [ ] API key authentication works
- [ ] CORS headers are set correctly

**Frontend Tests:**
- [ ] Activity feed page loads
- [ ] Infinite scroll loads more kills
- [ ] Ship type filter works (when integrated)
- [ ] Skeleton loaders appear during loading
- [ ] API calls include authentication header
- [ ] No console errors

**Integration Tests:**
- [ ] Frontend can fetch kills from backend
- [ ] Pagination works correctly
- [ ] Comments can be created and displayed
- [ ] Filters work correctly
- [ ] CORS doesn't block requests

### 9. Performance Verification 📊

- [ ] Database query times improved (check logs)
- [ ] Connection pool is being used (check metrics)
- [ ] Infinite scroll doesn't cause memory leaks
- [ ] Cached endpoints return quickly
- [ ] Initial page load is fast

### 10. Security Audit 🔒

- [ ] API key is strong (32+ characters)
- [ ] API key is not in version control
- [ ] API key required for protected endpoints
- [ ] Health checks don't require authentication
- [ ] CORS only allows trusted domains
- [ ] No sensitive data in logs
- [ ] Database passwords are strong
- [ ] Redis is secured (if exposed)

### 11. Monitoring Setup 📈

- [ ] Health check endpoints configured in load balancer
- [ ] Kubernetes liveness probe: `/live`
- [ ] Kubernetes readiness probe: `/ready`
- [ ] Log aggregation configured
- [ ] Error tracking configured (if using Sentry/etc)
- [ ] Performance monitoring configured

### 12. Documentation Review 📝

- [ ] README.md updated with new features
- [ ] IMPLEMENTATION_UPDATES.md reviewed
- [ ] MIGRATION_GUIDE.md available for team
- [ ] SESSION_SUMMARY.md archived
- [ ] Deployment guide updated if needed

### 13. Reverse Proxy Configuration 🌐

**Caddy Example:**
```caddy
tundragon.space {
    reverse_proxy localhost:12921
}

api.tundragon.space {
    reverse_proxy localhost:12922
}
```

- [ ] Frontend route configured
- [ ] Backend API route configured
- [ ] HTTPS/SSL certificates valid
- [ ] CORS origins match reverse proxy domains

### 14. Rollback Plan 🔄

- [ ] Previous version tagged in git
- [ ] Database backup created
- [ ] Rollback procedure documented
- [ ] Team knows how to revert changes

---

## Post-Deployment Verification ✅

### Immediate Checks (within 5 minutes)

- [ ] Frontend loads without errors
- [ ] Backend health checks respond
- [ ] Database connectivity confirmed
- [ ] Redis connectivity confirmed
- [ ] API authentication works
- [ ] No critical errors in logs

### Short-term Checks (within 1 hour)

- [ ] Users can view kills
- [ ] Infinite scroll works
- [ ] Comments can be added
- [ ] Dashboard displays correctly
- [ ] No memory leaks observed
- [ ] Response times acceptable

### Long-term Monitoring (within 24 hours)

- [ ] No errors in error tracking
- [ ] Database performance stable
- [ ] Cache hit rate acceptable
- [ ] API response times good
- [ ] User feedback is positive

---

## Troubleshooting Common Issues 🔧

### "Invalid or missing API key"
1. Check API_KEY in backend/.env matches frontend/.env.local
2. Restart both services after changing env files
3. Verify X-API-Key header is being sent

### "Cannot find module '@tanstack/react-query'"
1. Delete node_modules: `rm -rf node_modules`
2. Delete package-lock.json: `rm package-lock.json`
3. Reinstall: `npm install`

### Database connection fails
1. Check PostgreSQL is running
2. Verify credentials in .env
3. Check network connectivity
4. Review database logs

### Redis connection fails
1. Check Redis is running: `redis-cli ping`
2. Verify host/port in .env
3. Check network connectivity
4. Review Redis logs

### CORS errors
1. Verify origin is in allowed origins list
2. Check reverse proxy configuration
3. Verify NEXT_PUBLIC_API_URL is correct
4. Check browser developer console for exact error

---

## Emergency Contacts 📞

- **Database Admin:** [contact info]
- **DevOps Lead:** [contact info]
- **Backend Lead:** [contact info]
- **Frontend Lead:** [contact info]

---

## Rollback Procedure 🔴

If critical issues occur:

1. **Stop Services:**
   ```bash
   docker-compose down
   ```

2. **Restore Database:**
   ```bash
   # Restore from backup
   psql -U tdg_admin -d tdg_data < backup.sql
   ```

3. **Checkout Previous Version:**
   ```bash
   git checkout <previous-tag>
   ```

4. **Restart Services:**
   ```bash
   docker-compose up -d
   ```

5. **Verify:**
   - [ ] Services are running
   - [ ] Frontend loads
   - [ ] API responds
   - [ ] Database accessible

---

## Success Criteria ✨

Deployment is successful when:

- ✅ All health checks pass
- ✅ No critical errors in logs
- ✅ Users can access the application
- ✅ API authentication works
- ✅ Database queries are fast
- ✅ Infinite scroll works smoothly
- ✅ Comments can be added/viewed
- ✅ Response times < 500ms for most endpoints
- ✅ No memory leaks after 1 hour
- ✅ Team is satisfied with performance

---

**Deployment Date:** ___________________  
**Deployed By:** ___________________  
**Version:** v2.1.0  
**Sign-off:** ___________________  

---

*Keep this checklist for future reference and deployment audits.*
