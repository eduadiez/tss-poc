package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"

	"github.com/bnb-chain/tss/client"
	"github.com/bnb-chain/tss/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ipfs/go-log"
	"github.com/spf13/viper"
)

func main() {
	// Create a debug logger
	logger := log.Logger("main")
	// Define command-line flags
	var (
		home            = flag.String("home", "test1", "Home directory for configuration")
		vault           = flag.String("vault", "default", "Vault name")
		password        = flag.String("password", "123456789", "Password for the vault")
		channelId       = flag.String("channelId", "1116C145287", "Channel ID")
		channelPassword = flag.String("channelPassword", "123456789", "Channel Password")
		message         = flag.String("message", "123456789", "Message to sign")
	)

	flag.Parse()
	if err := common.ReadConfigFromHome(viper.GetViper(), false, *home, *vault, *password); err != nil {
		common.Panic(err)
	}

	common.TssCfg.ChannelId = *channelId
	common.TssCfg.ChannelPassword = *channelPassword
	common.TssCfg.Password = *password

	// Print as JSON
	jsonData, err := json.MarshalIndent(common.TssCfg, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling to JSON: %v\n", err)
		return
	}
	initLogLevel(common.TssCfg)
	logger.Debugf("TSS Configuration: %s", string(jsonData))

	messageHash := setEthereumMessage(*message)

	c := client.NewTssClient(&common.TssCfg, client.SignMode, false)
	fmt.Printf("Client type: %T\n", c)

	sig, err := c.Start()
	if err != nil {
		common.Panic(err)
	}
	fmt.Println("Signature: ", hex.EncodeToString(sig))

	sigWithV := make([]byte, 65)
	copy(sigWithV[:64], sig)
	sigWithV[64] = 0 // v should be 0 or 1 (or 27/28 for some libraries)

	pubKey, err := crypto.SigToPub(messageHash, sigWithV)
	if err != nil {
		logger.Fatalf("failed to recover pubkey: %v", err)
	}

	// 5. Compute the address
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	fmt.Println("Recovered address: ", recoveredAddr.Hex())
}

func setEthereumMessage(message string) []byte {

	// Create Ethereum-style message hash
	// Ethereum uses: keccak256("\x19Ethereum Signed Message:\n" + len(message) + message)
	// For simplicity, we'll use SHA-256 here, but you can use keccak-256 if needed

	// Create the Ethereum-style message prefix
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))

	// Combine prefix and message
	prefixedMessage := append([]byte(prefix), []byte(message)...)

	// Hash the message
	hash := crypto.Keccak256(prefixedMessage)

	fmt.Println("message: ", message, "hash: ", hex.EncodeToString(hash[:]))

	// Convert hex string to big.Int
	hashBigInt := new(big.Int)
	hashBigInt.SetString(string(hash), 16)

	// Store the hash in the config
	common.TssCfg.Message = hashBigInt.String()
	return hash
}

func initLogLevel(cfg common.TssConfig) {
	log.SetLogLevel("main", cfg.LogLevel)
	log.SetLogLevel("tss", cfg.LogLevel)
	log.SetLogLevel("tss-lib", cfg.LogLevel)
	log.SetLogLevel("srv", cfg.LogLevel)
	log.SetLogLevel("trans", cfg.LogLevel)
	log.SetLogLevel("p2p_utils", cfg.LogLevel)
	log.SetLogLevel("common", cfg.LogLevel)

	// libp2p loggers
	log.SetLogLevel("dht", "error")
	log.SetLogLevel("discovery", "error")
	log.SetLogLevel("swarm2", "error")
	log.SetLogLevel("stream-upgrader", "error")
}
