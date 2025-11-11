# Terraform State Management with SOPS

This project uses [SOPS](https://github.com/getsops/sops) (Secrets OPerationS) with [age](https://github.com/FiloSottile/age) encryption to securely store Terraform state files in Git.

## Why SOPS + age?

✅ **Encrypted at rest** - State files contain sensitive data (IPs, tokens)
✅ **Version controlled** - Track infrastructure changes over time
✅ **No infrastructure required** - No need for S3, PostgreSQL, or other backends
✅ **Perfect for homelabs** - Simple, secure, self-contained
✅ **FOSS** - Fully open source tools

## Prerequisites

### 1. Install age

**Debian/Ubuntu:**
```bash
sudo apt update
sudo apt install age
```

**macOS:**
```bash
brew install age
```

**Manual installation:**
```bash
# Download from https://github.com/FiloSottile/age/releases
wget https://github.com/FiloSottile/age/releases/download/v1.1.1/age-v1.1.1-linux-amd64.tar.gz
tar xzf age-v1.1.1-linux-amd64.tar.gz
sudo mv age/age age/age-keygen /usr/local/bin/
```

### 2. Install SOPS

**Debian/Ubuntu:**
```bash
# Download from https://github.com/getsops/sops/releases
wget https://github.com/getsops/sops/releases/download/v3.8.1/sops-v3.8.1.linux.amd64
sudo mv sops-v3.8.1.linux.amd64 /usr/local/bin/sops
sudo chmod +x /usr/local/bin/sops
```

**macOS:**
```bash
brew install sops
```

Verify installation:
```bash
age --version
sops --version
```

## Initial Setup

### 1. Generate Age Encryption Key

```bash
# Create SOPS directory
mkdir -p ~/.sops

# Generate a new age key pair
age-keygen -o ~/.sops/homelab-terraform.txt

# View the key (you'll need the public key)
cat ~/.sops/homelab-terraform.txt
```

Output will look like:
```
# created: 2025-11-11T12:34:56Z
# public key: age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
AGE-SECRET-KEY-1XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXx
```

**⚠️ IMPORTANT:**
- The line starting with `AGE-SECRET-KEY-1` is your **private key** - keep it secret!
- The line starting with `age1` is your **public key** - you'll use this in .sops.yaml
- **Backup this file** to a secure location (password manager, encrypted backup, etc.)
- If you lose this key, you **cannot decrypt** your state files!

### 2. Configure SOPS

```bash
cd terraform/proxmox-examples/docker-host

# Copy the example config
cp .sops.yaml.example .sops.yaml

# Edit and replace YOUR_AGE_PUBLIC_KEY_HERE with your public key from step 1
nano .sops.yaml
```

Your `.sops.yaml` should look like:
```yaml
creation_rules:
  - path_regex: \.tfstate$
    age: age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  - path_regex: \.secret$
    age: age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
  - path_regex: terraform\.tfvars$
    age: age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

### 3. Set Environment Variable (Optional but Recommended)

```bash
# Add to your ~/.bashrc or ~/.zshrc
echo 'export SOPS_AGE_KEY_FILE=~/.sops/homelab-terraform.txt' >> ~/.bashrc
source ~/.bashrc
```

This tells SOPS where to find your private key for decryption.

## Usage

### Option A: Automatic Wrapper Script (Recommended)

Use the `./scripts/tf` wrapper that handles encryption/decryption automatically:

```bash
# Initialize (first time)
./scripts/tf init

# Plan changes
./scripts/tf plan

# Apply changes (automatically encrypts after)
./scripts/tf apply

# Destroy infrastructure (automatically encrypts after)
./scripts/tf destroy

# View state
./scripts/tf show
```

The wrapper script:
1. Decrypts state files before running
2. Runs your terraform/tofu command
3. Encrypts state files after (if state was modified)

### Option B: Manual Encryption/Decryption

If you prefer manual control:

```bash
# 1. Decrypt state files
./scripts/tf-decrypt

# 2. Run terraform commands
tofu init
tofu plan
tofu apply

# 3. Encrypt state files
./scripts/tf-encrypt

# 4. Commit encrypted files to Git
git add *.enc
git commit -m "Update infrastructure"
git push
```

## Workflow Examples

### First Time Setup

```bash
cd terraform/proxmox-examples/docker-host

# 1. Configure your variables
cp terraform.tfvars.example terraform.tfvars
nano terraform.tfvars  # Add your API tokens, SSH keys, etc.

# 2. Initialize Terraform
./scripts/tf init

# 3. Plan infrastructure
./scripts/tf plan

# 4. Apply infrastructure
./scripts/tf apply

# 5. Encrypted state files are automatically created
# terraform.tfstate.enc now exists

# 6. Commit encrypted state to Git
git add terraform.tfstate.enc .sops.yaml.example
git commit -m "Add encrypted Terraform state"
git push
```

### Making Infrastructure Changes

```bash
# 1. Decrypt, apply changes, re-encrypt (all automatic)
./scripts/tf apply

# 2. Commit updated encrypted state
git add terraform.tfstate.enc
git commit -m "Update VM configuration"
git push
```

### Cloning on a New Machine

```bash
# 1. Clone the repository
git clone https://github.com/efigueroa/homelab.git
cd homelab/terraform/proxmox-examples/docker-host

# 2. Copy your age private key to the new machine
# (Securely transfer ~/.sops/homelab-terraform.txt)
mkdir -p ~/.sops
# Copy the key file here

# 3. Set up SOPS config
cp .sops.yaml.example .sops.yaml
# Edit with your public key

# 4. Decrypt state
./scripts/tf-decrypt

# 5. Now you can run terraform commands
./scripts/tf plan
```

## Security Best Practices

### DO ✅

- **Backup your age private key** to multiple secure locations
- **Use different keys** for different projects/environments
- **Commit `.sops.yaml.example`** to Git (without your actual key)
- **Commit encrypted `*.enc` files** to Git
- **Use the wrapper script** to avoid forgetting to encrypt

### DON'T ❌

- **Never commit `.sops.yaml`** with your actual key (it's in .gitignore)
- **Never commit unencrypted `.tfstate`** files (they're in .gitignore)
- **Never commit unencrypted `terraform.tfvars`** with secrets
- **Never share your private age key** publicly
- **Don't lose your private key** - you can't decrypt without it!

## File Structure

```
terraform/proxmox-examples/docker-host/
├── .gitignore                    # Ignores unencrypted files
├── .sops.yaml                    # Your SOPS config (NOT in Git)
├── .sops.yaml.example            # Template (in Git)
├── terraform.tfstate             # Unencrypted state (NOT in Git)
├── terraform.tfstate.enc         # Encrypted state (in Git) ✅
├── terraform.tfvars              # Your config with secrets (NOT in Git)
├── terraform.tfvars.enc          # Encrypted config (in Git) ✅
├── terraform.tfvars.example      # Template without secrets (in Git)
├── scripts/
│   ├── tf                        # Wrapper script
│   ├── tf-encrypt                # Manual encrypt
│   └── tf-decrypt                # Manual decrypt
└── STATE_MANAGEMENT.md           # This file
```

## Troubleshooting

### Error: "no key could decrypt the data"

**Cause:** SOPS can't find your private key

**Solution:**
```bash
# Set the key file location
export SOPS_AGE_KEY_FILE=~/.sops/homelab-terraform.txt

# Or add to ~/.bashrc permanently
echo 'export SOPS_AGE_KEY_FILE=~/.sops/homelab-terraform.txt' >> ~/.bashrc
```

### Error: "YOUR_AGE_PUBLIC_KEY_HERE"

**Cause:** You didn't replace the placeholder in `.sops.yaml`

**Solution:**
```bash
# Edit .sops.yaml and replace with your actual public key
nano .sops.yaml
```

### Error: "failed to get the data key"

**Cause:** The file was encrypted with a different key

**Solution:**
- Ensure you're using the same age key that encrypted the file
- If you lost the original key, you'll need to re-create the state by running `tofu import`

### Accidentally Committed Unencrypted State

**Solution:**
```bash
# Remove from Git history (DANGEROUS - coordinate with team if not solo)
git filter-branch --force --index-filter \
  'git rm --cached --ignore-unmatch terraform.tfstate' \
  --prune-empty --tag-name-filter cat -- --all

# Force push (only if solo or coordinated)
git push origin --force --all
```

### Lost Private Key

**Solution:**
- Restore from your backup (you made a backup, right?)
- If truly lost, you'll need to:
  1. Manually recreate infrastructure or import existing resources
  2. Generate a new age key
  3. Re-encrypt everything with the new key

## Advanced: Multiple Keys (Team Access)

If multiple people need access:

```yaml
# .sops.yaml
creation_rules:
  - path_regex: \.tfstate$
    age: >-
      age1person1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx,
      age1person2xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx,
      age1person3xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

Each person's private key can decrypt the files.

## Backup Strategy

### Recommended Backup Locations:

1. **Password Manager** (1Password, Bitwarden, etc.)
   ```bash
   # Copy the contents
   cat ~/.sops/homelab-terraform.txt
   # Store as a secure note in your password manager
   ```

2. **Encrypted USB Drive**
   ```bash
   # Copy to encrypted drive
   cp ~/.sops/homelab-terraform.txt /media/encrypted-usb/
   ```

3. **Encrypted Cloud Storage**
   ```bash
   # Encrypt with gpg before uploading
   gpg -c ~/.sops/homelab-terraform.txt
   # Upload homelab-terraform.txt.gpg to cloud
   ```

## Resources

- [SOPS Documentation](https://github.com/getsops/sops)
- [age Documentation](https://github.com/FiloSottile/age)
- [Terraform State Security](https://developer.hashicorp.com/terraform/language/state/sensitive-data)
- [OpenTofu Documentation](https://opentofu.org/docs/)

## Questions?

Common questions answered in this document:
- ✅ How do I set up SOPS? → See [Initial Setup](#initial-setup)
- ✅ How do I use it daily? → See [Option A: Automatic Wrapper](#option-a-automatic-wrapper-script-recommended)
- ✅ What if I lose my key? → See [Lost Private Key](#lost-private-key)
- ✅ How do I backup my key? → See [Backup Strategy](#backup-strategy)
- ✅ Can multiple people access? → See [Advanced: Multiple Keys](#advanced-multiple-keys-team-access)
