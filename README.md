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
   cd tss
   go mod tidy
   cd ..
   go mod tidy
   ```

4. **Build the binaries**:
   ```bash
   cd tss
   go build -o tss-cli main.go
   mv ./tss-cli ..
   cd ..
   go build -o tss-signer cmd/tss-signer/main.go
   ```

## 📁 Project Structure

```
tss-poc/
├── cmd/
│   └── tss-signer/          # Command-line interface
│       └── main.go
├── tsslib/                  # TSS library wrapper
│   └── tsslib.go
├── tss/                     # TSS submodule (forked BNB Chain TSS implementation)
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

The TSS system requires a multi-step configuration process to set up the distributed signing environment. Follow these steps in order:

### Step 1: Initialize Keystores

First, create the initial file structure and keystores for each participant:

```bash
# Initialize keystore for signer 1
./tss-cli init --home ./keystore/signer1 --vault_name "default" --moniker "signer1" --password "123456789"

# Initialize keystore for signer 2  
./tss-cli init --home ./keystore/signer2 --vault_name "default" --moniker "signer2" --password "123456789"

# Initialize keystore for signer 3
./tss-cli init --home ./keystore/signer3 --vault_name "default" --moniker "signer3" --password "123456789"
```

**What this does**: Creates the necessary directory structure and configuration files for each TSS participant.

### Step 2: Generate Channel ID

Create a secure communication channel for the TSS participants:

```bash
./tss-cli channel --channel_expire 30
```

**What this does**: Generates a unique channel ID that will be used by all participants to communicate securely during the TSS operations. The `--channel_expire 30` sets the channel to expire in 30 days.

**Note**: Save the generated channel ID - you'll need it for the next step.

### Step 3: Generate Distributed Keys

Run the key generation process on all participants simultaneously. This is a distributed process where all parties must participate:

```bash
# Run on signer 1 (replace CHANNEL_ID with the actual ID from step 2)
./tss-cli keygen --home ./keystore/signer1 --vault_name "default" --parties 3 --threshold 1 --password "123456789" --channel_password "123456789" --channel_id CHANNEL_ID

# Run on signer 2 (replace CHANNEL_ID with the actual ID from step 2)
./tss-cli keygen --home ./keystore/signer2 --vault_name "default" --parties 3 --threshold 1 --password "123456789" --channel_password "123456789" --channel_id CHANNEL_ID

# Run on signer 3 (replace CHANNEL_ID with the actual ID from step 2)
./tss-cli keygen --home ./keystore/signer3 --vault_name "default" --parties 3 --threshold 1 --password "123456789" --channel_password "123456789" --channel_id CHANNEL_ID
```

**Parameters explained**:
- `--parties 3`: Total number of participants in the TSS scheme
- `--threshold 1`: Minimum number of participants required to sign (1 means any single party can sign)
- `--channel_id`: The channel ID generated in step 2
- `--channel_password`: Password for the secure channel

**Important**: All three keygen commands must be run simultaneously for the distributed key generation to succeed.

### Configuration Directory Structure

After successful configuration, your directory structure should look like this:

```
keystore/
├── signer1/
│   └── default/          # Configuration for signer 1
│       ├── config.json
│       ├── keygen.json
│       └── ...
├── signer2/
│   └── default/          # Configuration for signer 2
│       ├── config.json
│       ├── keygen.json
│       └── ...
└── signer3/
    └── default/          # Configuration for signer 3
        ├── config.json
        ├── keygen.json
        └── ...
```

### Configuration Contents

Each configuration directory contains:
- **Party information**: Participant details and network addresses
- **Network settings**: P2P communication parameters
- **Cryptographic parameters**: TSS-specific security settings
- **Communication channels**: Secure channel configurations
- **Key shares**: Distributed key material (never complete keys)

### Security Notes

- **Passwords**: Use strong, unique passwords in production
- **Channel security**: The channel password should be shared securely among participants
- **Key generation**: This is a critical security operation - ensure all participants are trusted
- **Backup**: Safely backup the keystore directories after successful key generation

## 🎯 Usage

### Command-Line Interface

The `tss-signer` binary provides a comprehensive CLI for TSS operations:

```bash
./tss-signer [flags]
```

### Basic Usage

#### Message Signing

The TSS system allows any configured participant to sign messages independently. Here are different ways to sign messages:

##### Basic Message Signing

Sign a message using any of the configured participants:

```bash
# Sign with signer 1
./tss-signer -home=./keystore/signer1 -message="Hello World"

# Sign with signer 2  
./tss-signer -home=./keystore/signer2 -message="Hello World"

# Sign with signer 3
./tss-signer -home=./keystore/signer3 -message="Hello World"
```

**Note**: Since we configured with `--threshold 1`, two participants are needed to sign messages

##### Advanced Configuration

Use custom parameters for more control:

```bash
./tss-signer \
  -home=./keystore/signer1 \
  -vault=default \
  -password=123456789 \
  -channelId=YOUR_CHANNEL_ID \
  -channelPassword=123456789 \
  -message="Test message" \
  -logLevel=debug
```

##### JSON Output Format

Get structured output for programmatic use:

```bash
# Basic JSON output
./tss-signer -home=./keystore/signer1 -message="Hello World" -json

# JSON with debug logging
./tss-signer -home=./keystore/signer2 -message="Test message" -logLevel=debug -json
```

##### Message Signing Examples

```bash
# Simple message signing
./tss-signer -home=./keystore/signer1 -message="Hello World"

# Sign with custom vault
./tss-signer -home=./keystore/signer2 -vault=myvault -message="Custom message"

# Debug mode for troubleshooting
./tss-signer -home=./keystore/signer3 -message="Debug test" -logLevel=debug

# JSON output for automation
./tss-signer -home=./keystore/signer1 -message="API test" -json

# Show help
./tss-signer -help
```

