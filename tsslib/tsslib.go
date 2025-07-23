package tsslib

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss/client"
	"github.com/bnb-chain/tss/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipfs/go-log"
	"github.com/spf13/viper"
)

// TSSConfig holds all the configuration parameters for the TSS library
type TSSConfig struct {
	Home            string
	Vault           string
	Password        string
	ChannelID       string
	ChannelPassword string
	Message         string
	LogLevel        string
}

// TSSResult contains the results of the TSS operation
type TSSResult struct {
	Signature     string
	RecoveredAddr string
	MessageHash   string
}

// SignMessage performs TSS signing with the given configuration
func SignMessage(config TSSConfig) (*TSSResult, error) {
	// Create a debug logger
	logger := log.Logger("tsslib")

	// Read configuration from home
	if err := common.ReadConfigFromHome(viper.GetViper(), false, config.Home, config.Vault, config.Password); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	// Set configuration values
	common.TssCfg.ChannelId = config.ChannelID
	common.TssCfg.ChannelPassword = config.ChannelPassword
	common.TssCfg.Password = config.Password

	// Initialize log level
	initLogLevel(config.LogLevel)

	// Create Ethereum message hash
	messageHash := setEthereumMessage(config.Message)

	// Create and start TSS client
	c := client.NewTssClient(&common.TssCfg, client.SignMode, false)
	logger.Debugf("Client type: %T", c)

	// Get signature
	sig, err := c.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start TSS client: %w", err)
	}
	//fmt.Printf("sig: %s\n", hex.EncodeToString(sig))
	//fmt.Printf("sigRecovery: %s\n", hex.EncodeToString(sigRecovery))

	// Prepare signature with recovery byte
	//sigWithV := make([]byte, 65)
	//copy(sigWithV[:64], sig)
	//sigWithV[64] = 1 // v should be 0 or 1 (or 27/28 for some libraries)

	// Recover public key from signature
	pubKey, err := crypto.SigToPub(messageHash, sig)
	if err != nil {
		return nil, fmt.Errorf("failed to recover pubkey: %w", err)
	}

	// Compute the address
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)

	// Return results
	result := &TSSResult{
		Signature:     hex.EncodeToString(sig),
		RecoveredAddr: recoveredAddr.Hex(),
		MessageHash:   hex.EncodeToString(messageHash),
	}

	return result, nil
}

// setEthereumMessage creates an Ethereum-style message hash
func setEthereumMessage(message string) []byte {
	// Create the Ethereum-style message prefix
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))

	// Combine prefix and message
	prefixedMessage := append([]byte(prefix), []byte(message)...)

	// Hash the message
	hash := crypto.Keccak256(prefixedMessage)

	// Convert hash to big.Int and store in config
	hashBigInt := new(big.Int)
	hashBigInt.SetString(hex.EncodeToString(hash), 16)
	common.TssCfg.Message = hashBigInt.String()

	return hash
}

// initLogLevel sets up logging levels for various components
func initLogLevel(logLevel string) {
	log.SetLogLevel("tsslib", logLevel)
	log.SetLogLevel("tss", logLevel)
	log.SetLogLevel("tss-lib", logLevel)
	log.SetLogLevel("srv", logLevel)
	log.SetLogLevel("trans", logLevel)
	log.SetLogLevel("p2p_utils", logLevel)
	log.SetLogLevel("common", logLevel)

	// libp2p loggers
	log.SetLogLevel("dht", "error")
	log.SetLogLevel("discovery", "error")
	log.SetLogLevel("swarm2", "error")
	log.SetLogLevel("stream-upgrader", "error")
}
