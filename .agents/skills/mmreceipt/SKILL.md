```markdown
# mmreceipt Development Patterns

> Auto-generated skill from repository analysis

## Overview
This skill teaches the core development patterns and conventions used in the `mmreceipt` Go repository. You'll learn about file naming, import/export styles, commit message practices, and how to structure and run tests. This guide also provides suggested commands to streamline common workflows.

## Coding Conventions

### File Naming
- Use **camelCase** for file names.
  - Example: `receiptParser.go`, `userManager.go`

### Import Style
- Use **relative imports** within the project.
  - Example:
    ```go
    import "../utils"
    ```

### Export Style
- Use **named exports** for functions, types, and variables.
  - Example:
    ```go
    // In receiptParser.go
    func ParseReceipt(data string) (*Receipt, error) {
        // implementation
    }
    ```

### Commit Messages
- **Freeform** style, no strict prefixes.
- Average commit message length: ~35 characters.
  - Example: `fix bug in receipt parsing logic`

## Workflows

### Add a New Feature
**Trigger:** When implementing a new functionality  
**Command:** `/add-feature`

1. Create a new camelCase-named file if needed (e.g., `featureName.go`).
2. Write your Go code, using relative imports for internal packages.
3. Export functions/types that need to be accessed elsewhere using named exports.
4. Write or update corresponding test files (see Testing Patterns).
5. Commit your changes with a clear, concise message.

### Fix a Bug
**Trigger:** When resolving a bug in the codebase  
**Command:** `/fix-bug`

1. Locate the relevant file(s) using camelCase naming.
2. Apply the fix, maintaining code style and import conventions.
3. Update or add tests to cover the bug fix.
4. Commit with a descriptive message about the fix.

### Write or Update Tests
**Trigger:** When adding or modifying tests  
**Command:** `/test`

1. Create or update a test file matching the `*.test.*` pattern (e.g., `receiptParser.test.go`).
2. Write test cases for your functions.
3. Run tests using Go's testing tools (e.g., `go test`).
4. Ensure all tests pass before committing.

## Testing Patterns

- **Test File Naming:**  
  Test files follow the `*.test.*` pattern, such as `receiptParser.test.go`.
- **Framework:**  
  The specific testing framework is unknown, but standard Go testing practices are likely.
- **Example Test:**
  ```go
  // receiptParser.test.go
  package main

  import "testing"

  func TestParseReceipt(t *testing.T) {
      // test implementation
  }
  ```
- **Running Tests:**  
  Use Go's built-in test command:
  ```
  go test ./...
  ```

## Commands
| Command      | Purpose                                 |
|--------------|-----------------------------------------|
| /add-feature | Scaffold and implement a new feature    |
| /fix-bug     | Guide through fixing a bug              |
| /test        | Write, update, and run tests            |
```
