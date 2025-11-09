# Profilarr - Profile & Format Manager

Web-based management for Radarr/Sonarr custom formats and quality profiles with Git version control.

## Overview

**Profilarr** provides a user-friendly interface to manage your Radarr and Sonarr configurations:

- ✅ **Web UI**: Easy-to-use interface (no YAML editing)
- ✅ **Custom Formats**: Import and manage custom formats
- ✅ **Quality Profiles**: Create and sync quality profiles
- ✅ **Version Control**: Git-based configuration tracking
- ✅ **Community Configs**: Import from community databases
- ✅ **Conflict Resolution**: Handles configuration clashes
- ✅ **Multi-Instance**: Manage multiple Radarr/Sonarr instances
- ✅ **Local Preservation**: Keeps your custom changes during updates

## Quick Start

### 1. Deploy

```bash
cd ~/homelab/compose/media/automation/profilarr
docker compose up -d
```

### 2. Access Web UI

Go to: **https://profilarr.fig.systems**

### 3. Initial Setup

On first visit, you'll configure Profilarr:

**Add Radarr Instance:**
1. Click **Add Instance**
2. **Type**: Radarr
3. **Name**: Radarr Main
4. **URL**: `http://radarr:7878`
5. **API Key**: (from Radarr → Settings → General → API Key)
6. Click **Save**

**Add Sonarr Instance:**
1. Click **Add Instance**
2. **Type**: Sonarr
3. **Name**: Sonarr Main
4. **URL**: `http://sonarr:8989`
5. **API Key**: (from Sonarr → Settings → General → API Key)
6. Click **Save**

### 4. Import Custom Formats

**From Community Database:**
1. Click **Custom Formats**
2. Click **Import from Database**
3. Select formats you want (e.g., "Web Tier 01", "Bluray Tier 01")
4. Click **Import Selected**
5. Click **Sync to Radarr/Sonarr**

**From TRaSH Guides:**
1. Click **Custom Formats**
2. Click **Import from TRaSH**
3. Browse available formats
4. Select and import
5. Sync to your instances

## Features

### Custom Format Management

**What are Custom Formats?**

Custom formats score releases based on:
- Release group quality
- Source (Bluray, WEB, HDTV, etc.)
- Resolution
- Codecs
- Special editions

**Example scores:**
- **+10000**: Top-tier release groups (e.g., "FraMeSToR", "EPSiLON")
- **+5000**: High-quality groups
- **-10000**: Unwanted (BR-DISK, CAMs, low quality)

**Create Custom Format:**
1. Custom Formats → Create New
2. **Name**: `My Custom Format`
3. **Add Conditions**:
   - Release Group: `FraMeSToR|EPSiLON`
   - Source: `Bluray`
4. **Score**: 10000
5. Save and sync

### Quality Profile Management

**Create Quality Profile:**
1. Quality Profiles → Create New
2. **Name**: `HD Bluray + WEB`
3. **Add Qualities**:
   - Bluray-1080p (priority 1)
   - WEB 1080p (priority 2)
   - Bluray-720p (priority 3)
4. **Cutoff**: Bluray-1080p
5. **Custom Format Scoring**:
   - Enable custom format scores
   - Set minimum score: 0
   - Set upgrade until score: 10000
6. Save and sync

### Version Control

**Every change is tracked:**
- Automatic Git commits
- View change history
- Rollback to previous versions
- Compare configurations

**View History:**
1. Click **History** tab
2. See all changes with timestamps
3. Click any commit to view details
4. Click **Rollback** to restore

**Manual Commit:**
1. Make changes
2. Click **Commit Changes**
3. Add commit message
4. Save

### Sync to Instances

**Push to Radarr/Sonarr:**
1. Make changes in Profilarr
2. Click **Sync** button
3. Select instances to sync
4. Review changes preview
5. Click **Apply**

**Auto-sync:**
- Enable in Settings → Sync
- Choose sync interval
- Changes automatically push to instances

### Import/Export

**Export Configuration:**
1. Settings → Export
2. Download JSON
3. Store backup safely

