# issue-syncer

If you're anything like me, you have a habit of leaving TODOs, FIXMEs, and HACKs in your code. But keeping track of them can be a pain. That's where `issue-syncer` comes in!

`issue-syncer` is a tool that automatically synchronizes TODO comments in your code with GitHub issues. It scans your codebase for special comments (like `TODO`, `FIXME`, or `HACK`), and creates, updates, or closes GitHub issues to track them.

## Features

- Scans codebases for TODO, FIXME, and HACK comments
- Automatically creates GitHub issues for new comments
- Updates existing issues when comments change
- Closes issues when comments are removed
- Supports multiple programming languages
- Ignores specified directories (like node_modules, .git)
- Customizable search patterns and directories to skip

## Installation

```bash
go install github.com/alexdor/issue-syncer@latest
```

Or clone the repository and build manually:

```bash
git clone https://github.com/alexdor/issue-syncer.git
cd issue-syncer
go build
```

## Usage

### Add this to your workflow

```yaml
name: Sync TODOs
on:
  push:
    branches:
      - main

concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

permissions:
  issues: write
  contents: read

jobs:
  todo-sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Sync TODOs with Issues
        uses: alexdor/issue-syncer@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Command Line Interface

You can also run `issue-syncer` directly from the command line. Here’s how to use it:

```bash
issue-syncer [flags]
```

#### Required Environment Variables

For GitHub integration:

```bash
export GITHUB_TOKEN=your_github_token
export GITHUB_REPOSITORY=<username>/<repository>
```

#### Flags

- `-p, --path`: Path to the folder to scan (default: ".")
- `-w, --words`: Words to look for in comments (default: ["FIXME", "TODO", "HACK"])
- `-d, --dirs-to-skip`: Directories to skip (default: [".git", "node_modules", etc.])
- `-g, --use-gitignore`: Whether to use gitignore for skipping files (default: true)
- `-s, --storer`: Storer to use (default: "github")

#### Examples

Scan current directory with default settings:

```bash
issue-syncer
```

Scan a specific directory with custom comment markers:

```bash
issue-syncer --path ./src --words "TODO,FIXME,NOTE"
```

## Supported Languages

issue-syncer supports comments in many programming languages including:

- Go
- Python
- JavaScript/TypeScript
- Java
- C/C++
- Ruby
- PHP
- Rust
- HTML/CSS
- Shell scripts
- And many more!

## How It Works

1. The tool scans your codebase for specified comment patterns
2. For each comment found, it creates or updates a GitHub issue
3. The issue title is derived from the comment text
4. The issue body contains the file path, line number, and full comment text
5. Issues are tagged with "issue-syncer" and "auto-generated" labels
6. When a comment is removed from code, the corresponding issue is closed

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the [LICENSE](LICENSE) included in the repository.
