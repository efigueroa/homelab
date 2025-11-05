# Security Policy

## Supported Versions

This is a personal homelab configuration repository. The latest commit on `main` is always the supported version.

| Branch | Supported          |
| ------ | ------------------ |
| main   | :white_check_mark: |
| other  | :x:                |

## Security Considerations

### Secrets Management

**DO NOT commit secrets to this repository!**

- All passwords in `.env` files should use placeholder values (e.g., `changeme_*`)
- Real passwords should only be set in your local deployment
- Use environment variables or Docker secrets for sensitive data
- Never commit files containing real credentials

### Container Security

- All container images are scanned for vulnerabilities via GitHub Actions
- HIGH and CRITICAL vulnerabilities are reported in security scans
- Keep images up to date by pulling latest versions regularly
- Review security scan results before deploying

### Network Security

- All services are behind Traefik reverse proxy
- SSL/TLS is enforced via Let's Encrypt
- Internal services use isolated Docker networks
- SSO is enabled on most services via Tinyauth

### Authentication

- LLDAP provides centralized user management
- Tinyauth handles SSO authentication
- Services with built-in authentication are documented in README
- Change all default passwords before deployment

## Reporting a Vulnerability

If you discover a security vulnerability in this configuration:

1. **DO NOT** open a public issue
2. Contact the repository owner directly via GitHub private message
3. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Report

- Exposed secrets or credentials
- Insecure configurations
- Vulnerable container images (not already detected by CI)
- Authentication bypasses
- Network security issues

### What NOT to Report

- Issues with third-party services (report to their maintainers)
- Theoretical vulnerabilities without proof of concept
- Social engineering attempts

## Security Best Practices

### Before Deployment

1. **Change all passwords** in `.env` files
2. **Review** all service configurations
3. **Update** container images to latest versions
4. **Configure** firewall to only allow ports 80/443
5. **Enable** automatic security updates on host OS

### After Deployment

1. **Monitor** logs regularly for suspicious activity
2. **Update** services monthly (at minimum)
3. **Backup** data regularly
4. **Review** access logs
5. **Test** disaster recovery procedures

### Network Hardening

- Use a firewall (ufw, iptables, etc.)
- Only expose ports 80 and 443 to the internet
- Consider using a VPN for administrative access
- Enable fail2ban or similar intrusion prevention
- Use strong DNS providers with DNSSEC

### Container Hardening

- Run containers as non-root when possible
- Use read-only filesystems where applicable
- Limit container resources (CPU, memory)
- Enable security options (no-new-privileges, etc.)
- Regularly scan for vulnerabilities

## Automated Security Scanning

This repository includes automated security scanning:

- **Gitleaks**: Detects secrets in commits
- **Trivy**: Scans container images for vulnerabilities
- **YAML Linting**: Ensures proper configuration
- **Dependency Review**: Checks for vulnerable dependencies

Review GitHub Actions results before merging PRs.

## Compliance

This is a personal homelab configuration and does not claim compliance with any specific security standards. However, it follows general security best practices:

- Principle of least privilege
- Defense in depth
- Secure by default
- Regular updates and patching

## External Dependencies

Security of this setup depends on:

- Docker and Docker Compose security
- Container image maintainers
- Traefik security
- LLDAP security
- Host OS security

Always keep these dependencies up to date.

## Disclaimer

This configuration is provided "as is" without warranty. Use at your own risk. The maintainer is not responsible for any security incidents resulting from the use of this configuration.

## Additional Resources

- [Docker Security Best Practices](https://docs.docker.com/engine/security/)
- [Traefik Security Documentation](https://doc.traefik.io/traefik/https/overview/)
- [OWASP Container Security](https://cheatsheetseries.owasp.org/cheatsheets/Docker_Security_Cheat_Sheet.html)
