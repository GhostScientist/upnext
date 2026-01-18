# upnext

> The ultimate minimal todo CLI — fast, beautiful, calm.

---

## Philosophy

**"Do one thing beautifully."**

upnext should feel like a breath of fresh air compared to bloated productivity apps. Fast to invoke, pleasant to look at, zero friction.

### Core Principles

- **Instant** — sub-100ms response times
- **Calm** — soft colors, breathing room, not shouty
- **Playful** — subtle ASCII art that doesn't overwhelm
- **Focused** — your next action is always clear

---

## Visual Design

### Default List View

```
  ┌─────────────────────────────────────┐
  │           ✨ up next ✨              │
  └─────────────────────────────────────┘

  1. │░░░░░░░░░░│ Finish API documentation
  2. │░░░░░░    │ Review Dakota's PR
  3. │░░        │ Call mom back

  ───────────────────────────────────────
     3 items · added "Call mom" just now
```

### Empty State

```
    *  .  *
       *      ☀️
   .    *  .
  
  nothing up next — enjoy the calm
```

### Task Completed

```
  ✓ done! ──────────────
  
  "Finish API docs" archived
  2 items remaining
```

### Milestone Celebration (every 10 completions)

```
  ┌────────────────────────┐
  │  🎉 10 tasks crushed!  │
  │     keep flowing       │
  └────────────────────────┘
```

---

## Color Palette

Soft, accessible colors that work on both light and dark terminals:

| Element          | Color         | Hex       |
|------------------|---------------|-----------|
| Header accent    | Soft lavender | `#b4befe` |
| Item numbers     | Muted cyan    | `#89dceb` |
| Task text        | Warm white    | `#cdd6f4` |
| Done items       | Dim gray      | `#6c7086` |
| Success messages | Soft green    | `#a6e3a1` |
| Warnings         | Peach         | `#fab387` |

---

## Commands

### Primary Commands

| Command            | Action                    |
|--------------------|---------------------------|
| `upnext`           | Show your list            |
| `upnext add "task"`| Add a task                |
| `upnext done`      | Complete top item         |
| `upnext done 2`    | Complete item #2          |
| `upnext drop 3`    | Remove without completing |
| `upnext bump 3`    | Move item to top          |
| `upnext clear`     | Archive all done items    |

### Shorthand Aliases

For muscle-memory speed:

| Short    | Expands To      |
|----------|-----------------|
| `up`     | `upnext`        |
| `up a`   | `upnext add`    |
| `up d`   | `upnext done`   |
| `up b`   | `upnext bump`   |

### Flags

| Flag       | Effect                              |
|------------|-------------------------------------|
| `--plain`  | No colors/art (for scripting)       |
| `--json`   | Output as JSON                      |
| `--help`   | Show help                           |
| `--version`| Show version                        |

---

## Technical Architecture

### Language: Go

**Rationale:**
- Excellent cross-platform support
- Single binary distribution
- Fast startup times
- Charm ecosystem provides beautiful CLI primitives out of the box

### Dependencies

```
charm.sh/lipgloss    - styling and colors
spf13/cobra          - CLI command parsing
```

Alternatively, use only stdlib `flag` package to minimize dependencies further.

### Alternative: Rust

If maximum performance and minimal binary size are priorities:

```
crossterm           - terminal manipulation
clap                - argument parsing
serde + serde_json  - data serialization
directories         - XDG-compliant paths
```

---

## Data Storage

### Format: JSON

Simple, human-readable, easy to debug:

```json
{
  "version": 1,
  "items": [
    {
      "id": "a1b2c3",
      "text": "Finish API documentation",
      "created": "2025-01-17T10:00:00Z",
      "position": 0
    },
    {
      "id": "d4e5f6",
      "text": "Review PR",
      "created": "2025-01-17T09:30:00Z",
      "position": 1
    }
  ],
  "archive": [
    {
      "id": "g7h8i9",
      "text": "Setup CI pipeline",
      "created": "2025-01-16T14:00:00Z",
      "completed": "2025-01-17T08:00:00Z"
    }
  ],
  "stats": {
    "total_completed": 42,
    "streak_days": 5
  }
}
```

### File Location (XDG-compliant)

| Platform | Path                              |
|----------|-----------------------------------|
| Linux    | `~/.local/share/upnext/todos.json`|
| macOS    | `~/.local/share/upnext/todos.json`|
| Windows  | `%APPDATA%\upnext\todos.json`     |

---

## Project Structure

