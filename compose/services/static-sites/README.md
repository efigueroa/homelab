# Caddy Static Sites Server

Serves static websites for edfig.dev (professional), blog.edfig.dev (blog), and figgy.foo (experimental).

## Overview

**Caddy** is a modern web server with automatic HTTPS and simple configuration:

- ✅ **Static file serving** - HTML, CSS, JavaScript, images
- ✅ **Markdown rendering** - Write `.md` files, served as HTML automatically
- ✅ **Templates** - Dynamic content with Go templates
- ✅ **Directory browsing** - Beautiful file listing (figgy.foo)
- ✅ **Auto-compression** - Gzip for all responses
- ✅ **Zero-downtime reloads** - Config changes apply instantly

## Domain Strategy

### edfig.dev (Professional/Public)
- **Purpose**: Personal website, portfolio
- **URL**: https://edfig.dev or https://www.edfig.dev
- **SSO**: No (public site)
- **Content**: `/sites/edfig.dev/`

### blog.edfig.dev (Blog/Public)
- **Purpose**: Technical blog, articles
- **URL**: https://blog.edfig.dev
- **SSO**: No (public blog)
- **Content**: `/sites/blog.edfig.dev/`
- **Features**: Markdown auto-rendering, templates

### figgy.foo (Experimental/Private)
- **Purpose**: Testing, development, experiments
- **URL**: https://figgy.foo or https://www.figgy.foo
- **SSO**: Yes (protected by Tinyauth)
- **Content**: `/sites/figgy.foo/`
- **Features**: Directory browsing, templates

## Quick Start

### 1. Deploy

```bash
cd ~/homelab/compose/services/static-sites
docker compose up -d
```

### 2. Access Sites

- **edfig.dev**: https://edfig.dev
- **Blog**: https://blog.edfig.dev
- **Experimental**: https://figgy.foo (requires SSO login)

### 3. Verify

```bash
# Check container is running
docker ps | grep caddy-static

# Check logs
docker logs caddy-static

# Test sites
curl -I https://edfig.dev
curl -I https://blog.edfig.dev
```

## Directory Structure

```
static-sites/
├── compose.yaml          # Docker Compose + Traefik labels
├── Caddyfile            # Caddy configuration
├── .env                 # Environment variables
├── .gitignore           # Ignored files
├── README.md            # This file
└── sites/               # Site content (can be version controlled)
    ├── edfig.dev/
    │   ├── index.html
    │   ├── assets/
    │   │   ├── css/
    │   │   ├── js/
    │   │   └── images/
    │   └── ...
    ├── blog.edfig.dev/
    │   ├── index.html
    │   └── posts/
    │       ├── example-post.md    # Markdown posts
    │       └── ...
    └── figgy.foo/
        ├── index.html
        └── experiments/
            └── ...
```

## Managing Content

### Adding/Editing HTML

Simply edit files in the `sites/` directory:

```bash
# Edit main site
vim sites/edfig.dev/index.html

# Add new page
echo "<h1>About Me</h1>" > sites/edfig.dev/about.html

# Changes are live immediately (no restart needed!)
```

### Writing Blog Posts (Markdown)

Create `.md` files in `sites/blog.edfig.dev/posts/`:

```bash
# Create new post
cat > sites/blog.edfig.dev/posts/my-post.md << 'EOF'
# My New Blog Post

**Published:** January 10, 2025

This is my blog post content...

## Code Example

```bash
docker compose up -d
```

[Back to Blog](/)
EOF

# Access at: https://blog.edfig.dev/posts/my-post.md
# (renders as HTML automatically!)
```

**Markdown features:**
- Headers (`#`, `##`, `###`)
- **Bold**, *italic*, `code`
- Links, images
- Lists (ordered/unordered)
- Code blocks with syntax highlighting
- Tables
- Blockquotes

### Using Templates

Caddy supports Go templates for dynamic content:

**Example - Current time:**
```html
<!-- In any .html file under blog.edfig.dev -->
<p>Page generated at: {{.Now.Format "2006-01-02 15:04:05"}}</p>
```

**Example - Include header:**
```html
{{include "header.html"}}
<main>
  <h1>My Page</h1>
</main>
{{include "footer.html"}}
```

**Template variables:**
- `{{.Now}}` - Current time
- `{{.Req.URL}}` - Request URL
- `{{.Req.Host}}` - Request hostname
- `{{.Req.Method}}` - HTTP method
- `{{env "VARIABLE"}}` - Environment variable

