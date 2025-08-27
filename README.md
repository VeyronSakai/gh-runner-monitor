# gh-runner-monitor

A GitHub CLI extension that provides real-time monitoring of GitHub Actions self-hosted runners with a Terminal User Interface (TUI).

## Features

- 🔄 Real-time monitoring of self-hosted runners
- 📊 Display runner status (Idle, Active, Offline) with color coding
- 💼 Show currently executing jobs with execution time
- 🏢 Support for both repository and organization level monitoring
- ⌨️ Interactive TUI with keyboard navigation

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

## Architecture

This project follows Onion Architecture with the following layers:

- **Domain Layer**: Core entities and business rules
- **Use Cases Layer**: Application business logic
- **Infrastructure Layer**: External implementations (GitHub API)
- **Presentation Layer**: TUI components

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

### Running tests

```bash
go test ./...
```

## License

MIT