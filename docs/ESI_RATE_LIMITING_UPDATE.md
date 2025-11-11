# EVE ESI Rate Limiting Implementation

## Overview

Implementation of proper ESI (EVE Swagger Interface) rate limiting according to EVE Online API guidelines.

## Implementation Details

### ESI Error Manager (`backend/src/services/esiErrorManager.go`)

**Core Features:**
- Global singleton pattern with `GetGlobalESIManager()`
- Thread-safe rate limit tracking with mutex
- Tracks `X-ESI-Error-Limit-Remain` and `X-ESI-Error-Limit-Reset` headers
- Default: 100 errors remaining, updates from ESI responses
- `UpdateLimitsFromHeaders()`: Parses rate limit headers
- `CanMakeRequest()`: Checks if safe to make request
- `GetStatus()`: Returns current rate limit state for monitoring

### ESI Service (`backend/src/services/esi.go`)

**Centralized Request Handling:**
- `makeESIRequest(url, method)`: All ESI requests go through this helper
- Checks rate limit before each request
- Auto-waits when limit reached
- User-Agent: `EVE Ran Application - GitHub: tadeasf/eve-ran - Contact: github.com/tadeasf`
- Parses rate limit headers from responses
- Handles 420/520 status codes (rate limit errors)

**Functions Updated:**
- Universe: `FetchRegionIDs/Info`, `FetchSystemIDs/Info`, `FetchConstellationIDs/Info`
- Items: `FetchItemIDs/Info` (with pagination)
- Characters: `FetchCharacterInfo`, `SearchCharactersByName`
- Killmails: `FetchKillmailFromESI`

## Request Flow

1. **Pre-Request:** `CanMakeRequest()` checks rate limit, waits if needed
2. **Request:** Execute with proper User-Agent header
3. **Post-Request:** Parse `X-ESI-Error-Limit-*` headers, update manager
4. **Error Handling:** Decrement error count for 4xx/5xx, handle rate limit codes

## Benefits

- **Proactive Prevention:** Checks before requests, prevents hitting limits
- **Automatic Recovery:** Waits for reset when limit reached
- **ESI Compliant:** Proper headers and User-Agent
- **Centralized & Thread-Safe:** Single source of truth for rate limits
- **Observable:** Logs rate limit status for monitoring

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
