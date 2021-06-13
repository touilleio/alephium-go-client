package alephium

import (
	"context"
	"github.com/sqooba/go-common/logging"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"testing"
)

func TestCreateWalletE2E(t *testing.T) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "touilleio/alephium:v0.7.8",
		ExposedPorts: []string{"12973/tcp"},
		WaitingFor:   wait.ForListeningPort("12973/tcp"),
	}
	alephiumNode, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	defer alephiumNode.Terminate(ctx)
	port, err := alephiumNode.MappedPort(ctx, "12973/tcp")
	if err != nil {
		t.Error(err)
	}

	log := logging.NewLogger()
	logging.SetLogLevel(log, "debug")
	walletPassword := "dummy-password"
	alephiumClient, err := New("http://localhost:"+port.Port(), log)
	assert.Nil(t, err)

	newWallet, err := alephiumClient.CreateWallet("", walletPassword, true, "")
	assert.Nil(t, err)

	log.Printf("name: %s, mnemonic: %s\n", newWallet.Name, newWallet.Mnemonic)

	wallets, err := alephiumClient.GetWallets()
	assert.Nil(t, err)

	foundNewWallet := false
	for _, wallet := range wallets {
		if wallet.Name == newWallet.Name {
			foundNewWallet = true
		}
	}
	assert.True(t, foundNewWallet)

	walletAddresses, err := alephiumClient.GetWalletAddresses(newWallet.Name)
	assert.Nil(t, err)

	log.Printf("name: %s, activeAddress: %s, addresses: %s\n", newWallet.Name, walletAddresses.ActiveAddress, walletAddresses.Addresses)

	walletBalances, err := alephiumClient.GetWalletBalances(newWallet.Name)
	assert.Nil(t, err)

	totalBalance, ok := walletBalances.GetTotalBalance()
	assert.True(t, ok)

	log.Printf("name: %s, total balance: %d\n", newWallet.Name, totalBalance)

	ok, err = alephiumClient.DeleteWallet(newWallet.Name, walletPassword)
	assert.True(t, ok)

	restoredWallet, err := alephiumClient.RestoreWallet(walletPassword, newWallet.Mnemonic, newWallet.Name, false, "")
	assert.Nil(t, err)

	walletAddresses, err = alephiumClient.GetWalletAddresses(restoredWallet.Name)
	assert.Nil(t, err)
	assert.Contains(t, walletAddresses.Addresses, walletAddresses.ActiveAddress)
	log.Printf("name: %s, activeAddress: %s, addresses: %s\n", restoredWallet.Name, walletAddresses.ActiveAddress, walletAddresses.Addresses)

	activeAddress, err := alephiumClient.ChangeActiveAddress(restoredWallet.Name, walletAddresses.ActiveAddress)
	assert.Nil(t, err)
	assert.Contains(t, walletAddresses.Addresses, activeAddress)
	log.Printf("name: %s, new activeAddress: %s, addresses: %s\n", restoredWallet.Name, walletAddresses.ActiveAddress, walletAddresses.Addresses)

	derivedAddress, err := alephiumClient.DeriveNextAddress(restoredWallet.Name)
	assert.Nil(t, err)
	log.Printf("derived next address: %s\n", derivedAddress)

	locked, err := alephiumClient.LockWallet(restoredWallet.Name)
	assert.Nil(t, err)
	assert.True(t, locked)

	walletStatus, err := alephiumClient.GetWalletStatus(restoredWallet.Name)
	assert.Nil(t, err)
	assert.True(t, walletStatus.Locked)
	assert.Equal(t, restoredWallet.Name, walletStatus.Name)

	unlocked, err := alephiumClient.UnlockWallet(restoredWallet.Name, walletPassword)
	assert.Nil(t, err)
	assert.True(t, unlocked)

	walletStatus, err = alephiumClient.GetWalletStatus(restoredWallet.Name)
	assert.Nil(t, err)
	assert.False(t, walletStatus.Locked)
	assert.Equal(t, restoredWallet.Name, walletStatus.Name)
}
