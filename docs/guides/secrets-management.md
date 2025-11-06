# Secrets and Environment Variables Management

This guide explains how to properly configure and manage secrets in your homelab.

## Overview

Every service uses environment variables stored in `.env` files for configuration. This approach:
- ✅ Keeps secrets out of version control
- ✅ Makes configuration changes easy
- ✅ Follows Docker Compose best practices
- ✅ Provides clear examples of what each secret should look like

## Finding What Needs Configuration

### Search for Placeholder Values

All secrets that need changing are marked with `changeme_`:

```bash
# Find all files with placeholder secrets
grep -r "changeme_" ~/homelab/compose

# Output shows exactly what needs updating:
compose/core/lldap/.env:LLDAP_LDAP_USER_PASS=changeme_please_set_secure_password
compose/core/lldap/.env:LLDAP_JWT_SECRET=changeme_please_set_random_secret
compose/core/tinyauth/.env:LDAP_BIND_PASSWORD=changeme_please_set_secure_password
...
```

### Count What's Left to Configure

```bash
# Count how many secrets still need updating
grep -r "changeme_" ~/homelab/compose | wc -l

# Goal: 0
```

## Generating Secrets

Each `.env` file includes comments showing:
1. What the secret is for
2. How to generate it
3. What format it should be in

### Common Secret Types

#### 1. JWT Secrets (64 characters)

**Used by**: LLDAP, Vikunja, NextAuth

**Generate:**
```bash
openssl rand -hex 32
```

**Example output:**
```
a1b2c3d4e5f67890abcdef1234567890a1b2c3d4e5f67890abcdef1234567890
```

**Where to use:**
- `LLDAP_JWT_SECRET`
- `VIKUNJA_SERVICE_JWTSECRET`
- `NEXTAUTH_SECRET`
- `SESSION_SECRET`

#### 2. Database Passwords (32 alphanumeric)

**Used by**: Postgres, Immich, Vikunja, Linkwarden

**Generate:**
```bash
openssl rand -base64 32 | tr -d /=+ | cut -c1-32
```

**Example output:**
```
aB3dEf7HiJ9kLmN2oPqR5sTuV8wXyZ1
```

**Where to use:**
- `DB_PASSWORD` (Immich)
- `POSTGRES_PASSWORD` (Vikunja, Linkwarden)
- `VIKUNJA_DATABASE_PASSWORD`

#### 3. Strong Passwords (16+ characters, mixed)

**Used by**: LLDAP admin, service admin accounts

**Generate:**
```bash
# Option 1: Using pwgen (install: apt install pwgen)
pwgen -s 20 1

# Option 2: Using openssl
openssl rand -base64 20 | tr -d /=+

# Option 3: Manual (recommended for main admin password)
# Create something memorable but strong
# Example format: MyS3cur3P@ssw0rd!2024#HomeL@b
```

**Where to use:**
- `LLDAP_LDAP_USER_PASS`
- `LDAP_BIND_PASSWORD` (must match LLDAP_LDAP_USER_PASS!)

#### 4. API Keys / Master Keys (32 characters)

**Used by**: Meilisearch, various APIs

**Generate:**
```bash
openssl rand -hex 16
```

**Example output:**
```
f6g7h8i901234abcdef567890a1b2c3d
```

**Where to use:**
- `MEILI_MASTER_KEY`

## Service-Specific Configuration

### Core Services

#### LLDAP (`compose/core/lldap/.env`)

```bash
# Edit the file
cd ~/homelab/compose/core/lldap
nano .env
```

**Required secrets:**

```env
# Admin password - use a STRONG password you'll remember
# Example: MyS3cur3P@ssw0rd!2024#HomeL@b
LLDAP_LDAP_USER_PASS=changeme_please_set_secure_password

# JWT secret - generate with: openssl rand -hex 32
# Example: a1b2c3d4e5f67890abcdef1234567890a1b2c3d4e5f67890abcdef1234567890
LLDAP_JWT_SECRET=changeme_please_set_random_secret
```

