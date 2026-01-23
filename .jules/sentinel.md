## 2026-01-23 - Source Code Disclosure in Go FileServer
**Vulnerability:** The application was serving static files using `http.FileServer(http.Dir("./"))` from the project root. This exposed the entire source code (`main.go`), configuration (`go.mod`), and version control history (`.git/`) to any user who requested them.
**Learning:** The Go standard library `http.FileServer` does not filter files by extension or hidden status (dotfiles) by default. It faithfully serves whatever is in the directory.
**Prevention:** Always wrap `http.FileServer` with a middleware that enforces a whitelist of allowed file extensions (e.g., `.html`, `.css`, `.js`, images) and explicitly blocks sensitive files or directories.
