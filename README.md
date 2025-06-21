# 🔁 reloadr (Hot Reload CLI for Go, Node.js, and More)

A blazing fast, minimal hot-reload tool for developers — written in Go.  
Think of it as `nodemon`, but faster, simpler, and made for modern CLI workflows.

---

## 🚀 Features
- ✅ Watch file changes and auto-restart your app
- ✅ Supports any language: Go, Node.js, PHP, Python, etc.
- ⏱️ Easy setup with `reloadr.toml` (pending)
- ✅ Debounce built-in — no duplicate reloads
- ✅ Cross-platform: Linux, macOS, Windows
- ✅ Blazing fast — powered by Go’s native file system watchers

---

## 📦 Installation

```bash
go install github.com/testdrivenhuman/reloadr@latest
```

Make sure $GOPATH/bin is in your PATH:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```


---

## 🧪 Usage
```bash
reloadr
```
```bash
reloadr --exec "go run main.go"
```
```bash
reloadr --exec "node server.js" --ext ".js,.ts" --exclude ".git,node_modules"
```

---

## ⚙️ Config File (pending)

You can create a `hotreload.toml` in your root directory:

```bash
exec     = "go run main.go"
watch    = "."
exclude  = [".git", "node_modules", "vendor"]
ext      = [".go"]
debounce = 300
```

Now simply run:

```bash
reloadr
```

---

## 🧠 Example Use Cases

| Language | Command |
|----------|---------|
| **Go**     | `reloadr --exec "go run main.go"` |
| **Node.js**| `reloadr --exec "node index.js" --ext ".js,.ts"` |
| **Python** | `reloadr --exec "python app.py"` |
| **Docker** | `reloadr --exec "docker-compose up --build"` |

---

## 🔧 Flags

| Flag         | Description                             | Default              |
|--------------|-----------------------------------------|----------------------|
| `--exec`     | Command to run                          | `go run main.go`     |
| `--watch`    | Directory to watch                      | `.`                  |
| `--exclude`  | Comma-separated folders to ignore       | `.git,node_modules`  |
| `--ext`      | Extensions to watch (e.g. `.go`, `.js`) | `.go`                |
| `--debounce` | Delay before reload (in milliseconds)   | `300`                |
---

## 🤝 Contributing

PRs welcome! Feel free to suggest features like:
- Live browser reload (web dev)
- Pre/post command hooks
- Notify on crash
- Test runner integration

---

## ⚡️ License

This project is open source and freely available under the MIT license.
You are free to use, modify, and distribute it as you wish.
Pull requests and contributions are welcome!