# EVE ESI Rate Limiting Implementation (Updated 2025-11-22)

## Overview

Comprehensive update to support ESI's new floating window rate limiting with token-based bucket system, preventing 429 errors.

## Problem Addressed

Application was receiving 429 errors:
```
ERROR: ESI returned status 429 for killmail 67876058
```

## ESI Rate Limiting System (November 2025)

Per [ESI Documentation](https://developers.eveonline.com/docs/services/esi/rate-limiting/):

### Floating Window Rate Limiting
- **Token costs by status code:**
  - 2XX: 2 tokens
  - 3XX: 1 token (promotes If-Modified-Since)
  - 4XX: 5 tokens (discourages errors)
  - 5XX: 0 tokens (server errors)

### New Rate Limit Headers
- `X-Ratelimit-Group`: Route group identifier
- `X-Ratelimit-Limit`: Total tokens per window (e.g., "150/15m")
- `X-Ratelimit-Remaining`: Available tokens
- `X-Ratelimit-Used`: Tokens consumed by request
- `Retry-After`: Seconds to wait after 429

### Legacy Error Limiting
Still active on all routes:
- Max 100 non-2xx/3xx per minute
- Headers: `X-ESI-Error-Limit-Remain`, `X-ESI-Error-Limit-Reset`
- Returns 420 when exceeded

## Updated Implementation

### ESI Error Manager (`backend/src/services/esiErrorManager.go`)

**New Structures:**
```go
type ESIRateLimitInfo struct {
    group         string    // Route group
    limit         int       // Total tokens per window
    remaining     int       // Available tokens
    used          int       // Tokens used by last request
    windowSize    string    // e.g., "15m", "1h"
    retryAfter    time.Time // When to retry after 429
}
```

**Enhanced Features:**
- **Per-group tracking:** Separate limits for each route group
- **Token-based accounting:** Monitors tokens used/remaining
- **Dual limit tracking:** New rate limiting + legacy error limiting
- **Intelligent backoff:** Progressive slowdown as limits approach
  - Error limit < 20: 500ms delay
  - Route group < 10% remaining: 1 second delay
  - Route group < 30% remaining: 300ms delay

**Key Methods:**
- `UpdateLimitsFromHeaders()`: Parses ALL rate limit headers (new + legacy)
- `CanMakeRequest()`: Checks both systems, maintains safety buffers
- `ShouldBackoff()`: Returns whether to slow down and by how much
- `WaitForReset()`: Waits until rate limits reset with buffer
- `GetRateLimitInfo(group)`: Gets info for specific route group
- `LogStatus()`: Comprehensive logging of all limits

### ESI Service (`backend/src/services/esi.go`)

**New Function: `makeESIRequestWithRetry()`**
- **Automatic retries:** Up to 3 attempts for 429, 420, 5xx errors
- **Respect Retry-After:** Waits exactly as ESI instructs
- **Intelligent backoff:** Applies before requests when approaching limits
- **Better error handling:** Proper categorization and logging

**Enhanced Request Flow:**
```
1. Check CanMakeRequest() - verifies safety
2. Apply ShouldBackoff() - slows if approaching limits
3. Make request with proper headers
4. Update all rate limit info from response
5. Handle 429/420/5xx with retry + wait
6. Return response or fail after max retries
```

**Request Headers:**
```go
User-Agent: EVE Ran Application - GitHub: tadeasf/eve-ran - Contact: github.com/tadeasf
Accept: application/json
Cache-Control: no-cache
```

### Kill Enhancement (`backend/src/jobs/killEnhance.go`)

**Before:**
```go
resp, err := http.Get(url)  // Direct call, no rate limiting
```

**After:**
```go
kill, err := services.FetchKillmailFromESI(zkill.KillmailID, zkill.Hash)
// Centralized rate-limited service with retries
```

**Benefits:**
- Automatic rate limiting
- Built-in retry logic
- Consistent error handling
- Comprehensive logging

## Best Practices Implemented

Per [ESI Best Practices](https://developers.eveonline.com/docs/services/esi/rate-limiting/#best-practices):

✅ **Don't operate at the limit** - Buffers: 5 errors, 10% tokens  
✅ **Slow down when approaching** - Progressive backoff system  
✅ **Spread requests over time** - No bursting, intelligent delays  
✅ **Respect Retry-After** - Wait exactly as instructed on 429s  
✅ **Use proper User-Agent** - Identifies our application  
✅ **Respect cache times** - Cache-Control headers included

## Monitoring & Logging

System logs when:
- Error limit drops below 20
- Any route group drops below 30%
- 429 errors occur (with retry timing)
- Rate limit resets

**Example output:**
```
ESI Rate Limit Status - Error Remaining: 85, Reset: 2025-11-22T03:45:00Z
  Group killmails: 142/150 tokens remaining (15m window)
ESI Rate Limited: Retry after 15 seconds (at 2025-11-22T03:35:28Z)
```

## Production Recommendations

1. **Monitor Logs:** Watch for rate limit warnings
2. **Adjust Concurrency:** Reduce concurrent requests if limits hit frequently
3. **Respect Cache Headers:** Use `expires` and `last-modified` headers
4. **Consider:**
   - Circuit breaker pattern for failures
   - Metrics collection for rate limit status
   - Alert system when approaching limits
   - Dynamic concurrency based on rate limit state

## References

- [ESI Documentation](https://docs.esi.evetech.net/docs/esi_introduction.html)
- [ESI Guidelines](https://docs.esi.evetech.net/docs/guidelines.html)
- [ESI Support](https://www.eveonline.com/): #esi on Tweetfleet Slack