**Import Configuration:**
1. Settings → Import
2. Upload JSON file
3. Review changes
4. Apply

## Common Workflows

### 1. Set Up HD Quality Profiles

**Radarr - HD Bluray + WEB:**

1. Quality Profiles → Create New
2. Name: `HD Bluray + WEB`
3. Add qualities:
   - Bluray-1080p
   - WEB 1080p (WEBDL-1080p, WEBRip-1080p)
   - Bluray-720p
   - WEB 720p
4. Cutoff: Bluray-1080p
5. Upgrade until score: 10000

**Add custom formats:**
- Import: "Remux Tier 01" (+10000)
- Import: "Remux Tier 02" (+9000)
- Import: "WEB Tier 01" (+8000)
- Import: "BR-DISK" (-10000)
- Import: "LQ" (-10000)

6. Sync to Radarr

**Sonarr - WEB-1080p:**

Same process:
- WEB 1080p preferred
- HDTV-1080p fallback
- Import WEB tier custom formats
- Sync to Sonarr

### 2. Manage Multiple Instances

**Scenario:** Separate 1080p and 4K Radarr instances

**Setup:**
1. Add both instances:
   - Radarr 1080p (http://radarr:7878)
   - Radarr 4K (http://radarr-4k:7878)

2. Create separate profiles:
   - Profile 1: `HD Bluray + WEB` → Sync to Radarr 1080p
   - Profile 2: `UHD Bluray + WEB` → Sync to Radarr 4K

3. Import different custom formats for each

### 3. Import Community Configurations

**Popular community configs:**

1. Custom Formats → Community
2. Browse available configs:
   - **TRaSH Recommended**: Full TRaSH Guides setup
   - **Anime Optimized**: For anime content
   - **Remux Preferred**: Maximum quality
   - **Size Optimized**: Smaller file sizes

3. Click **Import**
4. Review settings
5. Sync to your instances

### 4. Create Anime Profile

**Sonarr Anime:**

1. Add custom formats:
   - "Anime Tier 01" (+10000)
   - "Anime Tier 02" (+9000)
   - "Uncensored" (+5000)
   - "Dual Audio" (+3000)

2. Create quality profile:
   - Name: `Anime`
   - Qualities: Bluray-1080p, WEB 1080p
   - Enable custom format scoring

3. Sync to Sonarr

## Integration with Recyclarr

**Use both together for best results:**

### Recyclarr
- **Automated** TRaSH Guides sync
- **Scheduled** updates (every 6 hours)
- **YAML** configuration
- **Set and forget**

### Profilarr
- **Manual** custom configurations
- **Web UI** for easy changes
- **Version control** for tracking
- **Quick tweaks** without editing files

### Workflow

1. **Recyclarr**: Base configuration from TRaSH Guides
2. **Profilarr**: Fine-tune and customize
3. **Sync**: Both can coexist (Profilarr takes precedence on conflicts)

**Example:**
- Recyclarr syncs TRaSH quality profiles every 6 hours
- You use Profilarr to add custom formats for your favorite release groups
- Both stay in sync; Profilarr preserves your customizations

## Advanced Features

### Conflict Resolution

**What are conflicts?**

When Profilarr and Radarr/Sonarr have different configurations.

**Resolve conflicts:**
1. Conflicts tab shows all conflicts
2. Choose resolution:
   - **Use Profilarr**: Override Radarr/Sonarr
   - **Use Instance**: Keep Radarr/Sonarr config
   - **Merge**: Combine both
3. Click **Resolve**

### Scheduling

**Automatic sync schedule:**

1. Settings → Sync
2. Enable **Auto-sync**
3. Set interval:
   - Every hour
   - Every 6 hours
   - Daily
4. Save

**One-way vs Two-way:**
- **One-way**: Profilarr → Instance (recommended)
- **Two-way**: Sync both directions (can cause conflicts)

### Webhooks

**Notify on changes:**

1. Settings → Webhooks
2. Add webhook URL (Discord, Slack, etc.)
3. Choose events:
   - Configuration changed
   - Sync completed
   - Conflict detected
4. Save

### API Access

Profilarr has an API for automation:

**Get API key:**
1. Settings → API
2. Generate API Key
3. Use in scripts/automation

**Example API calls:**
```bash
# List instances
curl http://profilarr.fig.systems/api/instances

# Trigger sync
curl -X POST http://profilarr.fig.systems/api/sync \
  -H "X-API-Key: your-api-key"

# Get custom formats
curl http://profilarr.fig.systems/api/custom-formats
```

## Troubleshooting

### Can't connect to Radarr/Sonarr

**Check network:**
```bash
# Test from Profilarr container
docker exec profilarr curl http://radarr:7878
docker exec profilarr curl http://sonarr:8989
```

**Verify:**
- Containers on same `homelab` network
- API keys are correct
- URLs use container names (not localhost)

### Changes not syncing

**Force sync:**
1. Click instance name
2. Click **Force Full Sync**
3. Check sync logs

**Check:**
- Instance is reachable
- API key has write permissions
- No conflicts blocking sync

### "Profile already exists"

**Conflict resolution:**
1. Go to Conflicts tab
2. Find the profile conflict
3. Choose resolution:
   - Merge with existing
   - Replace existing
   - Rename new profile

### Database locked

If you see "database is locked" errors:

**On Windows:**
- Don't use Windows filesystem for data volume
- Use Docker volume or WSL2 filesystem

**Fix:**
```bash
# Stop container
docker compose down

# Move data to Docker volume
docker volume create profilarr-data

# Update compose.yaml to use volume instead of bind mount
```

### UI not loading

**Check logs:**
```bash
docker logs profilarr
```

**Verify:**
- Container is running: `docker ps | grep profilarr`
- Traefik routing: `docker logs traefik | grep profilarr`
- Port 6868 accessible

## Best Practices

### Configuration Management

**Organize by purpose:**
- Create separate profiles for different quality tiers
- Use descriptive names (e.g., "HD Bluray Preferred", "4K Remux Only")
- Document custom formats with notes

### Version Control

**Commit regularly:**
- After major changes
- Before experimenting
- Use descriptive commit messages
- Tag stable configurations

### Backup

**Export regularly:**
```bash
# Export via UI
Settings → Export → Download

# Or backup data directory
cd ~/homelab/compose/media/automation/profilarr
tar czf profilarr-backup-$(date +%Y%m%d).tar.gz ./data
```

### Sync Strategy

**Recommended:**
- **Manual sync** for testing changes
- **Auto-sync** once you're confident
- **Preview** before applying large changes

### Instance Organization

**Name instances clearly:**
- `Radarr - 1080p`
- `Radarr - 4K`
- `Sonarr - TV Shows`
- `Sonarr - Anime`

## Getting API Keys

### Radarr API Key

1. Go to https://radarr.fig.systems
2. Settings → General
3. Security section
4. Copy **API Key**

### Sonarr API Key

1. Go to https://sonarr.fig.systems
2. Settings → General
3. Security section
4. Copy **API Key**

## Resource Usage

**Typical usage:**
- **RAM**: 100-200MB
- **CPU**: Low (spikes during sync)
- **Disk**: <500MB
- **Network**: Minimal

## Next Steps

1. ✅ Deploy Profilarr
2. ✅ Add Radarr/Sonarr instances
3. ✅ Import custom formats
4. ✅ Create quality profiles
5. ✅ Sync to instances
6. ✅ Test download quality
7. ⬜ Set up auto-sync
8. ⬜ Configure version control
9. ⬜ Export backup configuration

## Resources

- [Profilarr Website](https://dictionarry.dev)
- [Profilarr GitHub](https://github.com/Dictionarry-Hub/profilarr)
- [TRaSH Guides](https://trash-guides.info)
- [Radarr Wiki](https://wiki.servarr.com/radarr)
- [Sonarr Wiki](https://wiki.servarr.com/sonarr)

---

**Manage profiles with ease!** 🎛️
