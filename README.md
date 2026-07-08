<h1 align="center">arsenal-ng</h1>

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20WSL-lightgrey" alt="Platform">
  <img src="https://img.shields.io/github/license/hamedhdd/arsenal-ng?color=yellow" alt="License">
  <br>
  <img src="https://img.shields.io/badge/Tools-242-blueviolet?style=flat&logo=linux&logoColor=white" alt="Tools Count">
  <img src="https://img.shields.io/badge/Commands-2872-ff69b4?style=flat&logo=gnubash&logoColor=white" alt="Commands Count">
</p>

<p align="center">
  <b>🎯 Interactive pentest cheat sheet and command builder!</b>
</p>

<p align="center">
  Inspired by <a href="https://github.com/Orange-Cyberdefense/arsenal">arsenal</a>, rewritten from scratch with a focus on simplicity, speed and developer experience.
</p>

<p align="center">
  <strong>Arsenal-ng helps you find, build, and launch commands.</strong><br>
  You need to install the actual tools (nmap, ffuf, etc.) separately.
</p>

<p align="center">
  <img src="assets/basic-search.gif" alt="Basic Search Demo" width="800">
</p>

---

## 📦 Installation

### Option 1: Go Install

```bash
go install -v github.com/hamedhdd/arsenal-ng/cmd/arsenal-ng@latest
```
> Requires Go 1.22+ Ensure `$(go env GOPATH)/bin` is in your `$PATH`.


### Option 2: Build from Source Code

```bash
git clone https://github.com/hamedhdd/arsenal-ng.git
cd arsenal-ng
make build
./bin/arsenal-ng
```

### Alias (Optional)

You can create an alias for quick access (e.g., `a`):

**Zsh:**
```bash
echo "alias a='arsenal-ng'" >> ~/.zshrc
source ~/.zshrc
```

**Bash:**
```bash
echo "alias a='arsenal-ng'" >> ~/.bashrc
source ~/.bashrc
```
---

## 🖥️ Platform Support

