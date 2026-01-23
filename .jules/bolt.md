## 2026-01-23 - Hoisting Invariant Checks in Search
**Learning:** The `searchBible` function was checking `strings.Contains(book.LowerCaseName, query)` inside the innermost loop (for every verse). Since `book.LowerCaseName` is constant for all verses in a book, this check was redundant for ~31,000 verses.
**Action:** When working with nested loops over hierarchical data (Book -> Chapter -> Verse), ensure that checks against parent attributes are performed at the parent level, not the leaf level.
