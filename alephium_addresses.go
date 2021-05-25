package alephium

// GetAddressBalance
func (a *Client) GetAddressBalance(address string) (AddressUtxoBalance, error) {
	var addressBalance AddressUtxoBalance
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("addresses/"+address+"/balance").
		Receive(&addressBalance, &errorDetail)
	return addressBalance, relevantError(err, errorDetail)
}

// GetAddressGroup
func (a *Client) GetAddressGroup(address string) (AddressGroup, error) {
	var addressGroup AddressGroup
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("addresses/"+address+"/group").
		Receive(&addressGroup, &errorDetail)
	return addressGroup, relevantError(err, errorDetail)
}
