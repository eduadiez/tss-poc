package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/eduadiez/tss-poc/tsslib"
	"github.com/ipfs/go-log"
)

func main() {
	// Define command-line flags
	var (
		home            = flag.String("home", "test1", "Home directory for configuration")
		vault           = flag.String("vault", "default", "Vault name")
		password        = flag.String("password", "123456789", "Password for the vault")
		channelID       = flag.String("channelId", "1116C145287", "Channel ID")
		channelPassword = flag.String("channelPassword", "123456789", "Channel Password")
		message         = flag.String("message", "123456789", "Message to sign")
		logLevel        = flag.String("logLevel", "info", "Log level (debug, info, warn, error)")
		outputJSON      = flag.Bool("json", false, "Output results as JSON")
	)

	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s [flags]\n\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Flags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(flag.CommandLine.Output(), "\nExamples:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -home=test2 -vault=default -password=123456789 -message=\"Hello World\"\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "  %s -home=test3 -vault=myvault -password=mypassword -message=\"Test message\" -json\n", os.Args[0])
	}

	// Parse command-line flags
	flag.Parse()

	// Create configuration
	config := tsslib.TSSConfig{
		Home:            *home,
		Vault:           *vault,
		Password:        *password,
		ChannelID:       *channelID,
		ChannelPassword: *channelPassword,
		Message:         *message,
		LogLevel:        *logLevel,
	}

	// Initialize logger
	logger := log.Logger("tss-signer")
	log.SetLogLevel("tss-signer", *logLevel)

	logger.Infof("Starting TSS signing with home=%s, vault=%s, message=%s", *home, *vault, *message)

	// Perform TSS signing
	result, err := tsslib.SignMessage(config)
	if err != nil {
		logger.Errorf("Failed to sign message: %v", err)
		os.Exit(1)
	}

	// Output results
	if *outputJSON {
		// Output as JSON
		jsonResult := map[string]interface{}{
			"signature":     result.Signature,
			"recoveredAddr": result.RecoveredAddr,
			"messageHash":   result.MessageHash,
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
	}

	logger.Info("TSS signing completed successfully")
}
