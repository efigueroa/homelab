# Services Overview

Complete list of all services in the homelab with descriptions and use cases.

## Core Infrastructure (Required)

### Traefik
- **URL**: https://traefik.fig.systems
- **Purpose**: Reverse proxy with automatic SSL/TLS
- **Why**: Routes all traffic, manages Let's Encrypt certificates
- **Required**: ✅ Yes - Nothing works without this

### LLDAP
- **URL**: https://lldap.fig.systems
- **Purpose**: Lightweight LDAP directory for user management
- **Why**: Centralized user database for SSO
- **Required**: ✅ Yes (if using SSO)
- **Default Login**: admin / <your LLDAP_LDAP_USER_PASS>

### Tinyauth
- **URL**: https://auth.fig.systems
- **Purpose**: SSO forward authentication middleware
- **Why**: Single login for all services
- **Required**: ✅ Yes (if using SSO)

## Dashboard & Management

### Homarr
- **URL**: https://home.fig.systems
- **Purpose**: Service dashboard with auto-discovery
- **Why**: See all your services in one place, monitor status
- **Required**: ⬜ No, but highly recommended
- **Features**:
  - Auto-discovers Docker containers
  - Customizable widgets
  - Service status monitoring
  - Integration with media services

### Backrest
- **URL**: https://backup.fig.systems
- **Purpose**: Backup management with web UI (uses Restic)
- **Why**: Encrypted, deduplicated backups to Backblaze B2
- **Required**: ⬜ No, but critical for data safety
- **Features**:
  - Web-based backup management
  - Scheduled backups
  - File browsing and restore
  - Encryption at rest
  - S3-compatible storage support

## Media Services

### Jellyfin
- **URL**: https://flix.fig.systems
- **Purpose**: Media server (Netflix alternative)
- **Why**: Watch your movies/TV shows anywhere
- **Required**: ⬜ No
- **Features**:
  - Stream to any device
  - Hardware transcoding (with GPU)
  - Live TV & DVR
  - Mobile apps available
  - Subtitle support

### Immich
- **URL**: https://photos.fig.systems
- **Purpose**: Photo and video management (Google Photos alternative)
- **Why**: Self-hosted photo library with ML features
- **Required**: ⬜ No
- **Features**:
  - Face recognition (with GPU)
  - Object detection
  - Mobile apps with auto-upload
  - Timeline view
  - Album organization

### Jellyseerr
- **URL**: https://requests.fig.systems
- **Purpose**: Media request management
- **Why**: Let users request movies/shows
- **Required**: ⬜ No (only if using Sonarr/Radarr)
- **Features**:
  - Request movies and TV shows
  - Integration with Jellyfin
  - User permissions
  - Notification system

## Media Automation

### Sonarr
- **URL**: https://sonarr.fig.systems
- **Purpose**: TV show automation
- **Why**: Automatically download and organize TV shows
- **Required**: ⬜ No
- **Features**:
  - Episode tracking
  - Automatic downloading
  - Quality management
  - Calendar view

### Radarr
- **URL**: https://radarr.fig.systems
- **Purpose**: Movie automation
- **Why**: Automatically download and organize movies
- **Required**: ⬜ No
- **Features**:
  - Movie tracking
  - Automatic downloading
  - Quality profiles
  - Collection management

### SABnzbd
- **URL**: https://sabnzbd.fig.systems
- **Purpose**: Usenet downloader
- **Why**: Download from Usenet newsgroups
- **Required**: ⬜ No (only if using Usenet)
- **Features**:
  - Fast downloads
  - Automatic verification and repair
  - Category-based processing
  - Password support

### qBittorrent
- **URL**: https://qbt.fig.systems
- **Purpose**: BitTorrent client
- **Why**: Download torrents
- **Required**: ⬜ No (only if using torrents)
- **Features**:
  - Web-based UI
  - RSS support
  - Sequential downloading
  - IP filtering

## Productivity Services

### Linkwarden
- **URL**: https://links.fig.systems
- **Purpose**: Bookmark manager
- **Why**: Save and organize web links
- **Required**: ⬜ No
- **Features**:
  - Collaborative bookmarking
  - Full-text search
  - Screenshots and PDFs
  - Tags and collections
  - Browser extensions

### Vikunja
- **URL**: https://tasks.fig.systems
- **Purpose**: Task management (Todoist alternative)
- **Why**: Track tasks and projects
- **Required**: ⬜ No
- **Features**:
  - Kanban boards
  - Lists and sub-tasks
  - Due dates and reminders
  - Collaboration
  - CalDAV support

