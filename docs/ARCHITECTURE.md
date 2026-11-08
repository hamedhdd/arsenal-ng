# Arsenal-NG Architecture

## Overview

Arsenal-ng is a terminal-based command launcher and cheat sheet manager for penetration testing tools. Built with Go and the Bubble Tea TUI framework, it provides an interactive interface for browsing, searching, and executing command templates.

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                         User Input                           │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                    Terminal UI (Bubble Tea)                  │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │   Search    │  │   Command    │  │   Argument       │   │
│  │   View      │  │   List       │  │   Form           │   │
│  └─────────────┘  └──────────────┘  └──────────────────┘   │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                      Business Logic                          │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │   Loader    │  │   Search     │  │   Argument       │   │
│  │   (YAML)    │  │   Engine     │  │   Parser         │   │
│  └─────────────┘  └──────────────┘  └──────────────────┘   │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                      State & Storage                         │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │  Variables  │  │  Config Dir  │  │   Tool Check     │   │
│  │  (JSON)     │  │  (XDG)       │  │   (PATH)         │   │
│  └─────────────┘  └──────────────┘  └──────────────────┘   │
└────────────────────┬────────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│                   Output & Execution                         │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │  Terminal   │  │   Direct     │  │   Clipboard      │   │
│  │  Inject     │  │   Exec       │  │   Copy           │   │
│  └─────────────┘  └──────────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Package Structure

### `/cmd/arsenal-ng`
**Purpose**: Application entry point

**Responsibilities**:
- CLI flag parsing
- Initialization and orchestration
- Logging setup
- TUI lifecycle management
- User interaction flow (install prompts, action selection)

**Key Functions**:
- `main()`: Entry point
- `run()`: Main application flow
- `setupLogging()`: Initialize file logging
- `loadAndFilterCheats()`: Load and optionally filter cheat sheets
- `runTUI()`: Start the interactive UI
- `checkAndInstallTool()`: Verify tool availability and offer installation
- `handleCommandAction()`: Execute, inject, or copy the final command

### `/internal/loader`
**Purpose**: Load and search cheat sheets

**Responsibilities**:
- Parse embedded YAML files
- Validate cheat sheet structure
- Provide search functionality across all fields

**Key Files**:
- `loader.go`: Load cheat files using `//go:embed`
- `search.go`: Multi-field fuzzy search implementation
- `cheat-files/*.yaml`: Embedded command templates (242 tools, 2872 commands)

**Key Functions**:
- `Load()`: Parse all embedded YAML files
- `Search(query string, cheats []*Cheat)`: Filter cheats by search terms

**Data Flow**:
```
YAML Files (//go:embed)
    ↓
Parse & Validate
    ↓
[]Cheat struct
    ↓
Search Index
    ↓
Filtered Results
```

### `/internal/model`
**Purpose**: Core data structures

**Responsibilities**:
- Define cheat sheet structure
- Parse command arguments
- Build final commands from templates

**Key Structures**:
```go
type Cheat struct {
    Tool        string
    Tags        []string
    Actions     []Action
    InstallInfo *InstallCmd
}

type Action struct {
    Title   string
    Desc    string
    Command string
}

type Argument struct {
    Name    string
    Default string
    Value   string
}
```

**Key Functions**:
- `ParseArguments(command string)`: Extract `{{arg}}` placeholders
- `BuildCommand(command string, args []Argument)`: Replace placeholders with values
- `HasEmptyArgs(args []Argument)`: Check for unfilled required arguments

### `/internal/ui`
**Purpose**: Terminal user interface

**Responsibilities**:
- Render search interface
- Display command list with syntax highlighting
- Handle keyboard input
- Manage application state
- Provide argument input forms

**Key Files**:
- `app_model.go`: Main application state
- `app_new.go`: Initialize TUI app
- `app_state.go`: State management constants
- `update_*.go`: Event handlers for different views
- `view_*.go`: Rendering logic for different views
- `highlight_*.go`: Syntax highlighting and fuzzy matching
- `helpers_*.go`: Layout, string, and tag utilities
- `constants.go`: UI layout constants

**State Machine**:
```
┌──────────────┐
│  StateMain   │ ◄──────────┐
└──────┬───────┘            │
       │ Enter              │ Esc
       ▼                    │
┌──────────────┐            │
│  StateArgs   │────────────┘
└──────────────┘  Command Built
```

**Key Components**:
- **Search View**: Text input + filtered list
- **Info Box**: Command details and description
- **Tools View**: Paginated table of all available tools
- **Help View**: Keyboard shortcuts and usage guide

### `/internal/state`
**Purpose**: Variable persistence and management

**Responsibilities**:
- Store global variables (ip, domain, etc.)
- Persist variables to JSON file
- Auto-fill command arguments from variables
- Atomic file writes for data safety

**Key Files**:
- `store.go`: Variable CRUD operations
- `persist.go`: File I/O (load/save JSON)
- `format.go`: Display formatting

