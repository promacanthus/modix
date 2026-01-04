# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

See [AGENTS.md](AGENTS.md) for workflow.

---

## Project Overview

This is **Modix** (command: `modix`), a Go-based CLI tool for managing and switching between multiple LLM backends and large language models.

## Issue Tracking

We use **bd** (beads) for issue tracking instead of Markdown TODOs or external tools.

### Quick Reference

```bash
# Find ready work (no blockers)
bd ready --json

# Find ready work including future deferred issues
bd ready --include-deferred --json

# Create new issue
bd create "Issue title" -t bug|feature|task -p 0-4 -d "Description" --json

# Create issue with due date and defer (GH#820)
bd create "Task" --due=+6h              # Due in 6 hours
bd create "Task" --defer=tomorrow       # Hidden from bd ready until tomorrow
bd create "Task" --due="next monday" --defer=+1h  # Both

# Update issue status
bd update <id> --status in_progress --json

# Update issue with due/defer dates
bd update <id> --due=+2d                # Set due date
bd update <id> --defer=""               # Clear defer (show immediately)

# Link discovered work
bd dep add <discovered-id> <parent-id> --type discovered-from

# Complete work
bd close <id> --reason "Done" --json

# Show dependency tree
bd dep tree <id>

# Get issue details
bd show <id> --json

# Query issues by time-based scheduling (GH#820)
bd list --deferred              # Show issues with defer_until set
bd list --defer-before=tomorrow # Deferred before tomorrow
bd list --defer-after=+1w       # Deferred after one week from now
bd list --due-before=+2d        # Due within 2 days
bd list --due-after="next monday" # Due after next Monday
bd list --overdue               # Due date in past (not closed)
```

### Workflow

1. **Check for ready work**: Run `bd ready` to see what's unblocked
2. **Claim your task**: `bd update <id> --status in_progress`
3. **Work on it**: Implement, test, document
4. **Discover new work**: If you find bugs or TODOs, create issues:
   - `bd create "Found bug in auth" -t bug -p 1 --json`
   - Link it: `bd dep add <new-id> <current-id> --type discovered-from`
5. **Complete**: `bd close <id> --reason "Implemented"`
6. **Export**: Run `bd export -o .beads/issues.jsonl` before committing

### Issue Types

- `bug` - Something broken that needs fixing
- `feature` - New functionality
- `task` - Work item (tests, docs, refactoring)
- `epic` - Large feature composed of multiple issues
- `chore` - Maintenance work (dependencies, tooling)

### Priorities

- `0` - Critical (security, data loss, broken builds)
- `1` - High (major features, important bugs)
- `2` - Medium (nice-to-have features, minor bugs)
- `3` - Low (polish, optimization)
- `4` - Backlog (future ideas)

### Dependency Types

- `blocks` - Hard dependency (issue X blocks issue Y)
- `related` - Soft relationship (issues are connected)
- `parent-child` - Epic/subtask relationship
- `discovered-from` - Track issues discovered during work

Only `blocks` dependencies affect the ready work queue.

---

## Development Guidelines

### Code Standards

- **Go version**: 1.25+
- **Linting**: `golangci-lint run ./...`
- **Testing**: All new features need tests (`go test ./...`)
- **Documentation**: Update relevant .md files
- **CLI Standards**: Use Cobra framework for commands, Viper for configuration

### File Organization

```bash
modix/
├── cmd/modix/              # CLI commands and main entry point
│   ├── main.go            # Program entry point
│   └── commands/          # Individual command implementations
│       ├── add.go         # Add model command
│       ├── check.go       # Check configuration command
│       ├── init.go        # Initialize command
│       ├── list.go        # List models command
│       ├── remove.go      # Remove model command
│       ├── root.go        # Root command
│       ├── show.go        # Show vendor details command
│       ├── status.go      # Status command
│       ├── switch.go      # Switch model command
│       └── utils.go       # Utility functions
├── internal/               # Internal packages
│   └── config/           # Configuration management
│       ├── claude.go     # Claude-specific configuration
│       ├── config.go     # Core configuration structures
│       └── default.go    # Default configurations
├── .beads/               # Beads issue tracking system
├── .claude/              # Claude Code configuration
├── .github/              # GitHub Actions CI/CD
└── *.md                  # Documentation files
```

