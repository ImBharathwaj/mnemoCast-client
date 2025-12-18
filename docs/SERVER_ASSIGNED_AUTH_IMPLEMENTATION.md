# Server-Assigned Screen ID & Passkey Authentication

## Overview

This document outlines the implementation plan for changing from **client-generated screen ID + API key** to **server-assigned screen ID + passkey** authentication model.

## Current Flow (To Be Changed)

1. Screen generates its own UUID-based screen ID locally
2. User manually configures an API key
3. Screen registers with server using generated ID
4. All requests use `X-API-Key` header

## New Flow (To Implement)

1. **Server-side**: Admin registers a screen in the server, receives:
   - Unique Screen ID (assigned by server)
   - Passkey (assigned by server)
2. **Client-side**: User manually configures:
   - Screen ID (from server)
   - Passkey (from server)
3. **Client-side**: Screen connects/authenticates with server using Screen ID + Passkey
4. **All transactions**: Screen ID + Passkey used in every request (including heartbeat)

---

## Implementation Summary

### 1. Update Credentials Model

**File**: `internal/models/credentials.go`

**Changes**:
- Replace `APIKey` with `Passkey`
- Ensure `ScreenID` is stored (server-assigned, not generated)
- Add validation methods for passkey

```go
type Credentials struct {
    ScreenID string `json:"screenId"`        // Server-assigned screen ID
    Passkey  string `json:"passkey"`          // Server-assigned passkey
    // Remove: APIKey, AccessToken, RefreshToken, ExpiresAt
}
```

### 2. Update Identity Management

**File**: `internal/identity/manager.go`

**Changes**:
- Remove automatic UUID generation
- Screen ID should be manually configured (from server)
- Screen ID should be loaded from credentials, not generated
- Identity should be created/updated after successful server connection

**New Flow**:
- User configures Screen ID + Passkey
- Screen connects to server to validate credentials
- Server returns screen details (name, location, etc.)
- Screen saves identity with server-assigned ID

### 3. Update Ad Server Client

**File**: `internal/client/client.go`

**Changes**:
- Replace `apiKey` with `screenID` and `passkey`
- Update authentication headers:
  - Option A: `X-Screen-ID` + `X-Passkey` headers
  - Option B: `Authorization: Bearer <passkey>` + `X-Screen-ID` header
  - Option C: Basic Auth with screen ID as username and passkey as password
- Remove registration endpoint (screen is pre-registered on server)
- Add connection/authentication endpoint (e.g., `/api/v1/screens/{screenId}/connect` or `/api/v1/screens/authenticate`)

**New Methods**:
```go
// Connect authenticates with server using screen ID and passkey
func (c *Client) Connect(screenID, passkey string) (*models.Screen, error)

// All existing methods (Heartbeat, etc.) use screenID + passkey from client
```

### 4. Update Credentials Manager

**File**: `internal/credentials/manager.go`

**Changes**:
- Replace `SetAPIKey()` with `SetCredentials(screenID, passkey string)`
- Replace `GetAPIKey()` with `GetCredentials() (screenID, passkey string, error)`
- Update validation to check for both screen ID and passkey

### 5. Update Main Application

**File**: `cmd/screen/main.go`

**Changes**:
- Remove automatic screen ID generation
- Update interactive setup to ask for:
  - Screen ID (from server)
  - Passkey (from server)
- After configuration, attempt to connect to server
- On successful connection, load screen identity from server response
- Update all references from "API Key" to "Screen ID & Passkey"

**New Flow**:
```
1. Check if credentials exist
2. If not, prompt for Screen ID and Passkey
3. Save credentials (encrypted)
4. Connect to server with Screen ID + Passkey
5. On success, save server-returned screen details as identity
6. Start heartbeat system
```

### 6. Update Heartbeat System

**File**: `internal/heartbeat/scheduler.go`

**Changes**:
- No changes needed (uses client which will have screenID + passkey)
- Heartbeat will automatically use the configured credentials

### 7. Update Configuration

**File**: `internal/models/config.go`

**Changes**:
- Identity should be loaded from server after connection
- Default config should not generate screen ID

---

## API Endpoint Changes

### Current Endpoints (To Modify/Remove)

1. **POST `/api/v1/screens/register`** 
   - ❌ Remove or repurpose (screen is pre-registered)
   - ✅ Replace with connection/authentication endpoint

### New Endpoints (To Use)

