# Contributing to Awesome ClickUp CLI

Thank you for your interest in contributing! This document provides guidelines for contributing to the project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/awesome-clickup-cli.git`
3. Create a branch: `git checkout -b feature/your-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Commit: `git commit -m "Add your feature"`
7. Push: `git push origin feature/your-feature`
8. Open a Pull Request

## Development Setup

```bash
# Install dependencies
go mod download

# Build
go build -o awesome-clickup-cli .

# Run tests
go test ./...

# Test with your API key
./awesome-clickup-cli auth set-token YOUR_TOKEN
./awesome-clickup-cli doctor
```

## Code Style

- Follow standard Go conventions (`gofmt`, `go vet`)
- Keep functions focused and small
- Add tests for new functionality
- Update documentation for user-facing changes

## Pull Request Guidelines

- Keep PRs focused on a single change
- Update README.md if adding user-facing features
- Ensure all tests pass
- Add a clear description of what the PR does

## Reporting Issues

When reporting bugs, include:
- CLI version (`awesome-clickup-cli version`)
- Go version (`go version`)
- OS and architecture
- Steps to reproduce
- Expected vs actual behavior

## Feature Requests

Feature requests are welcome! Please:
- Check existing issues first
- Describe the use case
- Explain why existing features don't solve it

## Security

If you discover a security vulnerability, please email the maintainer directly instead of opening a public issue.

## License

By contributing, you agree that your contributions will be licensed under the Apache-2.0 license.
