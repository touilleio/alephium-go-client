package alephium

import (
	"fmt"
)

// GetUnconfirmedTransactions
func (a *AlephiumClient) GetUnconfirmedTransactions() error {
	return fmt.Errorf("not implemented yet")
}

// BuildTransaction
func (a *AlephiumClient) BuildTransaction(hash string) error {
	return fmt.Errorf("not implemented yet")
}

// SendTransaction
func (a *AlephiumClient) SendTransaction(transactionId string) error {
	return fmt.Errorf("not implemented yet")
}

type TransactionStatusRequestParams struct {
	TransactionId string `url:"txId"`
	FromGroup     int    `url:"fromGroup"`
	ToGroup       int    `url:"toGroup"`
}

// GetTransactionStatus
func (a *AlephiumClient) GetTransactionStatus(transactionId string, fromGroup int, toGroup int) (TransactionStatus, error) {

	var transactionStatus TransactionStatus
	var errorDetail ErrorDetail

	params := TransactionStatusRequestParams{
		TransactionId: transactionId,
		FromGroup:     fromGroup,
		ToGroup:       toGroup,
	}
	_, err := a.slingClient.New().Get("transactions/status").
		QueryStruct(params).Receive(&transactionStatus, &errorDetail)

	return transactionStatus, relevantError(err, errorDetail)
}
