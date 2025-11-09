# Recyclarr - TRaSH Guides Sync

Automatically sync TRaSH Guides recommendations to your Radarr and Sonarr instances.

## Overview

**Recyclarr** keeps your Radarr and Sonarr configurations in sync with [TRaSH Guides](https://trash-guides.info) best practices:

- ✅ **Quality Profiles**: Optimal quality settings
- ✅ **Custom Formats**: Release group scoring
- ✅ **Quality Definitions**: File size recommendations
- ✅ **Naming Formats**: Consistent file naming
- ✅ **Automated Sync**: Runs on schedule (every 6 hours by default)
- ✅ **Multi-Instance**: Support multiple Radarr/Sonarr instances

## Quick Start

### 1. Get API Keys

**Radarr API Key:**
```bash
# Visit https://radarr.fig.systems
# Settings → General → Security → API Key
```

**Sonarr API Key:**
```bash
# Visit https://sonarr.fig.systems
# Settings → General → Security → API Key
```

### 2. Create Configuration

```bash
cd ~/homelab/compose/media/automation/recyclarr

# Create config directory
mkdir -p config

# Copy example config
cp recyclarr.yml.example config/recyclarr.yml

# Edit with your API keys
nano config/recyclarr.yml
```

**Update these values:**
```yaml
radarr:
  radarr-main:
    api_key: YOUR_RADARR_API_KEY_HERE  # <- Replace this

sonarr:
  sonarr-main:
    api_key: YOUR_SONARR_API_KEY_HERE  # <- Replace this
```

### 3. Deploy

```bash
docker compose up -d
```

### 4. Verify

Check the logs to see sync results:
```bash
docker logs recyclarr
```

You should see:
```
Processing configuration...
Syncing Radarr: radarr-main
Syncing Sonarr: sonarr-main
Sync completed successfully
```

## Configuration

### Basic Configuration

The `config/recyclarr.yml` file controls what gets synced.

**Minimal example:**
```yaml
radarr:
  radarr-main:
    base_url: http://radarr:7878
    api_key: abc123...

sonarr:
  sonarr-main:
    base_url: http://sonarr:8989
    api_key: xyz789...
```

### Quality Profiles

**Radarr - HD Bluray + WEB:**
```yaml
radarr:
  radarr-main:
    quality_profiles:
      - name: HD Bluray + WEB
        upgrade:
          allowed: true
          until_quality: Bluray-1080p
        qualities:
          - name: Bluray-1080p
          - name: WEB 1080p
          - name: Bluray-720p
          - name: WEB 720p
```

**Sonarr - WEB-1080p:**
```yaml
sonarr:
  sonarr-main:
    quality_profiles:
      - name: WEB-1080p
        upgrade:
          allowed: true
          until_quality: WEB 1080p
        qualities:
          - name: WEB 1080p
          - name: HDTV-1080p
          - name: WEB 720p
```

### Custom Formats

Custom formats score releases based on quality, source, and release group.

**High-quality release groups (+10000 points):**
```yaml
custom_formats:
  - trash_ids:
      - 3a3ff47579026e76d6504ebea39390de # Remux Tier 01
      - 9f98181fe5a3fbeb0cc29340da2a468a # Remux Tier 02
    quality_profiles:
      - name: HD Bluray + WEB
        score: 10000
```

**Unwanted formats (-10000 points):**
```yaml
custom_formats:
  - trash_ids:
      - ed38b889b31be83ffc192888e2286d83 # BR-DISK
      - 90a6f9a284dff5103f6346090e6280c8 # LQ
      - b8cd450cbfa689c0259a01d9e29ba3d6 # 3D
    quality_profiles:
      - name: HD Bluray + WEB
        score: -10000
```

**Find trash_ids:**
- [TRaSH Radarr Custom Formats](https://trash-guides.info/Radarr/Radarr-collection-of-custom-formats/)
- [TRaSH Sonarr Custom Formats](https://trash-guides.info/Sonarr/sonarr-collection-of-custom-formats/)

### Sync Schedule

**Change sync frequency:**

Edit `.env`:
```env
# Every 6 hours (default)
CRON_SCHEDULE=0 */6 * * *

# Daily at 3 AM
CRON_SCHEDULE=0 3 * * *

# Every 12 hours
CRON_SCHEDULE=0 */12 * * *

# Every hour
CRON_SCHEDULE=0 * * * *
```

**Manual sync:**
```bash
docker exec recyclarr recyclarr sync
```

## Common Configurations

### 1. HD Quality (1080p)

Good balance of quality and file size.

**Radarr:**
- Bluray-1080p preferred
- WEB-1080p as fallback
- Scores high-quality release groups

**Sonarr:**
- WEB-1080p preferred
- HDTV-1080p as fallback

### 2. Maximum Quality (Remux)

Best possible quality, large file sizes.

```yaml
radarr:
  radarr-main:
    quality_profiles:
      - name: Remux-1080p
        upgrade:
          allowed: true
          until_quality: Remux-1080p
        qualities:
          - name: Remux-1080p
          - name: Bluray-1080p
```

### 3. 4K / UHD

For 4K content:

```yaml
radarr:
  radarr-4k:
    base_url: http://radarr:7878
    api_key: abc123...
    quality_profiles:
      - name: UHD Bluray + WEB
        upgrade:
          until_quality: Remux-2160p
        qualities:
          - name: Remux-2160p
          - name: Bluray-2160p
          - name: WEB 2160p
```

### 4. Anime

Special settings for anime:

```yaml
sonarr:
  sonarr-anime:
    base_url: http://sonarr:8989
    api_key: xyz789...
    quality_profiles:
      - name: Anime
        qualities:
          - name: Bluray-1080p
          - name: WEB 1080p
    custom_formats:
      - trash_ids:
          - 064af5f084a0a24458cc8ecd3220f93f # Uncensored
        quality_profiles:
          - name: Anime
            score: 10000
```

## Integration with Radarr/Sonarr

### How It Works

1. **Recyclarr reads** your `recyclarr.yml` configuration
2. **Connects to** Radarr/Sonarr via API
3. **Syncs settings:**
   - Creates/updates quality profiles
   - Adds custom formats with scores
   - Sets quality definitions
   - Configures naming formats
4. **Your instances** now use TRaSH Guides best practices

### What Gets Changed

**Recyclarr modifies:**
- Quality profile settings
- Custom format definitions
- Quality definitions (file sizes)
- Naming formats (if configured)

**Recyclarr does NOT touch:**
- Your media files
- Download client settings
- Indexer configurations
- Root folder locations
- Existing downloads/monitoring

### Multiple Instances

Run separate Radarr/Sonarr instances for different purposes:

```yaml
radarr:
  radarr-1080p:
    base_url: http://radarr:7878
    api_key: abc123...

  radarr-4k:
    base_url: http://radarr-4k:7878
    api_key: def456...

sonarr:
  sonarr-shows:
    base_url: http://sonarr:8989
    api_key: ghi789...

  sonarr-anime:
    base_url: http://sonarr-anime:8989
    api_key: jkl012...
```

## Troubleshooting

### Recyclarr won't start

**Check logs:**
```bash
docker logs recyclarr
```

**Common issues:**
- Missing API keys in `config/recyclarr.yml`
- Radarr/Sonarr not accessible at `base_url`
- Invalid YAML syntax

### "Unable to connect to Radarr/Sonarr"

**Verify connectivity:**
```bash
# From recyclarr container
docker exec recyclarr curl http://radarr:7878
docker exec recyclarr curl http://sonarr:8989
```

**Check:**
- Radarr/Sonarr containers are running
- Both on `homelab` network
- API keys are correct

### "Profile not found"

The quality profile name in `recyclarr.yml` must match exactly.

**Check existing profiles:**
1. Go to Radarr → Settings → Profiles
2. Note the exact profile name
3. Use that name in `recyclarr.yml`

Or let Recyclarr create the profile (it will if it doesn't exist).

### Changes not appearing

**Force a sync:**
```bash
docker exec recyclarr recyclarr sync --preview
docker exec recyclarr recyclarr sync
```

**Check:**
- Look for errors in logs
- Verify API key has write permissions
- Check Radarr/Sonarr system logs

### Invalid trash_id

If you see "invalid trash_id" errors:

1. Visit [TRaSH Guides](https://trash-guides.info)
2. Find the custom format you want
3. Copy the exact trash_id from the guide
4. Update `recyclarr.yml`

## Advanced Usage

### Preview Mode

See what would change without applying:

```bash
docker exec recyclarr recyclarr sync --preview
```

### Sync Specific Instance

```bash
# Sync only Radarr
docker exec recyclarr recyclarr sync radarr

# Sync only Sonarr
docker exec recyclarr recyclarr sync sonarr

# Sync specific instance
docker exec recyclarr recyclarr sync radarr-main
```

### Validate Configuration

```bash
docker exec recyclarr recyclarr config list
docker exec recyclarr recyclarr config check
```

### Debug Mode

Enable debug logging in `.env`:
```env
LOG_LEVEL=Debug
```

Then restart:
```bash
docker compose restart
```

## Best Practices

### Start Simple

1. **First sync:** Use default quality profiles
2. **Test:** Download a movie/show
3. **Verify:** Check quality and file size
4. **Adjust:** Modify scores and profiles as needed

### Quality Profile Strategy

**For most users:**
- Radarr: HD Bluray + WEB (1080p)
- Sonarr: WEB-1080p

**For quality enthusiasts:**
- Radarr: Remux-1080p or UHD Bluray + WEB
- Sonarr: Bluray-1080p

**For storage-conscious:**
- Lower minimum quality
- Add file size limits
- Score smaller releases higher

### Custom Format Scores

**Scoring guide:**
- **+10000**: Highest priority (preferred release groups)
- **+5000**: High priority (good quality)
- **+1000**: Moderate preference
- **-10000**: Never download (BR-DISK, CAM, etc.)
- **-5000**: Avoid unless no alternative

### Sync Frequency

**Recommended schedules:**
- **Every 6 hours**: Default, good balance
- **Daily**: If you rarely update settings
- **Every hour**: If actively tuning configurations

## TRaSH Guides Resources

- [TRaSH Guides Home](https://trash-guides.info)
- [Radarr Custom Formats](https://trash-guides.info/Radarr/)
- [Sonarr Custom Formats](https://trash-guides.info/Sonarr/)
- [Recyclarr Documentation](https://recyclarr.dev)

## Comparison: Recyclarr vs Profilarr

| Feature | Recyclarr | Profilarr |
|---------|-----------|-----------|
| **Interface** | CLI / Automated | Web UI |
| **Configuration** | YAML file | Web interface |
| **Source** | TRaSH Guides | Custom + Community |
| **Automation** | Scheduled sync | Manual/scheduled |
| **Version Control** | Via Git (manual) | Built-in |
| **Learning Curve** | Moderate | Easy |
| **Best For** | Set-and-forget | Active management |

**Use both together:**
- Recyclarr for automated TRaSH sync
- Profilarr for custom tweaks and management

## Next Steps

1. ✅ Deploy Recyclarr
2. ✅ Configure API keys
3. ✅ Run first sync
4. ✅ Check Radarr/Sonarr profiles
5. ✅ Test download quality
6. ⬜ Fine-tune custom format scores
7. ⬜ Set up Profilarr for additional management
8. ⬜ Monitor sync logs periodically

---

**Automate quality, maintain consistency!** 🎬
