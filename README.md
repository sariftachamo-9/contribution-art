# Contribution Art (sarif)

> A command-line tool developed by **Sarif Tachamo** that generates a pattern of Git commits across your repository to create visual text and artwork on your GitHub contribution heatmap graph.

---

## Description

**Contribution Art** transforms any given text string into a 7-row bitmap grid representing days of the week (Sunday to Saturday). It then generates backdated empty Git commits with precise timestamps (`GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE`), allowing developers to personalize their GitHub contribution graph with custom messages like **"SARIF"**.

The tool includes a built-in **dry-run mode** with colored terminal previews so you can inspect your heatmap before creating any commits.

---

## Screenshots / Demo

### Terminal Heatmap Preview for "SARIF" (ANSI Colored)

```text
Starting from: 2026-08-03

Heatmap preview (7 rows = Sun–Sat):

Sun   ████    ████   ████   █████  █████  
Mon  ██  ██  ██  ██  ██  ██   ██   ██     
Tue  ██      ██  ██  ██  ██   ██   ██     
Wed   ████   ██████  ████     ██   ████   
Thu      ██  ██  ██  ██ ██    ██   ██     
Fri  ██  ██  ██  ██  ██  ██   ██   ██     
Sat   ████   ██  ██  ██  ██ █████  ██     

[DRY RUN] No commits were created.
```

---

## Features

- **Text-to-Bitmap Engine**: Converts text strings into 7-row (Sunday to Saturday) bitmaps using a 7x5 font map.
- **Full Uppercase Font Support**: Out-of-the-box support for letters `A` through `Z` and spaces.
- **Configurable Intensity**: Adjust commits per day (1–10) to control contribution color depth.
- **Flexible Start Date**: Begins on the next Monday by default, or accepts custom dates (`YYYY-MM-DD`).
- **Dry-Run Preview**: Displays an ANSI colored preview in your terminal without modifying repository history.
- **Git Auto-Initialization**: Automatically initializes a `.git` repository if one does not exist.
- **Convenient Run Helper**: Includes `run.bat` launcher script and standalone binary compilation.
- **Zero External Dependencies**: Built strictly with standard Go libraries.

---

## Tech Stack

- **Language**: Go 1.21+ / Go 1.22+
- **VCS**: Git
- **Libraries**: Go Standard Library (`flag`, `fmt`, `os`, `os/exec`, `strings`, `time`)

---

## Project Architecture

```text
+------------------+     +------------------+     +------------------+
|    CLI Parser    | --> |  Grid Generator  | --> |  Preview Engine  |
|  (flag package)  |     |   (7x5 font map) |     |  (ANSI colours)  |
+------------------+     +------------------+     +------------------+
          |                       |                        |
          v                       v                        v
+------------------+     +------------------+     +------------------+
|  Date Calculator |     |   Git Committer  |     |  Output (stdout) |
|   (start date)   |     |  (exec.Command)  |     |    (terminal)    |
+------------------+     +------------------+     +------------------+
```

### Core Components
1. **CLI Parser**: Validates input flags (`-pattern`, `-start`, `-commits`, `-dry-run`).
2. **Date Calculator**: Computes start date (UTC) and aligns timeline to week boundaries.
3. **Grid Generator**: Maps string characters into a 7×N boolean grid with spacer columns.
4. **Preview Engine**: Renders ANSI background colors (`\x1b[42m` green, `\x1b[40m` dark) to terminal stdout.
5. **Git Committer**: Executes `git commit --allow-empty` with backdated environment variables.

---

## Installation

### Prerequisites
- [Go 1.21+](https://go.dev/dl/) installed (or portable Go SDK).
- [Git](https://git-scm.com/) installed and configured in your system `PATH`.

### Clone & Build
```bash
# Clone the repository
git clone https://github.com/sariftachamo-9/contribution-art.git
cd contribution-art

# Build the executable binary (sarif.exe)
go build -o sarif.exe main.go
```

---

## Configuration

| Flag | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `-pattern` | `string` | `"SARIF"` | Text string to render on contribution heatmap. |
| `-start` | `string` | `""` (next Monday) | Start date in `YYYY-MM-DD` format. |
| `-commits` | `int` | `3` | Number of commits per filled cell (range 1–10). |
| `-dry-run` | `bool` | `false` | Enable preview mode without creating any Git commits. |

---

## Usage

### 1. Run using Helper Launcher Script (`run.bat`)
```bash
# Preview mode (dry-run)
.\run.bat -pattern "SARIF" -dry-run

# Generate commits
.\run.bat -pattern "SARIF" -commits 5
```

### 2. Run using Standalone Binary (`sarif.exe`)
```bash
# Build the binary
go build -o sarif.exe main.go

# Execute directly
.\sarif.exe -pattern "SARIF" -dry-run
```

### 3. Run directly with Go CLI
```bash
# Preview mode
go run main.go -pattern "SARIF" -dry-run

# Custom pattern & start date
go run main.go -pattern "SARIF" -start 2025-06-01 -commits 4
```

### 4. Push Commits to GitHub
```bash
git remote add origin https://github.com/sariftachamo-9/contribution-art.git
git branch -M main
git push -u origin main
```

---

## Folder Structure

```text
contribution-art/
├── main.go            # CLI logic, font map, grid builder, preview & committer
├── main_test.go       # Unit and integration test suite
├── go.mod             # Go module definition
├── Makefile           # Build automation tasks (build, test, fmt, clean)
├── run.bat            # Helper launcher script for quick execution
├── .gitignore         # Git ignore rules for compiled binaries & installers
├── README.md          # Comprehensive user & developer documentation
├── prd.md             # Product Requirements Document
├── architecture.md    # High-Level Architecture specifications
├── design.md          # Detailed Technical Design specifications
├── rules.md           # Coding standards and development guidelines
├── phases.md          # Project development roadmap
└── memory.md          # Architecture decisions & knowledge base
```

---

## Security Features

- **Safe Dry-Run Guarantee**: Running with `-dry-run` guarantees zero changes to your filesystem or Git repository.
- **Isolated Empty Commits**: Uses `git commit --allow-empty` without modifying or deleting any working directory files.
- **Strict Timestamping**: Sets `GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE` per commit explicitly, avoiding environment leaks.
- **Local Operation Only**: Performs local Git actions only; does not execute auto-push commands without developer authorization.

---

## Testing

Execute unit tests to verify grid generation, date calculations, and Git integration:

```bash
# Run tests using Go toolchain
go test -v ./...

# Or use Makefile target
make test
```

---

## Future Improvements

- [ ] Support for lowercase letters and numeric digits (`0-9`).
- [ ] Integration with native Go Git library (`go-git`) to remove external binary invocation.
- [ ] Support loading custom font configurations from external JSON files.
- [ ] Multi-line text layout engine for long string patterns.
- [ ] Optional `--push` flag to push commits to remote repository automatically.

---

## License

Distributed under the MIT License. See `LICENSE` for more information.

---

## Author / Contact

- **Author**: Sarif Tachamo
- **Repository**: [contribution-art](https://github.com/sariftachamo-9/contribution-art.git)
