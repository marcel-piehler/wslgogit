package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: failed to get current directory: %v\n", err)
		os.Exit(1)
	}

	// Detect WSL distribution from UNC path if applicable
	distro, wslPath := parseUNCPath(cwd)

	// Build git command with properly escaped arguments
	gitArgs := make([]string, len(os.Args)-1)
	for i, arg := range os.Args[1:] {
		gitArgs[i] = escapeShellArg(arg)
	}

	// Construct and execute the WSL command
	var cmd *exec.Cmd
	if distro != "" {
		// UNC path: execute in specific WSL distribution
		gitCmd := fmt.Sprintf("cd '%s' && git %s", wslPath, strings.Join(gitArgs, " "))
		cmd = exec.Command("wsl", "-d", distro, "bash", "-c", gitCmd)
	} else {
		// Regular Windows path: execute in default WSL distribution
		unixPath := toUnixPath(cwd)
		gitCmd := fmt.Sprintf("cd '%s' && git %s", unixPath, strings.Join(gitArgs, " "))
		cmd = exec.Command("wsl", "bash", "-c", gitCmd)
	}

	// Pass through stdin/stdout/stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute and handle exit code
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

// escapeShellArg escapes an argument for safe use in bash shell.
// Arguments containing special characters are wrapped in single quotes,
// with any single quotes within the argument properly escaped.
func escapeShellArg(arg string) string {
	// If argument contains only safe characters, return as-is
	if regexp.MustCompile(`^[a-zA-Z0-9_\-\./]+$`).MatchString(arg) {
		return arg
	}

	// Escape single quotes: 'arg' -> 'arg'\''more'
	// This closes the quote, adds an escaped quote, and reopens the quote
	escaped := strings.ReplaceAll(arg, "'", "'\"'\"'")
	return fmt.Sprintf("'%s'", escaped)
}

// parseUNCPath extracts the WSL distribution name and path from UNC paths.
// Supports both \\wsl.localhost\distro\path and \\wsl$\distro\path formats.
//
// Example:
//   \\wsl.localhost\Ubuntu-24.04\home\user\repo
//   returns: ("Ubuntu-24.04", "/home/user/repo")
func parseUNCPath(path string) (distro, wslPath string) {
	// Match \\wsl.localhost\distro\path or \\wsl$\distro\path
	re := regexp.MustCompile(`^\\\\wsl(?:\.localhost|\$)\\([^\\]+)\\(.+)$`)
	matches := re.FindStringSubmatch(path)

	if len(matches) == 3 {
		distro = matches[1]
		// Convert Windows path separators to Unix
		wslPath = "/" + strings.ReplaceAll(matches[2], "\\", "/")
		return
	}

	return "", ""
}

// toUnixPath converts a Windows path to Unix/WSL path format.
// Drive letters are converted to /mnt/<drive> format.
//
// Example:
//   C:\Users\Marcel\repo -> /mnt/c/Users/Marcel/repo
func toUnixPath(path string) string {
	// Convert backslashes to forward slashes
	unixPath := filepath.ToSlash(path)

	// Convert drive letter (C:) to /mnt/c
	if len(unixPath) >= 2 && unixPath[1] == ':' {
		drive := strings.ToLower(string(unixPath[0]))
		unixPath = "/mnt/" + drive + unixPath[2:]
	}

	return unixPath
}