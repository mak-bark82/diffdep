# diffdep

> Compares dependency trees across branches to surface breaking version changes

---

## Installation

```bash
go install github.com/yourusername/diffdep@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/diffdep.git && cd diffdep && go build ./...
```

---

## Usage

Compare dependencies between your current branch and `main`:

```bash
diffdep --base main --head feature/my-branch
```

Compare two specific branches and output results as JSON:

```bash
diffdep --base v1.2.0 --head v2.0.0 --format json
```

**Example output:**

```
BREAKING  github.com/some/pkg       v1.4.2  →  v2.0.0
UPGRADED  github.com/another/lib    v0.8.0  →  v0.9.1
REMOVED   github.com/old/dep        v1.0.0  →  (removed)
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--base` | `main` | Base branch or tag to compare from |
| `--head` | `HEAD` | Target branch or tag to compare to |
| `--format` | `text` | Output format: `text` or `json` |
| `--breaking-only` | `false` | Show only breaking (major) version changes |

---

## Requirements

- Go 1.21+
- Git available in `$PATH`

---

## License

MIT © [yourusername](https://github.com/yourusername)