See [Caddy Templates Docs](https://caddyserver.com/docs/caddyfile/directives/templates)

### Directory Browsing (figgy.foo)

figgy.foo has directory browsing enabled:

```bash
# Add files to browse
cp some-file.txt sites/figgy.foo/experiments/

# Access: https://figgy.foo/experiments/
# Shows beautiful file listing with search!
```

## Adding New Sites

### Option 1: New Subdomain (same domain)

**Add to Caddyfile:**
```caddy
test.figgy.foo {
    root * /srv/test.figgy.foo
    file_server
    encode gzip
}
```

**Add Traefik labels to compose.yaml:**
```yaml
# test.figgy.foo
traefik.http.routers.figgy-test.rule: Host(`test.figgy.foo`)
traefik.http.routers.figgy-test.entrypoints: websecure
traefik.http.routers.figgy-test.tls.certresolver: letsencrypt
traefik.http.routers.figgy-test.service: caddy-static
traefik.http.routers.figgy-test.middlewares: tinyauth  # If SSO needed
```

**Create site directory:**
```bash
mkdir -p sites/test.figgy.foo
echo "<h1>Test Site</h1>" > sites/test.figgy.foo/index.html
```

**Reload (instant, no restart):**
```bash
# Caddy auto-reloads when Caddyfile changes!
# Just wait 1-2 seconds, then access https://test.figgy.foo
```

### Option 2: New Domain

Follow same process but use new domain name. Make sure DNS points to your server.

## Caddyfile Features

### Basic Site
```caddy
example.com {
    root * /srv/example
    file_server
}
```

### With Compression
```caddy
example.com {
    root * /srv/example
    file_server
    encode gzip zstd brotli
}
```

### With Caching
```caddy
example.com {
    root * /srv/example
    file_server

    @static {
        path *.css *.js *.jpg *.png *.gif *.ico
    }
    header @static Cache-Control "public, max-age=604800"
}
```

### With Redirects
```caddy
www.example.com {
    redir https://example.com{uri} permanent
}

example.com {
    root * /srv/example
    file_server
}
```

### With Custom 404
```caddy
example.com {
    root * /srv/example
    file_server
    handle_errors {
        rewrite * /404.html
        file_server
    }
}
```

### With Basic Auth (alternative to SSO)
```caddy
example.com {
    root * /srv/example
    basicauth {
        user $2a$14$hashedpassword
    }
    file_server
}
```

Generate hashed password:
```bash
docker exec caddy-static caddy hash-password --plaintext "mypassword"
```

## Traefik Integration

All sites route through Traefik:

```
Internet → DNS (*.edfig.dev, *.figgy.foo)
         ↓
    Traefik (SSL termination)
         ↓
    Tinyauth (SSO check for figgy.foo only)
         ↓
    Caddy (static file serving)
```

**SSL certificates:**
- Traefik handles Let's Encrypt
- Caddy receives plain HTTP on port 80
- Users see HTTPS

**SSO protection:**
- `edfig.dev` & `blog.edfig.dev`: No SSO (public)
- `figgy.foo`: SSO protected (private)

## Performance

### Caching

Static assets automatically cached:

```caddy
@static {
    path *.css *.js *.jpg *.jpeg *.png *.gif *.ico *.svg
}
header @static Cache-Control "public, max-age=604800, immutable"
```

- 7 days cache for images, CSS, JS
- Browsers won't re-request until expired

### Compression

All responses auto-compressed with gzip:

```caddy
encode gzip
```

- 70-90% size reduction for HTML/CSS/JS
- Faster page loads
- Lower bandwidth usage

### Performance Tips

1. **Optimize images**: Use WebP format, compress before uploading
2. **Minify CSS/JS**: Use build tools (optional)
3. **Use CDN**: For high-traffic sites (optional)
4. **Enable HTTP/2**: Traefik handles this automatically

## Monitoring

### Check Service Status

```bash
# Container status
docker ps | grep caddy-static

# Logs
docker logs caddy-static -f

# Resource usage
docker stats caddy-static
```

### Check Specific Site

```bash
# Test site is reachable
curl -I https://edfig.dev

# Test with timing
curl -w "@curl-format.txt" -o /dev/null -s https://edfig.dev

# Check SSL certificate
echo | openssl s_client -connect edfig.dev:443 -servername edfig.dev 2>/dev/null | openssl x509 -noout -dates
```

### Access Logs

Caddy logs to stdout (captured by Docker):

```bash
# View logs
docker logs caddy-static

# Follow logs
docker logs caddy-static -f

# Last 100 lines
docker logs caddy-static --tail 100
```

### Grafana Logs

All logs forwarded to Loki automatically:

**Query in Grafana** (https://logs.fig.systems):
```logql
{container="caddy-static"}
```

Filter by status code:
```logql
{container="caddy-static"} |= "404"
```

## Troubleshooting

### Site not loading

**Check container:**
```bash
docker ps | grep caddy-static
# If not running:
docker compose up -d
```

**Check logs:**
```bash
docker logs caddy-static
# Look for errors in Caddyfile or file not found
```

**Check DNS:**
```bash
dig +short edfig.dev
# Should point to your server IP
```

**Check Traefik:**
```bash
# See if Traefik sees the route
docker logs traefik | grep edfig
```

### 404 Not Found

**Check file exists:**
```bash
ls -la sites/edfig.dev/index.html
```

**Check path in Caddyfile:**
```bash
grep "root" Caddyfile
# Should show: root * /srv/edfig.dev
```

**Check permissions:**
```bash
# Files should be readable
chmod -R 755 sites/
```

### Changes not appearing

**Caddy auto-reloads**, but double-check:

```bash
# Check file modification time
ls -lh sites/edfig.dev/index.html

# Force reload (shouldn't be needed)
docker exec caddy-static caddy reload --config /etc/caddy/Caddyfile
```

**Browser cache:**
```bash
# Force refresh in browser: Ctrl+Shift+R (Linux/Win) or Cmd+Shift+R (Mac)
# Or open in incognito/private window
```

### Markdown not rendering

**Check templates enabled:**
```caddy
# In Caddyfile for blog.edfig.dev
blog.edfig.dev {
    templates  # <-- This must be present!
    # ...
}
```

**Check file extension:**
```bash
# Must be .md
mv post.txt post.md
```

**Test rendering:**
```bash
curl https://blog.edfig.dev/posts/example-post.md
# Should return HTML, not raw markdown
```

### SSO not working on figgy.foo

**Check middleware:**
```yaml
# In compose.yaml
traefik.http.routers.figgy-main.middlewares: tinyauth
```

**Check Tinyauth is running:**
```bash
docker ps | grep tinyauth
```

**Test without SSO:**
```bash
# Temporarily remove SSO to isolate issue
# Comment out middleware line in compose.yaml
# docker compose up -d
```

## Backup

### Backup Site Content

```bash
# Backup all sites
cd ~/homelab/compose/services/static-sites
tar czf sites-backup-$(date +%Y%m%d).tar.gz sites/

# Backup to external storage
scp sites-backup-*.tar.gz user@backup-server:/backups/
```

### Version Control (Optional)

Consider using Git for your sites:

```bash
cd sites/
git init
git add .
git commit -m "Initial site content"

# Add remote
git remote add origin git@github.com:efigueroa/sites.git
git push -u origin main
```

## Security

### Public vs Private

**Public sites** (`edfig.dev`, `blog.edfig.dev`):
- No SSO middleware
- Accessible to everyone
- Use for portfolio, blog, public content

**Private sites** (`figgy.foo`):
- SSO middleware enabled
- Requires LLDAP authentication
- Use for experiments, private content

### Content Security

**Don't commit:**
- API keys
- Passwords
- Private information
- Sensitive data

**Do commit:**
- HTML, CSS, JS
- Images, assets
- Markdown blog posts
- Public content

### File Permissions

```bash
# Sites should be read-only to Caddy
chmod -R 755 sites/
chown -R $USER:$USER sites/
```

## Resources

- [Caddy Documentation](https://caddyserver.com/docs/)
- [Caddyfile Tutorial](https://caddyserver.com/docs/caddyfile-tutorial)
- [Templates Documentation](https://caddyserver.com/docs/caddyfile/directives/templates)
- [Markdown Rendering](https://caddyserver.com/docs/caddyfile/directives/templates#markdown)

## Next Steps

1. ✅ Deploy Caddy static sites
2. ✅ Access edfig.dev, blog.edfig.dev, figgy.foo
3. ⬜ Customize edfig.dev with your content
4. ⬜ Write first blog post in Markdown
5. ⬜ Add experiments to figgy.foo
6. ⬜ Set up Git version control for sites
7. ⬜ Configure automated backups

---

**Serve static content, simply and securely!** 🌐
