package auth

import (
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func resetWalletAddress() {
	walletAddress = ""
	_ = os.Remove(addressFile)
}

func TestMain(m *testing.M) {
	resetWalletAddress()

	code := m.Run()

	resetWalletAddress()

	os.Exit(code)
}

func mockStdin(input string) (*os.File, func()) {
	tmpFile, err := os.CreateTemp("", "mockStdin")
	if err != nil {
		panic(err)
	}

	_, err = tmpFile.Write([]byte(input))
	if err != nil {
		panic(err)
	}

	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		panic(err)
	}

	originalStdin := os.Stdin
	os.Stdin = tmpFile

	return tmpFile, func() {
		os.Stdin = originalStdin
		tmpFile.Close()
		os.Remove(tmpFile.Name())
	}
}

func TestValidateAddress(t *testing.T) {
	assert := assert.New(t)

	validAddress := "0x" + common.Bytes2Hex(make([]byte, 20))
	err := validateAddress(validAddress)
	assert.NoError(err, "Expected no error for valid address")

	invalidAddress := "1x" + common.Bytes2Hex(make([]byte, 20))
	err = validateAddress(invalidAddress)
	assert.Error(err, "Expected error for invalid prefix")

	shortAddress := "0x1234"
	err = validateAddress(shortAddress)
	assert.Error(err, "Expected error for invalid length")
}

func TestSaveAndLoadWalletAddress(t *testing.T) {
	assert := assert.New(t)

	testAddress := "0x" + common.Bytes2Hex(make([]byte, 20))
	saveWalletAddress(testAddress)

	loadWalletAddress()

	assert.Equal(testAddress, walletAddress, "Expected wallet address to match the saved address")
}

func TestConnectWalletWithExistingAddress(t *testing.T) {
	assert := assert.New(t)

	testAddress := "0x" + common.Bytes2Hex(make([]byte, 20))
	saveWalletAddress(testAddress)

	mockInput := "n\n"
	_, cleanup := mockStdin(mockInput)
	defer cleanup()

	ConnectWallet()
	assert.Equal(testAddress, walletAddress, "Expected wallet address not to be overwritten")
}

func TestConnectWalletWithNewAddress(t *testing.T) {
	assert := assert.New(t)

	newAddress := "0x" + common.Bytes2Hex(make([]byte, 20))
	mockInput := newAddress + "\n"
	_, cleanup := mockStdin(mockInput)
	defer cleanup()

	ConnectWallet()
	assert.Equal(newAddress, walletAddress, "Expected wallet address to match the new address")
}

func TestGetWalletAddress(t *testing.T) {
	assert := assert.New(t)

	testAddress := "0x" + common.Bytes2Hex(make([]byte, 20))
	saveWalletAddress(testAddress)

	loadedAddress := GetWalletAddress()
	assert.Equal(testAddress, loadedAddress, "Expected wallet address to match the loaded address")
}
