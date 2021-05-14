package alephium

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

// GetWalletStatus
// Implementation will be updated with PR #200
func (a *AlephiumClient) GetWalletStatus(walletName string) (WalletInfo, error) {
	wallets, err := a.GetWallets()
	if err != nil {
		return WalletInfo{}, err
	}
	for _, w := range wallets {
		if w.Name == walletName {
			return w, nil
		}
	}
	return WalletInfo{}, ErrorDetail{Detail: walletName + " not found", Resource: walletName}
}

// LockWallet
func (a *AlephiumClient) LockWallet(walletName string) (bool, error) {

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/lock").
		Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

type WalletPasswordRequestBody struct {
	Password string `json:"password"`
}

// UnlockWallet
func (a *AlephiumClient) UnlockWallet(walletName string, password string) (bool, error) {

	body := WalletPasswordRequestBody{Password: password}

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
	var address Address
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/deriveNextAddress").
		Receive(&address, &errorDetail)

	return address, relevantError(err, errorDetail)
}

type AddressBodyRequest struct {
	Address string `json:"address"`
}

// ChangeActiveAddress
func (a *AlephiumClient) ChangeActiveAddress(walletName string, activeAddress string) (bool, error) {

	body := AddressBodyRequest{Address: activeAddress}

	//var address string
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Post("wallets/"+walletName+"/changeActiveAddress").
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}

// DeleteWallet
func (a *AlephiumClient) DeleteWallet(walletName string, walletPassword string) (bool, error) {

	body := WalletPasswordRequestBody{Password: walletPassword}

	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Delete("wallets/"+walletName).
		BodyJSON(body).Receive(nil, &errorDetail)

	return true, relevantError(err, errorDetail)
}
