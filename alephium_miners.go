package alephium

import (
	"fmt"
)

// StartMining
func (a *Client) StartMining() (bool, error) {
	return a.miningAction("start-mining")
}

// StopMining
func (a *Client) StopMining() (bool, error) {
	return a.miningAction("stop-mining")
}

type MiningActionRequestParams struct {
	Action string `url:"action"`
}

func (a *Client) miningAction(action string) (bool, error) {

	var errorDetail ErrorDetail
	var actionOk bool
	params := MiningActionRequestParams{
		Action: action,
	}
	_, err := a.slingClient.New().Post("miners").
		QueryStruct(params).Receive(&actionOk, &errorDetail)

	return actionOk, relevantError(err, errorDetail)
}

type UpdateMinersAddressesBodyParams struct {
	Addresses []string `json:"addresses"`
}

// UpdateMinersAddresses
func (a *Client) UpdateMinersAddresses(addresses []string) error {

	var errorDetail ErrorDetail
	params := UpdateMinersAddressesBodyParams{
		Addresses: addresses,
	}
	_, err := a.slingClient.New().Put("miners/addresses").
		BodyJSON(params).Receive(nil, &errorDetail)

	return relevantError(err, errorDetail)
}

// GetMinersAddresses
func (a *Client) GetMinersAddresses() (MinersAddresses, error) {

	var minersAddresses MinersAddresses
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("miners/addresses").
		Receive(&minersAddresses, &errorDetail)

	return minersAddresses, relevantError(err, errorDetail)
}

// GetBlockCandidate
func (a *Client) GetBlockCandidate() error {
	return fmt.Errorf("not implemented yet")
}

// SubmitNewBlock
func (a *Client) SubmitNewBlock() error {
	return fmt.Errorf("not implemented yet")
}
