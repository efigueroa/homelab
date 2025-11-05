# Backrest - Web UI for Restic Backups

Backrest provides a modern web interface for managing Restic backups with support for Backblaze B2 (S3-compatible) storage.

## Features

- **Web UI**: User-friendly interface for managing backups
- **Encrypted Backups**: Uses Restic for encrypted, deduplicated backups
- **Scheduling**: Built-in scheduler for automatic backups
- **Monitoring**: View backup status, logs, and statistics
- **Multiple Repositories**: Support for multiple backup destinations
- **Retention Policies**: Automatic cleanup of old backups
- **Notifications**: Email and webhook notifications for backup events

## Quick Start

### 1. Start the Service

```bash
cd compose/services/backrest
docker compose up -d
```

### 2. Access Web UI

Open your browser and navigate to:
- https://backup.fig.systems
- Or: https://backup.edfig.dev

Login with your SSO credentials (tinyauth).

### 3. Configure Backblaze B2 Repository

1. **Create B2 Bucket** (if not already done):
   - Go to https://secure.backblaze.com/b2_buckets.htm
   - Click "Create a Bucket"
   - Name: `homelab-backups` (or your choice)
   - Files: Private
   - Encryption: Server-Side (or Disabled - Backrest encrypts client-side)

2. **Create Application Key**:
   - Go to https://secure.backblaze.com/app_keys.htm
   - Click "Add a New Application Key"
   - Name: `backrest-homelab`
   - Access: Read and Write
   - Bucket: Select your backup bucket
   - Save the `keyID` and `applicationKey`

3. **Add Repository in Backrest**:
   - Click "Add Repository"
   - Repository Name: `B2 Immich Photos`
   - Storage Type: `S3-compatible storage`
   - Configuration:
     ```
     Endpoint: s3.us-west-002.backblazeb2.com
     Region: us-west-002
     Bucket: homelab-backups
     Path: /immich-photos
     Access Key ID: [your B2 keyID]
     Secret Access Key: [your B2 applicationKey]
     ```
   - Encryption Password: Set a strong password (SAVE THIS!)
   - Click "Initialize Repository"

### 4. Create Backup Plan

1. **Add Plan**:
   - Click "Add Plan"
   - Plan Name: `Immich Daily Backup`
   - Repository: Select your B2 repository
   - Paths to Backup:
     - `/backups/immich`
   - Exclude Patterns (optional):
     - `*.tmp`
     - `*.log`

2. **Schedule**:
   - Backup Schedule: `0 3 * * *` (3 AM daily)
   - Enable "Automatic Backups"

3. **Retention Policy**:
   - Keep Last: 7 daily backups
   - Keep Weekly: 4 weekly backups
   - Keep Monthly: 6 monthly backups
   - Keep Yearly: 2 yearly backups

4. **Notifications** (optional):
   - Configure email or webhook for backup status
   - Alert on failures

### 5. Run First Backup

Click "Run Now" to start your first backup immediately.

## Backup Locations

The service has access to these directories:

- `/backups/immich` - Immich photos (read-only)
- `/backups/homelab-config` - All compose configurations (read-only)

You can add more volumes in `compose.yaml` as needed.

## Monitoring

### View Backup Status

In the Backrest web UI:
- Dashboard shows all backup plans and their status
- Click on a plan to see backup history
- View logs for detailed information

### Check Repository Size

```bash
# Via web UI: Repository → Statistics
# Shows: Total size, deduplicated size, number of snapshots
```

### Verify Backups

Backrest has built-in verification:
1. Go to Repository → Verify
2. Click "Run Verification"
3. Check results for any errors

## Restore Files

### Via Web UI

1. Go to Plan → Snapshots
2. Select snapshot to restore
3. Click "Browse Files"
4. Select files/folders to restore
5. Choose restore location
6. Click "Restore"

### Via CLI (Advanced)

```bash
# List snapshots
docker exec backrest restic -r [repository] snapshots

# Restore specific snapshot
docker exec backrest restic -r [repository] restore [snapshot-id] --target /restore

# Restore specific file
docker exec backrest restic -r [repository] restore [snapshot-id] --target /restore --include /path/to/file
```

## Configuration Backup

### Backup Backrest Config

Your Backrest configuration (plans, schedules, repositories) is stored in:
- `./config/config.json`

**Important**: Backup this file! It contains your repository credentials (encrypted).

```bash
# Create backup
cp config/config.json config/config.json.backup

# Restore backup
cp config/config.json.backup config/config.json
docker compose restart
```

### Export Configuration

