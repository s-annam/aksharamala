# Documentation for Transliteration Scheme JSON Format

## Fields

- **`id`**:
  - Description: A unique identifier for the transliteration scheme.
  - Example: `"example_id"`

- **`name`**:
  - Description: A human-readable name for the transliteration scheme.
  - Example: `"Example Transliteration"`

- **`language`**:
  - Description: The target language of the transliteration scheme.
  - Examples: `"Telugu"`, `"Devanagari"`

- **`scheme`**:
  - Description: The transliteration or encoding scheme used.
  - Examples: `"ITRANS"`, `"Unicode"`

- **`metadata`**:
  - `virama`: Represents the virama character and its operational mode.
    - Examples: `"0x094D, smart"`, `"0x0, normal"`
  - `font_name`: The name of the associated font.
  - `font_size`: The font size in points.
  - `icon_enabled` and `icon_disabled`: Icons used for UI purposes.

- **`mappings`**:
  - Description: Defines transliteration rules.
  - Example:
    ```json
    {
      "lhs": "A",
      "rhs": ["\u0905"],
      "context": "vowel_start"
    }
    ```

- **`comments`**:
  - Description: Comments from the original AKT file or added for clarity.
  - Example: `"This is an example comment included for documentation purposes."`
