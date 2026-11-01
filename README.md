# gh-runner-monitor

A GitHub CLI extension that provides real-time monitoring of GitHub Actions self-hosted runners with a Terminal User Interface (TUI).

## Features

- 🔄 Real-time monitoring of self-hosted runners
- 📊 Display runner status (Idle, Active, Offline) with color coding
- 💼 Show currently executing jobs with execution time
- 🏢 Support for both repository and organization level monitoring
- ⌨️ Interactive TUI with keyboard navigation

<img width="904" height="195" alt="スクリーンショット 2025-11-03 16 14 13" src="https://github.com/user-attachments/assets/4d45ea0c-3374-4d16-a264-d478fdee290b" />

## Installation

```bash
gh extension install VeyronSakai/gh-runner-monitor
```

## Usage

### Monitor current repository
```bash
gh runner-monitor
```

### Monitor specific repository
```bash
gh runner-monitor --repo owner/repo
```

### Monitor organization
```bash
gh runner-monitor --org organization-name
```

### Monitor organization, limiting job lookups to specific repositories
```bash
gh runner-monitor --org organization-name --repos repo1,repo2
```

GitHub's REST API has no organization-wide endpoint for workflow runs, so in `--org` mode the jobs running on each runner are found by querying every repository in the organization. For organizations with many repositories, use `--repos` to reduce API calls (all organization runners are still shown).

### Custom update interval
```bash
gh runner-monitor --interval 10  # Update every 10 seconds
```

## Status Colors

- 🟢 **Green** - Idle: Runner is online and available
- 🟠 **Orange** - Active: Runner is executing a job
- ⚫ **Gray** - Offline: Runner is not connected

## Keyboard Shortcuts

- `↑/↓` or `j/k` - Navigate through runners
- `r` - Manual refresh
- `q` or `Ctrl+C` - Quit

## Development

### Prerequisites

- Go 1.20 or higher
- GitHub CLI (`gh`) installed and authenticated

### Building from source

```bash
git clone https://github.com/VeyronSakai/gh-runner-monitor.git
cd gh-runner-monitor
go build -o gh-runner-monitor
```

### Testing Locally

#### 1. Build and run directly
```bash
# Build the binary
go build -o gh-runner-monitor

# Run with help flag to see options
./gh-runner-monitor --help

# Monitor current repository
./gh-runner-monitor

# Monitor specific repository
./gh-runner-monitor --repo owner/repo
```

#### 2. Install as gh extension from local directory
```bash
# Install from current directory
gh extension install .

# Run as gh extension
gh runner-monitor

# Uninstall when done testing
gh extension remove runner-monitor
```

#### 3. Test with different configurations
```bash
# Monitor a public repository with runners
gh runner-monitor --repo actions/runner

# Monitor with custom refresh interval (10 seconds)
gh runner-monitor --interval 10

# Monitor organization (requires org access)
gh runner-monitor --org your-org-name
```

### Running tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

## License

MIT
