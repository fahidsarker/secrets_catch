# Secrets Catch

A lightweight Go tool designed to scan directories for sensitive configuration files (like `.env`, `secrets.json`, `config.yaml`) and securely back them up to a password-protected zip file.

## Why use this?

**Scrap all the envs and put them in a secure location so you don't delete them.**

When cleaning up projects, refactoring, or moving environments, it's easy to accidentally delete local configuration files that aren't checked into version control. This tool helps you aggregate all those scattered secrets into one encrypted archive for safekeeping.

## Features

- **Recursive Scan:** Scans the target directory and subdirectories.
- **Pattern Matching:** Uses glob patterns to find specific files.
- **Defaults:** Comes with default patterns for common secret files:

  - `**/.env*`
  - `**/config.yaml`
  - `**/config.json`
  - `**/secrets.*`
  - `**/firebase-*.json`
  - `**/*.pem`
  - `**/*.key`
  - `**/id_rsa*`
  - `**/credentials.json`

  And default ignore patterns:

  - `**/.git/**`
  - `**/node_modules/**`
  - `**/vendor/**`
  - `**/.idea/**`
  - `**/.vscode/**`
  - `**/dist/**`
  - `**/build/**`
  - `**/*.log`

- **Customizable:** Supports custom ignore (`-i`) and accept (`-a`) patterns.
- **Encrypted Backup:** Creates a password-protected zip file (AES-256) containing all found secrets.
- **Structure Preservation:** Maintains the original directory structure in the output zip.

## Usage

```bash
go run main.go -t <target_directory> -o <output_zip_file> -p <password> [flags]
```

### Flags

- `-t`: Target directory to scan (default: current directory `.`).
- `-o`: Output zip file path (Required).
- `-p`: Password for the zip file (Required).
- `-i`: Comma-separated glob patterns to ignore. **Overrides the default list.**
- `-a`: Comma-separated glob patterns to accept. **Overrides the default list.**

### Example

Scan the current directory for default secret files and save them to `secrets_backup.zip` with password `mysecretpass`:

```bash
go run main.go -t . -o ./secrets_backup.zip -p mysecretpass
```

## ⚠️ Warnings and Alerts

### 🔒 Security Alert

**This tool handles highly sensitive information.**

- **Strong Password:** Ensure you use a strong password for the zip file.
- **Secure Storage:** Store the output zip file in a secure location. Do not push it to a public repository.

### ⚠️ Operational Warnings

- **Overwrite Risk:** If the output file already exists, it will be overwritten only if `-w` is set.
- **Disk Usage:** Running this on very large directories without proper ignore patterns (like `node_modules` or `vendor`) might result in scanning a large number of files.
- **Hidden Files:** By default, it looks for `.env`, which is a hidden file. Ensure your glob patterns account for hidden files if you provide custom ones.
