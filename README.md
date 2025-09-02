<div align="center">
  <img src="https://github.com/user-attachments/assets/28ad7854-f6a7-4cd1-b904-4d7d72251f16" alt="morama logo" width="175" />
  <p><em>A CLI tool for managing your watched movies and dramas</em></p>
  <p style="margin: 0; line-height: 1;">
    <img src="https://img.shields.io/badge/Built_with-Go-00ADD8?style=for-the-badge&logo=go" alt="Go" />
    <img src="https://img.shields.io/badge/Database-SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white" alt="SQLite" />
    <img src="https://img.shields.io/badge/Development-2025.05~ing-9E7B6B?style=for-the-badge" alt="Development Period" />
  </p>
</div>

<br>

## Overview

**morama** is a simple command-line interface (CLI) application for recording and managing personal reviews and ratings for movies and dramas. Built in Go, it helps you keep track of what you’ve watched — and how you felt about it — all from your terminal.

![morama demo](assets/morama-demo.gif)

### Features
- Add reviews and star ratings for movies and dramas
- Browse and search your viewing history
- Filter by title, genre, or rating
- Edit or delete existing entries
- View yearly breakdowns and rating statistics

<br>

## Installation

### With Homebrew

```bash
# 1. Add the tap
brew tap kiku99/morama https://github.com/kiku99/morama

# 2. Install morama
brew install morama
```

<br>

## CLI Command Structure

```
morama
├── add [title]                   # Add a new entry
│   ├── --movie                   # Add as a movie
│   └── --drama                   # Add as a drama
│
├── list                          # View all records (grouped by year)
│
├── show [title]                  # Show details of a specific entry
│   ├── --movie                   # Specify movie
│   └── --drama                   # Specify drama
│
├── edit [title]                  # Edit an existing entry
│   ├── --id=<ID>                 # Target entry ID (required)
│   ├── --movie                   # Edit as a movie
│   └── --drama                   # Edit as a drama
│
├── delete                        # Delete entries
│   ├── --id=<ID>                 # Delete by ID
│   └── --all                     # Delete all records
│
├── stats                         # Show statistics
│
└── version                       # Show current version
```

<br>

## Examples

**Add a movie**

```bash
morama add "Inception" --movie
```

**Add a drama**

```bash
morama add "Hospital Playlist" --drama
```

**View all records**

```bash
morama list
```

**Show details of a movie**

```bash
morama show "Inception" --movie
```

**Edit a record by ID**

```bash
morama edit "Inception" --id=3 --movie
```

**Delete a record by ID**

```bash
morama delete --id=3
```

**Delete all records**

```bash
morama delete --all
```

**Show statistics**

```bash
morama stats
```

**Show version**

```bash
morama version
```

<br>

## License

This project is licensed under the MIT License.  
See the [LICENSE](LICENSE) file for full details.

<br>

## Contributing

Contributions are welcome!  
Feel free to open issues or submit pull requests to help improve **morama**.
 


<!-- Security scan triggered at 2025-09-02 16:17:35 -->