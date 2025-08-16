# CI/CD Workflows

This directory contains GitHub Actions workflows for continuous integration and deployment.

## Workflows

### CI (`ci.yml`)
Runs on every push and pull request to master/main branch.

**Jobs:**
- **Lint**: Code formatting, imports, and linting checks
- **Test**: Unit tests across multiple Go versions and OS platforms
- **Security**: Security scanning with gosec and govulncheck
- **Build**: Multi-platform builds (Linux, macOS, Windows)
- **Integration**: Integration test suite
- **Dependency Review**: Security review of dependencies (PRs only)

### Release (`release.yml`)
Runs when a version tag is pushed (e.g., `v1.0.0`).

**Features:**
- Automated packaging with `package.sh`
- GitHub release creation
- Binary artifact uploads

## Local Development

### Prerequisites
```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install goimports
go install golang.org/x/tools/cmd/goimports@latest
```

### Commands
```bash
# Run all checks locally
make lint

# Format code
make fmt

# Run tests
make test

# Build
make build
```

## Configuration

### golangci-lint (`.golangci.yml`)
Comprehensive linting configuration with:
- Code quality checks
- Security scanning
- Performance analysis
- Style enforcement

### Dependabot (`.github/dependabot.yml`)
Automated dependency updates:
- Go modules: Weekly updates
- GitHub Actions: Weekly updates
- Automatic PR creation with reviews

## Best Practices

1. **Always run `make lint` before committing**
2. **Ensure tests pass locally before pushing**
3. **Use conventional commit messages**
4. **Review Dependabot PRs promptly**
5. **Tag releases with semantic versioning (v1.0.0)**
