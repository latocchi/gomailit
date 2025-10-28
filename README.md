# gomailit

Lightweight CLI for sending emails from the terminal. Written in Go, gomailit sends email via Gmail REST API, authentication using OAuth2, supports sending email to multiple recipients, and supports attaching multiple attachments.

## Features
- Send to multiple recipients (supports `.txt` recipient lists)
- Attach multiple files or entire directories
- Inline or file-based email bodies
- Gmail OAuth2 authentication (no password handling)
- Simple flag-based configuration

## Installation

<!-- From source (requires Go 1.20+):
```bash
git clone https://github.com/latocchi/gomailit.git
cd gomailit
go build -o gomailit .
mv gomailit /usr/local/bin/
``` -->

```bash
go install github.com/latocchi/gomailit@latest
```

## Usage
### Basic setup
Authenticate your Gmail account to enable email sending:
```bash
gomailit setup google
```
This will:
   - Open a browser window for Google OAuth2 authentication.
   - Request the necessary Gmail API permissions.
   - Save your credentials locally for future use (typically under ~/.config/gomailit/). 

💡 Currently, only **Google Gmail** is supported. More providers will be added soon.

### Basic send
```bash
gomailit send --from alice@example.com --to bob@example.com \
    --subject "Hello" --body "This is a test"
```
| Flag        | Alias | Description                                                       |
| ----------- | ----- | ----------------------------------------------------------------- |
| `--to`      | `-t`  | Recipient (single email or `.txt` file with one address per line) |
| `--subject` | `-s`  | Email subject *(default: “No subject”)*                           |
| `--body`    | `-b`  | Inline body text or path to a `.txt` file *(default: "No body")*      |
| `--attach`  | `-a`  | One or more attachment files, or a directory path                 |


## Examples

### Send with attachment
```bash
gomailit send --to bob@example.com --subject "Files" --body "See attached" --attach report.pdf
```

### Send multiple attachments
```bash
gomailit send --to bob@example.com --subject "Files" --body "See attached" --attach report.pdf agenda.pdf
```

### Attach all files from a directory
```bash
gomailit send --to bob@example.com --subject "Files" --body "See attached" --attach ~/Documents/report/*
```

### Use file for email body
```bash
gomailit send --to bob@example.com --subject "Files" --body ~/Documents/body.txt --attach ~/Documents/report/*
```

### Send to multiple recipients via `.txt` file
```bash
gomailit send --to ~/Documents/recipients.txt --subject "Files" --body ~/Documents/body.txt --attach ~/Documents/report/*
```

Example of recipients.txt:
```text
recipient@example.com
recipient1@example.com
recipient2@example.com
recipient3@example.com
recipient4@example.com
recipient5@example.com
```

## Contributing

- Fork the repo and open a feature branch
- Open a pull request with clear description

## License

MIT — see LICENSE file for details.

Feedback and issues: open an issue on the repository.