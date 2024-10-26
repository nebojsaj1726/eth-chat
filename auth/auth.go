package auth

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

var walletAddress string

const addressFile = "address.json"

func ConnectWallet() {
	fmt.Println("Hello, welcome to Eth-Chat!")
	fmt.Println("Please connect to your Ethereum wallet:")

	loadWalletAddress()

	reader := bufio.NewReader(os.Stdin)
	if walletAddress != "" {
		fmt.Printf("You already have a connected address: %s\n", walletAddress)
		fmt.Print("Do you want to overwrite it? (y/n): ")
		confirmation, _ := reader.ReadString('\n')
		confirmation = strings.TrimSpace(confirmation)
		if confirmation != "y" {
			fmt.Println("Wallet address not changed.")
			return
		}
	}

	for {
		fmt.Print("Enter your wallet address (0x...): ")
		walletAddressInput, _ := reader.ReadString('\n')
		walletAddressInput = strings.TrimSpace(walletAddressInput)

		err := validateAddress(walletAddressInput)
		if err != nil {
			fmt.Println("Invalid address:", err)
			continue
		}

		walletAddress = walletAddressInput
		saveWalletAddress(walletAddress)
		fmt.Println("Wallet address is connected:", walletAddress)
		return
	}
}

func GetWalletAddress() string {
	if walletAddress == "" {
		loadWalletAddress()
	}
	return walletAddress
}

func loadWalletAddress() {
	file, err := os.Open(addressFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Error opening wallet file:", err)
		return
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&walletAddress); err != nil {
		fmt.Println("Error decoding wallet address:", err)
	}
}

func saveWalletAddress(address string) {
	file, err := os.Create(addressFile)
	if err != nil {
		fmt.Println("Error creating wallet file:", err)
		return
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(address); err != nil {
		fmt.Println("Error encoding wallet address:", err)
	}
}

func validateAddress(address string) error {
	if !strings.HasPrefix(address, "0x") {
		return errors.New("address must start with '0x'")
	}
	if len(address) != 42 {
		return errors.New("address must be 42 characters long, including '0x'")
	}

	if !common.IsHexAddress(address) {
		return errors.New("address contains invalid characters")
	}

	return nil
}
