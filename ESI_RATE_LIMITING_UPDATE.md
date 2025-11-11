# EVE ESI Rate Limiting Update

## Overview

This document details the updates made to implement proper ESI (EVE Swagger Interface) rate limiting handling according to the latest EVE Online API documentation.

## Changes Made

### 1. Updated `backend/src/services/esiErrorManager.go`

**Key Changes:**
- Added `UpdateLimitsFromHeaders(headers http.Header)` method to parse rate limit headers from ESI responses
- Implemented global singleton pattern with `GetGlobalESIManager()`
- Added `lastUpdate` field to track when rate limits were last updated
- Added `GetStatus()` method for debugging and monitoring
- Improved thread-safety with proper mutex usage

**New Headers Tracked:**
- `X-ESI-Error-Limit-Remain`: Number of errors remaining before hitting the limit
- `X-ESI-Error-Limit-Reset`: Seconds until the error limit resets

**Default Values:**
- Conservative default of 100 errors remaining
- Automatically updates from actual ESI responses

### 2. Updated `backend/src/services/esi.go`

**Key Changes:**
- Added `makeESIRequest(url, method)` helper function for centralized rate limiting
- Updated all ESI GET request functions to use the new helper
- Added proper User-Agent headers with contact information
- Implemented rate limit checking before making requests
- Added automatic waiting when rate limit is reached
- Improved error handling with status code checking

**Functions Updated:**
1. `FetchRegionIDs()` - Universe regions
2. `FetchRegionInfo()` - Individual region details
3. `FetchSystemIDs()` - Universe systems
4. `FetchSystemInfo()` - Individual system details
5. `FetchConstellationIDs()` - Universe constellations
6. `FetchConstellationInfo()` - Individual constellation details
7. `FetchItemIDs()` - Universe item types (with pagination)
8. `FetchItemInfo()` - Individual item details
9. `FetchConstellation()` - Constellation lookup
10. `FetchCharacterInfo()` - Character information

**Functions with Header Parsing Added:**
1. `FetchKillmailFromESI()` - Already used custom request handling, added header parsing
2. `SearchCharactersByName()` - POST request, added header parsing

**Rate Limit Features:**
- Automatic rate limit checking before each request
- Automatic waiting if error limit is reached
- Rate limit status logging for debugging
- Proper handling of 420 and 520 status codes (rate limit errors)
- User-Agent header: `EVE Ran Application - GitHub: tadeasf/eve-ran - Contact: github.com/tadeasf`

### 3. Updated `CLAUDE.md`

Added documentation about the new rate limiting implementation for future reference.

## How It Works

### Request Flow

1. **Before Request:**
   - Check `esiManager.CanMakeRequest()`
   - If rate limit reached, wait until reset time
   - Create HTTP request with proper headers

2. **During Request:**
   - Execute request via `esiClient.Do(req)`
   - Return response for processing

3. **After Request:**
   - Parse `X-ESI-Error-Limit-Remain` header
   - Parse `X-ESI-Error-Limit-Reset` header
   - Update global rate limit manager
   - Log current rate limit status

4. **On Error:**
   - Decrement error count for 4xx/5xx responses
   - Handle 420/520 specifically as rate limit errors
   - Return appropriate error message

### Rate Limit Manager

The `ESIErrorManager` is a global singleton that:
- Tracks errors remaining before hitting the limit
- Tracks when the limit will reset
- Provides thread-safe access to rate limit data
- Automatically waits when limit is reached
- Logs rate limit status for monitoring

## Benefits

1. **Prevents Rate Limiting:** Proactively checks before making requests
2. **Automatic Recovery:** Waits for reset when limit is reached
3. **Better Monitoring:** Logs rate limit status for debugging
4. **Compliant with ESI Guidelines:** Uses proper headers and User-Agent
5. **Centralized Management:** Single source of truth for rate limits
6. **Thread-Safe:** Handles concurrent requests properly

## Testing

The implementation has been tested by:
1. Building the backend successfully (`go build -o main ./src`)
2. Verifying all functions compile without errors
3. Ensuring backward compatibility with existing code

## Recommendations for Deployment

1. **Monitor Logs:** Watch for rate limit warnings in production
2. **Adjust Concurrency:** If rate limits are hit frequently, reduce concurrent ESI requests in `FetchAll*` functions
3. **Cache Aggressively:** Respect ESI's caching headers (`expires`, `last-modified`)
4. **Implement Retry Logic:** Consider adding exponential backoff for transient errors

## Further Improvements

Consider implementing:
1. Caching layer that respects ESI's `expires` and `last-modified` headers
2. Circuit breaker pattern for ESI failures
3. Metrics collection for rate limit status
4. Alert system when approaching error limit
5. Dynamic concurrency adjustment based on rate limit status

## References

- [ESI Documentation](https://docs.esi.evetech.net/docs/esi_introduction.html)
- [ESI Rate Limiting](https://docs.esi.evetech.net/docs/guidelines.html)
- [EVE ESI API](https://esi.evetech.net/)

## Contact

For issues or questions about this implementation:
- GitHub: https://github.com/tadeasf/eve-ran-monorepo
- ESI Support: #esi on Tweetfleet Slack
