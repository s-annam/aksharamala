# aks_converter

## Overview
The `aks_converter` tool provides powerful capabilities for converting `.akt` files that were built for Aksharamala Windows 2000/xp IME into `.aksj` format to be used with the Aksharamala Go project.

## Features
- Convert `.akt` files to `.aksj` format with normalized and structured output.
- Preserve existing comments and mappings during conversion.

## Future Enhancements
1. **Update-Only Mode**:
   - Update existing entries in-place, keeping their sections intact.
   - Add new mapping entries to appropriate sections.
   - Enable verbose output by default to log significant actions.
   - All the other sections and the metadata to be fully replaced.

2. **Verbose Mode**:
   - Log detailed information about sections, entries, and significant events.

3. **Dry-Run Mode**:
   - Simulate conversion without making changes, allowing users to review actions.

4. **Custom Section Mapping**:
   - Enable rules for mapping entries to specific sections based on prefixes or patterns.

5. **Interactive Mode**:
   - Prompt users to confirm adding or overwriting entries during conversion.

6. **Error Logging**:
   - Maintain a log file for warnings and errors encountered during conversion.

7. **Preserve History from File Comments**:
   - Extract historical information from AKT file comments (one line per update).
   - Append this history to the end of the `.aksj` file under a dedicated "history" section.
   - Clearly indicate that this history is sourced from the original AKT file.

8. **Skip Empty Section**:
   - Ignore any empty sections and skip to create them.