In Web UI:
1. Settings → Export Configuration
2. Save JSON file securely
3. Store encryption passwords separately

## Troubleshooting

### Cannot Access Web UI

Check container status:
```bash
docker compose logs backrest
docker compose ps
```

Verify Traefik routing:
```bash
docker logs traefik | grep backrest
```

### Backup Fails

1. **Check Logs**:
   - Web UI: Plan → View Logs
   - Or: `docker compose logs -f backrest`

2. **Verify B2 Credentials**:
   - Test connection in Repository settings
   - Ensure application key has read/write access

3. **Check Disk Space**:
   ```bash
   df -h
   docker exec backrest df -h /cache
   ```

### Repository Locked

If a backup is interrupted, the repository may be locked:

```bash
# Via Web UI: Repository → Unlock
# Or via CLI:
docker exec backrest restic -r [repository] unlock
```

### Slow Backups

1. **Enable Caching**: Already configured via `XDG_CACHE_HOME`
2. **Increase Upload Speed**: Check B2 endpoint is geographically close
3. **Exclude Unnecessary Files**: Add patterns to exclude list

## Security Considerations

### Encryption

- **Client-side**: All data encrypted before upload
- **Repository Password**: Required to access backups
- **Storage**: Store repository passwords in password manager

### Access Control

- **SSO Protected**: Web UI requires authentication via tinyauth
- **API Keys**: B2 application keys scoped to specific bucket
- **Read-Only Mounts**: Backup sources mounted read-only

### Best Practices

1. **Test Restores**: Regularly test restoring files
2. **Monitor Backups**: Check backup status weekly
3. **Verify Integrity**: Run verification monthly
4. **Secure Passwords**: Use strong, unique repository passwords
5. **Document Recovery**: Keep recovery procedures documented
6. **Offsite Storage**: B2 provides geographic redundancy

## Advanced Configuration

### Add More Backup Sources

Edit `compose.yaml` to add more volumes:

```yaml
volumes:
  - /path/to/backup:/backups/name:ro
```

Then create a new backup plan in the web UI.

### Multiple Repositories

Configure multiple destinations:
1. Primary: Backblaze B2
2. Secondary: Local NAS/USB drive
3. Archive: Another cloud provider

### Webhooks

Configure webhooks for monitoring:
1. Settings → Notifications
2. Add Webhook URL (e.g., Discord, Slack, Uptime Kuma)
3. Select events: Backup Success, Backup Failure

### Custom Retention

Fine-tune retention policies:
```
--keep-within 7d
--keep-within-daily 30d
--keep-within-weekly 90d
--keep-within-monthly 1y
--keep-within-yearly 5y
```

## Resource Usage

**Typical Usage:**
- CPU: Low (spikes during backup)
- Memory: ~200-500MB
- Disk: Cache grows over time (monitor)
- Network: Depends on backup size

**Monitoring Cache Size:**
```bash
du -sh compose/services/backrest/cache
```

Clean cache if needed (safe to delete - will rebuild):
```bash
rm -rf compose/services/backrest/cache/*
docker compose restart
```

## Backrest vs Duplicati

We chose Backrest over Duplicati because:
- **Modern**: Built on Restic (actively developed)
- **Performance**: Better deduplication and compression
- **Reliability**: Restic is battle-tested
- **Features**: More advanced scheduling and monitoring
- **UI**: Clean, responsive interface

## Cost Estimation

**Backblaze B2 Pricing (2024):**
- Storage: $0.006/GB/month
- Download: $0.01/GB (first 3x storage free)
- Upload: Free

**Example: 100GB Immich photos**
- Storage Cost: $0.60/month
- Download (3 restores/month): Free
- **Total: ~$0.60/month**

**With Deduplication:**
- First backup: 100GB
- Daily incrementals: ~1-5GB
- Monthly growth: ~20GB
- Avg monthly cost: ~$0.70

## Resources

- [Backrest Documentation](https://github.com/garethgeorge/backrest)
- [Restic Documentation](https://restic.readthedocs.io/)
- [Backblaze B2 Documentation](https://www.backblaze.com/b2/docs/)
- [S3-compatible API Guide](https://www.backblaze.com/b2/docs/s3_compatible_api.html)

## Next Steps

1. ✅ Configure B2 repository
2. ✅ Create backup plan for Immich
3. ⬜ Run initial backup
4. ⬜ Verify backup integrity
5. ⬜ Test restore procedure
6. ⬜ Set up notifications
7. ⬜ Add homelab-config backups
8. ⬜ Schedule monthly verification
