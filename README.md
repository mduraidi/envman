# envman

A portable, fully generic SDK environment manager for Windows (Go implementation).

## Features
- Manage and activate portable SDK environments for any toolchain (Python, Node.js, Java, Go, .NET, and more)
- No installation required for SDKs (xcopy/portable only)
- Supports global and per-project environments
- Generates activation scripts for PowerShell and CMD
- Extensible: add new SDKs/toolchains by dropping a YAML config in the `toolchains` folder
- All toolchain logic, prompts, and commands are defined in YAML—no Go code changes required

## Prerequisites
- Windows 11 (or Windows 10)
- [Go 1.22+](https://go.dev/dl/) (for building the tool)
- Portable SDKs for each language/toolchain you want to use
- Git (optional, for template .gitignore)

## Setup
1. **Clone or download this repository**
2. **Prepare SDKs:**
   - Download portable/xcopy versions of each SDK you want to use (e.g., Python, Node.js, Go, Java, .NET, etc.)
   - Place them in the following folder structure (edit `--sdk-root` if you use a different path):
     ```
     envman_sdks/
       python/3.12.1/
       nodejs/20.10.0/
       golang/1.22.3/
       java/17.0.2/
       dotnet/8.0.100/
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
- Prompts you to select a toolchain and version, then runs the steps defined in its YAML config.
- Creates `envman.json` in the current folder with your selection.
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
- Prompts you to choose a version for each toolchain (default is latest).
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
- By default, envman looks for SDKs in `envman_sdks` next to the executable.
- You can override this with the `--sdk-root` flag:
  ```
  envman init --sdk-root "D:/sdks"
  envman list --sdk-root "D:/sdks"
  # etc.
  ```

## Adding New Toolchains or SDKs
- To add a new toolchain, create a YAML config in the `toolchains/` folder (see below).
- To add a new SDK version, extract it to the appropriate folder under your SDKs root (e.g., `python/3.12.1/`).
- Run `envman list` to verify it is detected.

## Toolchain YAML Configs (DSL)
All toolchain logic is defined in YAML files in the `toolchains/` folder. Each YAML file describes the steps, prompts, commands, and environment variables for initializing and activating a project.

### YAML Structure
- `name`: Toolchain name (e.g., python, nodejs)
- `steps`: List of steps to run during `envman init` (prompt, select, run, file, message)
- `env_vars`: Mapping of environment variable names to values (for activation)

#### Supported Step Types
- `prompt`: Prompt the user for a value
- `select`: Prompt the user to select from options
- `run`: Run a command (with args)
- `file`: Write a file with given content
- `message`: Print a message
- `when`: (optional) Only run this step if the condition is true (supports `{var}` == "value")

### Example: Python Toolchain
```yaml
name: python
steps:
  - type: prompt
    var: venv_dir
    message: "Enter venv directory"
    default: venv
  - type: run
    command: "{sdk_root}/python/{version}/python.exe"
    args: ["-m", "venv", "{venv_dir}"]
  - type: file
    path: requirements.txt
    content: |
      # Add your dependencies here
  - type: message
    text: "Python venv created in {venv_dir}"
env_vars:
  PYTHON_VERSION: "{version}"
```

### Example: Node.js Toolchain
```yaml
name: nodejs
steps:
  - type: select
    var: pkg_manager
    message: "Choose package manager"
    options: ["npm", "yarn"]
    default: npm
  - type: run
    when: "{pkg_manager}" == "npm"
    command: "{sdk_root}/nodejs/{version}/npm.cmd"
    args: ["init", "-y"]
  - type: run
    when: "{pkg_manager}" == "yarn"
    command: "{sdk_root}/nodejs/{version}/yarn.cmd"
    args: ["init", "-y"]
  - type: message
    text: "Node.js project initialized with {pkg_manager}"
env_vars:
  NODE_VERSION: "{version}"
```

### Example: .NET Toolchain
```yaml
name: dotnet
steps:
  - type: run
    command: "{sdk_root}/dotnet/{version}/dotnet.exe"
    args: ["new", "console"]
  - type: message
    text: ".NET console project created"
env_vars:
  DOTNET_VERSION: "{version}"
```

### Example: Go Toolchain
```yaml
name: golang
steps:
  - type: prompt
    var: module_name
    message: "Enter Go module name"
    default: "{cwd_basename}"
  - type: run
    command: "{sdk_root}/golang/{version}/bin/go.exe"
    args: ["mod", "init", "{module_name}"]
  - type: message
    text: "Go module {module_name} initialized"
env_vars:
  GOLANG_VERSION: "{version}"
```

### Example: Java Toolchain
```yaml
name: java
steps:
  - type: select
    var: build_tool
    message: "Choose build tool"
    options: ["maven", "gradle"]
    default: maven
  - type: prompt
    var: group_id
    message: "Enter groupId"
    default: "com.example"
  - type: prompt
    var: artifact_id
    message: "Enter artifactId"
    default: "{cwd_basename}"
  - type: run
    when: "{build_tool}" == "maven"
    command: "mvn"
    args: ["archetype:generate", "-DgroupId={group_id}", "-DartifactId={artifact_id}", "-DinteractiveMode=false"]
  - type: run
    when: "{build_tool}" == "gradle"
    command: "gradle"
    args: ["init", "--type", "java-application"]
  - type: message
    text: "Java project initialized with {build_tool}"
env_vars:
  JAVA_VERSION: "{version}"
```

### Adding Your Own Toolchain
1. Copy one of the example YAMLs above to `toolchains/envman_<yourtool>.yaml`.
2. Edit the steps, commands, and env_vars as needed for your toolchain.
3. Place your portable SDK in the appropriate folder under your SDKs root.
4. Run `envman init` and select your toolchain!

## Troubleshooting
- If a command fails, check the error message for details.
- Make sure your SDKs are portable and have a `bin` directory.
- If activation scripts do not work, check your shell and PATH settings.
- You can fully customize or extend any toolchain by editing its YAML config.

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
