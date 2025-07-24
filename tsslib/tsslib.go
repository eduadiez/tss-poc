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
	Message         []byte
	LogLevel        string
}

// TSSResult contains the results of the TSS operation
type TSSResult struct {
	Signature     string
	RecoveredAddr string
	MessageHash   string
	SignedTx      string
	FromAddr      string
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

	// Convert hash to big.Int and store in config
	hashBigInt := new(big.Int)
	hashBigInt.SetString(hex.EncodeToString(config.Message), 16)
	common.TssCfg.Message = hashBigInt.String()

	// Create and start TSS client
	c := client.NewTssClient(&common.TssCfg, client.SignMode, false)
	logger.Debugf("Client type: %T", c)

	// Get signature
	sig, err := c.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start TSS client: %w", err)
	}

	// Recover public key from signature
	pubKey, err := crypto.SigToPub(config.Message, sig)
	if err != nil {
		logger.Errorf("failed to recover pubkey: %w", err)
		return nil, fmt.Errorf("failed to recover pubkey: %w", err)
	}

	// Compute the address
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	fmt.Printf("recoveredAddr: %s\n", recoveredAddr.Hex())

	// Return results
	result := &TSSResult{
		Signature:     hex.EncodeToString(sig),
		RecoveredAddr: recoveredAddr.Hex(),
		MessageHash:   hex.EncodeToString(config.Message),
	}

	return result, nil
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
