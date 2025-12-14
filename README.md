# Secrets Catch

A lightweight Go tool designed to scan directories for sensitive configuration files (like `.env`, `secrets.json`, `config.yaml`) and securely back them up to a specified output location.

## Why use this?

**Scrap all the envs and put them in a secure location so you don't delete them.**

When cleaning up projects, refactoring, or moving environments, it's easy to accidentally delete local configuration files that aren't checked into version control. This tool helps you aggregate all those scattered secrets into one folder for safekeeping.

## Features

- **Recursive Scan:** Scans the target directory and subdirectories.
- **Pattern Matching:** Uses glob patterns to find specific files.
- **Defaults:** Comes with default patterns for common secret files:
  - `**/.env`
  - `**/config.yaml`
  - `**/secrets.json`
  - `**/firebase-*.json`
- **Customizable:** Supports custom ignore (`-i`) and accept (`-a`) patterns.
- **Structure Preservation:** Maintains the original directory structure in the output folder.

## Usage

```bash
go run main.go -t <target_directory> -o <output_directory> [flags]
```

### Flags

- `-t`: Target directory to scan (default: current directory `.`).
- `-o`: Output directory where files will be copied (Required for safe operation).
- `-w`: Overwrite existing files in the output directory (default: false).
- `-i`: Comma-separated glob patterns to ignore (e.g., `**/node_modules/**,**/.git/**`).
- `-a`: Comma-separated glob patterns to accept. Overrides the default list.

### Example

Scan the current directory for default secret files and save them to a `backup_secrets` folder:

```bash
go run main.go -t . -o ./backup_secrets -i "**/node_modules/**"
```

## ⚠️ Warnings and Alerts

### 🔒 Security Alert

**This tool handles highly sensitive information.**

- **Secure the Output:** Ensure the output directory (`-o`) is stored in a secure location. Do not push the output directory to a public repository (add it to `.gitignore`).
- **Permissions:** Verify that the output directory has restricted file permissions so that only authorized users can access the backed-up secrets.

### ⚠️ Operational Warnings

- **Overwrite Risk:** If the output directory already contains files with the same names and paths, they will be overwritten without confirmation.
- **Disk Usage:** Running this on very large directories without proper ignore patterns (like `node_modules` or `vendor`) might result in scanning a large number of files.
- **Hidden Files:** By default, it looks for `.env`, which is a hidden file. Ensure your glob patterns account for hidden files if you provide custom ones.
