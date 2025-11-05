# Contributing Guide

Thank you for your interest in contributing to this homelab configuration! While this is primarily a personal repository, contributions are welcome.

## How to Contribute

### Reporting Issues

- Use the [bug report template](.github/ISSUE_TEMPLATE/bug-report.md) for bugs
- Use the [service request template](.github/ISSUE_TEMPLATE/service-request.md) for new services
- Search existing issues before creating a new one
- Provide as much detail as possible

### Submitting Changes

1. **Fork the repository**
2. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes** following the guidelines below
4. **Test your changes** locally
5. **Commit with clear messages**
   ```bash
   git commit -m "feat: add new service"
   ```
6. **Push to your fork**
   ```bash
   git push origin feature/your-feature-name
   ```
7. **Open a Pull Request** using the PR template

## Guidelines

### File Naming

- All Docker Compose files must be named `compose.yaml` (not `.yml`)
- Use lowercase with hyphens for service directories (e.g., `calibre-web`)
- Environment files must be named `.env`

### Docker Compose Best Practices

- Use version-pinned images when possible
- Include health checks for databases and critical services
- Use bind mounts for configuration, named volumes for data
- Set proper restart policies (`unless-stopped` or `always`)
- Include resource limits for production services

### Network Configuration

- All services must use the `homelab` network (marked as `external: true`)
- Services with multiple containers should use an internal network
- Example:
  ```yaml
  networks:
    homelab:
      external: true
    service_internal:
      name: service_internal
      driver: bridge
  ```

### Traefik Labels

All web services must include:

```yaml
labels:
  traefik.enable: true
  traefik.http.routers.service.rule: Host(`service.fig.systems`) || Host(`service.edfig.dev`)
  traefik.http.routers.service.entrypoints: websecure
  traefik.http.routers.service.tls.certresolver: letsencrypt
  traefik.http.services.service.loadbalancer.server.port: 8080
  # Optional SSO:
  traefik.http.routers.service.middlewares: tinyauth
```

### Environment Variables

- Use `.env` files for configuration
- Never commit real passwords
- Use `changeme_*` prefix for placeholder passwords
- Document all required environment variables
- Include comments explaining non-obvious settings

### Documentation

- Add service to README.md service table
- Include deployment instructions
- Document any special configuration
- Add comments to compose files explaining purpose
- Include links to official documentation

### Security

- Never commit secrets
- Scan compose files for vulnerabilities
- Use official or well-maintained images
- Enable SSO when appropriate
- Document security considerations

## Code Style

### YAML Style

- 2-space indentation
- No trailing whitespace
- Use `true/false` instead of `yes/no`
- Quote strings with special characters
- Follow yamllint rules in `.yamllint.yml`

### Commit Messages

Follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` New feature
- `fix:` Bug fix
- `docs:` Documentation changes
- `refactor:` Code refactoring
- `security:` Security improvements
- `chore:` Maintenance tasks

Examples:
```
feat: add jellyfin media server
fix: correct traefik routing for sonarr
docs: update README with new services
security: update postgres to latest version
```

## Testing

Before submitting a PR:

1. **Validate compose files**
   ```bash
   docker compose -f compose/path/to/compose.yaml config
   ```

2. **Check YAML syntax**
   ```bash
   yamllint compose/
   ```

3. **Test locally**
   ```bash
   docker compose up -d
   docker compose logs
   ```

4. **Check for secrets**
   ```bash
   git diff --cached | grep -i "password\|secret\|token"
   ```

5. **Run pre-commit hooks** (optional)
   ```bash
   pre-commit install
   pre-commit run --all-files
   ```

## Pull Request Process

1. Fill out the PR template completely
2. Ensure all CI checks pass
3. Request review if needed
4. Address review feedback
5. Squash commits if requested
6. Wait for approval and merge

## CI/CD Checks

Your PR will be automatically checked for:

- Docker Compose validation
- YAML linting
- Security scanning
- Secret detection
- Documentation completeness
- Traefik configuration
- Network setup
- File naming conventions

Fix any failures before requesting review.

## Adding a New Service

1. Choose the correct category:
   - `compose/core/` - Infrastructure (Traefik, auth, etc.)
   - `compose/media/` - Media-related services
   - `compose/services/` - Utility services

2. Create service directory:
   ```bash
   mkdir -p compose/category/service-name
   ```

3. Create `compose.yaml`:
   - Include documentation header
   - Add Traefik labels
   - Configure networks
   - Set up volumes
   - Add health checks if applicable

4. Create `.env` if needed:
   - Use placeholder passwords
   - Document all variables
   - Include comments

5. Update README.md:
   - Add to service table
   - Include URL
   - Document deployment

6. Test deployment:
   ```bash
   cd compose/category/service-name
   docker compose up -d
   docker compose logs -f
   ```

7. Create PR with detailed description

## Project Structure

```
homelab/
├── .github/
│   ├── workflows/        # CI/CD workflows
│   ├── ISSUE_TEMPLATE/   # Issue templates
│   └── pull_request_template.md
├── compose/
│   ├── core/            # Infrastructure services
│   ├── media/           # Media services
│   └── services/        # Utility services
├── README.md            # Main documentation
├── CONTRIBUTING.md      # This file
├── SECURITY.md          # Security policy
└── .yamllint.yml        # YAML linting config
```

## Getting Help

- Check existing issues and PRs
- Review the README.md
- Examine similar services for examples
- Ask in PR comments

## License

By contributing, you agree that your contributions will be licensed under the same terms as the repository.

## Code of Conduct

- Be respectful and professional
- Focus on constructive feedback
- Help others learn and improve
- Keep discussions relevant

## Questions?

Open an issue with the question label or comment on an existing PR/issue.

Thank you for contributing! 🎉
