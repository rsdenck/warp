# Contributing to WarpCTL

Thank you for your interest in contributing to WarpCTL! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to maintain a respectful and collaborative environment.

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- A GitHub account

### Setting Up Development Environment

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/warp.git
   cd warp
   ```

3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/rsdenck/warp.git
   ```

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Build the project:
   ```bash
   go build -o warpctl
   ```

## Development Workflow

### Creating a Branch

Always create a new branch for your work:

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or modifications

### Making Changes

1. Make your changes in your feature branch
2. Follow the Go coding standards
3. Add tests for new functionality
4. Ensure all tests pass:
   ```bash
   go test ./...
   ```

5. Format your code:
   ```bash
   go fmt ./...
   ```

### Commit Guidelines

Write clear, concise commit messages:

```
type: brief description

Detailed explanation of what changed and why.

Fixes #issue_number
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Test additions or modifications
- `chore`: Maintenance tasks

### Submitting Changes

1. Push your changes to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

2. Create a Pull Request on GitHub
3. Fill out the PR template with relevant information
4. Wait for review and address any feedback

## Code Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `go vet` to catch common mistakes
- Keep functions small and focused
- Write descriptive variable and function names

### Testing

- Write unit tests for new functionality
- Maintain or improve code coverage
- Test edge cases and error conditions
- Use table-driven tests where appropriate

Example test structure:
```go
func TestFeature(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        // test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Documentation

- Update README.md for user-facing changes
- Add godoc comments for exported functions
- Include examples in documentation
- Update CHANGELOG.md with your changes

## Line Endings

This project enforces LF line endings for all text files:

- Configure Git to use LF:
  ```bash
  git config core.autocrlf false
  git config core.eol lf
  ```

- The `.gitattributes` file ensures consistent line endings
- CI/CD will fail if CRLF line endings are detected

## Pull Request Process

1. Ensure your PR addresses a single concern
2. Update documentation as needed
3. Add tests for new functionality
4. Ensure all CI checks pass
5. Request review from maintainers
6. Address review feedback promptly

### PR Checklist

- [ ] Code follows project style guidelines
- [ ] Tests added/updated and passing
- [ ] Documentation updated
- [ ] Commit messages are clear
- [ ] Branch is up to date with main
- [ ] No merge conflicts
- [ ] CI/CD checks passing

## Reporting Issues

### Bug Reports

Include:
- Clear description of the issue
- Steps to reproduce
- Expected vs actual behavior
- Environment details (OS, Go version, etc.)
- Relevant logs or error messages

### Feature Requests

Include:
- Clear description of the feature
- Use case and benefits
- Proposed implementation (if any)
- Potential drawbacks or concerns

## Questions?

If you have questions about contributing:
- Open an issue with the `question` label
- Check existing issues and discussions
- Review the documentation

## License

By contributing to WarpCTL, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing to WarpCTL!