### FreshRSS
- **URL**: https://rss.fig.systems
- **Purpose**: RSS/Atom feed reader
- **Why**: Aggregate news and blogs
- **Required**: ⬜ No
- **Features**:
  - Web-based reader
  - Mobile apps via API
  - Filtering and search
  - Multi-user support

## Specialized Services

### LubeLogger
- **URL**: https://garage.fig.systems
- **Purpose**: Vehicle maintenance tracker
- **Why**: Track mileage, maintenance, costs
- **Required**: ⬜ No
- **Features**:
  - Service records
  - Fuel tracking
  - Cost analysis
  - Reminder system
  - Export data

### Calibre-web
- **URL**: https://books.fig.systems
- **Purpose**: Ebook library manager
- **Why**: Manage and read ebooks
- **Required**: ⬜ No
- **Features**:
  - Web-based ebook reader
  - Format conversion
  - Metadata management
  - Send to Kindle
  - OPDS support

### Booklore
- **URL**: https://booklore.fig.systems
- **Purpose**: Book tracking and reviews
- **Why**: Track reading progress and reviews
- **Required**: ⬜ No
- **Features**:
  - Reading lists
  - Progress tracking
  - Reviews and ratings
  - Import from Goodreads

### RSSHub
- **URL**: https://rsshub.fig.systems
- **Purpose**: RSS feed generator
- **Why**: Generate RSS feeds for sites without them
- **Required**: ⬜ No
- **Features**:
  - 1000+ source support
  - Custom routes
  - Filter and transform feeds

### MicroBin
- **URL**: https://paste.fig.systems
- **Purpose**: Encrypted pastebin with file upload
- **Why**: Share code snippets and files
- **Required**: ⬜ No
- **Features**:
  - Encryption support
  - File uploads
  - Burn after reading
  - Custom expiry
  - Password protection

### File Browser
- **URL**: https://files.fig.systems
- **Purpose**: Web-based file manager
- **Why**: Browse and manage media files
- **Required**: ⬜ No
- **Features**:
  - Upload/download files
  - Preview images and videos
  - Text editor
  - File sharing
  - User permissions

## Service Categories

### Minimum Viable Setup
Just want to get started? Deploy these:
1. Traefik
2. LLDAP
3. Tinyauth
4. Homarr

### Media Enthusiast Setup
For streaming media:
1. Core services (above)
2. Jellyfin
3. Sonarr
4. Radarr
5. qBittorrent
6. Jellyseerr

### Complete Homelab
Everything:
1. Core services
2. All media services
3. All productivity services
4. Backrest for backups

## Resource Requirements

### Light (2 Core, 4GB RAM)
- Core services
- Homarr
- 2-3 utility services

### Medium (4 Core, 8GB RAM)
- Core services
- Media services (without transcoding)
- Most utility services

### Heavy (6+ Core, 16GB+ RAM)
- All services
- GPU transcoding
- Multiple concurrent users

## Quick Deploy Checklist

**Before deploying a service:**
- ✅ Core infrastructure is running
- ✅ `.env` file configured with secrets
- ✅ DNS record created
- ✅ Understand what the service does
- ✅ Know how to configure it

**After deploying:**
- ✅ Check container is running: `docker ps`
- ✅ Check logs: `docker compose logs`
- ✅ Access web UI and complete setup
- ✅ Test SSO if applicable
- ✅ Add to Homarr dashboard

## Service Dependencies

```
Traefik (required for all)
├── LLDAP
│   └── Tinyauth
│       └── All SSO-protected services
├── Jellyfin
│   └── Jellyseerr
│       ├── Sonarr
│       └── Radarr
│           ├── SABnzbd
│           └── qBittorrent
├── Immich
│   └── Backrest (for backups)
└── All other services
```

## When to Use Each Service

### Use Jellyfin if:
- You have a movie/TV collection
- Want to stream from anywhere
- Have family/friends who want access
- Want apps on all devices

### Use Immich if:
- You want Google Photos alternative
- Have lots of photos to manage
- Want ML features (face recognition)
- Have mobile devices

### Use Sonarr/Radarr if:
- You watch a lot of TV/movies
- Want automatic downloads
- Don't want to manually search
- Want quality control

### Use Backrest if:
- You care about your data (you should!)
- Want encrypted cloud backups
- Have important photos/documents
- Want easy restore process

### Use Linkwarden if:
- You save lots of bookmarks
- Want full-text search
- Share links with team
- Want offline archives

### Use Vikunja if:
- You need task management
- Work with teams
- Want Kanban boards
- Need CalDAV for calendar integration

## Next Steps

1. Review which services you actually need
2. Start with core + 2-3 services
3. Deploy and configure each fully
4. Add more services gradually
5. Monitor resource usage

---

**Remember**: You don't need all services. Start small and add what you actually use!
