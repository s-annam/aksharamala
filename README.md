# aks.go - Transliteration for Indian Languages

## Overview
`aks.go` is a modern, extensible transliteration system for Indian languages, reimagining the original Aksharamala project (2000-2005) in Go. It provides developers with a robust, portable, and cloud-compatible toolkit to handle transliteration across languages and scripts efficiently.

## Purpose
Indian language data often exists in legacy formats or custom encodings, making it inaccessible to modern applications. `aks.go` bridges this gap by providing tools to convert and transliterate such data into standardized Unicode representations. This empowers developers to preserve and modernize historical content while supporting future innovations.

## Features
- **Transliteration Engine**:
  - Support for customizable mappings.
  - Handles multiple languages and scripts.
- **Legacy Compatibility**:
  - Converts old Aksharamala (`.akt`) files into JSON (`.aksj`) while preserving comments and structure.
- **Extensibility**:
  - JSON-driven configuration for easy customization.
  - Modular design enables new languages and scripts to be added seamlessly.
- **Smart Processing**:
  - Intelligent virama handling (with support for various modes that are helpful for Indic).
  - Optional logging and verbose modes for debugging.

## Quick Start
### Prerequisites
- Go 1.20+
- A working environment for Go projects.

### Installation
```bash
# Clone the repository
$ git clone https://github.com/s-annam/aks.go.git

# Navigate to the directory
$ cd aks.go

# Build the project
$ go build ./cmd/aksharamala
```

### Usage
To convert old `.akt` files to `.aksj` format:
```bash
go run ./cmd/akt_converter convert -input myfile.akt -output myfile.aksj
```
For dry-run mode:
```bash
go run ./cmd/akt_converter convert -input myfile.akt -output myfile.aksj -dry-run
```

## Architecture
1. **Transliteration Core**:
   - Parses and processes mappings.
   - Handles transliteration with contextual awareness.
2. **Keymap Store**:
   - Manages transliteration schemes and metadata.
   - Validates mappings and supports efficient lookup.
3. **Utilities**:
   - Supports metadata extraction and comment normalization.
4. **Logger**:
   - Provides configurable logging (debug/production modes).

## Roadmap
### Completed
- Basic `.akt` to `.aksj` conversion.
- Transliteration with configurable mappings.

### Future Enhancements
* **AI-Driven AKSJ Creation:** Implement machine learning models to analyze .akt files and generate .aksj mappings intelligently, reducing manual effort.
* **Dynamic Language Support:** Add runtime language detection and mapping support to handle multi-script input dynamically.
* **API Integration:** Develop REST APIs for transliteration tasks, enabling integration with other tools and systems.

See the the list of issues and planned enhancements in the Issues section that are more immediate. Please reaach out if you are interested in contributing.

## History
`aks.go` builds upon the legacy Aksharamala project, originally developed as an Indic transliteration system for Windows 2000/XP. By transitioning to Go, `aks.go` modernizes the original concepts and makes them accessible via cloud-based APIs and modern development environments.

For a detailed history, see [history.md](docs/history.md).

## Contributing
Contributions are welcome! Please read [CONTRIBUTING.md](docs/CONTRIBUTING.md) for details on our process.

## License
This project is licensed under the GNU Affero General Public License (AGPL-3.0-or-later). See [LICENSE](LICENSE) for details.

## Acknowledgments
Special thanks to the community and contributors of the original Aksharamala project for inspiring this modernization effort.