**Generate and update:**
```bash
# Generate JWT secret
echo "LLDAP_JWT_SECRET=$(openssl rand -hex 32)"

# Choose a strong password for LLDAP_LDAP_USER_PASS
# Write it down - you'll need it for Tinyauth too!
```

#### Tinyauth (`compose/core/tinyauth/.env`)

```bash
cd ~/homelab/compose/core/tinyauth
nano .env
```

**Required secrets:**

```env
# MUST match LLDAP_LDAP_USER_PASS from lldap/.env
LDAP_BIND_PASSWORD=changeme_please_set_secure_password

# Session secret - generate with: openssl rand -hex 32
SESSION_SECRET=changeme_please_set_random_session_secret
```

**⚠️ CRITICAL**: `LDAP_BIND_PASSWORD` must exactly match `LLDAP_LDAP_USER_PASS`!

```bash
# Generate session secret
echo "SESSION_SECRET=$(openssl rand -hex 32)"
```

### Media Services

#### Immich (`compose/media/frontend/immich/.env`)

```bash
cd ~/homelab/compose/media/frontend/immich
nano .env
```

**Required secrets:**

```env
# Database password - generate with: openssl rand -base64 32 | tr -d /=+ | cut -c1-32
DB_PASSWORD=changeme_please_set_secure_password
```

```bash
# Generate
echo "DB_PASSWORD=$(openssl rand -base64 32 | tr -d /=+ | cut -c1-32)"
```

### Utility Services

#### Linkwarden (`compose/services/linkwarden/.env`)

```bash
cd ~/homelab/compose/services/linkwarden
nano .env
```

**Required secrets:**

```env
# NextAuth secret - generate with: openssl rand -hex 32
NEXTAUTH_SECRET=changeme_please_set_random_secret_key

# Postgres password - generate with: openssl rand -base64 32 | tr -d /=+ | cut -c1-32
POSTGRES_PASSWORD=changeme_please_set_secure_postgres_password

# Meilisearch master key - generate with: openssl rand -hex 16
MEILI_MASTER_KEY=changeme_please_set_meili_master_key
```

```bash
# Generate all three
echo "NEXTAUTH_SECRET=$(openssl rand -hex 32)"
echo "POSTGRES_PASSWORD=$(openssl rand -base64 32 | tr -d /=+ | cut -c1-32)"
echo "MEILI_MASTER_KEY=$(openssl rand -hex 16)"
```

#### Vikunja (`compose/services/vikunja/.env`)

```bash
cd ~/homelab/compose/services/vikunja
nano .env
```

**Required secrets:**

```env
# Database password (used in two places - must match!)
VIKUNJA_DATABASE_PASSWORD=changeme_please_set_secure_password
POSTGRES_PASSWORD=changeme_please_set_secure_password  # Same value!

# JWT secret - generate with: openssl rand -hex 32
VIKUNJA_SERVICE_JWTSECRET=changeme_please_set_random_jwt_secret
```

**⚠️ CRITICAL**: Both password fields must match!

```bash
# Generate
DB_PASS=$(openssl rand -base64 32 | tr -d /=+ | cut -c1-32)
echo "VIKUNJA_DATABASE_PASSWORD=$DB_PASS"
echo "POSTGRES_PASSWORD=$DB_PASS"
echo "VIKUNJA_SERVICE_JWTSECRET=$(openssl rand -hex 32)"
```

## Automated Configuration Script

Create a script to generate all secrets at once:

