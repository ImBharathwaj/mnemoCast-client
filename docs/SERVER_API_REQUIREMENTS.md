# Server API Requirements

This document outlines what the server must return to the client for successful authentication and heartbeat transactions.

## Authentication Headers

All requests from the client include:
- `X-Screen-ID: {screen-id}` - The server-assigned screen ID
- `X-Passkey: {passkey}` - The server-assigned passkey
- `Content-Type: application/json`

---

## 1. Connection/Authentication Endpoint

### Request

**Endpoint:** `POST /api/v1/screens/{screenId}/connect`

**Headers:**
```
X-Screen-ID: d31f2fe7-16f3-4842-8db7-4b67868ecdc6
X-Passkey: LutpCadofzdKxT7aBzYhr5dsWmMDrzMKe97We1bAXak=
Content-Type: application/json
```

**Body:** Empty (or null)

### Server Response (Success)

**Status Code:** `200 OK`

**Response Body:** JSON Screen object
```json
{
  "id": "d31f2fe7-16f3-4842-8db7-4b67868ecdc6",
  "name": "Chennai Airport Screen 1",
  "country": "India",
  "city": "Chennai",
  "area": "Airport",
  "venueType": "airport",
  "timezone": "Asia/Kolkata",
  "width": 1920,
  "height": 1080,
  "isAudible": false,
  "isOnline": true,
  "lastSeen": "2025-12-19T00:07:52Z",
  "classification": 1,
  "createdAt": "2025-12-18T10:00:00Z",
  "updatedAt": "2025-12-19T00:07:52Z"
}
```

**Required Fields:**
- `id` (string) - Screen ID
- `name` (string) - Screen name
- `isOnline` (boolean) - Online status
- `createdAt` (string, ISO 8601) - Creation timestamp
- `updatedAt` (string, ISO 8601) - Last update timestamp

**Optional Fields:**
- `country`, `city`, `area`, `venueType`, `timezone`
- `width`, `height`, `isAudible`
- `lastSeen` (string, ISO 8601, nullable)
- `classification` (integer)

### Server Response (Error Cases)

**401 Unauthorized** - Invalid screen ID or passkey
```json
{
  "error": "Invalid credentials",
  "message": "The supplied screen ID or passkey is invalid"
}
```

**404 Not Found** - Screen ID does not exist
```json
{
  "error": "Screen not found",
  "message": "Screen with ID d31f2fe7-16f3-4842-8db7-4b67868ecdc6 not found"
}
```

**405 Method Not Allowed** - Wrong HTTP method
- Currently the client uses `POST`, but server only accepts `OPTIONS`
- **Fix:** Server should accept `POST` method for this endpoint

---

## 2. Heartbeat Endpoint

### Request

**Endpoint:** `PUT /api/v1/screens/{screenId}/heartbeat`

**Headers:**
```
X-Screen-ID: d31f2fe7-16f3-4842-8db7-4b67868ecdc6
X-Passkey: LutpCadofzdKxT7aBzYhr5dsWmMDrzMKe97We1bAXak=
Content-Type: application/json
```

**Request Body:**
```json
{
  "status": "online",
  "timestamp": "2025-12-19T00:07:52Z"
}
```

### Server Response (Success)

**Status Code:** `200 OK` or `204 No Content`

**Option 1: 200 OK with Response Body**
```json
{
  "message": "Heartbeat received",
  "timestamp": "2025-12-19T00:07:52Z"
}
```

**Option 2: 204 No Content** (Preferred - no body needed)
- Empty response body
- Client accepts both 200 and 204

### Server Response (Error Cases)

**401 Unauthorized** - Invalid credentials
```json
{
  "error": "Unauthorized",
  "message": "Invalid screen ID or passkey"
}
```

**403 Forbidden** - Credentials valid but not authorized
```json
{
  "error": "Forbidden",
  "message": "The supplied authentication is not authorized to access this resource"
}
```

**404 Not Found** - Screen ID does not exist
```json
{
  "error": "Screen not found",
  "message": "Screen with ID d31f2fe7-16f3-4842-8db7-4b67868ecdc6 not found"
}
```

---

## Current Issues

Based on the terminal output, the server is currently returning:

1. **Connection Endpoint:**
   - ❌ `405 Method Not Allowed` - Server only accepts `OPTIONS`
   - ✅ **Fix:** Server should accept `POST` method

