package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"

	"github.com/eduadiez/tss-poc/tsslib"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ipfs/go-log"
)

// SigningMode represents the type of signing operation
type SigningMode string

const (
	MessageMode SigningMode = "message"
	TxMode      SigningMode = "tx"
)

func main() {
	// Define command-line flags
	var (
		// Common flags
		home            = flag.String("home", "test1", "Home directory for configuration")
		vault           = flag.String("vault", "default", "Vault name")
		password        = flag.String("password", "123456789", "Password for the vault")
		channelID       = flag.String("channelId", "1116C145287", "Channel ID")
		channelPassword = flag.String("channelPassword", "123456789", "Channel Password")
		logLevel        = flag.String("logLevel", "info", "Log level (debug, info, warn, error)")
		outputJSON      = flag.Bool("json", false, "Output results as JSON")

		// Mode selection
		mode = flag.String("mode", "message", "Signing mode: 'message' or 'tx'")

		// Message mode flags
		message = flag.String("message", "123456789", "Message to sign (for message mode)")

		// Transaction mode flags
		to        = flag.String("to", "", "Recipient address (for tx mode)")
		value     = flag.String("value", "0", "Amount in wei (for tx mode)")
		gasLimit  = flag.Uint64("gasLimit", 21000, "Gas limit (for tx mode)")
		gasPrice  = flag.String("gasPrice", "20000000000", "Gas price in wei (for tx mode)")
		nonce     = flag.Uint64("nonce", 0, "Transaction nonce (for tx mode)")
		data      = flag.String("data", "", "Transaction data (hex, for tx mode)")
		chainID   = flag.Uint64("chainId", 1, "Chain ID (for tx mode)")
		rpcURL    = flag.String("rpc", "", "Ethereum RPC URL (e.g., http://localhost:8545)")
		send      = flag.Bool("send", false, "Send the signed transaction (requires -rpc)")
		etherscan = flag.String("etherscan", "https://sepolia.etherscan.io", "Etherscan URL (default: https://sepolia.etherscan.io)")
	)

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  # Message signing:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -mode=message -message=\"Hello World\"\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -home=test2 -vault=default -message=\"Test message\" -json\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  \n")
		fmt.Fprintf(flag.CommandLine.Output(), "  # Transaction signing:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000 -gasLimit=21000 -gasPrice=20000000000 -nonce=5 -chainId=1\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  \n")
		fmt.Fprintf(flag.CommandLine.Output(), "  # Transaction signing and sending:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -mode=tx -to=0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b6 -value=1000000000000000000 -rpc=http://localhost:8545 -send\n", os.Args[0])
	}

	// Parse command-line flags
	flag.Parse()

	// Validate mode
	signingMode := SigningMode(*mode)
	if signingMode != MessageMode && signingMode != TxMode {
		fmt.Fprintf(os.Stderr, "Error: Invalid mode '%s'. Must be 'message' or 'tx'\n", *mode)
		os.Exit(1)
	}

	// Validate transaction sending parameters
	if *send && *rpcURL == "" {
		fmt.Fprintf(os.Stderr, "Error: -rpc is required when -send is specified\n")
		os.Exit(1)
	}

	// Initialize logger
	logger := log.Logger("tss-signer")
	log.SetLogLevel("tss-signer", *logLevel)

	var config tsslib.TSSConfig
	var result *tsslib.TSSResult
	var err error

	switch signingMode {
	case MessageMode:
		logger.Infof("Starting TSS message signing with home=%s, vault=%s, message=%s", *home, *vault, *message)

		config = tsslib.TSSConfig{
			Home:            *home,
			Vault:           *vault,
			Password:        *password,
			ChannelID:       *channelID,
			ChannelPassword: *channelPassword,
			Message:         getEthereumMessage(*message),
			LogLevel:        *logLevel,
		}

		result, err = tsslib.SignMessage(config)

	case TxMode:
		// Validate transaction parameters
		if *to == "" {
			fmt.Fprintf(os.Stderr, "Error: 'to' address is required for transaction mode\n")
			os.Exit(1)
		}

		// Parse transaction parameters
		toAddr := common.HexToAddress(*to)
		valueBigInt := new(big.Int)
		valueBigInt.SetString(*value, 10)
		gasPriceBigInt := new(big.Int)
		gasPriceBigInt.SetString(*gasPrice, 10)

		// Create transaction
		tx := types.NewTransaction(
			*nonce,
			toAddr,
			valueBigInt,
			*gasLimit,
			gasPriceBigInt,
			common.FromHex(*data),
		)

		// Create transaction hash
		signer := types.NewEIP155Signer(big.NewInt(int64(*chainID)))

		txHash := signer.Hash(tx)

		logger.Infof("Starting TSS transaction signing with home=%s, vault=%s", *home, *vault)
		logger.Infof("Transaction: to=%s, value=%s, gasLimit=%d, gasPrice=%s, nonce=%d, chainId=%d",
			toAddr.Hex(), valueBigInt.String(), *gasLimit, gasPriceBigInt.String(), *nonce, *chainID)
		logger.Infof("Signer chain ID: %d", signer.ChainID().Int64())

		config = tsslib.TSSConfig{
			Home:            *home,
			Vault:           *vault,
			Password:        *password,
			ChannelID:       *channelID,
			ChannelPassword: *channelPassword,
			Message:         txHash.Bytes(), // Use transaction hash as message
			LogLevel:        *logLevel,
		}

		result, err = tsslib.SignMessage(config)

		// If successful, reconstruct the signed transaction
		if err == nil {
			// Parse signature
			sigBytes := common.FromHex(result.Signature)
			if len(sigBytes) != 65 {
				err = fmt.Errorf("invalid signature length: expected 65 bytes, got %d", len(sigBytes))
			} else {

				// Try using the original TSS signature first
				signedTx, err := tx.WithSignature(signer, sigBytes)
				if err != nil {
					logger.Errorf("Failed to create signed transaction with original signature: %v", err)
				} else {
					logger.Infof("Successfully created transaction with original TSS signature")
				}

				if err == nil {
					// Check the chain ID of the signed transaction
					logger.Infof("Signed transaction chain ID: %d", signedTx.ChainId().Int64())

					// Encode signed transaction
					encodedTx, err := rlp.EncodeToBytes(signedTx)
					if err != nil {
						logger.Errorf("Failed to encode signed transaction: %v", err)
					} else {
						result.SignedTx = common.Bytes2Hex(encodedTx)

						// Extract the 'from' address from the signed transaction
						// Use the same signer that was used to create the transaction
						fromAddr, err := types.Sender(signer, signedTx)
						if err != nil {
							logger.Errorf("Failed to extract 'from' address: %v", err)
							// Fallback to recovered address if extraction fails
							result.FromAddr = result.RecoveredAddr
							logger.Warnf("Using recovered address as fallback: %s", result.FromAddr)
						} else {
							result.FromAddr = fromAddr.Hex()
							logger.Infof("Transaction 'from' address: %s", result.FromAddr)

							// Verify that addresses match
							if result.FromAddr != result.RecoveredAddr {
								logger.Warnf("Address mismatch! Recovered: %s, From: %s", result.RecoveredAddr, result.FromAddr)
							} else {
								logger.Infof("Address verification successful: %s", result.FromAddr)
							}
						}

						if *send {
							txHash, err := sendTransaction(*rpcURL, encodedTx, logger)
							if err != nil {
								logger.Errorf("Failed to send transaction: %v", err)
							} else {
								fmt.Printf("Transaction sent! Hash: %s\n", txHash)
								fmt.Printf("View on Etherscan: %s/tx/%s\n", *etherscan, txHash)
							}
						}
					}
				}
			}
		}
	}

	// Output results
	if *outputJSON {
		// Output as JSON
		jsonResult := map[string]interface{}{
			"signature":     result.Signature,
			"recoveredAddr": result.RecoveredAddr,
			"messageHash":   result.MessageHash,
			"signedTx":      result.SignedTx,
		}
		if signingMode == TxMode && result.FromAddr != "" {
			jsonResult["fromAddr"] = result.FromAddr
		}

		jsonData, err := json.MarshalIndent(jsonResult, "", "  ")
		if err != nil {
			logger.Errorf("Failed to marshal result to JSON: %v", err)
			os.Exit(1)
		}
		fmt.Println(string(jsonData))
	} else {
		// Output as human-readable format
		fmt.Printf("Signature: %s\n", result.Signature)
		fmt.Printf("Recovered Address: %s\n", result.RecoveredAddr)
		fmt.Printf("Message Hash: %s\n", result.MessageHash)
		if signingMode == TxMode {
			fmt.Printf("Signed Tx: %s\n", result.SignedTx)
		}
	}

	logger.Info("TSS signing completed successfully")
}

// sendTransaction sends the raw signed transaction to the Ethereum node and returns the transaction hash.
func sendTransaction(rpcURL string, rawTx []byte, logger log.StandardLogger) (string, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return "", fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}
	defer client.Close()

	ctx := context.Background()
	tx := new(types.Transaction)
	if err := rlp.DecodeBytes(rawTx, tx); err != nil {
		return "", fmt.Errorf("failed to decode raw transaction: %w", err)
	}

	err = client.SendTransaction(ctx, tx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %w", err)
	}
	logger.Infof("Transaction sent! Hash: %s", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

// setEthereumMessage creates an Ethereum-style message hash
func getEthereumMessage(message string) []byte {
	// Create the Ethereum-style message prefix
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))

	// Combine prefix and message
	prefixedMessage := append([]byte(prefix), []byte(message)...)

	// Hash the message
	hash := crypto.Keccak256(prefixedMessage)

	return hash
}
