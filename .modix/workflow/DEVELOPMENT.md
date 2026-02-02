# Development Guidelines

## Code Standards

- **Go version**: 1.25+
- **Linting**: `golangci-lint run ./...`
- **Testing**: All new features need tests (`go test ./...`)
- **Documentation**: Update relevant .md files
- **CLI Standards**: Use Bubbletea framework for TUI

## File Organization

```bash
modix/
├── cmd/                  # CLI commands
├── internal/             # Internal packages
├── .claude/              # Claude Code configuration
├── .github/              # GitHub Actions CI/CD
├── .modix/               # Modix configuration
└── *.md                  # Documentation files
```

## Before Committing

1. **Run tests**: `go test ./...`
2. **Run linter**: `golangci-lint run ./...`
3. **Export issues**: `bd export -o .beads/issues.jsonl`
4. **Update docs**: If you changed behavior, update README.md or other docs
5. **Git add both**: `git add .beads/issues.jsonl <your-changes>`

## Git Workflow

```bash
# Make changes
git add <files>

# Commit
git commit -m "Your message"

# After pull
git pull
```

Or use the git hooks in `examples/git-hooks/` for automation.

## Building and Testing

```bash
# Build the project
go build -o modix ./cmd/modix

# Test all packages
go test ./...

# Test with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run locally
./modix init
./modix list
./modix status
```

## Release Process (Maintainers)

1. Update version in go.mod (if applicable)
2. Update CHANGELOG.md (if exists)
3. Run full test suite: `go test ./...`
4. Tag release: `git tag v0.x.0`
5. Push tag: `git push origin v0.x.0`
6. GitHub Actions handles CI/CD automatically
