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
		Image:        "touilleio/alephium:v0.7.5",
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
}
