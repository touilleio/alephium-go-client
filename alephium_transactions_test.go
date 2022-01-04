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

func TestTransactionE2E(t *testing.T) {

	log := logging.NewLogger()
	_ = logging.SetLogLevel(log, "debug")

	ctx, _ := context.WithCancel(context.Background())
	alephiumNode, port, walletFolder, err := setupContainer(ctx)
	defer tearDownContainer(ctx, alephiumNode, walletFolder)

	walletPassword := "dummy-password"
	//walletName := "test-wallet"
	alephiumClient, err := NewWithApiKey("http://localhost:"+port.Port(), TestApiKey, log)
	assert.Nil(t, err)

	sync, err := alephiumClient.WaitUntilSyncedWithAtLeastOnePeer(context.Background())
	assert.Nil(t, err)
	assert.True(t, sync)

	genesisWallet, err := alephiumClient.RestoreWallet(walletPassword, TestGenesisWalletMnemonics, TestGenesisWalletName, true, "")
	assert.Nil(t, err)

	_, err = alephiumClient.UnlockWallet(genesisWallet.Name, walletPassword, "")
	assert.Nil(t, err)

	walletAddress, err := alephiumClient.GetWalletAddresses(genesisWallet.Name)
	assert.Nil(t, err)

	walletAddressDetail, err := alephiumClient.GetWalletAddressDetail(genesisWallet.Name, walletAddress.ActiveAddress)
	assert.Nil(t, err)

	//log.Infof("http://localhost:"+port.Port())
	//time.Sleep(24 * time.Hour)

	amount1, _ := ALPHFromALPHString("2.60")
	amount2, _ := ALPHFromALPHString("4.44")
	unsignedTx, err := alephiumClient.BuildTransaction(walletAddressDetail.PublicKey, []TransactionDestination{
		{
			Address: "16FnqysnYf7qE6Xx1ZFeCixYFUwNKATTvRAArh3SD7w3S",
			Amount:  amount1,
		},
		{
			Address: "1AjSsNMLZwqgN7VSisVn5ZFESXaBb25ydyR41AXTK1Xvk",
			Amount:  amount2,
		},
	})
	assert.Nil(t, err)
	log.Infof("%s txId=%s (%d -> %d)", unsignedTx.UnsignedTx, unsignedTx.TxId, unsignedTx.FromGroup, unsignedTx.ToGroup)

	signature, err := alephiumClient.Sign(genesisWallet.Name, unsignedTx.TxId)
	assert.Nil(t, err)
	log.Infof("%s", signature)

	tx, err := alephiumClient.SubmitTransaction(unsignedTx.UnsignedTx, signature)
	assert.Nil(t, err)
	log.Infof("%s (%d -> %d)", tx.TransactionId, tx.FromGroup, tx.ToGroup)

	//ok, err := alephiumClient.WaitForTransactionConfirmed(context.Background(), tx.TransactionId, tx.FromGroup, tx.ToGroup)
	//assert.Nil(t, err)
	//assert.True(t, ok)
}