**Key Functions**:
- `Set(key, value string)`: Store a variable
- `Get(key string)`: Retrieve a variable
- `Unset(key string)`: Remove a variable
- `ApplyToCommand(args []Argument)`: Pre-fill arguments from variables
- `SaveToFile()`: Persist to `~/.config/arsenal-ng/variables.json`

**Storage Location**:
```
Linux/macOS: ~/.config/arsenal-ng/variables.json
Windows:     %LOCALAPPDATA%\arsenal-ng\variables.json
```

### `/internal/config`
**Purpose**: Configuration and path management

**Responsibilities**:
- Determine config directory (XDG_CONFIG_HOME support)
- Provide paths for variables, logs
- Application metadata (version, name)

**Key Files**:
- `app.go`: Version and app name constants
- `path.go`: Path resolution and directory creation
- `theme.go`: Color scheme definitions

**Key Functions**:
- `GetConfigDir()`: Returns `~/.config/arsenal-ng/` (or XDG override)
- `GetVariablesPath()`: Returns path to `variables.json`
- `GetLogPath()`: Returns path to log file

### `/internal/toolmgr`
**Purpose**: Tool installation and availability checking

**Responsibilities**:
- Check if tools are installed (PATH lookup)
- Detect available package managers
- Generate installation commands
- Provide user-friendly installation prompts

**Key Files**:
- `check.go`: Tool availability verification
- `install.go`: Installation command generation

**Key Functions**:
- `IsInstalled(tool string)`: Check if binary exists in PATH
- `GetToolBinary(tool string)`: Normalize tool name
- `DetectPackageManager()`: Find apt/dnf/brew/pacman
- `GetInstallCommand(install *InstallCmd)`: Generate install command

**Package Manager Priority**:
1. Language-specific: go, pipx, pip, cargo
2. System package manager: apt, dnf, pacman, brew
3. Git clone
4. Manual instructions

### `/internal/output`
**Purpose**: Command output handling

**Responsibilities**:
- Inject commands into terminal input buffer (TIOCSTI)
- Fallback to stdout when TIOCSTI unavailable
- Platform-specific behavior (Linux/macOS/Windows)

**Key Files**:
- `terminal.go`: Linux and macOS TIOCSTI implementation
- `terminal_windows.go`: Windows stub (prints to stdout)

**Platform Support**:
- **Linux**: TIOCSTI (requires kernel 6.2+ flag or CAP_SYS_ADMIN)
- **macOS**: TIOCSTI (native support)
- **Windows**: Fallback to stdout only

## Data Flow

### 1. Application Startup

```
main()
  ↓
setupLogging()
  ↓
loadAndFilterCheats()
  ├─ loader.Load() → Parse YAML files
  └─ filterCheatsByTag() → Apply --profile filter
  ↓
runTUI()
  ├─ ui.New() → Initialize Bubble Tea app
  └─ program.Run() → Start interactive UI
```

### 2. User Search Flow

```
User types "nmap scan"
  ↓
update_search.go: handleSearchInput()
  ↓
loader.Search("nmap scan", allCheats)
  ├─ Split query into terms: ["nmap", "scan"]
  ├─ Check each cheat's fields (tool, title, desc, tags, command)
  └─ Return matching cheats
  ↓
view_search.go: renderList()
  ├─ Apply syntax highlighting
  ├─ Apply fuzzy highlight on matching terms
  └─ Render paginated list
```

### 3. Command Building Flow

```
User selects command: "nmap -sV {{ip}}"
  ↓
model.ParseArguments()
  ├─ Extract: [{Name: "ip", Default: "", Value: ""}]
  └─ Return arguments list
  ↓
state.ApplyToCommand()
  ├─ Check if variable "ip" exists (e.g., "10.10.10.10")
  └─ Pre-fill: [{Name: "ip", Default: "", Value: "10.10.10.10"}]
  ↓
ui.StateArgs: Render form
  ↓
User edits/confirms
  ↓
model.BuildCommand()
  ├─ Replace "{{ip}}" → "10.10.10.10"
  └─ Return: "nmap -sV 10.10.10.10"
  ↓
checkAndInstallTool()
  ├─ toolmgr.IsInstalled("nmap")
  └─ If not installed → offerToolInstallation()
  ↓
handleCommandAction()
  ├─ [1] Execute: exec.Command("sh", "-c", command)
  ├─ [2] Inject: output.ToTerminal(command)
  └─ [3] Copy: clipboard.WriteAll(command)
```

### 4. Variable Persistence Flow

```
User types: set ip=10.10.10.10
  ↓
state.Set("ip", "10.10.10.10")
  ├─ Validate variable name (alphanumeric + _ -)
  ├─ Store in memory: map[string]string{"ip": "10.10.10.10"}
  └─ state.SaveToFile()
      ├─ Atomic write to temp file
      ├─ JSON marshal
      └─ Rename to variables.json
  ↓
Next session: loader starts
  ↓
state.LoadFromFile()
  ├─ Read ~/.config/arsenal-ng/variables.json
  └─ Restore variables into memory
```

