# Example Blog Post

**Published:** January 10, 2025
**Tags:** #homelab #docker #traefik

---

## Introduction

This is an example blog post written in Markdown. Caddy automatically renders `.md` files as HTML!

## Why Markdown?

Markdown is perfect for writing blog posts because:

1. **Simple syntax** - Easy to write and read
2. **Fast** - No build step required
3. **Portable** - Works everywhere
4. **Clean** - Focus on content, not formatting

## Code Examples

Here's some example code:

```bash
# Deploy a service
cd ~/homelab/compose/services/example
docker compose up -d

# Check logs
docker logs example-service -f
```

## Features

### Supported Elements

- **Bold text**
- *Italic text*
- `Code snippets`
- [Links](https://edfig.dev)
- Lists (ordered and unordered)
- Code blocks with syntax highlighting
- Blockquotes
- Tables

### Example Table

| Service | URL | Purpose |
|---------|-----|---------|
| Traefik | traefik.fig.systems | Reverse Proxy |
| Sonarr | sonarr.fig.systems | TV Automation |
| Radarr | radarr.fig.systems | Movie Automation |

## Blockquote Example

> "The best way to predict the future is to invent it."
> — Alan Kay

## Conclusion

This is just an example post. Delete this file and create your own posts in the `posts/` directory!

Each `.md` file will be automatically rendered when accessed via the browser.

---

[← Back to Blog](/)