| Platform | Status | Notes |
|----------|--------|-------|
| **Linux** | ✅ Fully Supported | Requires kernel 6.2+ configuration for terminal prefill (see [Troubleshooting](#-troubleshooting)) |
| **macOS** | ✅ Fully Supported | Works out of the box, no additional configuration needed |
| **WSL** | ✅ Supported | Runs the Linux binary under Windows. Same kernel 6.2+ note applies |
| **Windows (native)** | ❌ Not Supported | Use Linux, macOS, or WSL |

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| ⚡ **Instant Startup** | Single binary, no dependencies, launches in milliseconds |
| 🔍 **Smart Search** | Multi-word fuzzy search across tool names, titles, tags, descriptions, and commands |
| 🎨 **Syntax Highlighting** | Commands are color-coded with syntax highlighting for better readability |
| 🏷️ **Colored Tags** | Each tag has a consistent, distinct color based on hash for quick visual identification |
| 📝 **Simple YAML Format** | Easy to maintain and extend cheatsheets |
| 🔧 **Argument System** | Support for `{{arg}}` and `{{arg\|default}}` placeholders with auto-completion |
| 🌐 **Global Variables** | Set once, use everywhere - variables auto-fill in all commands |
| 📊 **Tools View** | Browse all available tools with command counts in a paginated table |
| 💡 **Command Hints** | Interactive hints for special commands (`set`, `unset`, `variables`, `tools`) |
| ❓ **Built-in Help** | Press `?` for comprehensive help screen with all shortcuts |
| � **Command Builder** | Builds complete commands with filled arguments ready to copy and execute |
| 🛠️ **242 Tool References** | Cheat sheets for nmap, ffuf, gobuster, impacket, bloodhound, and 237 more (tools not included - install separately) |

---

## 🚀 Usage

### What Arsenal-ng Does

**Arsenal-ng is an interactive command builder with automatic tool installation.** It helps you:
- 🔍 Find the right command syntax for 242 pentesting tools
- 📝 Fill in command arguments with an interactive form
- 🛠️ **Automatically detect missing tools and offer to install them**
- 🚀 **Execute the final command directly**, inject it to your terminal, or copy it to the clipboard

**New in this version:**
- Arsenal-ng now checks if tools are installed before showing commands
- If a tool is missing, it offers to install it automatically
- Supports apt, dnf, pacman, brew, pip, pipx, go, cargo, and git installations
- You confirm before any installation happens

### Quick Start

```bash
# Launch the application
arsenal-ng

# The TUI will open with all available commands
# Use arrow keys to navigate, type to search, Enter to select
```

When you select a command and fill in arguments:
1. Arsenal-ng checks if the tool is installed
2. **If missing:** Prompts you to install it automatically
3. **If installed:** Prints the command to your terminal

### Basic Workflow

1. **Search for a command**: Type keywords (e.g., `nmap scan`, `ffuf`)
2. **Navigate results**: Use arrow keys to browse matching commands
3. **Select command**: Press Enter on the desired command
4. **Fill arguments**: If the command has `{{placeholders}}`, fill them in
5. **Proceed**: Choose to execute directly, inject to terminal, or copy to clipboard

### Example Session

```bash
# 1. Launch arsenal-ng
arsenal-ng

# 2. Search for "403bypasser" and select a command

# 3. Fill in the URL: https://example.com/admin

# 4. Arsenal-ng checks if 403bypasser is installed:
⚠️  Tool '403bypasser' is not installed.

Install with: git clone https://github.com/yunemse48/403bypasser && cd 403bypasser && pip install -r requirements.txt

Would you like to install it now? [y/N]: y

# 5. Tool installs automatically

✅ 403bypasser installed successfully!

# 6. Command is ready! Choose how to proceed:
🚀 Command ready: 403bypasser -u https://example.com/admin

How would you like to proceed?
  [1] Execute directly
  [2] Inject into terminal (TIOCSTI)
  [3] Copy to clipboard
  [4] Cancel

Choice [1]: 1

# 7. The command executes directly!
```

### Important Notes

- **Automatic tool detection** — Arsenal-ng checks if tools exist before showing commands
- **Optional auto-installation** — You're prompted before any installation happens
- **Multiple package managers supported** — apt, dnf, pacman, brew, pip, pipx, go, cargo, git
- **242 tool references included** — Installation info available for commonly used tools
- Global variables (like `set ip=10.10.10.10`) auto-fill across all commands in your session

### Keyboard Shortcuts

#### Main Search View
| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate up/down in the list |
| `Ctrl+P` / `Ctrl+N` | Navigate list (vim-style) |
| `PgUp` / `PgDown` | Jump one page up/down |
| `Enter` | Select highlighted command or execute special command |
| `Esc` / `Ctrl+C` | Exit application |
| `?` | Show help screen |
| `q` | Quit (in some views) |

#### Argument Input View
| Key | Action |
|-----|--------|
| `Tab` / `↓` | Move to next argument field |
| `Shift+Tab` / `↑` | Move to previous argument field |
| `Enter` | Execute command with filled arguments |
| `Esc` | Go back to search view |

#### Tools View

<p align="center">
  <img src="assets/tools-view.gif" alt="Tools View Demo" width="800">
</p>

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate table rows |
| `←` / `→` or `h` / `l` | Change page (previous/next) |
| `Enter` / `Esc` | Go back to search view |



#### Help Views

<p align="center">
  <img src="assets/help-screen.gif" alt="Help Screen Demo" width="800">
</p>

| Key | Action |
|-----|--------|
| `Enter` / `Esc` | Go back to search view |

---

## 🌐 Global Variables

Set variables once and reuse them across all commands in your session. Variables automatically pre-fill argument fields in commands.

### Special Commands

| Command | Description |
|---------|-------------|
| `set key=value` | Set a global variable (e.g., `set ip=10.10.10.10`) |
| `unset key` | Remove a global variable |
| `variables` | List all currently set variables |
| `tools` | Show all available tools with command counts |
| `help` | Display comprehensive help screen |

### How it works

<p align="center">
  <img src="assets/variables.gif" alt="Variables Demo" width="800">
</p>

1. Type `set ip=10.10.10.10` and press Enter
2. Select any command with `{{ip}}` placeholder
3. The `ip` field will be **pre-filled** automatically with `10.10.10.10`
4. You can still edit the pre-filled value if needed
5. Variables persist throughout your session until you exit or unset them

### Example Workflow

```bash
# Set common variables at the start of your session
set ip=10.10.10.10
set domain=corp.local
set user=administrator

# Now select any command - arguments will auto-fill!
# Before: nmap -sV {{ip}}
# After:  nmap -sV 10.10.10.10  (auto-filled!)

# View all variables
variables

# Remove a variable when done
unset ip
```

### Common Variables

| Variable | Description |
|----------|-------------|
| `ip` | Target IP address |
| `domain` | Target domain |
| `user` | Username |
| `pass` | Password |
| `hash` | NTLM hash |
| `url` | Target URL |
| `port` | Port number |
| `lhost` | Local host (your IP) |
| `lport` | Local port |
| `wordlist` | Wordlist path |
| `output` | Output file name |

---

## 📄 Cheat File Format

Add your own commands by creating YAML files in `internal/loader/cheat-files/`:

```yaml
tool: mytool
tags: [recon, web, custom]

actions:
  - title: mytool - basic scan
    desc: Performs a basic scan against the target
    command: "mytool scan {{target}}"

  - title: mytool - scan with multiple arguments
    desc: Advanced scan with IP, port, and output file
    command: "mytool scan -t {{ip}} -p {{port|443}} -o {{output|scan.log}}"
```

### Argument Syntax

| Syntax | Description | Example |
|--------|-------------|---------|
| `{{arg}}` | Required argument (user must fill) | `{{ip}}` |
| `{{arg\|default}}` | Argument with default value (can be edited) | `{{port\|8080}}` |

### File Structure

- **tool**: Tool name (e.g., `nmap`, `ffuf`, `gobuster`)
- **tags**: Array of tags for categorization (e.g., `[recon, scan, network]`)
- **actions**: List of command entries
  - **title**: Display name shown in the list
  - **desc**: Optional description (shown in info box)
  - **command**: Command template with `{{placeholders}}`

### Tips

- Use descriptive titles that make it easy to find commands
- Add relevant tags for better searchability
- Use default values for commonly used options
- Global variables will auto-fill matching argument names

---

## 🔧 Development

### Prerequisites

- Go 1.22 or higher
- Make (optional, for using Makefile)

### Building from Source

```bash
git clone https://github.com/hamedhdd/arsenal-ng.git
cd arsenal-ng
make build
# Binary will be in ./bin/arsenal-ng
```

Or build directly with Go:

```bash
go build -o bin/arsenal-ng ./cmd/arsenal-ng
```

### Makefile Targets

| Command | Description |
|---------|-------------|
| `make build` | Build the binary into `bin/arsenal-ng` |
| `make run` | Run directly without producing a binary |
| `make test` | Run all tests |
| `make vet` | Run `go vet` on all packages |
| `make lint` | Run `golangci-lint` (requires installation) |
| `make clean` | Remove build artifacts |
| `make install` | Install into `$GOPATH/bin` |

### Adding New Tools

1. Create a new YAML file in `internal/loader/cheat-files/`
2. Follow the format above (see [Cheat File Format](#-cheat-file-format))
3. Rebuild: `make build`

### Testing

After making changes, test your changes:

```bash
# Run the app directly without building
make run

# Or build and test the binary
make build
./bin/arsenal-ng

# Run tests (once available)
make test

# Vet code for common issues
make vet
```

### Contributing Guidelines

- Follow Go code style conventions
- Run `make vet` before submitting PRs
- Keep YAML files organized and well-documented
- Add descriptive titles and tags to commands
- Test your changes before submitting PRs

---

## ⚠️ Troubleshooting

### Terminal Prefill Not Working (Linux kernel 6.2+)

On Linux kernel 6.2+, arsenal-ng cannot inject commands into your terminal due to TIOCSTI being restricted by default.

**Current behavior:** When you select a command, the app exits and the command is **printed to stdout** on a new line. You need to manually copy and execute it.

**To enable terminal prefill** (command appears in your input buffer ready to edit), you have two options:

#### Option 1: Enable TIOCSTI Globally
The TIOCSTI ioctl is disabled by default in newer Linux kernels for security reasons.

```bash
# Temporary (current session only)
sudo sysctl -w dev.tty.legacy_tiocsti=1

# Permanent (survives reboot)
echo "dev.tty.legacy_tiocsti=1" | sudo tee /etc/sysctl.d/99-tiocsti.conf
sudo sysctl --system
```

#### Option 2: Grant CAP_SYS_ADMIN Capability
Instead of modifying system-wide settings, you can grant the `CAP_SYS_ADMIN` capability specifically to the arsenal-ng binary.
CAP_SYS_ADMIN is powerful and virtually equivalent to `root` access. Use this method only if you fully understand the risks.

```Bash
# Ensure 'setcap' is installed (Debian/Kali/Ubuntu)
sudo apt-get install libcap2-bin

# Grant the required capability to the binary
sudo setcap "cap_sys_admin+ep" $(which arsenal-ng)
```

## 🤝 Contributing

This project is **open source** and contributions are welcome!

### How to Contribute

- 🔧 **Add a tool**: Create a YAML file in `internal/loader/cheat-files/` and submit a PR to [hamedhdd/arsenal-ng](https://github.com/hamedhdd/arsenal-ng)
- 🐛 **Report bugs**: [Open an issue](https://github.com/hamedhdd/arsenal-ng/issues) with details about the problem
- 💡 **Suggest features**: [Share your ideas](https://github.com/hamedhdd/arsenal-ng/issues) for improvements
- 📝 **Improve documentation**: Help make the README and code comments better
- ⭐ **Star the project**: Show your support at [hamedhdd/arsenal-ng](https://github.com/hamedhdd/arsenal-ng)!

### Adding Cheat Sheets

The easiest way to contribute is by adding new cheat sheet YAML files:

1. Fork [https://github.com/hamedhdd/arsenal-ng](https://github.com/hamedhdd/arsenal-ng)
2. Add your YAML file(s) to `internal/loader/cheat-files/`
3. Follow the existing format and style
4. Test your changes locally
5. Submit a pull request to `https://github.com/hamedhdd/arsenal-ng`

See [Cheat File Format](#-cheat-file-format) for details.

---

## 📋 Changelog

### Project Structure & Testing Improvements (2026-07-08)

#### 🧪 Testing

- **Unit test coverage** — Added 38 comprehensive unit tests across core packages:
  - `internal/model/argument_test.go`: 14 tests covering ParseArguments, BuildCommand, and validation helpers
  - `internal/loader/search_test.go`: 6 tests for search functionality (multi-term, case sensitivity, field matching)
  - `internal/loader/loader_test.go`: 4 tests validating all 2872 embedded cheat files parse correctly
  - `internal/state/variables_test.go`: 14 tests covering CRUD operations, persistence, and command integration
- **CI quality gate** — GitHub Actions workflow now runs `go vet` and `go test` before build step

#### 🏗️ Code Organization

- **State package refactored** — Split `variables.go` (290 lines) into logical files:
  - `store.go`: Global struct, CRUD operations, command integration
  - `persist.go`: Atomic file I/O (LoadFromFile, SaveToFile)
  - `format.go`: Display helper (FormatList)
- **Expanded Makefile** — Added 7 targets for improved developer experience:
  - `make run`: Run directly without building
  - `make test`: Run all tests
  - `make vet`: Run go vet
  - `make lint`: Run golangci-lint
  - `make clean`: Remove build artifacts
  - `make install`: Install into $GOPATH/bin
- **README updated** — Added Makefile targets table and updated contributing guidelines

### Security & Stability Fixes (2026-07-08)

#### 🔒 Security

- **Variable name validation** — `set` commands now enforce `^[a-zA-Z0-9_-]+$` (max 64 chars). Malformed names (null bytes, path separators, special chars) are rejected with a clear error message before reaching the JSON store.
- **`XDG_CONFIG_HOME` path validation** — The environment variable is now only accepted if it is an absolute path. Relative or empty values fall back to the default `~/.config/arsenal-ng`, preventing environment-based redirection of the config store.
- **Pinned GoReleaser CI version** — Changed `version: latest` to `version: "~> v2"` in the release workflow, preventing a compromised upstream `latest` tag from running arbitrary code in the build pipeline.

#### 🐛 Bug Fixes

- **Terminal output formatting** — Fixed TIOCSTI fallback on Linux 6.2+. When TIOCSTI fails, commands now print on a new line (was: appeared on same line as prompt, making them hard to read).
- **TIOCSTI silent failure fixed** — On Linux kernel 6.2+ where `TIOCSTI` is restricted by default, the app previously swallowed the command silently. It now detects the `EPERM` errno on the first byte, logs a clear message with the fix (`sysctl -w dev.tty.legacy_tiocsti=1`), and falls back to printing the command to stdout so it is never lost.
- **Windows native support dropped** — Removed the Windows stub. The app targets Linux and macOS only, where proper terminal prefill via TIOCSTI is available.

#### 🧹 Code Quality

- **Removed deprecated `rand.Seed`** — Dropped the `rand.Seed(time.Now().UnixNano())` call deprecated since Go 1.20. The global RNG is auto-seeded in modern Go versions.

---

<p align="center">
  Made with ❤️ and <a href="https://github.com/hamedhdd">HamedHD.</a>
</p>
