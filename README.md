# Diary CLI

A dead simple diary CLI written in Go

## Features
- **AES-256-GCM Encryption**: Securely encrypts entries so they cannot be read without your master password.
- **TUI & Neovim Integration**: Beautiful terminal user interfaces for selecting and interacting with your past entries, seamlessly dropping you into your `$EDITOR` (defaults to Neovim).
- **Markdown & Frontmatter**: Parses YAML frontmatter for tags, mood, and weather. Renders beautiful Markdown directly in your terminal.
- **Statistics**: View your writing streaks and word counts.
- **Full-Text Search**: Search through the decrypted body of all your entries in memory.
- **Password Caching**: Temporarily caches your session for 15 minutes of uninterrupted writing.

## Install and basic usage
install with `go install` in your terminal.  
then usage is simple:
```bash
# Initial setup
diary setup # this will prompt for you to create a password

# Write a new entry (or open today's existing one)
diary new

# list, browse, and read entries
diary read

# View statistics
diary stats

# View other options
diary help
```
