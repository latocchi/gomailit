# gomailit

Lightweight CLI for sending emails from the terminal. Written in Go, gomailit supports SMTP sending, templated messages, attachments, and CSV-based batch sends.

## Features
- Attachments and inline files
- Plain flags

## Installation

<!-- From source (requires Go 1.20+):
```bash
git clone https://github.com/latocchi/gomailit.git
cd gomailit
go build -o gomailit .
mv gomailit /usr/local/bin/
``` -->

```bash
go install github.com/youruser/gomailit/cmd/gomailit@latest
```

## Usage
### Basic setup
Authenticate your Gmail account to enable email sending:
```bash
gomail setup google
```
This will:
   - Open a browser window to authenticate with your google account.
   - Request the necessary Gmail API permissions.
   - Save your credentials locally for future use.

💡 Currently, only **Google Gmail** is supported. More providers will be added soon.

### Basic send
```bash
gomailit send --from alice@example.com --to bob@example.com \
    --subject "Hello" --body "This is a test"
```
```bash
Common flags
- --to          -t      Recipient (can be a single recipient or .txt file of recipients)
- --subject     -s      Email subject (Default: No subject)
- --body        -b      Inline body text (Default: No body)
- --attach      -a      Attachments (can be a directory or single file)
```

## Examples

### Send with attachment
```bash
gomailit send --to bob@example.com --subject "Files" --body "See attached" --attach report.pdf
```

### Send with attachments
```bash
gomailit send --to bob@example.com --subject "Files" --body "See attached" --attach report.pdf agenda.pdf
```

### Send with attachments on directory
```bash
gomailit send --to bob@example.com --subject "Files" --body "See attached" --attach ~/Documents/report/*
```

### Send with file as body
```bash
gomailit send --to bob@example.com --subject "Files" --body ~/Documents/body.txt --attach ~/Documents/report/*
```

### Send with multiple recipients in .txt file
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