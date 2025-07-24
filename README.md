# TSS POC - Threshold Signature Scheme Implementation

A comprehensive implementation of Threshold Signature Scheme (TSS) for Ethereum using the BNB Chain TSS library. This project provides both a library and a command-line interface for distributed signing of Ethereum messages and transactions

## 🚀 Features

- **TSS Library**: Reusable library for TSS operations
- **Message Signing**: Sign arbitrary messages using TSS
- **Transaction Signing**: Sign Ethereum transactions with TSS
- **Transaction Broadcasting**: Send signed transactions to Ethereum networks
- **Multi-Network Support**: Support for Ethereum mainnet, testnets, and other EVM chains
- **Address Recovery**: Verify recovered addresses from signatures
- **JSON Output**: Structured output for programmatic use
- **Comprehensive Logging**: Detailed logging for debugging and monitoring
- **Command-line Interface**: Easy-to-use CLI for TSS operations

## 📋 Prerequisites

- Go 1.19 or higher
- TSS configuration files (see Configuration section)

## 🛠️ Installation

1. **Clone the repository with submodules**:
   ```bash
   git clone --recursive https://github.com/eduadiez/tss-poc.git
   cd tss-poc
   ```

2. **If you already cloned without submodules, initialize them**:
   ```bash
   git submodule update --init --recursive
   ```

3. **Install dependencies**:
   ```bash
   go mod download
   ```

4. **Build the binary**:
   ```bash
   go build -o tss-signer cmd/tss-signer/main.go
   ```

## 📁 Project Structure

```
tss-poc/
├── cmd/
│   └── tss-signer/          # Command-line interface
│       └── main.go
├── tsslib/                  # TSS library
│   └── tsslib.go
├── tss/                     # TSS submodule (BNB Chain TSS implementation)
│   ├── client/
│   ├── common/
│   ├── server/
│   └── ...
├── keystore/                # TSS keystore configurations
│   ├── test1/
│   │   └── default/
│   ├── test2/
│   │   └── default/
│   └── test3/
│       └── default/
├── go.mod
├── go.sum
├── .gitmodules              # Git submodule configuration
└── README.md
```

## ⚙️ Configuration

The TSS system requires configuration files for each participant. These should be placed in the appropriate keystore directories:

- `keystore/test1/default/` - Configuration for test environment 1
- `keystore/test2/default/` - Configuration for test environment 2  
- `keystore/test3/default/` - Configuration for test environment 3

Each configuration includes:
- Party information
- Network settings
- Cryptographic parameters
- Communication channels

## 🎯 Usage

### Command-Line Interface

The `tss-signer` binary provides a comprehensive CLI for TSS operations:

```bash
./tss-signer [flags]
```

### Basic Usage

#### Message Signing

```bash
# Basic message signing with default parameters
./tss-signer -message="Hello World"

# With custom configuration
./tss-signer \
  -home=test2 \
  -vault=default \
  -password=123456789 \
  -channelId=1116C145287 \
  -channelPassword=123456789 \
  -message="Test message" \
  -logLevel=debug

# Output as JSON
./tss-signer -message="Hello World" -json
```

### Command-Line Flags

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

### Examples

#### Basic usage
```bash
./tss-signer -message="Hello World"
```

#### Debug mode with JSON output
```bash
./tss-signer -message="Test message" -logLevel=debug -json
```

#### Custom configuration
```bash
./tss-signer \
  -home=test3 \
  -vault=myvault \
  -password=mypassword \
  -message="Custom message"
```

#### Show help
```bash
./tss-signer -help
```

## 📊 Output Format

### Human-Readable Output
```
Signature: 0x568b17caef811c510b413e6f20dd05f3306f85c79de667f6be07ed35421a23a35d184eb7483efe48bdfecc6f24d85051b96b1bcdcd9a4c3c75fa4b0e2d5f84ce01
Recovered Address: 0xf6844377aE73B4Ae396A75405807f862E2f220d4
Message Hash: 0x1234567890abcdef...
Configuration:
{
  "p2p": {
    "listen": "/ip4/0.0.0.0/tcp/8080"
  },
  "log_level": "info",
  ...
}
```

### JSON Output
```json
{
  "signature": "0x568b17caef811c510b413e6f20dd05f3306f85c79de667f6be07ed35421a23a35d184eb7483efe48bdfecc6f24d85051b96b1bcdcd9a4c3c75fa4b0e2d5f84ce01",
  "recoveredAddr": "0xf6844377aE73B4Ae396A75405807f862E2f220d4",
  "messageHash": "0x1234567890abcdef...",
  "config": {
    "p2p": {
      "listen": "/ip4/0.0.0.0/tcp/8080"
    },
    "log_level": "info"
  }
}
```

## 🔧 Library Usage

### Import the Library

```go
import "github.com/eduadiez/tss-poc/tsslib"
```

### Create Configuration

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

### Sign a Message

```go
result, err := tsslib.SignMessage(config)
if err != nil {
    log.Fatalf("Failed to sign message: %v", err)
}

fmt.Printf("Signature: %s\n", result.Signature)
fmt.Printf("Recovered Address: %s\n", result.RecoveredAddr)
fmt.Printf("Message Hash: %s\n", result.MessageHash)
```

## 🌐 Library API

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

## 🔍 Troubleshooting

### Common Issues

#### 1. TSS Configuration Not Found
```
Error: failed to read config
```
**Solution**: Verify that TSS configuration files exist in the specified keystore directory.

#### 2. Submodule Issues
```
Error: cannot find package github.com/bnb-chain/tss
```
**Solution**: Initialize and update the submodules:
```bash
git submodule update --init --recursive
```

#### 3. Address Inconsistency
If you're getting different recovered addresses for different messages, this indicates a signature recovery issue. The library now uses the TSS public key directly to ensure consistent address recovery.

### Debug Mode

Enable debug logging for detailed troubleshooting:

```bash
./tss-signer -message="Test message" -logLevel=debug
```

## 🔐 Security Considerations

- **Private Keys**: TSS eliminates the need for single private keys by distributing signing authority
- **Threshold Security**: Messages require a threshold of participants to sign
- **No Single Point of Failure**: No single party can sign messages alone
- **Key Management**: Keys are never reconstructed in a single location
- **Consistent Address Recovery**: Uses TSS public key for reliable address recovery

## 📝 Examples

### Complete Workflow Examples

#### 1. Message Signing Workflow
```bash
# Sign a message
./tss-signer -message="Hello from TSS!" -home=test1 -json

# Verify the signature (using external tools)
# The recovered address should match the expected signer
```

#### 2. Debug Mode for Troubleshooting
```bash
# Sign with debug logging
./tss-signer -message="Test message" -logLevel=debug -home=test2
```

#### 3. Custom Configuration
```bash
# Use custom keystore and parameters
./tss-signer \
  -home=test3 \
  -vault=myvault \
  -password=mypassword \
  -channelId=MYCHANNEL123 \
  -channelPassword=mychannelpass \
  -message="Custom message" \
  -json
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [BNB Chain TSS Library](https://github.com/bnb-chain/tss-lib) - Core TSS implementation
- [Ethereum Go Client](https://github.com/ethereum/go-ethereum) - Ethereum integration
- [IPFS Go Log](https://github.com/ipfs/go-log) - Logging framework

## 📞 Support

For support and questions:
- Create an issue in the GitHub repository
- Check the troubleshooting section above
- Review the debug logs with `-logLevel=debug`

---

**Note**: This is a proof-of-concept implementation. For production use, ensure proper security audits and testing. 