### Before Committing

1. **Run tests**: `go test ./...`
2. **Run linter**: `golangci-lint run ./...`
3. **Export issues**: `bd export -o .beads/issues.jsonl`
4. **Update docs**: If you changed behavior, update README.md or other docs
5. **Git add both**: `git add .beads/issues.jsonl <your-changes>`

### Git Workflow

```bash
# Make changes
git add <files>

# Export beads issues
bd export -o .beads/issues.jsonl
git add .beads/issues.jsonl

# Commit
git commit -m "Your message"

# After pull
git pull
bd import -i .beads/issues.jsonl  # Sync SQLite cache
```

Or use the git hooks in `examples/git-hooks/` for automation.

## Current Project Status

Run `bd stats` to see overall progress.

### Active Areas

- **Core CLI**: Go CLI tool with Cobra framework implemented
- **Configuration Management**: Complete configuration system with vendor-based organization
- **Multi-Provider Support**: Support for 8+ LLM providers (Anthropic, DeepSeek, Alibaba, ByteDance, Moonshot AI, Kuaishou, MiniMax, ZHIPU AI)
- **CLI Commands**: Core commands implemented (init, list, check, add, remove, show, status, switch, path)
- **Enhanced Features**: Color-coded output, special Anthropic handling, health checks
- **CI/CD**: Automated builds and releases with GitHub Actions
- **Documentation**: Comprehensive CLI documentation and usage examples

### Project Goals

- **Stable CLI**: Complete CLI implementation with all planned commands
- **API Integration**: Integrate with various LLM provider APIs
- **GUI Interface**: Future GUI configuration interface
- **Performance**: Model performance benchmarking and optimization
- **IDE Integration**: Integration with popular IDEs and editors

## Common Tasks

### Adding a New CLI Command

1. Create command file in `cmd/modix/commands/`
2. Implement the command using Cobra framework
3. Add command to root command in `cmd/modix/commands/root.go`
4. Add `--json` flag for programmatic use
5. Add tests in `cmd/modix/commands/*_test.go`
6. Update CLI help and documentation

### Adding New LLM Provider

1. Update configuration schema in `internal/config/config.go`
2. Add vendor definition in `internal/config/default.go`
3. Add model definitions if needed
4. Update CLI commands to handle new vendor
5. Add tests for new functionality
6. Update README.md with new provider information

### Configuration Management

1. Update configuration structures in `internal/config/config.go`
2. Add validation logic
3. Update CLI commands that interact with configuration
4. Add tests for configuration handling
5. Update documentation

### Adding Examples

1. Create directory in `examples/` (if exists)
2. Add README.md explaining the example
3. Include working code
4. Link from main README.md if applicable
5. Mention in relevant documentation

## Questions?

- Check existing issues: `bd list`
- Look at recent commits: `git log --oneline -20`
- Read the docs: README.md, AGENTS.md
- Check CLI help: `modix --help` or `modix <command> --help`
- Create an issue if unsure: `bd create "Question: ..." -t task -p 2`

## Important Files

- **README.md** - Main project documentation (keep this updated!)
- **AGENTS.md** - AI agent usage guidelines
- **.beads/issues.jsonl** - Current issue tracking data
- **go.mod/go.sum** - Go module and dependency definitions
- **cmd/modix/main.go** - CLI entry point
- **internal/config/** - Configuration management system

## Pro Tips for Agents

- Always use `--json` flags for programmatic use of CLI commands
- Link discoveries with `discovered-from` to maintain context
- Check `bd ready` before asking "what next?"
- Export to JSONL before committing (or use git hooks)
- Use `bd dep tree` to understand complex dependencies
- Priority 0-1 issues are usually more important than 2-4
- Modix commands support color-coded output for better UX
- Anthropic models are pre-configured in Claude Code (special handling)
- Use `modix check claude-code` to validate Claude Code integration

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

---

**Remember**: If you find the workflow confusing or have ideas for improvement, create an issue with your feedback.

Happy coding!
