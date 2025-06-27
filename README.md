# te: Terminal Text Editor

`te` is a simple terminal-based text editor written in Go, inspired by modal editors like Vim. It supports basic editing, navigation, and file operations in both Linux and Windows terminals.

---

## Features

- Open and edit text files in the terminal
- Modal editing: Command mode and Edit mode
- Basic navigation (arrow keys, Home/End)
- Insert, delete, and append lines
- Save changes to file
- Status bar with file and cursor info

---

## Building

### Prerequisites

- [Go](https://golang.org/dl/) installed

### Build for Linux and Windows

Run the provided build script:

```bash
./build.sh
```

This will generate:
- `te` (Linux binary)
- `te.exe` (Windows binary)

---

## Usage

```bash
./te [filename]
```

If no filename is provided, it opens or creates `text.txt`.

---

## Controls

### Modes

- **Command Mode** (default): For navigation and commands
- **Edit Mode**: For inserting text

### Switching Modes

- Press [`i`](te.go ) to enter **Edit Mode**
- Press `Esc` to return to **Command Mode**

### Editing

- **Insert text**: In Edit Mode, type normally
- **Delete character**: `Backspace`
- **Insert new line**: `Enter`
- **Insert tab**: `Tab` (inserts spaces)

### Navigation

- **Arrow keys**: Move cursor
- **Home**: Go to first line
- **End**: Go to last visible line

### File Operations

- **Save**: Press `w` in Command Mode
- **Quit**: Press `Ctrl+C`

---

## Status Bar

Displays:
- Current mode (EDIT/CMD)
- Filename and line count
- Save status (modified/saved)
- Cursor position

---

## Notes

- The editor uses [termbox-go](https://github.com/nsf/termbox-go) for terminal UI.
- Some features (undo, copy/paste) are placeholders for future development.

---

## License

MIT License

---

**Enjoy editing in your terminal!**
