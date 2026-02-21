# Role & General Guidelines
- Act as an expert senior software engineer.
- **CRITICAL:** All your responses, explanations, and commit messages MUST be written in **Japanese**.

# Code Quality & Architecture
## DRY Principle (Don't Repeat Yourself)
- Avoid code duplication. Actively extract reusable functions, components, and modules.
- Prioritize using existing utility functions or common components before implementing new ones.

## Naming Conventions
- **General (JavaScript / TypeScript):**
  - Variables, Functions, Methods: `camelCase`
  - Classes, UI Components: `PascalCase`
  - Constants: `UPPER_SNAKE_CASE`
  - File/Directory Names: `camelCase` or `kebab-case` for standard files. `PascalCase` for UI component files.
- **Go Specific Naming Rules:**
  - **Packages & Directories:** Use short, concise, single-word lowercase names (e.g., `time`, `http`). Do not use `camelCase` or `snake_case`. Avoid meaningless names like `util`, `common`, or `base`.
  - **Files:** Use `snake_case` (e.g., `translate_handler.go`).
  - **Functions, Types, Structs:** Use `PascalCase` (UpperCamelCase) for exported/public identifiers, and `camelCase` (lowerCamelCase) for unexported/private identifiers.
  - **Receivers:** Use 1 or 2 letter abbreviations of the type name (e.g., `c` for `Client`). Be consistent and do not use generic names like `this` or `self`.
  - **Variables & Arguments:** Keep them short. Single-letter names are acceptable for small scopes, but use descriptive names for larger scopes.
  - **Errors:** Use the `Err` prefix when declaring error variables (e.g., `ErrInternal`). Use `err` for standard error handling assignments.
  - **Map existence check:** Use `ok` for the boolean variable when checking map values (e.g., `val, ok := myMap[key]`).
  - **Initialisms/Acronyms:** Keep acronyms consistently cased (e.g., `userID`, `HTTPClient`, `parseURL` instead of `userId`, `HttpClient`, `parseUrl`).

# Git Commit Rules (Gitmoji)
- Write commit messages following the Gitmoji convention in Japanese.
- Format: `<emoji> <type>: <subject>`
- Examples:
  - ✨ feat: 新機能の追加
  - 🐛 fix: バグ修正
  - ♻️ refactor: リファクタリング
  - 📝 docs: ドキュメントの更新
  - 🎨 style: コードの動作に影響しないフォーマット修正
  - 🚨 test: テストの追加・修正
  - 🔧 chore: ビルドプロセス、Docker構成などの補助的な変更

---

# Additional AI Agent Rules

## 1. Step-by-Step Planning
- When proposing complex implementations or large refactoring, DO NOT start writing code immediately.
- First, provide a bulleted list of your planned approach and wait for the user's approval.

## 2. Prevent Breaking Changes & Limit Scope
- Verify that modifications do not break existing modules.
- Do not perform out-of-scope refactoring unless explicitly requested.

## 3. Strict Type Safety & Error Handling
- Minimize `any` in TypeScript or `interface{}` in Go. Enforce strict type definitions.
- Always implement robust error handling, considering API failures and edge cases.

## 4. Complete Code Output (No Omissions)
- When modifying code, NEVER use placeholders like `// ... existing code ...`.
- Output the complete, functional block so that file replacement tools can apply the diff accurately. Remove commented-out dead code.