<div align="center">
  <img src="https://github.com/user-attachments/assets/6d71dac4-88e2-43e6-824a-fd1e48f2031e" alt="morama logo" width="80" />
  <h1 style="margin-top: 0em;">morama</h1>
  <p><em>A CLI tool for managing your watched movies and dramas</em></p>
  <p>
    <img src="https://img.shields.io/badge/개발언어-Go-00ADD8?style=for-the-badge&logo=go" alt="Go" />
    <img src="https://img.shields.io/badge/개발기간-2024.05~2024.06-9E7B6B?style=for-the-badge" alt="개발기간" />
  </p>
</div>


## Overview

**morama** is a simple command-line interface (CLI) application for recording and managing your personal reviews and ratings for movies and TV dramas. Built with Go, morama helps you keep track of what you've watched and how you felt about it — all from your terminal.

morama lets you:
- Record reviews and star ratings for movies and dramas.
- Browse and search your personal watch history.
- Categorize and filter entries by title, genre, or rating.
- Easily update or delete entries as your opinions evolve.

<br>

## Install

### Using Homebrew (macOS/Linux)

```bash
# 1. Add the tap
brew tap kiku99/morama https://github.com/kiku99/morama

# 2. Install morama
brew install morama
```

### Manual Installation

1. Download the latest release from [GitHub Releases](https://github.com/kiku99/morama/releases)
2. Extract and move the binary to your PATH

<br>

## Usage
새로운 감상 기록 추가
```
morama add "inception" --movie
```

모든 기록 조회
```
morama list
```

특정 제목 상세 조회
```
morama show "inception" --movie
```

기록 수정
```
morama edit "inception" --movie
```

기록 삭제
```
morama delete --id=3
```

### 주석 추가




