## **AKT-to-JSON Converter Documentation**

---

### **Overview**
The AKT-to-JSON Converter processes AKT files, a structured format for transliteration mappings, and outputs structured JSON. This tool supports regular sections, pseudo-sections, multiple LHS entries, and metadata extraction.

---

### **Features**
1. **Metadata Parsing**:
   - Extracts metadata from AKT files, such as `#id`, `#name`, `#language`, etc.
   - Metadata is included at the top level of the JSON.

2. **Section Handling**:
   - Supports regular sections (e.g., `#vowels#`) and pseudo-sections (e.g., `//=*= vowels =*=`).
   - Pseudo-sections are treated as real sections and dynamically mapped to categories (e.g., `vowels-no-virama`).

3. **Multiple LHS Entries**:
   - LHS-only lines are attached to the previous mapping as additional LHS entries.

4. **`#end` Marker**:
   - Parsing stops when the `#end` marker is encountered.
   - Lines following `#end` are ignored.

5. **Robust JSON Output**:
   - Optional metadata fields are omitted if empty.
   - Comments from AKT files are included as file-level or section-level comments.

---

### **Input Format**
#### **Metadata**
Metadata lines start with `#` and contain key-value pairs:

```
#id = example_id#
#name = Example Transliteration#
```

#### **Regular Sections**
Sections are defined with a header like `#vowels#`:

```
#vowels#

अ   a[v]
आ   A[v]
```

#### **Pseudo-Sections**
Pseudo-sections are defined within comments:

```
#others#

// =*= vowels =*=
अ   a[v]
आ   A[v]

// =*= DIGITS =*=
०   0 
१   1
```

#### **Multiple LHS Entries**
Lines with only an LHS (no RHS) are attached to the previous mapping:

```
aa  0x0906 0x093E
A
i   0x0907 0x093F
```

#### **End Marker**
   - Parsing stops when the `#end` marker is encountered.

```
#end
```

#### **Robust JSON Output**:
- Optional metadata fields are omitted if empty.

---

### **Common Pitfalls**
- Ensure that metadata lines are correctly formatted.
- Be cautious with section markers.

### **Output Format**
#### **Top-Level Fields**
- `id`: Unique identifier for the transliteration scheme.
- `name`: Human-readable name of the scheme.
- `categories`: Contains mappings grouped by sections.

#### **Categories**
Each section or pseudo-section is represented as a category:

```json
"categories": {
  "vowels": {
    "mappings": [
      { "lhs": ["aa", "A"], "rhs": ["0x0906", "0x093E"] },
      { "lhs": ["i"], "rhs": ["0x0907", "0x093F"] }
    ]
  }
}
```

### Examples
#### Input AKT

```
#id = example_id#
#name = Example Transliteration#

#vowels#
अ   a[v]
आ   A[v]

// =*= DIGITS =*=
०   0
१   1

#end

#others#
.   0x0964
```

#### Generated JSON

```json
{
  "id": "example_id",
  "name": "Example Transliteration",
  "categories": {
    "vowels": {
      "mappings": [
        { "lhs": ["अ"], "rhs": ["a[v]"] },
        { "lhs": ["आ"], "rhs": ["A[v]"] }
      ]
    },
    "digits": {
      "mappings": [
        { "lhs": ["०"], "rhs": ["0"] },
        { "lhs": ["१"], "rhs": ["1"] }
      ]
    }
  }
}
```