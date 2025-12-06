# wslgogit

A minimal Git wrapper for using WSL Git from Windows applications.

## Problem

Windows GUI tools (Sublime Merge, Fork, VS Code) cannot directly use Git installed in WSL when working with repositories in the WSL filesystem (`\\wsl.localhost\...`). Existing solutions have bugs with argument escaping or don't support UNC paths.

## Solution

wslgogit is a simple wrapper that:
- Detects the WSL distribution from UNC paths
- Converts Windows paths to Unix paths
- Properly escapes all arguments
- Passes through stdin/stdout/stderr

## Installation

Download the appropriate binary from [Releases](../../releases):
- `wslgogit-amd64.exe` for x64 systems
- `wslgogit-arm64.exe` for ARM64 systems

Place it in a directory (e.g., `C:\tools\wslgogit.exe`)

Or build from source:
```bash
make build
```

Binaries will be in the `build/` directory.

## Usage

### Sublime Merge
```json
{
    "git_binary": "C:\\tools\\wslgogit.exe"
}
```

### VS Code
```json
{
    "git.path": "C:\\tools\\wslgogit.exe"
}
```

### Fork

Preferences → Git → Custom Git Instance → `C:\tools\wslgogit.exe`

## Requirements

- Windows 10/11 with WSL
- Git installed in WSL (`wsl sudo apt install git`)

## How it works

Converts commands like:
```
git commit -m "message"
```

Into:
```
wsl -d Ubuntu-24.04 bash -c "cd '/home/user/repo' && git commit -m 'message'"
```

The WSL distribution is automatically detected from UNC paths like `\\wsl.localhost\Ubuntu-24.04\...`
