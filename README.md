# TSS POC Library and Binary

This project provides a library and binary for Threshold Signature Scheme (TSS) operations using the BNB Chain TSS implementation.

## Structure

- `tsslib/` - The main library containing TSS functionality
- `cmd/tss-signer/` - A command-line binary that uses the library
- `main.go` - Original main file (now simplified to use the library)

## Library Usage

### Import the library

```go
import "github.com/eduadiez/tss-poc/tsslib"
```

### Create configuration

```go
config := tsslib.TSSConfig{
    Home:            "test1",
    Vault:           "default", 
    Password:        "123456789",
    ChannelID:       "1116C145287",
    ChannelPassword: "123456789",
    Message:         "Hello World",
    LogLevel:        "info",
}
```

### Sign a message

```go
result, err := tsslib.SignMessage(config)
if err != nil {
    log.Fatalf("Failed to sign message: %v", err)
}

fmt.Printf("Signature: %s\n", result.Signature)
fmt.Printf("Recovered Address: %s\n", result.RecoveredAddr)
fmt.Printf("Message Hash: %s\n", result.MessageHash)
```

## Binary Usage

### Build the binary

```bash
go build -o tss-signer cmd/tss-signer/main.go
```

### Run with default parameters

```bash
./tss-signer
```

### Run with custom parameters

```bash
./tss-signer \
  -home=test2 \
  -vault=default \
  -password=123456789 \
  -channelId=1116C145287 \
  -channelPassword=123456789 \
  -message="Hello World" \
  -logLevel=debug
```

### Output as JSON

```bash
./tss-signer -message="Test message" -json
```

### Show help

```bash
./tss-signer -help
```

## Command Line Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `home` | `test1` | Home directory for configuration |
| `vault` | `default` | Vault name |
| `password` | `123456789` | Password for the vault |
| `channelId` | `1116C145287` | Channel ID |
| `channelPassword` | `123456789` | Channel Password |
| `message` | `123456789` | Message to sign |
| `logLevel` | `info` | Log level (debug, info, warn, error) |
| `json` | `false` | Output results as JSON |

## Examples

### Basic usage
```bash
./tss-signer -message="Hello World"
```

### Debug mode with JSON output
```bash
./tss-signer -message="Test message" -logLevel=debug -json
```

### Custom configuration
```bash
./tss-signer \
  -home=test3 \
  -vault=myvault \
  -password=mypassword \
  -message="Custom message"
```

## Library API

### TSSConfig

```go
type TSSConfig struct {
    Home            string
    Vault           string
    Password        string
    ChannelID       string
    ChannelPassword string
    Message         string
    LogLevel        string
}
```

### TSSResult

```go
type TSSResult struct {
    Signature       string
    RecoveredAddr   string
    MessageHash     string
    ConfigJSON      string
}
```

### SignMessage

```go
func SignMessage(config TSSConfig) (*TSSResult, error)
```

Performs TSS signing with the given configuration and returns the signature, recovered address, message hash, and configuration JSON.

## Error Handling

The library returns detailed error messages wrapped with context. Always check for errors:

```go
result, err := tsslib.SignMessage(config)
if err != nil {
    log.Fatalf("TSS signing failed: %v", err)
}
```

## Logging

The library uses structured logging with different log levels. Set the `LogLevel` in the configuration to control verbosity:

- `debug` - Most verbose, includes configuration details
- `info` - Standard information messages
- `warn` - Warning messages only
- `error` - Error messages only 