## Configuration Files

### Cheat Sheet Format (YAML)

```yaml
tool: nmap
tags: [recon, network, scan]
install:
  apt: nmap
  brew: nmap
  dnf: nmap
  pacman: nmap

actions:
  - title: nmap - TCP SYN scan
    desc: Fast TCP scan of top 1000 ports
    command: "nmap -sS {{ip}}"

  - title: nmap - Service version detection
    desc: Detect service versions on target
    command: "nmap -sV -p {{port|1-65535}} {{ip}}"
```

### Variables Storage (JSON)

```json
{
  "ip": "10.10.10.10",
  "domain": "corp.local",
  "user": "administrator",
  "lhost": "192.168.1.100"
}
```

## Extension Points

### Adding a New Cheat Sheet

1. Create `internal/loader/cheat-files/newtool.yaml`
2. Follow the YAML schema
3. Rebuild: `go build ./cmd/arsenal-ng`

### Adding a New UI View

1. Create view state constant in `app_state.go`
2. Add `update_newview.go` for event handling
3. Add `view_newview.go` for rendering
4. Update `Update()` and `View()` in main app

### Adding a New Package Manager

1. Edit `internal/toolmgr/install.go`
2. Add detection in `DetectPackageManager()`
3. Add install command format in `GetInstallCommand()`
4. Update cheat YAML files with new package field

## Testing Strategy

### Unit Tests
- **Model**: Argument parsing, command building
- **Loader**: YAML validation, search algorithm
- **State**: Variable CRUD, persistence, atomic writes
- **Config**: Path resolution, XDG validation
- **ToolMgr**: Tool detection, command generation
- **UI Helpers**: String truncation, word wrap, color hashing

### Integration Tests
- End-to-end cheat loading
- Search across all embedded files
- Variable persistence round-trip

### Manual Testing
- TUI interaction flows
- Platform-specific behavior (TIOCSTI)
- Terminal size handling

## Performance Considerations

1. **Startup Time**
   - All cheat files are embedded using `//go:embed`
   - Lazy parsing on startup (~0.2s for 242 files)
   - No external file I/O

2. **Search Performance**
   - Linear scan through all cheats (acceptable for <3000 items)
   - No indexing (trade-off: simplicity vs speed)
   - Multi-field search with short-circuit evaluation

3. **Memory Usage**
   - All cheats loaded in memory (~5MB)
   - UI state minimal (current view, search, cursor)
   - Variables stored in memory + JSON file

## Security Considerations

1. **Variable Name Validation**
   - Regex: `^[a-zA-Z0-9_-]+$`
   - Max length: 64 characters
   - Prevents path traversal (no `/` or `..`)

2. **XDG_CONFIG_HOME Validation**
   - Only absolute paths accepted
   - Falls back to `~/.config` if relative

3. **Command Execution**
   - User confirmation before install
   - No automatic execution without selection
   - Shell injection risk: user is responsible for input

4. **File Permissions**
   - Config dir: 0755
   - Variables file: 0644 (default)
   - Log file: readable by user only

## Future Enhancements

### Planned Features
1. **Command History**: Track recently used commands
2. **Favorites System**: Star frequently used commands
3. **Configuration File**: User-defined themes and preferences
4. **Export Functionality**: Generate bash scripts from commands
5. **Plugin System**: Load custom cheat files from external directories

### Potential Improvements
1. **Search Index**: Build inverted index for faster search
2. **Command Templates**: Reusable multi-command workflows
3. **Remote Cheat Sheets**: Fetch from Git repositories
4. **Syntax Validation**: Validate command syntax before execution
5. **Command History Persistence**: Save to history.json

## Troubleshooting

### Build Fails on Windows
- **Cause**: Missing `terminal_windows.go` stub
- **Solution**: Create stub file with `//go:build windows` constraint

### TIOCSTI Not Working (Linux 6.2+)
- **Cause**: Kernel security restriction
- **Solution**: `sysctl -w dev.tty.legacy_tiocsti=1` or use CAP_SYS_ADMIN

### Variables Not Persisting
- **Cause**: Permission error on config directory
- **Solution**: Check `~/.config/arsenal-ng/` permissions

### Tool Not Found After Install
- **Cause**: Binary not in PATH
- **Solution**: Verify installation and refresh PATH

## Contributing

See `CONTRIBUTING.md` for detailed contribution guidelines.

### Quick Start for Contributors
1. Fork the repository
2. Create feature branch: `git checkout -b feature/new-cheat`
3. Add cheat file: `internal/loader/cheat-files/mytool.yaml`
4. Test locally: `make run`
5. Run tests: `make test`
6. Submit PR

---

**Maintainer**: HamedHD  
**License**: See LICENSE file  
**Last Updated**: 2026-11-08
