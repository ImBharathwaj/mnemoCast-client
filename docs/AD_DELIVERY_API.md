# Ad Delivery API - Server Response Format

This document specifies the exact JSON response format that the server must return for the ad delivery endpoint.

## Endpoint

**GET** `/api/v1/screens/{screenId}/ads/deliver`

**Headers:**
```
X-Screen-Id: {screenId}
X-Screen-Passkey: {passkey}
Content-Type: application/json
```

---

## Response Format

### Success Response (200 OK)

When ads are available, the server should return:

```json
{
  "ads": [
    {
      "id": "ad-12345",
      "title": "Summer Sale 2025",
      "type": "image",
      "contentUrl": "https://example.com/ads/summer-sale.jpg",
      "duration": 30,
      "startTime": "2025-12-19T10:00:00Z",
      "endTime": "2025-12-19T18:00:00Z",
      "priority": 1,
      "metadata": {
        "campaignId": "campaign-001",
        "targetAudience": "all"
      }
    },
    {
      "id": "ad-67890",
      "title": "New Product Launch",
      "type": "video",
      "contentUrl": "https://example.com/ads/product-launch.mp4",
      "duration": 60,
      "startTime": "2025-12-19T12:00:00Z",
      "endTime": "2025-12-19T20:00:00Z",
      "priority": 2,
      "metadata": {
        "campaignId": "campaign-002",
        "targetAudience": "premium"
      }
    }
  ],
  "playlistId": "playlist-abc123",
  "updatedAt": "2025-12-19T01:24:10Z"
}
```

### No Ads Available (204 No Content)

When no ads are available, the server should return:

**Status Code:** `204 No Content`

**Body:** Empty (no response body)

**OR**

**Status Code:** `200 OK`

**Body:**
```json
{
  "ads": [],
  "updatedAt": "2025-12-19T01:24:10Z"
}
```

---

## Field Specifications

### Root Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `ads` | `array` | Yes | Array of ad objects (can be empty) |
| `playlistId` | `string` | No | Associated playlist identifier |
| `updatedAt` | `string` (ISO 8601) | Yes | Timestamp when ads were last updated |

### Ad Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `id` | `string` | Yes | Unique ad identifier |
| `title` | `string` | No | Ad title/name |
| `type` | `string` | Yes | Ad type: `image`, `video`, `html`, `text`, etc. |
| `contentUrl` | `string` | Yes | URL to ad content (image, video, etc.) |
| `duration` | `integer` | No | Display duration in seconds |
| `startTime` | `string` (ISO 8601) | No | Scheduled start time |
| `endTime` | `string` (ISO 8601) | No | Scheduled end time |
| `priority` | `integer` | No | Display priority (higher = more important) |
| `metadata` | `object` | No | Additional metadata (key-value pairs) |

---

## Example Responses

### Example 1: Single Image Ad

```json
{
  "ads": [
    {
      "id": "ad-001",
      "title": "Promotional Banner",
      "type": "image",
      "contentUrl": "https://cdn.example.com/ads/banner-001.jpg",
      "duration": 15,
      "priority": 1
    }
  ],
  "playlistId": "playlist-001",
  "updatedAt": "2025-12-19T01:24:10Z"
}
```

### Example 2: Multiple Ads with Scheduling

```json
{
  "ads": [
    {
      "id": "ad-morning",
      "title": "Morning Special",
      "type": "image",
      "contentUrl": "https://cdn.example.com/ads/morning.jpg",
      "duration": 20,
      "startTime": "2025-12-19T06:00:00Z",
      "endTime": "2025-12-19T12:00:00Z",
      "priority": 2,
      "metadata": {
        "category": "food",
        "location": "lobby"
      }
    },
    {
      "id": "ad-afternoon",
      "title": "Afternoon Deal",
      "type": "video",
      "contentUrl": "https://cdn.example.com/ads/afternoon.mp4",
      "duration": 45,
      "startTime": "2025-12-19T12:00:00Z",
      "endTime": "2025-12-19T18:00:00Z",
      "priority": 1,
      "metadata": {
        "category": "retail",
        "location": "all"
      }
    }
  ],
  "playlistId": "playlist-daily",
  "updatedAt": "2025-12-19T01:24:10Z"
}
```

