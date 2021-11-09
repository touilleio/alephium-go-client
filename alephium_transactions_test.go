package alephium

import (
	"context"
	"github.com/sqooba/go-common/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactions(t *testing.T) {

	log := logging.NewLogger()
	logging.SetLogLevel(log, "debug")
	alephiumClient, err := New("http://localhost:12973", log)
	assert.Nil(t, err)

	txId := "2822de6f3d8754ce217593585bb9a69f00f5ed3c421fa2198a561cbee02f2584"
	fromGroup := 1
	toGroup := 1

	transactionStatus, err := alephiumClient.GetTransactionStatus(txId, fromGroup, toGroup)
	log.Debugf("transactionStatus is %v", transactionStatus)

	ok, err := alephiumClient.WaitForTransactionConfirmed(context.Background(), txId, fromGroup, toGroup)
	assert.Nil(t, err)
	assert.True(t, ok)
}
