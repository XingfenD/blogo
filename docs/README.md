# Blogo -- A simple blog engine by go


[![License](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Report Card](https://goreportcard.com/badge/github.com/XingfenD/blogo)](https://goreportcard.com/report/github.com/XingfenD/blogo)

Blogo is a minimalist blog engine built with Go, featuring Markdown support and SQLite database.

## Features

- 📝 Markdown Articles
- 🏷️ Categories & Tags System
- 📆 Timeline Archiving
- 🎨 Responsive Theme
- ⚡ Fast Rendering
- 🔒 File-based Storage

## Quick Start

### Prerequisites
- Go 1.24+
- SQLite3

### Installation
```bash
# Clone repository
git clone https://github.com/XingfenD/blogo.git

# Enter project
cd blogo

# Install dependencies
go mod tidy

# Start server
go run main.go
```

## Configuration

Edit `config.toml`:

```toml
[basic]
port2listen = 8080         # Server port
base_url = 'http://localhost:8080' # Site URL
root_path = 'website'      # Resource directory

[user]
name = "Your Name"         # Author name
avatar_url = "/img/avatar.png" # Avatar path
description = "Personal Blog" # Site description

# See config_example.toml for more options
```

## Project Structure
```plaintext
blogo/
├── website/             # Frontend
│   ├── template/        # HTML templates
│   ├── static/          # Static assets
│   └── data/            # Database
├── module/              # Go modules
│   ├── router/          # Routing
│   ├── sqlite/          # Database
│   └── tpl/             # Templates
└── config.toml          # Configuration
```

## Tech Stack

- Backend: Go 1.24
- Database: SQLite3
- Templating: Go html/template
- Markdown: Blackfriday
- Frontend: HTML5/CSS3

## License

Licensed under the [Mozilla Public License 2.0](https://opensource.org/licenses/MPL-2.0)

This project utilizes some icons from the ByteDance Icon Library.

## Todo

- [ ] Implement admin page
- [ ] Document complete
