package alephium

import (
	"github.com/sqooba/go-common/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransactions(t *testing.T) {

	log := logging.NewLogger()
	logging.SetLogLevel(log, "debug")
	alephiumClient, err := New("http://localhost:12973", log)
	assert.Nil(t, err)

	//txId := "56f35934e4005fff7f6c6529848c8d6314f9705d49ac8e0df6f565c64f4cdb37"
	txId := "1113eabfeae718adc55fa12462fdeab38b30677b9493f7a0782015ca9afe6f96"
	fromGroup := 1
	toGroup := 1

	transactionStatus, err := alephiumClient.GetTransactionStatus(txId, fromGroup, toGroup)
	log.Debugf("transactionStatus is %v", transactionStatus)

	err = alephiumClient.WaitForTransactionConfirmed(txId, fromGroup, toGroup)
	assert.Nil(t, err)
}
