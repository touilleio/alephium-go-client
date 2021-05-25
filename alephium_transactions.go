package alephium

import (
	"fmt"
	"time"
)

const (
	txConfirmed = "confirmed"
)

// GetUnconfirmedTransactions
func (a *Client) GetUnconfirmedTransactions() error {
	return fmt.Errorf("not implemented yet")
}

// BuildTransaction
func (a *Client) BuildTransaction(hash string) error {
	return fmt.Errorf("not implemented yet")
}

// SendTransaction
func (a *Client) SendTransaction(transactionId string) error {
	return fmt.Errorf("not implemented yet")
}

type TransactionStatusRequestParams struct {
	TransactionId string `url:"txId"`
	FromGroup     int    `url:"fromGroup"`
	ToGroup       int    `url:"toGroup"`
}

// GetTransactionStatus
func (a *Client) GetTransactionStatus(transactionId string, fromGroup int, toGroup int) (TransactionStatus, error) {

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

// WaitForTransactionConfirmed
func (a *Client) WaitForTransactionConfirmed(transactionId string, fromGroup int, toGroup int) error {
	return a.WaitForTransactionStatus(txConfirmed, transactionId, fromGroup, toGroup)
}

// WaitForTransactionStatus
func (a *Client) WaitForTransactionStatus(status string, transactionId string, fromGroup int, toGroup int) error {
	txStatus := "unknown"
	for {
		tx, err := a.GetTransactionStatus(transactionId, fromGroup, toGroup)
		if err != nil {
			return err
		}
		txStatus = tx.Type
		if txStatus == status {
			return nil
		} else {
			a.log.Debugf("Tx %s not %s yet, sleeping %s", transactionId, status, a.sleepTime)
			time.Sleep(a.sleepTime)
		}
	}
}
