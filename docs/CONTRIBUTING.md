# Contributing to Ham Radio Cloud

Thank you for your interest in contributing to Ham Radio Cloud! This document provides guidelines and instructions for contributing to the project.

---

## Table of Contents

1. [Code of Conduct](#code-of-conduct)
2. [Getting Started](#getting-started)
3. [Development Workflow](#development-workflow)
4. [Coding Standards](#coding-standards)
5. [Commit Guidelines](#commit-guidelines)
6. [Pull Request Process](#pull-request-process)
7. [Testing](#testing)
8. [Documentation](#documentation)

---

## Code of Conduct

We are committed to providing a welcoming and inclusive environment for all contributors. Please:

- Be respectful and considerate
- Focus on constructive feedback
- Assume good intentions
- Help create a positive community

Remember: **73 - Best Regards** 📻

---

## Getting Started

### Prerequisites

- Go 1.23+
- Node.js 20+
- Docker & Docker Compose
- Git

### Fork and Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/ham-radio-cloud.git
cd ham-radio-cloud

# Add upstream remote
git remote add upstream https://github.com/nitefawkes/ham-radio-cloud.git
```

### Set Up Development Environment

```bash
# Install dependencies
make install

# Copy environment file
cp .env.example .env

# Start development environment
make dev
```

See [GETTING_STARTED.md](./GETTING_STARTED.md) for detailed setup instructions.

---

## Development Workflow

### 1. Create a Branch

```bash
# Fetch latest changes
git fetch upstream
git checkout main
git merge upstream/main

# Create a feature branch
git checkout -b feature/your-feature-name
```

### Branch Naming Convention

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions/updates
- `chore/` - Maintenance tasks

Examples:
- `feature/adif-import`
- `fix/lotw-sync-timeout`
- `docs/api-endpoints`

### 2. Make Changes

- Write clean, readable code
- Follow coding standards (see below)
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Backend tests
cd backend && go test -v ./...

# Frontend type check
cd frontend && npm run check

# Lint
make lint
```

### 4. Commit

```bash
git add .
git commit -m "feat: add ADIF import functionality"
```

See [Commit Guidelines](#commit-guidelines) below.

### 5. Push and Create Pull Request

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

---

## Coding Standards

### Backend (Go)

- **Formatting:** Use `gofmt` and `gofmt -s`
- **Linting:** Code must pass `go vet`
- **Naming:** Follow Go naming conventions (camelCase for private, PascalCase for public)
- **Error Handling:** Always handle errors explicitly
- **Comments:** Document exported functions and types

```go
// Good
func CreateQSO(qso *models.QSO) error {
    if qso == nil {
        return fmt.Errorf("qso cannot be nil")
    }
    // ...
}

// Bad
func createqso(q *models.QSO) error {
    // no validation
    // no documentation
}
```

### Frontend (TypeScript/Svelte)

- **Formatting:** Use Prettier (configured)
- **Linting:** Code must pass ESLint
- **Types:** Use TypeScript types, avoid `any`
- **Components:** Keep components focused and composable
- **Naming:** PascalCase for components, camelCase for functions/variables

```typescript
// Good
export interface QSO {
	id: string;
	callsign: string;
	frequency: number;
}

// Bad
export interface qso {
	id: any;
	callsign: any;
}
```

### Database

- **Migrations:** Use versioned migration files
- **Indexes:** Add indexes for frequently queried columns
- **Naming:** Use snake_case for tables and columns

---

## Commit Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/).

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation only
- `style:` - Code style (formatting, missing semi-colons, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance tasks

### Examples

```bash
feat(qso): add ADIF import functionality

Implement ADIF file parser and bulk QSO import endpoint.
Supports ADIF 3.1.0 format.

Closes #123

---

fix(auth): resolve OAuth token expiration issue

JWT tokens were not refreshing correctly. Updated token
validation middleware to check expiration properly.

Fixes #456

---

docs(api): update propagation endpoint documentation

Add examples for band conditions endpoint.
```

---

## Pull Request Process

### Before Submitting

1. ✅ Tests pass (`make test`)
2. ✅ Linting passes (`make lint`)
3. ✅ Code is formatted (`make format`)
4. ✅ Documentation updated (if applicable)
5. ✅ Branch is up to date with `main`

### PR Template

When creating a PR, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
How was this tested?

## Screenshots (if applicable)

## Checklist
- [ ] Tests pass
- [ ] Linting passes
- [ ] Documentation updated
- [ ] Follows coding standards
```

### Review Process

- At least one maintainer approval required
- CI must pass
- No merge conflicts
- Discussion and requested changes addressed

---

## Testing

### Backend Tests

```bash
cd backend
go test -v ./...

# With coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Frontend Tests

```bash
cd frontend
npm run check     # Type check
npm run test      # Unit tests (when implemented)
```

### Integration Tests

```bash
# Start test environment
docker-compose -f deployments/docker/docker-compose.yml up -d

# Run integration tests
# (To be implemented in Phase 6)
```

---

## Documentation

### When to Update Documentation

- Adding new features
- Changing API endpoints
- Modifying configuration
- Updating dependencies

### Documentation Files

- `README.md` - Project overview
- `docs/GETTING_STARTED.md` - Setup instructions
- `docs/API.md` - API documentation
- `docs/ARCHITECTURE.md` - Architecture overview
- `docs/PROJECT_STATUS.md` - Project status
- Code comments - Inline documentation

### API Documentation

Update `docs/API.md` when:
- Adding new endpoints
- Changing request/response formats
- Modifying authentication

---

## Issue Reporting

### Bug Reports

Include:
- Description of the bug
- Steps to reproduce
- Expected behavior
- Actual behavior
- Environment (OS, Go version, Node version)
- Logs/screenshots

### Feature Requests

Include:
- Use case
- Proposed solution
- Alternative solutions considered
- Additional context

---

## Community

- **GitHub Issues:** Bug reports and feature requests
- **GitHub Discussions:** General questions and discussions
- **Pull Requests:** Code contributions

---

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.

---

## Questions?

- Check existing issues and discussions
- Read the documentation
- Ask in GitHub Discussions

---

**Thank you for contributing to Ham Radio Cloud!**

*73 de W1AW* 📻