1. **POST `/api/v1/screens/{screenId}/connect`** or **POST `/api/v1/screens/authenticate`**
   - **Purpose**: Authenticate screen with server using screen ID + passkey
   - **Request Headers**:
     - `X-Screen-ID: <screenId>`
     - `X-Passkey: <passkey>`
     - OR `Authorization: Bearer <passkey>` + `X-Screen-ID: <screenId>`
   - **Response**: Full screen details (name, location, etc.)
   - **Status Codes**:
     - `200 OK`: Authentication successful, returns screen data
     - `401 Unauthorized`: Invalid screen ID or passkey
     - `404 Not Found`: Screen ID not found

2. **POST `/api/v1/screens/{screenId}/heartbeat`**
   - **Purpose**: Send heartbeat (unchanged functionality)
   - **Request Headers**: Same authentication as connect endpoint
   - **Response**: Heartbeat acknowledgment

---

## Security Considerations

1. **Passkey Storage**: 
   - ✅ Already encrypted using AES-256-GCM (no changes needed)
   - ✅ Stored in `~/.mnemocast/credentials.json.enc`

2. **Screen ID Storage**:
   - ✅ Stored in encrypted credentials file
   - ✅ Also stored in identity.json (after server connection)

3. **Network Security**:
   - ✅ All requests use HTTPS (if configured)
   - ✅ Passkey never sent in URL or logs

---

## Migration Path

### For Existing Installations

1. **Option A**: Clear existing identity and credentials, reconfigure with server-assigned values
2. **Option B**: Provide migration script to convert API key to passkey format (if server supports both)

### Recommended Approach

- Clear existing configuration and start fresh with server-assigned credentials
- This ensures clean state and proper server-side registration

---

## Implementation Checklist

### Phase 1: Data Model Updates
- [ ] Update `Credentials` struct (replace APIKey with Passkey)
- [ ] Update credential validation methods
- [ ] Update credential storage/loading

### Phase 2: Client Updates
- [ ] Update `Client` struct (screenID + passkey instead of apiKey)
- [ ] Update authentication headers
- [ ] Replace registration with connection/authentication endpoint
- [ ] Update all request methods to use new auth

### Phase 3: Identity Management
- [ ] Remove automatic UUID generation
- [ ] Load screen ID from credentials
- [ ] Update identity after server connection
- [ ] Save server-returned screen details

### Phase 4: User Interface
- [ ] Update main application prompts
- [ ] Change "API Key" to "Screen ID & Passkey"
- [ ] Add connection/authentication flow
- [ ] Update status messages

### Phase 5: Testing
- [ ] Test credential configuration
- [ ] Test server connection
- [ ] Test heartbeat with new auth
- [ ] Test error handling (invalid credentials)

---

## Example User Flow

```
$ ./bin/screen

🖥️  MnemoCast Screen System
============================
Config directory: /home/user/.mnemocast

📋 Loading screen identity...
⚠️  Screen identity not found

🔐 Checking credentials...
⚠️  Credentials: Not configured

Would you like to configure credentials now? (y/n): y

Enter Screen ID (from server): screen-abc123xyz
Enter Passkey (from server): ****************

✅ Credentials saved successfully!

🌐 Connecting to ad server...
   Authenticating with server...
   ✅ Connected successfully!
   Screen Name: Chennai Airport Screen 1
   Location: Chennai, Airport (airport)

❤️  Starting heartbeat scheduler...
   ✅ Heartbeat system active

[System running...]
```

---

## Benefits of This Approach

1. **Centralized Control**: Server controls all screen IDs
2. **Security**: Passkey is server-generated and can be rotated
3. **Audit Trail**: Server knows exactly which screens are registered
4. **Simplified Client**: No need for client-side ID generation
5. **Better Management**: Server can track and manage all screens centrally

---

## Questions to Clarify with Backend Team

1. **Authentication Method**: Which header format should be used?
   - `X-Screen-ID` + `X-Passkey`?
   - `Authorization: Bearer <passkey>` + `X-Screen-ID`?
   - Basic Auth?

2. **Connection Endpoint**: What is the exact endpoint for authentication?
   - `/api/v1/screens/{screenId}/connect`?
   - `/api/v1/screens/authenticate`?
   - Something else?

3. **Response Format**: What does the connection endpoint return?
   - Full screen object with all details?
   - Just confirmation?

4. **Error Handling**: What status codes for invalid credentials?
   - `401 Unauthorized`?
   - `403 Forbidden`?
   - `404 Not Found`?

---

## Next Steps

1. ✅ Review this implementation plan
2. ⏳ Confirm authentication method with backend team
3. ⏳ Implement Phase 1 (Data Model Updates)
4. ⏳ Implement Phase 2 (Client Updates)
5. ⏳ Implement Phase 3 (Identity Management)
6. ⏳ Implement Phase 4 (User Interface)
7. ⏳ Test end-to-end flow
8. ⏳ Update documentation

