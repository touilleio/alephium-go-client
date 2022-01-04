package alephium

import (
	"context"
	"time"
)

const (
	TxConfirmed = "confirmed"
)

// GetUnconfirmedTransactions gets the list of unconfirmed transactions
func (a *Client) GetUnconfirmedTransactions() error {

	var errorDetail ErrorDetail

	_, err := a.slingClient.New().Get("transactions/unconfirmed").
		Receive(nil, &errorDetail)

	return relevantError(err, errorDetail)
}

type BuildTransactionBodyRequest struct {
	FromPublicKey string                   `json:"fromPublicKey"`
	Destinations  []TransactionDestination `json:"destinations"`
}

type TransactionDestination struct {
	Address string `json:"address"`
	Amount  ALPH   `json:"amount"`
}

type UnsignedTransaction struct {
	UnsignedTx string `json:"unsignedTx"`
	TxId       string `json:"txId"`
	FromGroup  int    `json:"fromGroup"`
	ToGroup    int    `json:"toGroup"`
}

// BuildTransaction builds an unsigned transaction
func (a *Client) BuildTransaction(publicKey string, destinations []TransactionDestination) (UnsignedTransaction, error) {

	var unsignedTx UnsignedTransaction
	var errorDetail ErrorDetail

	body := BuildTransactionBodyRequest{
		FromPublicKey: publicKey,
		Destinations:  destinations,
	}

	_, err := a.slingClient.New().Post("transactions/build").
		BodyJSON(body).Receive(&unsignedTx, &errorDetail)

	return unsignedTx, relevantError(err, errorDetail)
}

type SubmitTransactionBodyRequest struct {
	UnsignedTx string `json:"unsignedTx"`
	Signature  string `json:"signature"`
}

// SubmitTransaction submit a previously built and signed transaction
func (a *Client) SubmitTransaction(unsignedTxId string, signature string) (Transaction, error) {

	var tx Transaction
	var errorDetail ErrorDetail

	params := SubmitTransactionBodyRequest{
		UnsignedTx: unsignedTxId,
		Signature:  signature,
	}
	_, err := a.slingClient.New().Post("transactions/submit").
		BodyJSON(params).Receive(&tx, &errorDetail)

	return tx, relevantError(err, errorDetail)
}

type TransactionStatusRequestParams struct {
	TransactionId string `url:"txId"`
	FromGroup     int    `url:"fromGroup"`
	ToGroup       int    `url:"toGroup"`
}

// GetTransactionStatus gets the status of a given transaction
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

// WaitForTransactionConfirmed waits until the transaction is confirmed
func (a *Client) WaitForTransactionConfirmed(ctx context.Context, transactionId string, fromGroup int, toGroup int) (bool, error) {
	return a.WaitForTransactionStatus(ctx, TxConfirmed, transactionId, fromGroup, toGroup)
}

// WaitForTransactionStatus waits until the transaction is in a given status
func (a *Client) WaitForTransactionStatus(ctx context.Context, status string, transactionId string, fromGroup int, toGroup int) (bool, error) {
	txStatus := "unknown"
	for {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:

		}
		tx, err := a.GetTransactionStatus(transactionId, fromGroup, toGroup)
		if err != nil {
			return false, err
		}
		txStatus = tx.Type
		if txStatus == status {
			return true, nil
		} else {
			a.log.Debugf("Tx %s not %s yet, sleeping %s", transactionId, status, a.sleepTime)
			time.Sleep(a.sleepTime)
		}
	}
}