### Example 3: No Ads Available

```json
{
  "ads": [],
  "updatedAt": "2025-12-19T01:24:10Z"
}
```

---

## Error Responses

### 401 Unauthorized

```json
{
  "error": "Unauthorized",
  "message": "Invalid screen ID or passkey"
}
```

### 403 Forbidden

```json
{
  "error": "Forbidden",
  "message": "The supplied authentication is not authorized to access this resource"
}
```

### 404 Not Found

```json
{
  "error": "Screen not found",
  "message": "Screen with ID d31f2fe7-16f3-4842-8db7-4b67868ecdc6 not found"
}
```

---

## Client Behavior

### On Success (200 OK)

1. Client receives JSON response
2. Parses `ads` array
3. Saves ads to `~/.mnemocast/ads/current_ads.json`
4. Logs ad count and details
5. Ready for display/processing

### On No Content (204 No Content)

1. Client receives empty response
2. Creates empty ads array
3. Saves to filesystem (with empty array)
4. Logs "No ads available"

### On Error (401/403/404)

1. Client logs error
2. Retries according to retry configuration
3. Does not update stored ads on failure

---

## Data Types

### Timestamps

All timestamps should be in **ISO 8601** format:
- Format: `YYYY-MM-DDTHH:MM:SSZ` or `YYYY-MM-DDTHH:MM:SS+00:00`
- Example: `2025-12-19T01:24:10Z`

### Duration

- Type: `integer`
- Unit: **seconds**
- Example: `30` = 30 seconds

### Priority

- Type: `integer`
- Higher numbers = higher priority
- Example: `1` = normal, `2` = high, `3` = urgent

### Content URL

- Type: `string`
- Format: Full URL (http:// or https://)
- Should be accessible by the client
- Example: `https://cdn.example.com/ads/image.jpg`

---

## Server Implementation Checklist

- [ ] Return `200 OK` with ads array (even if empty)
- [ ] Return `204 No Content` when no ads (optional, client handles both)
- [ ] Include `updatedAt` timestamp in ISO 8601 format
- [ ] Include `playlistId` if available
- [ ] Validate `X-Screen-Id` and `X-Screen-Passkey` headers
- [ ] Return proper error codes (401, 403, 404) for failures
- [ ] Ensure all required fields (`id`, `type`, `contentUrl`) are present
- [ ] Use ISO 8601 format for all timestamps
- [ ] Support empty ads array (no ads available)

---

## Client Storage Format

The client stores the response in `~/.mnemocast/ads/current_ads.json`:

```json
{
  "fetchedAt": "2025-12-19T01:24:10Z",
  "playlistId": "playlist-abc123",
  "updatedAt": "2025-12-19T01:24:10Z",
  "ads": [
    {
      "id": "ad-12345",
      "title": "Summer Sale 2025",
      "type": "image",
      "contentUrl": "https://example.com/ads/summer-sale.jpg",
      "duration": 30,
      "priority": 1
    }
  ],
  "adsCount": 1
}
```

Note: The client adds `fetchedAt` and `adsCount` fields for tracking.

---

## Summary

**Minimum Required Response:**
```json
{
  "ads": [],
  "updatedAt": "2025-12-19T01:24:10Z"
}
```

**Full Response with Ads:**
```json
{
  "ads": [
    {
      "id": "string (required)",
      "type": "string (required)",
      "contentUrl": "string (required)",
      "title": "string (optional)",
      "duration": "integer (optional)",
      "startTime": "ISO 8601 (optional)",
      "endTime": "ISO 8601 (optional)",
      "priority": "integer (optional)",
      "metadata": "object (optional)"
    }
  ],
  "playlistId": "string (optional)",
  "updatedAt": "ISO 8601 (required)"
}
```

