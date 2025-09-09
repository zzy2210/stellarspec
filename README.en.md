# StellarSpec

[English](README.en.md) | [简体中文](README.md)

StellarSpec is a local code review tool powered by LLMs, written in Go and built on top of the Eino framework. It detects Git changes in your repo and asks an AI model to review them, helping you spot potential issues and improve code quality.

## Features Overview

Available today
- LLM-based code review
- Concurrency with a fixed 10 workers
- Git integration: detect working tree and staged changes vs HEAD
- Language recognition for 20+ file types by extension
- i18n: Chinese/English prompts and report templates
- Config management: persist API server/model/key/language locally
- Markdown report: append per-file results into `code-review.md`

Planned / in progress
- Configurable concurrency: `--max-pool` (flag present, not wired; default 10)
- Review a specific commit: `--commit-id` (not wired)
- Thinking chain output: `--thinking-chain` (not wired)
- Custom prompt via file: `--prompt-file` (not wired)
- Scope filtering for single file/subdirectory (current flow reviews repo changes only)

## Quick Start

### Build & Install

```bash
git clone https://github.com/your-username/stellarspec.git
cd stellarspec

# Build (artifact at build/stellar)
make build

# Optional: install to PATH
sudo cp build/stellar /usr/local/bin/stellar
```

### Configure

```bash
# API server
stellar --set-apiserver https://api.siliconflow.cn/v1/

# Model name
stellar --set-model deepseek-chat

# API key
stellar --set-key sk-xxxxxxxxxxxxx

# Language (default: zh)
stellar --set-lang zh
stellar --set-lang en
```

Config is saved to `$HOME/.stellarspec/cnf`.

### Usage

```bash
# Review changes in current Git repo (against HEAD)
./build/stellar review

# If installed
stellar review

# Help (or make run)
stellar --help
```

Notes
- Currently reviews “repo changes” only; `review [path]` is treated as the repo root to open.
- Precise review for a single file/subdirectory is planned, not yet wired.
- Output report `code-review.md` is created/updated in the working directory.

### Config Options (placeholders)

Flags below are present in CLI but not wired into the engine yet:

```bash
# Concurrency (default 10)
stellar review . --max-pool 20

# Review a specific commit
stellar review --commit-id 1bacd3f

# Thinking chain output
stellar review --thinking-chain

# Custom prompt file
stellar review --prompt-file custom_prompt.txt
```

## Advanced

### Concurrency

Fixed at 10 for now. `--max-pool` will enable custom limits once wired.

### Thinking Chain

`--thinking-chain` is planned; it will surface more detailed reasoning (off by default).

### Custom Prompt

`--prompt-file` is planned to customize review focus and prompts.

## Report

Results are appended to `code-review.md` per file and include:

- File: path, detected language, timestamp
- Review result: the model’s conclusion text (template respects language)

Structured sections like “Issues/Suggestions/Fixes” are carried in model output; the template does not enforce them yet.

## Architecture

```
stellarspec/
├── cmd/                    # Cobra CLI entry
│   └── stellarspec.go
├── internal/
│   ├── model/
│   │   └── conf/          # INI config I/O
│   └── reviewer/          # diff collection / concurrency / reporting
├── build/                 # build artifacts (git-ignored)
├── go.mod
├── go.sum
└── README.md
```

## Contributing

PRs are welcome. Keep changes focused and small, and ensure `make build` succeeds.

## License

MIT

## Thanks

- Eino, Cobra, go-git, go-diff

## Support

- File an Issue or start a Discussion
