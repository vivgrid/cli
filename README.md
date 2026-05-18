# yc

The command line client for [Vivgrid](https://vivgrid.com), the Enterprise-grade AI Agent Platform.

## Getting Started

Install the compiled binary:

```sh
curl "https://bina.egoist.dev/vivgrid/yc" | sh
```

## Build from source:

```sh
make

sudo cp ./bin/yc /usr/local/bin/
```

## Usage

```
yc --help
```

### Overview

The `yc` CLI tool allows you to manage globally deployed Serverless LLM Functions on [Vivgrid](https://vivgrid.com). It provides a complete workflow for developing, deploying, and monitoring LLM Tools in serverless.

### Configuration

You can configure `yc` using a configuration file or command-line flags:

**Configuration file** (`yc.yml`):
```yaml
api: http://127.0.0.1:9040
secret: your_app_secret
tool: my_llm_function_tool
```

**Environment variable for config file location**:
```bash
export YC_CONFIG_FILE=/path/to/your/yc.yml
```

### Global Flags

- `--api string`: REST API endpoint (default "http://127.0.0.1:9040")
- `--secret string`: App secret for authentication
- `--tool string`: Serverless LLM Function name (default "my_first_llm_tool")

### Commands

#### General Commands

##### `yc deploy <source>`

One-command deployment that chains: upload → remove → create

**Examples:**
```bash
# Deploy current directory
yc deploy .
```

##### `yc upload <source>`

Upload and compile your source code to the vivgrid platform.

**Supported source formats:**
- Directories - Will be automatically zipped (respects .gitignore)
- `.zip` files - Pre-packaged zip archive
- `.go` files - Single Go source file

**Examples:**
```bash
# Upload a directory (auto-zips with exclusions)
yc upload ./my-function-dir
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

##### `yc create`

Create and start a serverless deployment from previously uploaded code.

**Examples:**
```bash
# Create deployment
yc create

# Create with environment variables
yc create --env DATABASE_URL=postgres://... --env API_KEY=secret
```

**Flags:**
- `--env key=value`: Set environment variables (can be used multiple times)

##### `yc remove`

Delete the current serverless deployment.

```bash
yc remove
```

#### Monitoring & Observability

##### `yc status`

Show the current status of your serverless deployment.

```bash
yc status
```

**Output includes:**
- Deployment status
- Start time
- Global deployments

##### `yc logs`

Observe serverless logs in real-time.

```bash
yc logs
```

**Flags:**
- `--tail int`: Number of log lines to tail (default 20)

#### Utility Commands

##### `yc version`

Show the current version of the yc CLI tool.

```bash
yc version
```

## Docs

For more detailed documentation, visit the [Vivgrid Developer Docs](https://docs.vivgrid.com).
