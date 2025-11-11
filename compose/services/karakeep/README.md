# Karakeep - Bookmark Everything App

AI-powered bookmark manager for links, notes, images, and PDFs with automatic tagging and full-text search.

## Overview

**Karakeep** (previously known as Hoarder) is a self-hostable bookmark-everything app:

- ✅ **Bookmark Everything**: Links, notes, images, PDFs
- ✅ **AI-Powered**: Automatic tagging and summarization
- ✅ **Full-Text Search**: Find anything instantly with Meilisearch
- ✅ **Web Archiving**: Save complete webpages (full page archive)
- ✅ **Browser Extensions**: Chrome and Firefox support
- ✅ **Mobile Apps**: iOS and Android apps available
- ✅ **Ollama Support**: Use local AI models (no cloud required!)
- ✅ **OCR**: Extract text from images
- ✅ **Self-Hosted**: Full control of your data

## Quick Start

### 1. Configure Secrets

```bash
cd ~/homelab/compose/services/karakeep

# Edit .env and update:
# - NEXTAUTH_SECRET (generate with: openssl rand -base64 36)
# - MEILI_MASTER_KEY (generate with: openssl rand -base64 36)
nano .env
```

### 2. Deploy

```bash
docker compose up -d
```

### 3. Access

Go to: **https://links.fig.systems**

**First-time setup:**
1. Create your admin account
2. Start bookmarking!

## Features

### Bookmark Types

**1. Web Links**
- Save any URL
- Automatic screenshot capture
- Full webpage archiving
- Extract title, description, favicon
- AI-generated summary and tags

**2. Notes**
- Quick text notes
- Markdown support
- AI-powered categorization
- Full-text searchable

**3. Images**
- Upload images directly
- OCR text extraction (if enabled)
- AI-based tagging
- Image search

**4. PDFs**
- Upload PDF documents
- Full-text indexing
- Searchable content

### AI Features

Karakeep can use AI to automatically:
- **Tag** your bookmarks
- **Summarize** web content
- **Extract** key information
- **Organize** by category

**Three AI options:**

**1. Ollama (Recommended - Local & Free)**
```env
# In .env, uncomment:
OLLAMA_BASE_URL=http://ollama:11434
INFERENCE_TEXT_MODEL=llama3.2:3b
INFERENCE_IMAGE_MODEL=llava:7b
```

**2. OpenAI**
```env
OPENAI_API_KEY=sk-...
OPENAI_BASE_URL=https://api.openai.com/v1
INFERENCE_TEXT_MODEL=gpt-4o-mini
```

**3. OpenRouter (multiple providers)**
```env
OPENAI_API_KEY=sk-or-v1-...
OPENAI_BASE_URL=https://openrouter.ai/api/v1
INFERENCE_TEXT_MODEL=anthropic/claude-3.5-sonnet
```

### Web Archiving

Karakeep saves complete web pages for offline viewing:
- **Full HTML archive**
- **Screenshots** of the page
- **Extracted text** for search
- **Works offline** - view archived pages anytime

### Search

Powered by Meilisearch:
- **Instant** full-text search
- **Fuzzy matching** - finds similar terms
- **Filter by** type, tags, dates
- **Search across** titles, content, tags, notes

### Browser Extensions