#### Transaction Signing

The TSS system can also sign Ethereum transactions using the `-mode=tx` parameter:

##### Basic Transaction Signing

```bash
# Sign a transaction with signer 1
./tss-signer -home=./keystore/signer1 -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000

# Sign a transaction with signer 2
./tss-signer -home=./keystore/signer2 -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000
```

##### Transaction Signing with Custom Parameters

```bash
./tss-signer \
  -home=./keystore/signer1 \
  -mode=tx \
  -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 \
  -value=1000000000000000000 \
  -gasLimit=21000 \
  -gasPrice=20000000000 \
  -nonce=5 \
  -chainId=1 \
  -data=0x1234567890abcdef
```

##### Transaction Signing and Broadcasting

```bash
# Sign and send transaction
./tss-signer \
  -home=./keystore/signer1 \
  -mode=tx \
  -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 \
  -value=1000000000000000000 \
  -chainId=11155111 \
  -gasPrice=1174319 \
  -rpc=https://eth-sepolia.public.blastapi.io \
  -send
```

##### Transaction Signing Examples

```bash
# Basic transaction signing
./tss-signer -home=./keystore/signer1 -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000

# Sign with custom gas parameters
./tss-signer -home=./keystore/signer2 -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000 -gasLimit=50000 -gasPrice=30000000000

# Sign with JSON output
./tss-signer -home=./keystore/signer3 -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000 -json

# Debug mode for transaction signing
./tss-signer -home=./keystore/signer1 -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000 -logLevel=debug
```

**Note**: In transaction mode (`-mode=tx`), the `-to` address is required, and the system automatically creates and signs the transaction hash.

### Command-Line Flags

| Parameter | Default | Description |
|-----------|---------|-------------|
| `home` | `test1` | Home directory containing keystore configuration |
| `vault` | `default` | Vault name within the keystore |
| `password` | `123456789` | Password for the vault |
| `channelId` | `1116C145287` | Channel ID for TSS communication |
| `channelPassword` | `123456789` | Password for the secure channel |
| `mode` | `message` | Signing mode: `message` for messages, `tx` for transactions |
| `message` | `123456789` | Message to sign (for message mode) |
| `logLevel` | `info` | Log level (debug, info, warn, error) |
| `json` | `false` | Output results as JSON format |

#### Transaction Mode Flags (use with `-mode=tx`)

| Parameter | Default | Description |
|-----------|---------|-------------|
| `to` | - | Recipient address (required for tx mode) |
| `value` | `0` | Amount in wei (for tx mode) |
| `gasLimit` | `21000` | Gas limit (for tx mode) |
| `gasPrice` | `20000000000` | Gas price in wei (for tx mode) |
| `nonce` | `0` | Transaction nonce (for tx mode) |
| `data` | - | Transaction data in hex (for tx mode) |
| `chainId` | `1` | Chain ID (for tx mode) |
| `rpc` | - | Ethereum RPC URL (e.g., http://localhost:8545) |
| `send` | `false` | Send the signed transaction (requires -rpc) |
| `etherscan` | `https://sepolia.etherscan.io` | Etherscan URL for transaction viewing |

## 📊 Output Format

### Human-Readable Output

#### Message Mode Output
```
Signature: 0x568b17caef811c510b413e6f20dd05f3306f85c79de667f6be07ed35421a23a35d184eb7483efe48bdfecc6f24d85051b96b1bcdcd9a4c3c75fa4b0e2d5f84ce01
Recovered Address: 0xf6844377aE73B4Ae396A75405807f862E2f220d4
Message Hash: 0x1234567890abcdef...
```

#### Transaction Mode Output
```
Signature: 0x568b17caef811c510b413e6f20dd05f3306f85c79de667f6be07ed35421a23a35d184eb7483efe48bdfecc6f24d85051b96b1bcdcd9a4c3c75fa4b0e2d5f84ce01
Recovered Address: 0xf6844377aE73B4Ae396A75405807f862E2f220d4
Message Hash: 0x1234567890abcdef...
Signed Tx: 0xf86c8085174876e800830186a094742d35cc6634c0532925a3b8d4c9db96c4b4d8b687038d7ea4c68000801ca0...
From Address: 0xf6844377aE73B4Ae396A75405807f862E2f220d4
```

### JSON Output

#### Message Mode JSON
```json
{
  "signature": "0x568b17caef811c510b413e6f20dd05f3306f85c79de667f6be07ed35421a23a35d184eb7483efe48bdfecc6f24d85051b96b1bcdcd9a4c3c75fa4b0e2d5f84ce01",
  "recoveredAddr": "0xf6844377aE73B4Ae396A75405807f862E2f220d4",
  "messageHash": "0x1234567890abcdef..."
}
```

#### Transaction Mode JSON
```json
{
  "signature": "0x568b17caef811c510b413e6f20dd05f3306f85c79de667f6be07ed35421a23a35d184eb7483efe48bdfecc6f24d85051b96b1bcdcd9a4c3c75fa4b0e2d5f84ce01",
  "recoveredAddr": "0xf6844377aE73B4Ae396A75405807f862E2f220d4",
  "messageHash": "0x1234567890abcdef...",
  "signedTx": "0xf86c8085174876e800830186a094742d35cc6634c0532925a3b8d4c9db96c4b4d8b687038d7ea4c68000801ca0...",
  "fromAddr": "0xf6844377aE73B4Ae396A75405807f862E2f220d4"
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