```
upnext/
├── cmd/
│   └── upnext/
│       └── main.go          # entry point
├── internal/
│   ├── commands/
│   │   ├── add.go
│   │   ├── done.go
│   │   ├── list.go
│   │   ├── drop.go
│   │   ├── bump.go
│   │   └── clear.go
│   ├── store/
│   │   ├── store.go         # storage interface
│   │   └── json.go          # JSON file implementation
│   ├── model/
│   │   └── todo.go          # data structures
│   └── ui/
│       ├── colors.go        # color palette
│       ├── render.go        # list rendering
│       ├── art.go           # ASCII art moments
│       └── messages.go      # success/error messages
├── go.mod
├── go.sum
├── Makefile
├── README.md
└── LICENSE
```

---

## Installation

### Homebrew (macOS/Linux)

```bash
brew install upnext
```

### Go Install

```bash
go install github.com/yourusername/upnext@latest
```

### Binary Download

Pre-built binaries for:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64)

### Shell Alias Setup

Add to `~/.bashrc`, `~/.zshrc`, or equivalent:

```bash
alias up="upnext"
```

---

## Roadmap

### Phase 1: MVP ✓

- [ ] Add tasks with `upnext add "task"`
- [ ] List tasks with beautiful colored output
- [ ] Complete tasks with `upnext done [n]`
- [ ] Delete tasks with `upnext drop [n]`
- [ ] Reorder with `upnext bump [n]`
- [ ] Persistent JSON storage
- [ ] Cross-platform binary releases
- [ ] Basic ASCII art states

### Phase 2: Polish

- [ ] Relative timestamps ("added 2h ago")
- [ ] Milestone celebrations (every 10 completions)
- [ ] Shell completions (bash, zsh, fish, PowerShell)
- [ ] `--plain` flag for scripting/piping
- [ ] `--json` flag for machine-readable output
- [ ] Configurable colors via config file

### Phase 3: Quality of Life

- [ ] `upnext edit 2` — edit task text inline
- [ ] `upnext undo` — restore last completed/dropped item
- [ ] Archive auto-pruning (configurable retention)
- [ ] Stats command (`upnext stats`)

### Phase 4: Maybe Never

Features intentionally out of scope to maintain minimalism:

- Tags/contexts
- Due dates and reminders
- Multiple lists
- Syncing/cloud storage
- Collaboration features

---

## Design Decisions

### Why No Due Dates?

upnext is about **right now**. What's your next action? Due dates add cognitive overhead and push the tool toward calendar territory. Keep it simple.

### Why No Categories/Tags?

Categories encourage over-organization. If you need categories, you might need a different tool. upnext is a single stream of "what's next."

### Why JSON Over SQLite?

- Human-readable and editable
- No dependencies
- Trivial to backup (it's just a file)
- Fast enough for hundreds of items
- Easy to debug

### Why Go Over Node/Python?

- Single binary, no runtime required
- Fast startup (critical for CLI tools)
- Excellent cross-compilation
- Charm ecosystem is Go-native

---

## Configuration (Optional)

`~/.config/upnext/config.toml`:

```toml
[display]
# Use plain output (no colors/unicode)
plain = false

# Show relative timestamps
relative_time = true

# Max items to show (0 = unlimited)
max_display = 10

[archive]
# Days to keep completed items (0 = forever)
retention_days = 30

[theme]
# Override default colors
accent = "#b4befe"
success = "#a6e3a1"
```

---

## Usage Examples

### Daily Workflow

```bash
# Morning: check what's up
$ up
  ┌─────────────────────────────────────┐
  │           ✨ up next ✨              │
  └─────────────────────────────────────┘

  1. │░░░░░░░░░░│ Review PR from Seth
  2. │░░░░░░    │ Prep slides for Monday
  3. │░░        │ Email Julian re: deployment

# Add something that came up
$ up a "Call back vendor about license"

# Knock out the top item
$ up d
  ✓ done! ──────────────
  "Review PR from Seth" archived

# Promote something urgent
$ up b 3
  ↑ bumped! ──────────────
  "Email Julian" moved to top
```

### Scripting Integration

```bash
# Check if there are any tasks
if upnext --plain | grep -q .; then
  echo "You have tasks!"
fi

# Get count of tasks
upnext --json | jq '.items | length'

# Add task from another script
upnext add "$(date): Deploy completed"
```

---

## Contributing

1. Fork the repository
2. Create a feature branch
3. Keep changes minimal and focused
4. Ensure cross-platform compatibility
5. Add tests for new functionality
6. Submit a pull request

---

## License

MIT License — Use it, fork it, make it yours.

---

## Acknowledgments

- [Charm](https://charm.sh) for making beautiful CLIs possible
- [Catppuccin](https://catppuccin.com) color palette inspiration
- Everyone who believes simple tools can be delightful
