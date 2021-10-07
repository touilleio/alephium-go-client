package alephium

import (
	"strings"
)

// GetWallets
func (a *Client) GetWallets() ([]WalletInfo, error) {
	var wallets []WalletInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets").
		Receive(&wallets, &errorDetail)

	return wallets, relevantError(err, errorDetail)
}

type CreateWalletRequestBody struct {
	Password           string `json:"password"`
	IsMiner            bool   `json:"isMiner"`
	WalletName         string `json:"walletName"`
	MnemonicPassphrase string `json:"mnemonicPassphrase,omitempty"`
	mMnemonicSize      int    `json:"mnemonicSize,omitempty"`
}

// CreateWallet
func (a *Client) CreateWallet(walletName string, password string, isMiner bool, mnemonicPassphrase string) (WalletCreate, error) {

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
	Password           string `json:"password"`
	Mnemonic           string `json:"mnemonic"`
	IsMiner            bool   `json:"isMiner,omitempty"`
	WalletName         string `json:"walletName,omitempty"`
	MnemonicPassphrase string `json:"mnemonicPassphrase,omitempty"`
}

// RestoreWallet
func (a *Client) RestoreWallet(password string, mnemonic string, walletName string,
	isMiner bool, mnemonicPassphrase string) (Wallet, error) {

	body := RestoreWalletRequestBody{
		Password:           password,
		Mnemonic:           mnemonic,
		WalletName:         walletName,
		IsMiner:            isMiner,
		MnemonicPassphrase: mnemonicPassphrase,
	}

	var wallet Wallet
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Put("wallets").
		BodyJSON(body).Receive(&wallet, &errorDetail)

	return wallet, relevantError(err, errorDetail)
}

// GetWalletStatus
func (a *Client) GetWalletStatus(walletName string) (WalletInfo, error) {
	var walletInfo WalletInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName).
		Receive(&walletInfo, &errorDetail)
	return walletInfo, relevantError(err, errorDetail)
}

// LockWallet
func (a *Client) LockWallet(walletName string) (bool, error) {

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/lock").
		Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

type WalletPasswordRequestBody struct {
	Password string `json:"password"`
}

// UnlockWallet
func (a *Client) UnlockWallet(walletName string, password string) (bool, error) {

	body := WalletPasswordRequestBody{Password: password}

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/unlock").
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

// GetWalletBalances
func (a *Client) GetWalletBalances(walletName string) (WalletBalances, error) {

	var walletBalances WalletBalances
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/balances").
		Receive(&walletBalances, &errorDetail)

	return walletBalances, relevantError(err, errorDetail)
}

// GetWalletAddresses
func (a *Client) GetWalletAddresses(walletName string) (WalletAddresses, error) {

	var walletAddresses WalletAddresses
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/addresses").
		Receive(&walletAddresses, &errorDetail)

	return walletAddresses, relevantError(err, errorDetail)
}

type TransferRequest struct {
	Destinations []TransferDestination `json:"destinations"`
}

type TransferDestination struct {
	Address string          `json:"address"`
	Amount  ALF             `json:"amount"`
	Tokens  []TransferToken `json:"tokens,omitempty"`
}

type TransferToken struct {
	Id     string `json:"id"`
	Amount string `json:"amount"`
}

// Transfer
func (a *Client) Transfer(walletName string, address string, amount ALF) (Transaction, error) {

	// TODO: run sanity check on address and amount
	body := TransferRequest{Destinations: []TransferDestination{{Address: address, Amount: amount}}}

	var transaction Transaction
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/transfer").
		BodyJSON(body).Receive(&transaction, &errorDetail)

	return transaction, relevantError(err, errorDetail)
}

type SweepAllRequest struct {
	Address string `json:"toAddress"`
}

// SweepAll
func (a *Client) SweepAll(walletName string, toAddress string) (Transaction, error) {

	// TODO: run sanity check on address
	body := SweepAllRequest{Address: toAddress}

	var transaction Transaction
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/sweep-all").
		BodyJSON(body).Receive(&transaction, &errorDetail)

	return transaction, relevantError(err, errorDetail)
}

// DeriveNextAddress
func (a *Client) DeriveNextAddress(walletName string) (Address, error) {
	var address Address
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/derive-next-address").
		Receive(&address, &errorDetail)

	return address, relevantError(err, errorDetail)
}

type AddressBodyRequest struct {
	Address string `json:"address"`
}

// ChangeActiveAddress
func (a *Client) ChangeActiveAddress(walletName string, activeAddress string) (bool, error) {

	body := AddressBodyRequest{Address: activeAddress}

	//var address string
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/change-active-address").
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

// DeleteWallet
func (a *Client) DeleteWallet(walletName string, walletPassword string) (bool, error) {

	body := WalletPasswordRequestBody{Password: walletPassword}

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Delete("wallets/"+walletName).
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

// CheckWalletExist is a convenience function which checks if the wallet exists,
// since this information is based on the error string returned by the API call.
// TODO: should the "not found" exception being typed?
func (a *Client) CheckWalletExist(walletName string) (bool, error) {
	_, err := a.GetWalletStatus(walletName)
	if err != nil {
		if strings.HasPrefix(walletName+" not found", err.Error()) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func GetAddressesAsString(walletAddresses []WalletAddress) []string {
	addresses := make([]string, 0, len(walletAddresses))
	for _, wa := range walletAddresses {
		addresses = append(addresses, wa.Address)
	}
	return addresses
}
