# Telegram Music Monitor

A self-hosted tool that monitors Telegram group chats for music links (Spotify, YouTube, etc.), extracts metadata, enriches it with MusicBrainz data, and automatically adds artists to Lidarr.

## Features

- **Multi-Platform Support**: Detects and processes links from Spotify and YouTube
- **Automatic Metadata Extraction**: Extracts song, album, and artist information
- **MusicBrainz Integration**: Enriches metadata with canonical artist information
- **Lidarr Integration**: Automatically adds artists to your Lidarr instance
- **SQLite Database**: Tracks all detected music and metadata
- **REST API**: Query statistics and metrics
- **Web Dashboard**: Beautiful real-time dashboard to view activity
- **Docker Support**: Easy deployment with Docker Compose

## Architecture

```
Telegram Groups → Monitor → Extract Metadata → Enrich with MusicBrainz → Add to Lidarr
                              ↓
                         SQLite Database
                              ↓
                         REST API + Dashboard
```

## Prerequisites

Before you begin, you'll need to obtain API credentials for the following services:

### 1. Telegram API Credentials

1. Go to https://my.telegram.org
2. Log in with your phone number
3. Go to "API development tools"
4. Create a new application
5. Note down your `API ID` and `API Hash`

### 2. Spotify API Credentials

1. Go to https://developer.spotify.com/dashboard
2. Log in or create an account
3. Click "Create an App"
4. Note down your `Client ID` and `Client Secret`

### 3. YouTube Data API Key

1. Go to https://console.cloud.google.com
2. Create a new project (or select an existing one)
3. Enable the "YouTube Data API v3"
4. Go to "Credentials" and create an API key
5. Note down your API key

### 4. Lidarr

You need a running Lidarr instance. You'll need:
- Lidarr URL (e.g., `http://localhost:8686`)
- Lidarr API Key (found in Lidarr Settings → General)

## Installation

### Using Docker Compose (Recommended)

1. Clone this repository:
```bash
git clone https://github.com/efigueroa/telegram-music-monitor.git
cd telegram-music-monitor
```

2. Edit `docker-compose.yml` and replace the placeholder values with your actual credentials:
```yaml
environment:
  - TELEGRAM_API_ID=YOUR_API_ID
  - TELEGRAM_API_HASH=YOUR_API_HASH
  - TELEGRAM_PHONE=+1234567890
  - SPOTIFY_CLIENT_ID=YOUR_SPOTIFY_CLIENT_ID
  - SPOTIFY_CLIENT_SECRET=YOUR_SPOTIFY_CLIENT_SECRET
  - YOUTUBE_API_KEY=YOUR_YOUTUBE_API_KEY
  - LIDARR_URL=http://lidarr:8686
  - LIDARR_API_KEY=YOUR_LIDARR_API_KEY
```

3. Create the data directory:
```bash
mkdir -p data
```

4. Build and start the container:
```bash
docker-compose up -d
```

5. **First Run - Authentication**:

   On the first run, you'll need to authenticate with Telegram. Check the logs:
   ```bash
   docker-compose logs -f music-monitor
   ```

   When you see "Enter code:", you'll need to enter the authentication code sent to your Telegram app:
   ```bash
   docker exec -it telegram-music-monitor sh
   # Then enter the code when prompted
   ```

   After successful authentication, the session will be saved and you won't need to do this again.

6. Access the dashboard at http://localhost:8080

### Building from Source

1. Install Go 1.21 or later

2. Clone the repository:
```bash
git clone https://github.com/efigueroa/telegram-music-monitor.git
cd telegram-music-monitor
```

3. Install dependencies:
```bash
go mod download
```

4. Set environment variables (or create a `.env` file):
```bash
export TELEGRAM_API_ID=your_api_id
export TELEGRAM_API_HASH=your_api_hash
export TELEGRAM_PHONE=+1234567890
export SPOTIFY_CLIENT_ID=your_client_id
export SPOTIFY_CLIENT_SECRET=your_client_secret
export YOUTUBE_API_KEY=your_youtube_key
export LIDARR_URL=http://localhost:8686
export LIDARR_API_KEY=your_lidarr_key
export DATABASE_PATH=./data/music_monitor.db
export SERVER_PORT=8080
```

5. Run the application:
```bash
go run main.go
```

## Usage

### Dashboard

Access the web dashboard at http://localhost:8080

The dashboard shows:
- **Recent Music**: Latest detected music links with metadata
- **Top Artists**: Most frequently shared artists
- **Top Albums**: Most frequently shared albums
- **Group Activity**: Activity breakdown by Telegram group
- **Platform Stats**: Distribution of music platforms (Spotify, YouTube)

### API Endpoints

All API endpoints return JSON responses.

#### Get Recent Music
```
GET /api/metrics/recent?limit=50
```

