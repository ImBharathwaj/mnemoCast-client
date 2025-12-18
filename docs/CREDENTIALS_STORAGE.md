# Credentials Storage Location

## Where Credentials Are Stored

The **Screen ID** and **Passkey** are stored securely in an **encrypted file**, not in plain text JSON files.

### Storage Files

```
~/.mnemocast/
├── config.json              # Application configuration (NO credentials)
├── identity.json            # Screen identity (NO credentials)
├── credentials.json.enc     # ✅ ENCRYPTED credentials (Screen ID + Passkey)
└── .encryption_key          # Encryption key (0600 permissions)
```

### Why Not in config.json or identity.json?

1. **Security**: Credentials are sensitive and must be encrypted
2. **Separation**: Configuration and identity are separate from authentication
3. **Best Practice**: Credentials should never be stored in plain text

### How It Works

1. **When you configure credentials:**
   - Screen ID and Passkey are encrypted using AES-256-GCM
   - Encrypted data is saved to `credentials.json.enc`
   - Encryption key is stored in `.encryption_key` (with restricted permissions)

2. **When the app starts:**
   - Reads `credentials.json.enc`
   - Decrypts using `.encryption_key`
   - Loads Screen ID and Passkey into memory
   - **You don't need to configure again!**

### Verification

To verify credentials are stored:

```bash
# Check if credentials file exists
ls -la ~/.mnemocast/credentials.json.enc

# Check file size (should not be empty)
stat ~/.mnemocast/credentials.json.enc

# Verify encryption key exists
ls -la ~/.mnemocast/.encryption_key
```

### File Permissions

- `credentials.json.enc`: `0600` (read/write for owner only)
- `.encryption_key`: `0600` (read/write for owner only)

### Do You Need to Configure Every Time?

**NO!** Once you configure credentials:
- They are saved to `credentials.json.enc`
- They persist between application restarts
- You only need to configure again if:
  - You delete `credentials.json.enc`
  - You want to change the Screen ID or Passkey
  - The encryption key is lost/corrupted

### Current Status

From your terminal output, credentials ARE being loaded:
```
✅ Credentials: Configured
   Screen ID: d31f2fe7-16f3-4842-8db7-4b67868ecdc6
   Passkey: Lutp...Xak=
```

This means:
- ✅ Credentials file exists
- ✅ Credentials are being decrypted successfully
- ✅ Screen ID and Passkey are loaded into memory

### Troubleshooting

If credentials are not persisting:

1. **Check file permissions:**
   ```bash
   ls -la ~/.mnemocast/credentials.json.enc
   ls -la ~/.mnemocast/.encryption_key
   ```

2. **Verify files exist:**
   ```bash
   test -f ~/.mnemocast/credentials.json.enc && echo "Exists" || echo "Missing"
   test -f ~/.mnemocast/.encryption_key && echo "Exists" || echo "Missing"
   ```

3. **Check if app can read:**
   - Ensure you're running as the same user who created the files
   - Check file ownership: `ls -la ~/.mnemocast/`

### Security Note

- **Never** commit `credentials.json.enc` or `.encryption_key` to version control
- These files are in `.gitignore` by default
- If you need to backup, ensure secure storage
- If compromised, regenerate credentials on the server

