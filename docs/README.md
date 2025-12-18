# 🖥️ MnemoCast Screen System

Digital Display Screen System with identity, credentials, heartbeat, and ad server integration.

## 🎯 Features

- ✅ **Screen Identity** - Unique UUID-based screen identification
- ✅ **Configuration Management** - Secure configuration file handling
- ✅ **Credentials Management** - Secure API key storage with AES-256-GCM encryption
- ⏭️ **Ad Server Integration** - Registration and communication (Phase 3)
- ⏭️ **Heartbeat System** - Regular status updates (Phase 4)

## 🚀 Quick Start

### Build

```bash
go build -o bin/screen ./cmd/screen
```

### Run

```bash
./bin/screen
```

This will:
1. Create `~/.mnemocast/` directory if it doesn't exist
2. Generate or load screen identity
3. Load or create default configuration
4. Check for credentials (prompts to set API key if missing)
5. Display screen information and status

## 📁 Project Structure

```
mnemoCast-client/
├── cmd/
│   └── screen/
│       └── main.go              # Entry point
├── internal/
│   ├── identity/                 # Identity management
│   │   ├── generator.go         # UUID generation
│   │   └── manager.go           # Identity operations
│   ├── credentials/             # Credentials management
│   │   ├── manager.go          # Credential operations
│   │   └── storage.go          # Secure storage
│   ├── heartbeat/               # Heartbeat (Phase 4)
│   ├── client/                  # Ad server client (Phase 3)
│   ├── config/                  # Configuration
│   │   └── loader.go            # Config loading/saving
│   └── models/                  # Data models
│       ├── identity.go
│       ├── credentials.go
│       └── config.go
└── pkg/
    └── storage/                  # Storage utilities
        ├── encryption.go        # AES-256-GCM encryption
        └── keygen.go            # Key generation
```

## 📝 Configuration

Configuration is stored in `~/.mnemocast/`:

- `identity.json` - Screen identity (ID, name, location)
- `config.json` - Application configuration
- `credentials.json.enc` - Encrypted credentials (AES-256-GCM)
- `.encryption_key` - Encryption key (0600 permissions)

## 🔧 Development

### Current Status

- ✅ **Phase 1 Complete:** Identity & Configuration
- ✅ **Phase 2 Complete:** Credentials Management
- ⏭️ **Phase 3 Next:** Ad Server Client
- ⏭️ **Phase 4:** Heartbeat System

### Running Tests

```bash
go test ./...
```

## 📚 Documentation

- **Implementation Plan:** `SCREEN_SYSTEM_PLAN.md`
- **Architecture:** See plan document

## 🔐 Security

- Configuration files use restricted permissions (0600)
- Credentials encrypted with AES-256-GCM
- Encryption key stored separately with restricted permissions
- No plain-text credential storage
- No sensitive data in logs

## 🛠️ Requirements

- Go 1.22+
- Backend running at `http://10.42.0.1:8080` (for Phase 3+)

---

**Status:** Phase 1 & 2 Complete | Phase 3 Next
