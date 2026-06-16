# viv

The command line client for [Vivgrid](https://vivgrid.com), the Enterprise-grade AI Agent Platform.

## Getting Started

Install the compiled binary:

```sh
curl "https://bina.egoist.dev/vivgrid/cli?file=viv" | sh
```

## Build from source:

```sh
make

sudo cp ./bin/viv /usr/local/bin/
```

## Usage

```
viv --help
```

### Overview

The `viv` CLI tool allows you to manage globally deployed Serverless LLM Functions on [Vivgrid](https://vivgrid.com). It provides a complete workflow for developing, deploying, and monitoring LLM Tools in serverless.

### Configuration

You can configure `viv` using a configuration file or command-line flags:

**Configuration file** (`vivgrid.yml`):
```yaml
secret: your_app_secret
tool: my_first_llm_tool
```

**Environment variable for config file location**:
```bash
export VIV_CONFIG_FILE=/path/to/your/vivgrid.yml
```

### Global Flags

- `--secret string`: App secret for authentication
- `--tool string`: Serverless LLM Function name (Optional, default "my_first_llm_tool")
- `--api string`: REST API endpoint (Optional, default "https://hosting.vivgrid.com")

### Commands

#### General Commands

##### `viv deploy <source>`

One-command deployment that chains: upload → remove → create

**Examples:**
```bash
# Deploy current directory
viv deploy .
```

##### `viv upload <source>`

Upload and compile your source code to the vivgrid platform.

**Supported source formats:**
- Directories - Will be automatically zipped (respects .gitignore)
- `.zip` files - Pre-packaged zip archive
- `.go` files - Single Go source file

**Examples:**
```bash
# Upload a directory (auto-zips with exclusions)
viv upload ./my-function-dir
```

**Auto-exclusions when uploading directories:**
- `.git/` - Git repository directory
- `.vscode/` - VS Code settings
- `.DS_Store` - macOS system files
- `.env` - Environment files
- Files matching patterns in `.gitignore`

**Flags:**
- `--env key=value`: Set environment variables (can be used multiple times)

#### Deployment Management

##### `viv create`

Create and start a serverless deployment from previously uploaded code.

**Examples:**
```bash
# Create deployment
viv create

# Create with environment variables
viv create --env DATABASE_URL=postgres://... --env API_KEY=secret
```

**Flags:**
- `--env key=value`: Set environment variables (can be used multiple times)

##### `viv remove`

Delete the current serverless deployment.

```bash
viv remove
```

#### Monitoring & Observability

##### `viv status`

Show the current status of your serverless deployment.

```bash
viv status
```

**Output includes:**
- Deployment status
- Start time
- Global deployments

##### `viv logs`

Observe serverless logs in real-time.

```bash
viv logs
```

**Flags:**
- `--tail int`: Number of log lines to tail (default 20)

#### Utility Commands

##### `viv version`

Show the current version of the viv CLI tool.

```bash
viv version
```

## Docs

For more detailed documentation, visit the [Vivgrid Developer Docs](https://docs.vivgrid.com).
