package alephium

import (
	"strings"
)

// GetWallets returns the list of wallet present on the full node
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
	MnemonicSize       int    `json:"mnemonicSize,omitempty"`
}

// CreateWallet creates a new wallet, generating mnemonic while doing so
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

// RestoreWallet creates a wallet with provided mnemonics (unlike CreateWallet which generates new mnemonics)
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

// GetWalletStatus returns the status of a given wallet
func (a *Client) GetWalletStatus(walletName string) (WalletInfo, error) {
	var walletInfo WalletInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName).
		Receive(&walletInfo, &errorDetail)
	return walletInfo, relevantError(err, errorDetail)
}

// LockWallet locks a given wallet. Returns false if the wallet was already locked.
func (a *Client) LockWallet(walletName string) (bool, error) {

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/lock").
		Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

type WalletPasswordRequestBody struct {
	Password string `json:"password"`
	MnemonicPassphrase string `json:"mnemonicPassphrase,omitempty"`
}

// UnlockWallet unlocks wallet with the provided password and optional passphrase.
// Returns true if the wallet got successfully unlocked, false when the wallet was already unlocked
func (a *Client) UnlockWallet(walletName string, password string, mnemonicPassphrase string) (bool, error) {

	body := WalletPasswordRequestBody{
		Password: password,
		MnemonicPassphrase: mnemonicPassphrase,
	}

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/unlock").
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

// GetWalletBalances returns the balance of all the addresses inside the wallet.
func (a *Client) GetWalletBalances(walletName string) (WalletBalances, error) {

	var walletBalances WalletBalances
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/balances").
		Receive(&walletBalances, &errorDetail)

	return walletBalances, relevantError(err, errorDetail)
}

// GetWalletAddresses lists all the addresses from a wallet
func (a *Client) GetWalletAddresses(walletName string) (WalletAddresses, error) {

	var walletAddresses WalletAddresses
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/addresses").
		Receive(&walletAddresses, &errorDetail)

	return walletAddresses, relevantError(err, errorDetail)
}

// GetWalletAddressDetail returns detailed info about a specific address of a wallet
func (a *Client) GetWalletAddressDetail(walletName string, address string) (AddressDetailResponse, error) {

	var addressDetailResponse AddressDetailResponse
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/addresses/"+address).
		Receive(&addressDetailResponse, &errorDetail)

	return addressDetailResponse, relevantError(err, errorDetail)
}

type AddressDetailResponse struct {
	Address   string `json:"address"`
	PublicKey string `json:"publicKey"`
	Group     int    `json:"group"`
}

type TransferRequest struct {
	Destinations []TransferDestination `json:"destinations"`
}

type 	TransferDestination struct {
	Address string `json:"address"`
	Amount  ALPH   `json:"amount"`
}

type TransferToken struct {
	Id     string `json:"id"`
	Amount string `json:"amount"`
}

// Transfer transfers ALPH from one wallet to a given address
func (a *Client) Transfer(walletName string, address string, amount ALPH) (Transaction, error) {

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

// SweepAll transfers all (unlocked) ALPH from a wallet to another address
func (a *Client) SweepAll(walletName string, toAddress string) (Transaction, error) {

	// TODO: run sanity check on address
	body := SweepAllRequest{Address: toAddress}

	var transaction Transaction
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/sweep-all").
		BodyJSON(body).Receive(&transaction, &errorDetail)

	return transaction, relevantError(err, errorDetail)
}

type RevealMnemonicRequest	 struct {
	Password string `json:"password"`
}

type RevealMnemonicResponse struct {
	Mnemonic string `json:"mnemonic"`
}

// RevealWalletMnemonic reveals your mnemonic. Please use with caution!!
func (a *Client) RevealWalletMnemonic(walletName string, password string) (string, error) {

	body := RevealMnemonicRequest{Password: password}

	var response RevealMnemonicResponse
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Get("wallets/"+walletName+"/reveal-mnemonic").
		BodyJSON(body).Receive(&response, &errorDetail)

	return response.Mnemonic, relevantError(err, errorDetail)
}
type SignRequest struct {
	Data string `json:"data"`
}

type SignResponse struct {
	Signature string `json:"signature"`
}

// Sign signs the given data and returns the signature
func (a *Client) Sign(walletName string, data string) (string, error) {

	body := SignRequest{Data: data}

	var response SignResponse
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/sign").
		BodyJSON(body).Receive(&response, &errorDetail)

	return response.Signature, relevantError(err, errorDetail)
}

// DeriveNextAddress derives the next address
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

// ChangeActiveAddress changes the active address of the wallet. Has no effect on non-miner wallet.
func (a *Client) ChangeActiveAddress(walletName string, activeAddress string) (bool, error) {

	body := AddressBodyRequest{Address: activeAddress}

	//var address string
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/change-active-address").
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

// DeleteWallet deletes a wallet.
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

type MinerWalletAddresses struct {
	Addresses []WalletAddress `json:"addresses"`
}

// GetMinerWalletAddresses lists all the addresses from a miner wallet
func (a *Client) GetMinerWalletAddresses(walletName string) ([]MinerWalletAddresses, error) {

	var minerAddresses []MinerWalletAddresses
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("wallets/"+walletName+"/miner-addresses").
		Receive(&minerAddresses, &errorDetail)

	return minerAddresses, relevantError(err, errorDetail)
}

// DeriveNextMinerAddresses derives the next miner address
func (a *Client) DeriveNextMinerAddresses(walletName string) ([]WalletAddress, error) {
	var addresses []WalletAddress
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/derive-next-miner-addresses").
		Receive(&addresses, &errorDetail)

	return addresses, relevantError(err, errorDetail)
}
