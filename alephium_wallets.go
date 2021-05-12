package alephium

import (
	"fmt"
)

// GetWallets
func (a *AlephiumClient) GetWallets() ([]WalletInfo, error) {
	var wallets []WalletInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets").
		Receive(&wallets, &errorDetail)

	return wallets, relevantError(err, errorDetail)
}

type CreateWalletRequestBody struct {
	Password           string `json:"password"`
	IsMiner            bool   `json:"isMiner"`
	WalletName         string `json:"walletName,omitempty"`
	MnemonicPassphrase string `json:"mnemonicPassphrase,omitempty"`
	mMnemonicSize      int    `json:"mnemonicSize,omitempty"`
}

// CreateWallet
func (a *AlephiumClient) CreateWallet(walletName string, password string, isMiner bool, mnemonicPassphrase string) (WalletCreate, error) {

	body := CreateWalletRequestBody{
		Password:           password,
		IsMiner:            isMiner,
		WalletName:         walletName,
		MnemonicPassphrase: mnemonicPassphrase,
	}

	var wallet WalletCreate
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets").
		BodyJSON(body).Receive(&wallet, &errorDetail)

	return wallet, relevantError(err, errorDetail)
}

type RestoreWalletRequestBody struct {
	Password string `json:"password"`
	Mnemonic string `json:"mnemonic"`
}

// RestoreWallet
func (a *AlephiumClient) RestoreWallet(password string, mnemonic string) (Wallet, error) {

	body := RestoreWalletRequestBody{
		Password: password,
		Mnemonic: mnemonic,
	}

	var wallet Wallet
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Put("wallets").
		BodyJSON(body).Receive(&wallet, &errorDetail)

	return wallet, relevantError(err, errorDetail)
}

// LockWallet
func (a *AlephiumClient) LockWallet(walletName string) (bool, error) {

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/lock").
		Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

type UnlockWalletRequestBody struct {
	Password string `json:"password"`
}

// UnlockWallet
func (a *AlephiumClient) UnlockWallet(walletName string, password string) (bool, error) {

	body := UnlockWalletRequestBody{Password: password}

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/unlock").
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

// GetWalletBalances
func (a *AlephiumClient) GetWalletBalances(walletName string) (WalletBalances, error) {

	var walletBalances WalletBalances
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/balances").
		Receive(&walletBalances, &errorDetail)

	return walletBalances, relevantError(err, errorDetail)
}

// GetWalletAddresses
func (a *AlephiumClient) GetWalletAddresses(walletName string) (WalletAddresses, error) {

	var walletAddresses WalletAddresses
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/addresses").
		Receive(&walletAddresses, &errorDetail)

	return walletAddresses, relevantError(err, errorDetail)
}

type TransferRequest struct {
	Address string `json:"address"`
	Amount  string `json:"amount"`
}

// Transfer
func (a *AlephiumClient) Transfer(walletName string, address string, amount string) (Transaction, error) {

	// TODO: run sanity check on address and amount
	body := TransferRequest{Address: address, Amount: amount}

	var transaction Transaction
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/transfer").
		BodyJSON(body).Receive(&transaction, &errorDetail)

	return transaction, relevantError(err, errorDetail)
}

// DeriveNextAddress
func (a *AlephiumClient) DeriveNextAddress(walletName string) (Address, error) {
	return Address{}, fmt.Errorf("not implemented yet")
}

// ChangeActiveAddress
func (a *AlephiumClient) ChangeActiveAddress(walletName string) (Address, error) {
	return Address{}, fmt.Errorf("not implemented yet")
}
