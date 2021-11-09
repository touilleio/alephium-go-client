package alephium

// GetAddressBalance returns the balance (computed) of the address
func (a *Client) GetAddressBalance(address string, utxosLimit int) (AddressUtxoBalance, error) {
	params := &UtxosLimit{}
	if utxosLimit > 0 {
		params.UtxosLimit = utxosLimit
	}
	var addressBalance AddressUtxoBalance
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("addresses/"+address+"/balance").
		QueryStruct(params).Receive(&addressBalance, &errorDetail)
	return addressBalance, relevantError(err, errorDetail)
}

// GetAddressGroup returns the group of the address
func (a *Client) GetAddressGroup(address string) (AddressGroup, error) {
	var addressGroup AddressGroup
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("addresses/"+address+"/group").
		Receive(&addressGroup, &errorDetail)
	return addressGroup, relevantError(err, errorDetail)
}

type AddressUtxosList struct {
	Utxos []struct {
		Ref struct {
			Hint int    `json:"hint"`
			Key  string `json:"key"`
		} `json:"ref"`
		Amount string `json:"amount"` // ALPH
		Tokens []struct {
			Id     string `json:"id"`
			Amount string `json:"amount"` // ALPH
		} `json:"tokens"`
		LockTime       int64  `json:"lockTime"`
		AdditionalData string `json:"additionalData"`
	} `json:"utxos"`
}

// GetAddressUtxos returns the UTXOs of the address
func (a *Client) GetAddressUtxos(address string, utxosLimit int) (AddressUtxosList, error) {
	params := &UtxosLimit{}
	if utxosLimit > 0 {
		params.UtxosLimit = utxosLimit
	}
	var utxosList AddressUtxosList
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("addresses/"+address+"/utxos").
		QueryStruct(params).Receive(&utxosList, &errorDetail)
	return utxosList, relevantError(err, errorDetail)
}

type UtxosLimit struct {
	UtxosLimit int `url:"utxosLimit,omitempty"`
}
