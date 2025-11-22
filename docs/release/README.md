# GitHub Actions Workflows

This directory contains GitHub Actions workflows for CI/CD automation.

## Workflows

### `ci.yml`
Continuous Integration workflow that runs on:
- Pull requests to main/master
- Pushes to main/master

**What it does:**
- Builds the Go application
- Builds the Docker image
- Verifies image size
- Tests the Docker container
- Validates all endpoints

### `release.yml`
Release workflow that runs on:
- Version tags (e.g., `v1.0.0`)
- Manual workflow dispatch

**What it does:**
- Builds the Docker image with optimizations
- Pushes to Docker Hub with version and latest tags
- Creates a GitHub release

**Required Secrets:**
- `DOCKER_USERNAME` - Docker Hub username
- `DOCKER_PASSWORD` - Docker Hub access token

## Setting up Secrets

1. Go to your repository Settings → Secrets and variables → Actions
2. Add the following secrets:
   - `DOCKER_USERNAME`: Your Docker Hub username
   - `DOCKER_PASSWORD`: Your Docker Hub access token (not password)

## Creating a Release

### Automatic (via tag):
```bash
git tag v1.0.0
git push origin v1.0.0
```

### Manual (via GitHub UI):
1. Go to Actions → Release workflow
2. Click "Run workflow"
3. Enter the version (e.g., `1.0.0`)
4. Click "Run workflow"

