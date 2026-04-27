# envchain-cli

> A CLI tool to manage and inject environment variables per project using encrypted local storage.

---

## Installation

```bash
go install github.com/yourusername/envchain-cli@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/envchain-cli.git
cd envchain-cli
go build -o envchain .
```

---

## Usage

**Set a variable for a project namespace:**

```bash
envchain set myproject API_KEY=supersecret DB_PASSWORD=hunter2
```

**Run a command with injected environment variables:**

```bash
envchain run myproject -- go run main.go
```

**List stored namespaces:**

```bash
envchain list
```

**Remove a variable:**

```bash
envchain unset myproject API_KEY
```

Variables are encrypted at rest using AES-256 and stored locally in `~/.envchain/store`.

---

## How It Works

1. Variables are stored encrypted in a local keystore file.
2. When you run a command with `envchain run`, variables are decrypted and injected into the subprocess environment.
3. Nothing is ever written to `.env` files or shell history.

---

## License

MIT © [yourusername](https://github.com/yourusername)