# Aksharamala (aks.go)

Aksharamala is a versatile transliteration framework designed to convert text 
across multiple languages and scripts using customizable keymaps. This project 
builds upon the legacy of the original Aksharamala, introduced in 1989 as an 
undergraduate project and later expanded in 2000 with a C++ implementation.

This Go-based implementation modernizes the transliteration process with:
- JSON-based keymaps for flexibility and compatibility.
- Command-line tools for managing and applying transliteration schemes.

## License

This project is licensed under the GNU Affero General Public License (AGPL),
version 3 or later. See the [LICENSE](./LICENSE) file for details.

## History

- **1989**: Original Aksharamala developed as an undergraduate project. This
  was written in Turbo Pascal, to generate characters on a CRT and dot-matrix
  printers to type in Telugu. The solution was certainly matching that era.
- **2000-2003**: Built a brand new C++ version for Windows as a keyboard hook
  (or colloquially referred to as IME, Input Method Editor) and sold under the
  banner Deshweb.com Pvt. Ltd.
- **2025**: Reimagined as an open-source Go project with a focus on modularity
  and modern workflows. Largely reusing the concepts of transliteration while
  taking support from ChatGPT, Claude and Gemini to develop code using LLMs.
  The architecture otherwise has been evolved (for example, taking AKT text
  files into a JSON based format etc.) to fit today's needs.

### Contributors
Special acknowledgment to collaborators and contractors who contributed to the
development of the two original Aksharamala projects.

## Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/s-annam/aks.go.git
   cd aks.go
   ```

2. Build the tools:

   ```bash
    go build ./cmd/aksharamala
    go build ./cmd/akt_converter
    ```

3. Run the tools:

   ```bash
    ./aksharamala
    ./akt_converter
    ```

## Commands

- **aksharamala**: Main transliteration tool for converting text using keymaps.
- **akt_converter**: Converts `.akt` files to `.aksj` JSON format.

---

#### **3. Directory Structure**
Explain the directory structure for better navigation:

## Directory Structure

```graphql
    cmd/ 
        ├── aksharamala/ 
        │       └── main.go # Main transliteration command
        ├── akt_converter/ 
        │       └── main.go # AKT-to-JSON converter command
    docs/ 
        └── akt_converter.md # Documentation for AKT-to-JSON converter
    examples/ 
        ├── example.akt # Example AKT file
        └── example.aksj # Example JSON keymap file
    internal/ 
        ├── akt_converter/ 
        │       └── utils.go # Utility functions for AKT conversion 
        ├── keymap/ 
        │       ├── keymap_store.go # Keymap management logic
        │       └── transliteration.go # Core transliteration logic 
        ├── types/ 
        │       └── scheme.go # Shared types like TransliterationScheme
    keymaps/ 
        ├── hindi.aksj # Hindi JSON keymap 
        ├── telugu.aksj # Telugu JSON keymap 
        └── marathi.aksj # Marathi JSON keymap
```

## Usage

### aksharamala
Run the main transliteration tool with input text:

    ```bash
    ./aksharamala -keymap hindi -input "example text"
    ```

### akt_converter
Convert a .akt file to .aksj format:

    ```bash
    ./akt_converter input.akt output.aksj
    ```

---

#### **5. Contribution Guidelines**
Encourage contributions and provide instructions:

```markdown
## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bug fix:
   ```bash
   git checkout -b my-feature
   ```
3. Commit your changes
    ```bash
    git commit -m "Add my feature"
    ```
4. Push to your fork and submit a pull request.
For major changes, please open an issue to discuss your proposal.
Let me know if you’d like further refinements or additions! 😊
