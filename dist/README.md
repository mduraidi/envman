# envman

A portable SDK environment manager for Windows (Go implementation).

## Features
- Manage and activate portable SDK environments for dotnet, java, python, nodejs, golang
- No installation required for SDKs (xcopy/portable only)
- Supports global and per-project environments
- Generates activation scripts for PowerShell and CMD
- Extensible: add new SDKs/versions by dropping them in the SDKs folder

## Prerequisites
- Windows 11 (or Windows 10)
- [Go 1.22+](https://go.dev/dl/) (for building the tool)
- Portable SDKs for each language (see below)
- Git (optional, for template .gitignore)

## Setup
1. **Clone or download this repository**
2. **Prepare SDKs:**
   - Download portable/xcopy versions of each SDK you want to use (dotnet, python, nodejs, golang, java)
   - Place them in the following folder structure (edit `--sdk-root` if you use a different path):
     ```
     envman_sdks/
       dotnet/8.0.100/
       python/3.12.1/
       nodejs/20.10.0/
       golang/1.22.3/
       java/17.0.2/
     ```
   - Each version folder should contain the SDK's `bin` directory and all required files for that SDK to run portably.
3. **Build envman:**
   - Open a terminal in the `envman` folder
   - Run:
     ```
     go build -o envman.exe
     ```
   - This will produce `envman.exe` in the current directory.
4. **(Optional) Add envman.exe to your PATH**
   - So you can run `envman` from any folder.

## Usage
### Initialize a new environment
```
cd myproject
envman init
```
- Creates `envman.json` in the current folder with the latest SDK versions found in your SDKs root.
- Optionally copies template files (e.g., `.gitignore`, `README.md`) if present in the `templates` folder.

### List available SDKs and versions
```
envman list
```
- Shows all SDKs and versions found in your SDKs root.

### Interactively select SDK versions
```
envman select
```
- Prompts you to choose a version for each SDK (default is latest).
- Updates `envman.json`.

### Generate activation scripts
```
envman activate
```
- Generates `activate.bat` (CMD) and `activate.ps1` (PowerShell) in the current folder.
- These scripts set up your shell to use the selected SDKs.
- Run the appropriate script in your shell to activate the environment.

### Deactivate the environment
```
envman deactivate
```
- (Currently, just a reminder. To fully deactivate, close your shell or manually restore your PATH.)

### Switch between local and global environments
```
envman use local   # Use the envman.json in the current folder
envman use global  # Use the global envman.json (see below)
```
- Copies the selected config to `envman.json` in the current folder.

## Global Environment
- To set up a global environment, create a folder `%USERPROFILE%\.envman` and place an `envman.json` there.
- Use `envman use global` to switch to it in any project.

## Customizing SDKs Root
- By default, envman looks for SDKs in `../envman_sdks` relative to your project.
- You can override this with the `--sdk-root` flag:
  ```
  envman init --sdk-root "D:/sdks"
  envman list --sdk-root "D:/sdks"
  # etc.
  ```

## Adding New SDKs or Versions
- Download the portable SDK and extract it to the appropriate folder under your SDKs root.
- The folder name should match the language and version (e.g., `dotnet/8.0.100/`).
- Run `envman list` to verify it is detected.

## Troubleshooting
- If a command fails, check the error message for details.
- Make sure your SDKs are portable and have a `bin` directory.
- If activation scripts do not work, check your shell and PATH settings.

## Example Workflow
```
# Prepare SDKs in envman_sdks/ as described above
cd myproject
envman init
envman select   # (optional, to pick specific versions)
envman activate
# In your shell, run: .\activate.ps1  (PowerShell) or activate.bat (CMD)
# Now your shell uses the selected SDKs
```

## Uninstallation
- To remove envman, simply delete the `envman.exe` file and any generated scripts/configs.

## License
MIT

---