```bash
#!/bin/bash
# save as: ~/homelab/generate-secrets.sh

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${YELLOW}Homelab Secrets Generator${NC}\n"

echo "This script will help you generate secure secrets for your homelab."
echo "You'll need to manually copy these values into the respective .env files."
echo ""

# LLDAP
echo -e "${GREEN}=== LLDAP (compose/core/lldap/.env) ===${NC}"
echo "LLDAP_JWT_SECRET=$(openssl rand -hex 32)"
echo "LLDAP_LDAP_USER_PASS=<choose-a-strong-password-manually>"
echo ""

# Tinyauth
echo -e "${GREEN}=== Tinyauth (compose/core/tinyauth/.env) ===${NC}"
echo "LDAP_BIND_PASSWORD=<same-as-LLDAP_LDAP_USER_PASS-above>"
echo "SESSION_SECRET=$(openssl rand -hex 32)"
echo ""

# Immich
echo -e "${GREEN}=== Immich (compose/media/frontend/immich/.env) ===${NC}"
echo "DB_PASSWORD=$(openssl rand -base64 32 | tr -d /=+ | cut -c1-32)"
echo ""

# Linkwarden
echo -e "${GREEN}=== Linkwarden (compose/services/linkwarden/.env) ===${NC}"
echo "NEXTAUTH_SECRET=$(openssl rand -hex 32)"
echo "POSTGRES_PASSWORD=$(openssl rand -base64 32 | tr -d /=+ | cut -c1-32)"
echo "MEILI_MASTER_KEY=$(openssl rand -hex 16)"
echo ""

# Vikunja
VIKUNJA_PASS=$(openssl rand -base64 32 | tr -d /=+ | cut -c1-32)
echo -e "${GREEN}=== Vikunja (compose/services/vikunja/.env) ===${NC}"
echo "VIKUNJA_DATABASE_PASSWORD=$VIKUNJA_PASS"
echo "POSTGRES_PASSWORD=$VIKUNJA_PASS  # Must match above!"
echo "VIKUNJA_SERVICE_JWTSECRET=$(openssl rand -hex 32)"
echo ""

echo -e "${YELLOW}Done! Copy these values into your .env files.${NC}"
echo ""
echo "Don't forget to:"
echo "1. Choose a strong LLDAP_LDAP_USER_PASS manually"
echo "2. Use the same password for LDAP_BIND_PASSWORD in tinyauth"
echo "3. Save all secrets in a password manager"
```

**Usage:**
```bash
chmod +x ~/homelab/generate-secrets.sh
~/homelab/generate-secrets.sh > secrets.txt

# Review and copy secrets
cat secrets.txt

# Keep this file safe or delete after copying to .env files
```

## Security Best Practices

### 1. Use a Password Manager

Store all secrets in a password manager:
- **1Password**: Great for teams
- **Bitwarden**: Self-hostable option
- **KeePassXC**: Offline, open-source

Create an entry for each service with:
- Service name
- URL
- All secrets from `.env` file
- Admin credentials

### 2. Never Commit Secrets

The repository `.gitignore` already excludes `.env` files, but double-check:

```bash
# Verify .env files are ignored
git status

# Should NOT show any .env files
```

### 3. Backup Your Secrets

```bash
# Create encrypted backup of all .env files
cd ~/homelab
tar czf env-backup-$(date +%Y%m%d).tar.gz $(find compose -name ".env")

# Encrypt with GPG
gpg -c env-backup-$(date +%Y%m%d).tar.gz

# Store encrypted file safely
mv env-backup-*.tar.gz.gpg ~/backups/

# Delete unencrypted tar
rm env-backup-*.tar.gz
```

### 4. Rotate Secrets Regularly

Change critical secrets periodically:
- **Admin passwords**: Every 90 days
- **JWT secrets**: Every 180 days
- **Database passwords**: When personnel changes

### 5. Limit Secret Access

- Don't share raw secrets over email/chat
- Use password manager's sharing features
- Delete shared secrets when no longer needed

## Verification

### Check All Secrets Are Set

```bash
# Should return 0 (no changeme_ values left)
grep -r "changeme_" ~/homelab/compose | wc -l
```

### Test Service Startup

```bash
# Start a service and check for password errors
cd ~/homelab/compose/core/lldap
docker compose up -d
docker compose logs

# Should NOT see:
# - "invalid password"
# - "authentication failed"
# - "secret not set"
```

### Verify SSO Works

