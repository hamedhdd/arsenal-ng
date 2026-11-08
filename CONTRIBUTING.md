# Contributing to Arsenal-NG

First off, thank you for considering contributing to arsenal-ng! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [How Can I Contribute?](#how-can-i-contribute)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)
- [Adding Cheat Sheets](#adding-cheat-sheets)
- [Project Structure](#project-structure)

## Code of Conduct

This project follows a simple code of conduct:

- **Be respectful**: Treat everyone with respect and kindness
- **Be constructive**: Provide helpful feedback and suggestions
- **Be collaborative**: Work together to improve the project
- **Be professional**: Maintain a professional demeanor in all interactions

## How Can I Contribute?

### 1. Adding Cheat Sheets (Easiest!)

The most common and welcome contribution is adding new cheat sheets for pentesting tools.

**What we need:**
- Command templates for security tools
- Clear descriptions and use cases
- Proper categorization with tags
- Installation instructions (if available)

See [Adding Cheat Sheets](#adding-cheat-sheets) below for detailed instructions.

### 2. Reporting Bugs

Found a bug? Please open an issue with:

- **Title**: Clear, descriptive summary
- **Environment**: OS, Go version, terminal emulator
- **Steps to reproduce**: Detailed steps to trigger the bug
- **Expected behavior**: What you expected to happen
- **Actual behavior**: What actually happened
- **Screenshots**: If applicable

### 3. Suggesting Enhancements

Have an idea? Open an issue with:

- **Title**: Clear description of the feature
- **Use case**: Why is this needed?
- **Proposed solution**: How should it work?
- **Alternatives**: Other approaches you considered

### 4. Code Contributions

Want to contribute code? Great! See the [Development Workflow](#development-workflow) section.

## Getting Started

### Prerequisites

- **Go 1.22+**: Required for building
- **Git**: For version control
- **Make**: Optional but recommended
- **golangci-lint**: Optional, for linting

### Setup Development Environment

1. **Fork the repository**
   ```bash
   # Click "Fork" on GitHub, then:
   git clone https://github.com/YOUR_USERNAME/arsenal-ng.git
   cd arsenal-ng
   ```

2. **Add upstream remote**
   ```bash
   git remote add upstream https://github.com/hamedhdd/arsenal-ng.git
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Verify setup**
   ```bash
   make test
   make build
   ./bin/arsenal-ng
   ```

## Development Workflow

### 1. Create a Feature Branch

```bash
# Update your fork
git checkout main
git pull upstream main

# Create feature branch
git checkout -b feature/your-feature-name
```

**Branch naming conventions:**
- `feature/`: New features (e.g., `feature/add-history`)
- `fix/`: Bug fixes (e.g., `fix/search-crash`)
- `cheat/`: New cheat sheets (e.g., `cheat/add-nuclei`)
- `docs/`: Documentation updates (e.g., `docs/improve-readme`)
- `refactor/`: Code refactoring (e.g., `refactor/ui-helpers`)

### 2. Make Your Changes

- Write clean, readable code
- Follow Go conventions and idioms
- Add tests for new functionality
- Update documentation if needed
- Keep commits focused and atomic

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run specific package tests
go test -v ./internal/loader/...

# Check for issues
make vet

# Run linter (if installed)
make lint
```

### 4. Commit Your Changes

```bash
# Stage your changes
git add .

# Commit with a clear message
git commit -m "feat: add command history tracking"
```

**Commit message format:**
```
<type>: <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only
- `style`: Code style changes (formatting, no logic change)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Example:**
```
feat: add command history tracking

- Add history.json storage in config directory
- Implement history special command
- Track last 100 executed commands with timestamps
- Add tests for history persistence

Closes #42
```

### 5. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/your-feature-name
```

Then open a Pull Request on GitHub.

## Coding Standards

### Go Code Style

Follow standard Go conventions:

```go
// ✅ Good: Clear function names, documented
// GetConfigDir returns the application configuration directory.
// It respects XDG_CONFIG_HOME on Linux/macOS and falls back to
// the user's home directory if not set.
func GetConfigDir() (string, error) {
    // Implementation
}

// ❌ Bad: No documentation, unclear naming
func gcd() (string, error) {
    // Implementation
}
```

### Code Organization

- **Keep functions small**: Max 50 lines per function
- **Single responsibility**: One function = one task
- **Use constants**: Avoid magic numbers
- **Error handling**: Always handle errors explicitly
- **Documentation**: Add godoc comments for exported items

### File Organization

```go
// Package declaration with description
package mypackage

// Imports grouped: stdlib, external, internal
import (
    "fmt"
    "strings"
    
    "github.com/external/package"
    
    "github.com/hamedhdd/arsenal-ng/internal/model"
)

// Constants
const (
    MaxItems = 100
    DefaultTimeout = 30
)

// Types
type MyStruct struct {
    // fields
}

// Functions
func NewMyStruct() *MyStruct {
    // implementation
}
```

## Testing Guidelines

### Writing Tests

```go
func TestFunctionName_Scenario(t *testing.T) {
    // Arrange
    input := "test input"
    expected := "expected output"
    
    // Act
    result := FunctionName(input)
    
    // Assert
    if result != expected {
        t.Errorf("FunctionName(%q) = %q, expected %q", input, result, expected)
    }
}
```

### Test Coverage Goals

- **Critical packages**: >80% coverage (model, loader, state)
- **Business logic**: >70% coverage (toolmgr, config)
- **UI helpers**: >60% coverage (ui/helpers_*)
- **UI components**: Best effort (TUI testing is complex)

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Specific package
go test -v ./internal/model/

# Specific test
go test -v -run TestParsArguments ./internal/model/
```

## Pull Request Process

### Before Submitting

- [ ] Code follows project style guidelines
- [ ] All tests pass: `make test`
- [ ] No linting errors: `make vet`
- [ ] Added tests for new functionality
- [ ] Updated documentation if needed
- [ ] Commit messages are clear and descriptive

### PR Description Template

```markdown
## Description
Brief description of what this PR does.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update
- [ ] New cheat sheet(s)

## Testing
Describe how you tested your changes:
- [ ] Unit tests added/updated
- [ ] Manual testing performed
- [ ] Tested on: Linux / macOS / Windows

## Checklist
- [ ] My code follows the code style of this project
- [ ] I have added tests that prove my fix/feature works
- [ ] All new and existing tests pass
- [ ] I have updated the documentation accordingly
```

### Review Process

1. **Automated Checks**: CI runs tests and linting
2. **Code Review**: Maintainer reviews your code
3. **Feedback**: Address any requested changes
4. **Approval**: Once approved, PR will be merged
5. **Merge**: Squash and merge to main branch

## Adding Cheat Sheets

This is the easiest and most valuable contribution!

### Step-by-Step Guide

1. **Create a new YAML file**
   ```bash
   cd internal/loader/cheat-files/
   touch mytool.yaml
   ```

2. **Follow the template**
   ```yaml
   tool: mytool
   tags: [category1, category2, category3]
   
   install:
     apt: mytool
     brew: mytool
     pip: mytool
     go: github.com/user/mytool@latest
     git: git clone https://github.com/user/mytool
   
   actions:
     - title: mytool - brief description
       desc: Longer description of what this command does
       command: "mytool --flag {{argument}}"
     
     - title: mytool - another command
       desc: Description of the second command
       command: "mytool --option {{arg1}} --target {{arg2|default_value}}"
   ```

3. **Choose appropriate tags**
   
   Common tags:
   - **recon**: Reconnaissance and information gathering
   - **scan**: Port/vulnerability scanning
   - **web**: Web application testing
   - **exploit**: Exploitation tools
   - **post**: Post-exploitation
   - **enum**: Enumeration
   - **network**: Network tools
   - **wireless**: Wireless testing
   - **password**: Password attacks
   - **forensics**: Digital forensics
   - **crypto**: Cryptography
   - **stealth**: Stealth and evasion
   - **mobile**: Mobile security
   - **cloud**: Cloud security
   - **api**: API testing

4. **Test locally**
   ```bash
   make build
   ./bin/arsenal-ng
   # Search for your tool and test all commands
   ```

5. **Submit PR**
   ```bash
   git add internal/loader/cheat-files/mytool.yaml
   git commit -m "cheat: add mytool command templates"
   git push origin cheat/add-mytool
   ```

### Cheat Sheet Best Practices

**Do:**
- ✅ Use clear, descriptive titles
- ✅ Add helpful descriptions for complex commands
- ✅ Include default values for optional arguments
- ✅ Group related commands in one file
- ✅ Add installation instructions when available
- ✅ Use common argument names (ip, port, user, pass, etc.)

**Don't:**
- ❌ Add destructive commands without clear warnings
- ❌ Include credentials or sensitive data
- ❌ Use overly complex multi-line commands
- ❌ Duplicate existing cheat sheets
- ❌ Add non-security-related tools

### Example: High-Quality Cheat Sheet

```yaml
tool: ffuf
tags: [web, fuzzing, recon]

install:
  go: github.com/ffuf/ffuf@latest
  apt: ffuf
  brew: ffuf

actions:
  - title: ffuf - directory fuzzing
    desc: Discover directories using a wordlist
    command: "ffuf -w {{wordlist|/usr/share/wordlists/dirb/common.txt}} -u {{url}}/FUZZ"

  - title: ffuf - subdomain fuzzing
    desc: Enumerate subdomains using DNS
    command: "ffuf -w {{wordlist}} -u {{url}} -H 'Host: FUZZ.{{domain}}'"

  - title: ffuf - POST data fuzzing
    desc: Fuzz POST parameters
    command: "ffuf -w {{wordlist}} -X POST -d 'username=FUZZ&password={{pass}}' -u {{url}}"

  - title: ffuf - filter by status code
    desc: Only show results with specific status codes
    command: "ffuf -w {{wordlist}} -u {{url}}/FUZZ -fc {{codes|404,403}}"
```

## Project Structure

Understanding the codebase:

```
arsenal-ng/
├── cmd/
│   └── arsenal-ng/        # Main application entry point
│       └── main.go
├── internal/
│   ├── config/           # Configuration and paths
│   ├── loader/           # Cheat file loading and search
│   │   └── cheat-files/  # YAML cheat sheets
│   ├── model/            # Data structures
│   ├── output/           # Command output handling
│   ├── state/            # Variable persistence
│   ├── toolmgr/          # Tool installation
│   └── ui/               # Terminal UI
├── docs/
│   └── ARCHITECTURE.md   # Architecture documentation
├── .github/
│   └── workflows/        # CI/CD pipelines
├── Makefile             # Build automation
├── go.mod               # Go dependencies
└── README.md            # User documentation
```

## Getting Help

### Resources

- **Documentation**: Check `docs/ARCHITECTURE.md` for technical details
- **Issues**: Search existing issues for similar problems
- **Discussions**: Use GitHub Discussions for questions

### Communication

- **GitHub Issues**: Bug reports and feature requests
- **GitHub Discussions**: Questions and general discussion
- **Pull Requests**: Code contributions and reviews

## Recognition

Contributors are recognized in:

- GitHub contributors page
- Release notes (for significant contributions)
- README acknowledgments (for major features)

---

## Quick Reference

### Common Commands

```bash
# Development
make build         # Build binary
make run          # Run directly
make test         # Run all tests
make vet          # Check for issues
make lint         # Run linter

# Testing
go test ./...                    # All tests
go test -v ./internal/model/     # Specific package
go test -coverprofile=cover.out  # With coverage

# Formatting
gofmt -s -w .     # Format all files
go mod tidy       # Clean dependencies
```

### Need Help?

- 📖 Read the [Architecture Documentation](docs/ARCHITECTURE.md)
- 🐛 [Open an Issue](https://github.com/hamedhdd/arsenal-ng/issues/new)
- 💬 Start a [Discussion](https://github.com/hamedhdd/arsenal-ng/discussions)

---

**Thank you for contributing to arsenal-ng!** 🎉

Your contributions help security professionals worldwide work more efficiently.
