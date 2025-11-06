# Homelab Documentation

Welcome to the homelab documentation! This folder contains comprehensive guides for setting up, configuring, and maintaining your self-hosted services.

## 📚 Documentation Structure

### Quick Start
- [Getting Started](./getting-started.md) - First-time setup walkthrough
- [Quick Reference](./quick-reference.md) - Common commands and URLs

### Configuration
- [Environment Variables & Secrets](./guides/secrets-management.md) - How to configure secure secrets
- [DNS Configuration](./guides/dns-setup.md) - Setting up domain names
- [SSL/TLS Certificates](./guides/ssl-certificates.md) - Let's Encrypt configuration
- [GPU Acceleration](./guides/gpu-setup.md) - NVIDIA GPU setup for Jellyfin and Immich

### Services
- [Service Overview](./services/README.md) - All available services
- [SSO Configuration](./services/sso-setup.md) - Single Sign-On with LLDAP and Tinyauth
- [Media Stack](./services/media-stack.md) - Jellyfin, Sonarr, Radarr setup
- [Backup Solutions](./services/backup.md) - Backrest configuration

### Troubleshooting
- [Common Issues](./troubleshooting/common-issues.md) - Frequent problems and solutions
- [FAQ](./troubleshooting/faq.md) - Frequently asked questions
- [Debugging Guide](./troubleshooting/debugging.md) - How to diagnose problems

### Operations
- [Maintenance](./operations/maintenance.md) - Regular maintenance tasks
- [Updates](./operations/updates.md) - Updating services
- [Backups](./operations/backups.md) - Backup and restore procedures
- [Monitoring](./operations/monitoring.md) - Service monitoring

## 🚀 Quick Links

### First Time Setup
1. [Prerequisites](./getting-started.md#prerequisites)
2. [Configure Secrets](./guides/secrets-management.md)
3. [Setup DNS](./guides/dns-setup.md)
4. [Deploy Services](./getting-started.md#deployment)

### Common Tasks
- [Add a new service](./guides/adding-services.md)
- [Generate secure passwords](./guides/secrets-management.md#generating-secrets)
- [Enable GPU acceleration](./guides/gpu-setup.md)
- [Backup configuration](./operations/backups.md)
- [Update a service](./operations/updates.md)

### Troubleshooting
- [Service won't start](./troubleshooting/common-issues.md#service-wont-start)
- [SSL certificate errors](./troubleshooting/common-issues.md#ssl-errors)
- [SSO not working](./troubleshooting/common-issues.md#sso-issues)
- [Can't access service](./troubleshooting/common-issues.md#access-issues)

## 📖 Documentation Conventions

Throughout this documentation:
- `command` - Commands to run in terminal
- **Bold** - Important concepts or UI elements
- `https://service.fig.systems` - Example URLs
- ⚠️ - Warning or important note
- 💡 - Tip or helpful information
- ✅ - Verified working configuration

## 🔐 Security Notes

Before deploying to production:
1. ✅ Change all passwords in `.env` files
2. ✅ Configure DNS records
3. ✅ Verify SSL certificates are working
4. ✅ Enable backups
5. ✅ Review security settings

## 🆘 Getting Help

If you encounter issues:
1. Check [Common Issues](./troubleshooting/common-issues.md)
2. Review [FAQ](./troubleshooting/faq.md)
3. Check service logs: `docker compose logs servicename`
4. Review the [Debugging Guide](./troubleshooting/debugging.md)

## 📝 Contributing to Documentation

Found an error or have a suggestion? Documentation improvements are welcome!
- Keep guides clear and concise
- Include examples and code snippets
- Test all commands before documenting
- Update the table of contents when adding new files

## 🔄 Last Updated

This documentation is automatically maintained and reflects the current state of the homelab repository.
