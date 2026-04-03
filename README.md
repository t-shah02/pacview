# pacview

Terminal UI for browsing **installed** Arch Linux packages. It runs `pacman -Qi`, parses the output, and shows a scrollable table with search, dependency context, and keyboard-driven navigation ([Bubble Tea](https://github.com/charmbracelet/bubbletea) + [bubbles](https://github.com/charmbracelet/bubbles)).

## Requirements

- **pacman** on `PATH` (Arch Linux or another environment where `pacman -Qi` works)
- **Go 1.26+** (see `go.mod`) to build from source
- A terminal with reasonable size (the UI resizes with the window)

## Build and run

```bash
make build    # writes bin/pacview
./bin/pacview
```

Or without Make:

```bash
go build -o bin/pacview .
./bin/pacview
```

Run from source:

```bash
make run
# or: go run .
```

Tests:

```bash
make test
```

## Local install script

`scripts/install-local.sh` builds with `make build`, copies the binary to `~/.local/bin/pacview`, and appends a managed `alias pacview=...` block to your shell profile (chosen from `$SHELL`: zsh → `~/.zshrc`, bash → `~/.bashrc` or `~/.bash_profile`, else `~/.profile`). Safe to run again; it skips the alias block if already present.

```bash
./scripts/install-local.sh
source ~/.zshrc   # or the profile path the script prints
```

## Using the UI

| Key | Action |
|-----|--------|
| `↑` / `↓` (also `j` / `k` in the table) | Move selection |
| `/` | Focus the search bar |
| `esc` | Leave the search bar (focus returns to the table) |
| `f` | Narrow the table to packages listed in **Required by** for the highlighted row (resolved against the full install set) |
| `b` | Pop that filter and return to the previous scope |
| `q` or `ctrl+c` | Quit |

Search matches case-insensitively against name, description, version, install date, depends, and required-by text. The footer shows how many rows are visible and whether you are in a narrowed scope.

## How it works

1. `internal/utils` runs `pacman -Qi` and parses records into `PacmanPackage` (name, description, version, install date, depends, required-by).
2. `internal/ui` renders the Bubble Tea table and search field.

Logging from the utils layer goes to **stdout** via `log/slog` (see `internal/utils/logger.go`).

## Project layout

```
main.go              # entrypoint
cmd/                 # Cobra root command
internal/utils/      # pacman invocation + parsing + models
internal/utils_test/ # external tests (`package utils_test`)
internal/ui/         # TUI (table + search)
scripts/install-local.sh
```

## License

Unspecified; add a `LICENSE` file if you publish this repository.
