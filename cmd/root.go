package cmd

import (
	"eth-chat/auth"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "eth-chat",
	Short: "Eth-Chat is a decentralized messaging application",
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to an Ethereum wallet",
	Run: func(cmd *cobra.Command, args []string) {
		auth.ConnectWallet()
	},
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check the current connection status",
	Run: func(cmd *cobra.Command, args []string) {
		walletAddress := auth.GetWalletAddress()
		if walletAddress == "" {
			fmt.Println("No wallet connected.")
		} else {
			fmt.Printf("Wallet connected: %s\n", walletAddress)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
	rootCmd.AddCommand(statusCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
