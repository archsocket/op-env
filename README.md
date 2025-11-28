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

```
go install github.com/archsocket/op-env@latest
```

## Usage

| Flag      | Description                                                                                                    |
| --------- | -------------------------------------------------------------------------------------------------------------- |
| `--token` | 1Password service account token. Can also be set via the `OP_SERVICE_ACCOUNT_TOKEN` environment variable.      |
| `--vault` | Name or ID of a 1Password vault to export. Can be specified multiple times. Defaults to all accessible vaults. |
| `--file`  | Output filename (default: `.env`)                                                                              |

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
op-env --file app.env
```

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
git clone https://github.com/archsocket/op-env.git
cd op-env
go build
```

Run locally:
```bash
go run main.go --token $OP_SERVICE_ACCOUNT_TOKEN
```