1. Start LLDAP and Tinyauth
2. Access protected service (e.g., https://tasks.fig.systems)
3. Should redirect to auth.fig.systems
4. Login with LLDAP credentials
5. Should redirect back to service

If this works, your LLDAP ↔ Tinyauth passwords match! ✅

## Common Mistakes

### ❌ Using Weak Passwords

**Don't:**
```env
LLDAP_LDAP_USER_PASS=password123
```

**Do:**
```env
LLDAP_LDAP_USER_PASS=MyS3cur3P@ssw0rd!2024#HomeL@b
```

### ❌ Mismatched Passwords

**Don't:**
```env
# In lldap/.env
LLDAP_LDAP_USER_PASS=password1

# In tinyauth/.env
LDAP_BIND_PASSWORD=password2  # Different!
```

**Do:**
```env
# In lldap/.env
LLDAP_LDAP_USER_PASS=MyS3cur3P@ssw0rd!2024#HomeL@b

# In tinyauth/.env
LDAP_BIND_PASSWORD=MyS3cur3P@ssw0rd!2024#HomeL@b  # Same!
```

### ❌ Using Same Secret Everywhere

**Don't:**
```env
# Same secret in multiple places
LLDAP_JWT_SECRET=abc123
NEXTAUTH_SECRET=abc123
SESSION_SECRET=abc123
```

**Do:**
```env
# Unique secret for each
LLDAP_JWT_SECRET=a1b2c3d4e5f67890...
NEXTAUTH_SECRET=f6g7h8i9j0k1l2m3...
SESSION_SECRET=x9y8z7w6v5u4t3s2...
```

### ❌ Forgetting to Update Both Password Fields

In Vikunja `.env`, both must match:
```env
# Both must be the same!
VIKUNJA_DATABASE_PASSWORD=aB3dEf7HiJ9kLmN2oPqR5sTuV8wXyZ1
POSTGRES_PASSWORD=aB3dEf7HiJ9kLmN2oPqR5sTuV8wXyZ1
```

## Troubleshooting

### "Authentication failed" in Tinyauth

**Cause**: LDAP_BIND_PASSWORD doesn't match LLDAP_LDAP_USER_PASS

**Fix**:
```bash
# Check LLDAP password
grep LLDAP_LDAP_USER_PASS ~/homelab/compose/core/lldap/.env

# Check Tinyauth password
grep LDAP_BIND_PASSWORD ~/homelab/compose/core/tinyauth/.env

# They should be identical!
```

### "Invalid JWT" errors

**Cause**: JWT_SECRET is too short or invalid format

**Fix**:
```bash
# Regenerate with proper length
openssl rand -hex 32

# Update in .env file
```

### "Database connection failed"

**Cause**: Database password mismatch

**Fix**:
```bash
# Check both password fields match
grep -E "(POSTGRES_PASSWORD|DATABASE_PASSWORD)" compose/services/vikunja/.env

# Both should be identical
```

## Next Steps

Once all secrets are configured:
1. ✅ [Deploy services](../getting-started.md#step-6-deploy-services)
2. ✅ [Configure SSO](../services/sso-setup.md)
3. ✅ [Set up backups](../operations/backups.md)
4. ✅ Store secrets in password manager
5. ✅ Create encrypted backup of .env files

## Reference

### Quick Command Reference

```bash
# Generate 64-char hex
openssl rand -hex 32

# Generate 32-char password
openssl rand -base64 32 | tr -d /=+ | cut -c1-32

# Generate 32-char hex
openssl rand -hex 16

# Find all changeme_ values
grep -r "changeme_" compose/

# Count remaining secrets to configure
grep -r "changeme_" compose/ | wc -l

# Backup all .env files (encrypted)
tar czf env-files.tar.gz $(find compose -name ".env")
gpg -c env-files.tar.gz
```

### Secret Types Quick Reference

| Secret Type | Command | Example Length | Used By |
|-------------|---------|----------------|---------|
| JWT Secret | `openssl rand -hex 32` | 64 chars | LLDAP, Vikunja, NextAuth |
| Session Secret | `openssl rand -hex 32` | 64 chars | Tinyauth |
| DB Password | `openssl rand -base64 32 \| tr -d /=+ \| cut -c1-32` | 32 chars | Postgres, Immich |
| API Key | `openssl rand -hex 16` | 32 chars | Meilisearch |
| Admin Password | Manual | 16+ chars | LLDAP admin |

---

**Remember**: Strong, unique secrets are your first line of defense. Take the time to generate them properly! 🔐
