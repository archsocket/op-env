<p align="center">
    <img src="./assets/icon-128.png">
</p>

# op-env

op-env is a CLI tool that converts items from one or more 1Password vaults into a .env (dotenv) file.

It simplifies integrating 1Password with application deployments by exporting vault contents as environment variables.

* [Get started with 1Password Service Accounts](https://developer.1password.com/docs/service-accounts/get-started/)

## Features

* Fetches items from specified 1Password vaults, or all accessible vaults by default

* Converts item titles and field names into environment-variable-safe keys

* Supports multiline values and escapes quotes/newlines

* Outputs a valid .env file

## Installation

### a) [Go Install](https://go.dev/doc/install)
```bash
go install github.com/archsocket/op-env@latest
```

### b) [Direct Binary Download](https://github.com/archsocket/op-env/releases)
```
https://github.com/archsocket/op-env/releases/download/{VERSION}/op-env_{OS}_{ARCH}
```

Path parameters:
* `version`:
    * See [Tags](https://github.com/archsocket/op-env/tags) for possible values.
* `os`:
    * Possible values: `linux`, `windows`, `darwin`
* `arch`:
    * Possible values: `amd64`, `arm64`

Example with linux amd64:
```bash
wget https://github.com/archsocket/op-env/releases/download/v1.1.0/op-env_linux_amd64
chmod +x op-env_linux_amd64
sudo mv op-env_linux_amd64 /usr/local/bin/op-env
```

### c) [Docker](https://github.com/archsocket/op-env/pkgs/container/op-env)
```bash
docker pull ghcr.io/archsocket/op-env:latest
```

## Usage

| Flag      | Description                                                                                               |
| --------- | --------------------------------------------------------------------------------------------------------- |
| `--token` | 1Password service account token. Can also be set via the `OP_SERVICE_ACCOUNT_TOKEN` env variable.         |
| `--vault` | Name or ID of a 1Password vault to export. Can be used multiple times. Defaults to all accessible vaults. |
| `--out`   | Output filename. Use `-` to write to stdout. (default: `.env`)                                            |

## Examples

### Export all accessible vaults into `.env`:
```bash
export OP_SERVICE_ACCOUNT_TOKEN="your-token-here"
op-env
```

### Export specific vaults:
```bash
op-env --vault "Engineering" --vault "Production"
```

### Write output to a custom file
```bash
op-env --out app.env
```

### Write output to stdout
```bash
op-env --out -
```

### Docker

Unix
```bash
docker run --rm \
    -e OP_SERVICE_ACCOUNT_TOKEN=$OP_SERVICE_ACCOUNT_TOKEN \
    -v "$PWD:/out" \
    ghcr.io/archsocket/op-env:latest \
    op-env \
    --out ./out/.env
```
</details>

Windows
```bash
docker run --rm \
    -e "OP_SERVICE_ACCOUNT_TOKEN=$OP_SERVICE_ACCOUNT_TOKEN" \
    -v "$(pwd -W):/out" \
    ghcr.io/archsocket/op-env:latest \
    op-env \
    --out ./out/.env
```
</details>

## How It Works

### Key Formatting

All keys are formatted using the following rules:
* Convert spaces to _
* Remove any characters not matching [a-zA-Z0-9_]
* Convert to UPPERCASE

Example:
```
"My API Key" → MY_API_KEY
```

Fields inside an item are namespaced under the item key:
```
Item title: "Database"
Field title: "username"
→ DATABASE_USERNAME
```

### Value Handling

* Newlines (\n) are escaped as \\n

* Quotes (") are escaped as \\"

* Values are wrapped in quotes in the dotenv output

## Output Example

A vault item such as:
* Title: API Credentials
* Notes: Some internal notes
* Fields:
    * Key: abc123
    * Secret: def456

Would produce:
```
API_CREDENTIALS="Some internal notes"
API_CREDENTIALS_KEY="abc123"
API_CREDENTIALS_SECRET="def456"
```

## Development

Clone the repo and build:
```bash
git clone git@github.com:archsocket/op-env.git
cd op-env
go build
```

Run locally:
```bash
go run main.go --token $OP_SERVICE_ACCOUNT_TOKEN
```
