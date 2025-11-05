## Description

<!-- Provide a brief description of what this PR does -->

## Type of Change

<!-- Mark the relevant option with an "x" -->

- [ ] New service addition
- [ ] Service configuration update
- [ ] Bug fix
- [ ] Documentation update
- [ ] Security fix
- [ ] Infrastructure change

## Changes Made

<!-- List the main changes in this PR -->

-
-
-

## Checklist

<!-- Mark completed items with an "x" -->

### General
- [ ] All compose files use `compose.yaml` (not `.yml`)
- [ ] Code follows Docker Compose best practices
- [ ] Changes tested locally
- [ ] Documentation updated (README.md)

### Services (if applicable)
- [ ] Service added to correct category (core/media/services)
- [ ] Proper network configuration (homelab + internal if needed)
- [ ] Volumes properly configured
- [ ] Environment variables use `.env` file or are documented

### Traefik & SSL (if applicable)
- [ ] Traefik labels configured correctly
- [ ] Uses `websecure` entrypoint
- [ ] Let's Encrypt cert resolver configured
- [ ] Both domains configured (`fig.systems` and `edfig.dev`)
- [ ] SSO middleware applied (if appropriate)

### Security
- [ ] No secrets committed in `.env` files
- [ ] Placeholder passwords use `changeme_*` format
- [ ] No sensitive data in compose files
- [ ] Container runs as non-root user (where possible)

### Documentation
- [ ] Service added to README.md service table
- [ ] Deployment instructions added/updated
- [ ] Configuration requirements documented
- [ ] Comments added to compose file explaining purpose

## Testing

<!-- Describe how you tested these changes -->

```bash
# Commands used to test:


# Expected behavior:


# Actual behavior:

```

## Screenshots (if applicable)

<!-- Add screenshots of the service running, configuration, etc. -->

## Related Issues

<!-- Link any related issues: Fixes #123, Closes #456 -->

## Additional Notes

<!-- Any additional context, breaking changes, migration notes, etc. -->

---

## For Reviewers

<!-- Automatically checked by CI/CD -->

- [ ] All CI checks pass
- [ ] Docker Compose validation passes
- [ ] YAML linting passes
- [ ] Security scans pass
- [ ] No security vulnerabilities introduced