Returns the most recently detected music links.

#### Get Top Artists
```
GET /api/metrics/top-artists?limit=10
```

Returns the most frequently detected artists.

#### Get Top Albums
```
GET /api/metrics/top-albums?limit=10
```

Returns the most frequently detected albums.

#### Get Group Activity
```
GET /api/metrics/group-activity
```

Returns activity statistics by Telegram group.

#### Get Platform Stats
```
GET /api/metrics/platform-stats
```

Returns statistics by platform (Spotify, YouTube, etc.).

#### Health Check
```
GET /api/health
```

Returns `{"status": "ok"}` if the service is running.

## How It Works

1. **Message Monitoring**: The application runs as a Telegram userbot and monitors all messages in groups it's a member of.

2. **Link Detection**: When a message contains a Spotify or YouTube link, it's detected and processed.

3. **Metadata Extraction**:
   - For Spotify: Uses Spotify API to get track/album/artist information
   - For YouTube: Uses YouTube Data API and parses video titles for artist/song names

4. **MusicBrainz Enrichment**: Searches MusicBrainz for the artist to get canonical name and MusicBrainz ID.

5. **Database Storage**: All detected music is stored in SQLite with full metadata.

6. **Lidarr Integration**: If a MusicBrainz ID is found, the artist is automatically added to Lidarr (if not already present).

## Configuration

All configuration is done via environment variables:

| Variable | Required | Description |
|----------|----------|-------------|
| `TELEGRAM_API_ID` | Yes | Telegram API ID from my.telegram.org |
| `TELEGRAM_API_HASH` | Yes | Telegram API Hash from my.telegram.org |
| `TELEGRAM_PHONE` | Yes | Phone number in international format (e.g., +1234567890) |
| `SPOTIFY_CLIENT_ID` | Yes | Spotify Client ID |
| `SPOTIFY_CLIENT_SECRET` | Yes | Spotify Client Secret |
| `YOUTUBE_API_KEY` | Yes | YouTube Data API key |
| `LIDARR_URL` | Yes | Lidarr instance URL |
| `LIDARR_API_KEY` | Yes | Lidarr API key |
| `DATABASE_PATH` | No | Path to SQLite database (default: /data/music_monitor.db) |
| `SERVER_PORT` | No | API server port (default: 8080) |

## Database

The application uses SQLite to store all data. The database is automatically created on first run.

### Tables

- **detected_music**: Stores all detected music links with metadata
- **group_configs**: Stores per-group configuration (future use for enabling/disabling monitoring)

### Location

By default, the database is stored at `/data/music_monitor.db` inside the container, which maps to `./data/music_monitor.db` on the host.

## Supported Platforms

### Spotify
- Single tracks
- Albums
- Playlists (processes all tracks)
- Artist pages

### YouTube
- Individual videos
- Playlists (processes up to 50 videos)

## Limitations

- YouTube title parsing may not always correctly identify artist/song (depends on video title format)
- Playlists are processed but currently only the first track/video is added to the database (future enhancement will process all items)
- MusicBrainz has a rate limit of 1 request per second (the application respects this automatically)
- YouTube API has daily quota limits (10,000 units per day by default)

## Troubleshooting

### Authentication Issues

If you see "Not authenticated" in the logs:
1. Make sure your Telegram credentials are correct
2. Connect to the container and enter the code: `docker exec -it telegram-music-monitor sh`
3. Check that the session file has write permissions

### Lidarr Connection Issues

If artists aren't being added to Lidarr:
1. Check that Lidarr URL is accessible from the container
2. Verify the API key is correct
3. Ensure Lidarr has at least one root folder, quality profile, and metadata profile configured

### API Rate Limits

- **Spotify**: No strict rate limits for client credentials flow
- **YouTube**: 10,000 quota units per day (50 video requests = ~50-100 units)
- **MusicBrainz**: 1 request per second (automatically handled)

## Development

### Project Structure

```
telegram-music-monitor/
├── api/              # HTTP API handlers and dashboard
├── config/           # Configuration loading
├── database/         # SQLite database layer
├── lidarr/          # Lidarr API client
├── musicbrainz/     # MusicBrainz API client
├── platforms/       # Platform-specific clients (Spotify, YouTube)
├── telegram/        # Telegram userbot
├── web/             # Web dashboard (embedded)
├── main.go          # Application entry point
├── Dockerfile       # Multi-stage Docker build
└── docker-compose.yml
```

### Code Comments

The code is extensively commented to help Go beginners understand how everything works. Each function, struct, and important code block has explanatory comments.

### Building

```bash
# Build the Docker image
docker build -t telegram-music-monitor .

# Or build the Go binary directly
go build -o music-monitor .
```

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions, please open an issue on GitHub.