**Install extensions:**
- [Chrome Web Store](https://chromewebstore.google.com/detail/karakeep/kbkejgonjhbmhcaofkhdegeoeoemgkdm)
- [Firefox Add-ons](https://addons.mozilla.org/en-US/firefox/addon/karakeep/)

**Configure extension:**
1. Install extension
2. Click extension icon
3. Enter server URL: `https://links.fig.systems`
4. Login with your credentials
5. Save bookmarks from any page!

### Mobile Apps

**Download apps:**
- [iOS App Store](https://apps.apple.com/app/karakeep/id6479258022)
- [Android Google Play](https://play.google.com/store/apps/details?id=app.karakeep.mobile)

**Setup:**
1. Install app
2. Open app
3. Enter server: `https://links.fig.systems`
4. Login
5. Bookmark on the go!

## Configuration

### Basic Settings

**Disable public signups:**
```env
DISABLE_SIGNUPS=true
```

**Set max file size (100MB default):**
```env
MAX_ASSET_SIZE_MB=100
```

**Enable OCR for multiple languages:**
```env
OCR_LANGS=eng,spa,fra,deu
```

### Ollama Integration

**Prerequisites:**
1. Deploy Ollama service (see `compose/services/ollama/`)
2. Pull models: `docker exec ollama ollama pull llama3.2:3b`

**Enable in Karakeep:**
```env
# In karakeep/.env
OLLAMA_BASE_URL=http://ollama:11434
INFERENCE_TEXT_MODEL=llama3.2:3b
INFERENCE_IMAGE_MODEL=llava:7b
INFERENCE_LANG=en
```

**Restart:**
```bash
docker compose restart
```

**Recommended models:**
- **Text**: llama3.2:3b (fast, good quality)
- **Images**: llava:7b (vision model)
- **Advanced**: llama3.3:70b (slower, better results)

### Advanced Settings

**Custom logging:**
```env
LOG_LEVEL=debug  # Options: debug, info, warn, error
```

**Custom data directory:**
```env
DATADIR=/custom/path
```

**Chrome timeout (for slow sites):**
```env
# Add to compose.yaml environment section
BROWSER_TIMEOUT=60000  # 60 seconds
```

## Usage Workflows

### 1. Bookmark a Website

**Via Browser:**
1. Click Karakeep extension
2. Bookmark opens automatically
3. AI generates tags and summary
4. Edit tags/notes if needed
5. Save

**Via Mobile:**
1. Open share menu
2. Select Karakeep
3. Bookmark saved

**Manually:**
1. Open Karakeep
2. Click "+" button
3. Paste URL
4. Click Save

### 2. Quick Note

1. Open Karakeep
2. Click "+" → "Note"
3. Type your note
4. AI auto-tags
5. Save

### 3. Upload Image

1. Click "+" → "Image"
2. Upload image file
3. OCR extracts text (if enabled)
4. AI generates tags
5. Save

### 4. Search Everything

**Simple search:**
- Type in search box
- Results appear instantly

**Advanced search:**
- Filter by type (links, notes, images)
- Filter by tags
- Filter by date range
- Sort by relevance or date

### 5. Organize with Tags

**Auto-tags:**
- AI generates tags automatically
- Based on content analysis
- Can be edited/removed

**Manual tags:**
- Add your own tags
- Create tag hierarchies
- Color-code tags

**Tag management:**
- Rename tags globally
- Merge duplicate tags
- Delete unused tags

## Browser Extension Usage

### Quick Bookmark

1. **Visit any page**
2. **Click extension icon** (or keyboard shortcut)
3. **Automatically saved** with:
   - URL
   - Title
   - Screenshot
   - Full page archive
   - AI tags and summary

### Save Selection

1. **Highlight text** on any page
2. **Right-click** → "Save to Karakeep"
3. **Saves as note** with source URL

### Save Image

1. **Right-click image**
2. Select "Save to Karakeep"
3. **Image uploaded** with AI tags

## Mobile App Features

- **Share from any app** to Karakeep
- **Quick capture** - bookmark in seconds
- **Offline access** to archived content
- **Search** your entire collection
- **Browse by tags**
- **Dark mode** support

## Data Management

### Backup

**Important data locations:**
```bash
compose/services/karakeep/
├── data/           # Uploaded files, archives
└── meili_data/     # Search index
```

**Backup script:**
```bash
#!/bin/bash
cd ~/homelab/compose/services/karakeep
tar czf karakeep-backup-$(date +%Y%m%d).tar.gz ./data ./meili_data
```

### Export

**Export bookmarks:**
1. Settings → Export
2. Choose format:
   - JSON (complete data)
   - HTML (browser-compatible)
   - CSV (spreadsheet)
3. Download

### Import

**Import from other services:**
1. Settings → Import
2. Select source:
   - Browser bookmarks (HTML)
   - Pocket
   - Raindrop.io
   - Omnivore
   - Instapaper
3. Upload file
4. Karakeep processes and imports

## Troubleshooting

### Karakeep won't start

**Check logs:**
```bash
docker logs karakeep
docker logs karakeep-chrome
docker logs karakeep-meilisearch
```

**Common issues:**
- Missing `NEXTAUTH_SECRET` in `.env`
- Missing `MEILI_MASTER_KEY` in `.env`
- Services not on `karakeep_internal` network

### Bookmarks not saving

**Check chrome service:**
```bash
docker logs karakeep-chrome
```

**Verify chrome is accessible:**
```bash
docker exec karakeep curl http://karakeep-chrome:9222
```

**Increase timeout:**
```env
# Add to .env
BROWSER_TIMEOUT=60000
```

### Search not working

**Rebuild search index:**
```bash
# Stop services
docker compose down

# Remove search data
rm -rf ./meili_data

# Restart (index rebuilds automatically)
docker compose up -d
```

**Check Meilisearch:**
```bash
docker logs karakeep-meilisearch
```

### AI features not working

**With Ollama:**
```bash
# Verify Ollama is running
docker ps | grep ollama

# Test Ollama connection
docker exec karakeep curl http://ollama:11434

# Check models are pulled
docker exec ollama ollama list
```

**With OpenAI/OpenRouter:**
- Verify API key is correct
- Check API balance/credits
- Review logs for error messages

### Extension can't connect

**Verify server URL:**
- Must be `https://links.fig.systems`
- Not `http://` or `localhost`

**Check CORS:**
```env
# Add to .env if needed
CORS_ALLOW_ORIGINS=https://links.fig.systems
```

**Clear extension data:**
1. Extension settings
2. Logout
3. Clear extension storage
4. Login again

### Mobile app issues

**Can't connect:**
- Use full HTTPS URL
- Ensure server is accessible externally
- Check firewall rules

**Slow performance:**
- Check network speed
- Reduce image quality in app settings
- Enable "Low data mode"

## Performance Optimization

### For Large Collections (10,000+ bookmarks)

**Increase Meilisearch RAM:**
```yaml
# In compose.yaml, add to karakeep-meilisearch:
deploy:
  resources:
    limits:
      memory: 2G
    reservations:
      memory: 1G
```

**Optimize search index:**
```env
# In .env
MEILI_MAX_INDEXING_MEMORY=1048576000  # 1GB
```

### For Slow Archiving

**Increase Chrome resources:**
```yaml
# In compose.yaml, add to karakeep-chrome:
deploy:
  resources:
    limits:
      memory: 1G
      cpus: '1.0'
```

**Adjust timeouts:**
```env
BROWSER_TIMEOUT=90000  # 90 seconds
```

### Database Maintenance

**Vacuum (compact) database:**
```bash
# Karakeep uses SQLite by default
docker exec karakeep sqlite3 /data/karakeep.db "VACUUM;"
```

## Comparison with Linkwarden

| Feature | Karakeep | Linkwarden |
|---------|----------|------------|
| **Bookmark Types** | Links, Notes, Images, PDFs | Links only |
| **AI Tagging** | Yes (Ollama/OpenAI) | No |
| **Web Archiving** | Full page + Screenshot | Screenshot only |
| **Search** | Meilisearch (fuzzy) | Meilisearch |
| **Browser Extension** | Yes | Yes |
| **Mobile Apps** | iOS + Android | No official apps |
| **OCR** | Yes | No |
| **Collaboration** | Personal focus | Team features |
| **Database** | SQLite | PostgreSQL |

**Why Karakeep?**
- More bookmark types
- AI-powered organization
- Better mobile support
- Lighter resource usage (SQLite vs PostgreSQL)
- Active development

## Resources

- [Official Website](https://karakeep.app)
- [Documentation](https://docs.karakeep.app)
- [GitHub Repository](https://github.com/karakeep-app/karakeep)
- [Demo Instance](https://try.karakeep.app)
- [Chrome Extension](https://chromewebstore.google.com/detail/karakeep/kbkejgonjhbmhcaofkhdegeoeoemgkdm)
- [Firefox Extension](https://addons.mozilla.org/en-US/firefox/addon/karakeep/)

## Next Steps

1. ✅ Deploy Karakeep
2. ✅ Create admin account
3. ✅ Install browser extension
4. ✅ Install mobile app
5. ⬜ Deploy Ollama for AI features
6. ⬜ Import existing bookmarks
7. ⬜ Configure AI models
8. ⬜ Set up automated backups

---

**Bookmark everything, find anything!** 🔖