2. **Heartbeat Endpoint:**
   - ❌ `403 Forbidden` - Authentication not authorized
   - ✅ **Fix:** Server should validate `X-Screen-ID` and `X-Passkey` headers correctly

---

## Server Implementation Checklist

### Connection Endpoint (`POST /api/v1/screens/{screenId}/connect`)

- [ ] Accept `POST` method (not just `OPTIONS`)
- [ ] Read `X-Screen-ID` header
- [ ] Read `X-Passkey` header
- [ ] Validate screen ID exists in database
- [ ] Validate passkey matches the screen
- [ ] Return `200 OK` with Screen JSON object on success
- [ ] Return `401 Unauthorized` if credentials invalid
- [ ] Return `404 Not Found` if screen ID doesn't exist
- [ ] Update `lastSeen` timestamp in database
- [ ] Set `isOnline` to `true` in database

### Heartbeat Endpoint (`PUT /api/v1/screens/{screenId}/heartbeat`)

- [ ] Accept `PUT` method
- [ ] Read `X-Screen-ID` header
- [ ] Read `X-Passkey` header
- [ ] Validate screen ID exists in database
- [ ] Validate passkey matches the screen
- [ ] Parse request body (status, timestamp)
- [ ] Update `lastSeen` timestamp in database
- [ ] Set `isOnline` to `true` in database
- [ ] Return `200 OK` or `204 No Content` on success
- [ ] Return `401 Unauthorized` if credentials invalid
- [ ] Return `403 Forbidden` if authorization fails
- [ ] Return `404 Not Found` if screen ID doesn't exist

---

## Authentication Flow

```
1. Client sends: POST /api/v1/screens/{id}/connect
   Headers: X-Screen-ID, X-Passkey
   
2. Server validates:
   - Screen ID exists
   - Passkey matches
   
3. Server responds:
   - 200 OK + Screen JSON (success)
   - 401 Unauthorized (invalid credentials)
   - 404 Not Found (screen not found)
   
4. Client receives response and updates local identity
```

## Heartbeat Flow

```
1. Client sends: PUT /api/v1/screens/{id}/heartbeat
   Headers: X-Screen-ID, X-Passkey
   Body: {"status": "online", "timestamp": "..."}
   
2. Server validates:
   - Screen ID exists
   - Passkey matches
   
3. Server updates:
   - lastSeen = current timestamp
   - isOnline = true
   
4. Server responds:
   - 200 OK or 204 No Content (success)
   - 401/403/404 (error)
   
5. Client continues sending heartbeats every 30 seconds
```

---

## Example Server Response (Scala/Pekko)

### Connection Endpoint

```scala
POST("/api/v1/screens/:screenId/connect") {
  val screenId = params("screenId")
  val passkey = request.header("X-Passkey").getOrElse("")
  
  // Validate credentials
  screenService.authenticate(screenId, passkey) match {
    case Some(screen) =>
      // Update last seen
      screenService.updateLastSeen(screenId)
      
      // Return screen object
      Ok(screen.toJson)
    case None =>
      Unauthorized(Json.obj("error" -> "Invalid credentials"))
  }
}
```

### Heartbeat Endpoint

```scala
PUT("/api/v1/screens/:screenId/heartbeat") {
  val screenId = params("screenId")
  val passkey = request.header("X-Passkey").getOrElse("")
  
  // Validate credentials
  screenService.authenticate(screenId, passkey) match {
    case Some(_) =>
      // Update heartbeat
      screenService.updateHeartbeat(screenId)
      
      // Return success (204 No Content)
      NoContent
    case None =>
      Unauthorized(Json.obj("error" -> "Invalid credentials"))
  }
}
```

---

## Summary

**To keep transactions alive, the server must:**

1. ✅ Accept `POST` for `/api/v1/screens/{id}/connect`
2. ✅ Accept `PUT` for `/api/v1/screens/{id}/heartbeat`
3. ✅ Validate `X-Screen-ID` and `X-Passkey` headers
4. ✅ Return `200 OK` or `204 No Content` for successful heartbeats
5. ✅ Update `lastSeen` and `isOnline` in database
6. ✅ Return proper error codes (401, 403, 404) for failures

The client will automatically retry failed requests and continue sending heartbeats every 30 seconds when the server responds correctly.

