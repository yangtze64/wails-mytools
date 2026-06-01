# MyTools

MyTools is a desktop toolbox application built with Wails 3 + Go + Vue 3, integrating a variety of daily development and text processing tools. All data is stored locally — no internet connection required.

[中文](README.md)

## Features

### Text Tools

| Tool | Description |
|------|-------------|
| Text Dedupe | Remove duplicate lines |
| Random String | Generate passwords and tokens |
| Text Diff | View added and removed lines |
| JSON Tool | Format, minify, and validate JSON |
| Markdown Editor | Live editing and preview |
| Base64 Codec | Encode and decode text and files |

### Developer Tools

| Tool | Description |
|------|-------------|
| Timestamp Converter | Convert between time and Unix timestamp |
| Cron Expression Generator | Visual cron expression builder |
| JWT Parser | Decode JSON Web Tokens |
| URL Codec | URL encoding, decoding, and parsing |
| Bookmark Manager | Save and open frequently used websites |
| Secret Manager | Locally encrypt and store passwords and keys |

### Image Tools

| Tool | Description |
|------|-------------|
| QR Code Generator | Convert text and links to PNG |
| QR Code Decoder | Extract QR code content from images |
| Image Compressor | Compress images with adjustable quality |

### Diagram Tools

| Tool | Description |
|------|-------------|
| Mermaid Editor | Live preview of Mermaid diagrams |
| Mind Map | Generate mind maps from indented text |
| Markmap | Generate mind maps from Markdown |

## Getting Started

### Prerequisites

- Go 1.25+
- Node.js 18+
- pnpm (recommended via Corepack: `corepack enable`)
- Wails 3 CLI (`go install github.com/wailsapp/wails/v3/cmd/wails3@latest`)

### Installation

```sh
git clone https://github.com/your-username/wails-mytools.git
cd wails-mytools
cd frontend && pnpm install && cd ..
task dev
```

Alternatively, use the Wails CLI directly:

```sh
wails3 dev -config ./build/config.yml
```

## Common Commands

```sh
task dev          # Start development mode
task build        # Build the application
task package      # Package the application (generate installer)
```

Frontend commands (run inside `frontend/`):

```sh
pnpm run dev      # Start frontend dev server
pnpm run build    # Build frontend for production
```

## Directory Structure

```
wails-mytools/
├── main.go                    # Wails application entry point
├── internal/
│   └── services/              # Go backend services (exported to frontend)
├── frontend/
│   ├── src/
│   │   ├── api/               # Wails bindings wrapper layer
│   │   ├── components/        # Reusable Vue components
│   │   ├── composables/       # Vue composables
│   │   ├── config/            # Tool registry and other configs
│   │   ├── styles/            # Global styles
│   │   ├── types/             # Shared TypeScript types
│   │   └── views/             # Page-level Vue components
│   ├── bindings/              # Auto-generated Wails TypeScript bindings
│   └── package.json
├── build/
│   ├── config.yml             # Wails 3 build configuration
│   ├── Taskfile.yml           # Common build tasks
│   ├── darwin/                # macOS build assets
│   ├── windows/               # Windows build assets
│   └── linux/                 # Linux build assets
├── Taskfile.yml               # Root Taskfile
├── .github/
│   └── workflows/             # GitHub Actions workflows
│       ├── ci.yml             # CI checks (type check, lint)
│       └── release.yml        # Auto build & release
└── README.md
```

## Development Conventions

- After adding a Go service, run `wails3 generate bindings -clean=true -ts` to update frontend bindings.
- Use pnpm for all frontend dependency management, development, and builds — do not use npm.
- Do not import `frontend/bindings` directly in pages; wrap calls in `frontend/src/api/` first.
- Keep the structure simple while the number of tools is small; introduce Vue Router only when multi-page navigation, history, or deep linking is needed.
- Introduce Pinia only when shared state becomes complex; avoid premature global state management.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Framework | Wails 3 |
| Backend | Go 1.25+ |
| Frontend | Vue 3 + TypeScript |
| UI Library | Element Plus |
| Package Manager | pnpm |
| Build Tool | Vite 5 |
| Diagrams | Mermaid |
| Mind Map | mind-elixir / markmap |
| QR Code | qrcode / jsqr |

## License

MIT
