package alephium

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/sqooba/go-common/logging"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func setupContainer(ctx context.Context) (testcontainers.Container, nat.Port, string, error) {

	configFile, err := filepath.Abs("./user-dev-standalone.conf")
	if err != nil {
		return nil, "", "", err
	}
	walletFolder, err := ioutil.TempDir("./test-data", "wallets")
	if err != nil {
		return nil, "", "", err
	}
	walletFolder, err = filepath.Abs(walletFolder)
	err = os.Chmod(walletFolder, 0777)
	if err != nil {
		return nil, "", "", err
	}

	req := testcontainers.ContainerRequest{
		Image:        "alephium/alephium:v1.0.0",
		ExposedPorts: []string{"12973/tcp"},
		WaitingFor:   wait.ForListeningPort("12973/tcp"),
		BindMounts: map[string]string{
			configFile: "/alephium-home/.alephium/user.conf",
			walletFolder: "/alephium-home/.alephium-wallets",
		},
	}
	alephiumNode, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, "", "", err
	}
	port, err := alephiumNode.MappedPort(ctx, "12973/tcp")
	if err != nil {
		return nil, "", "", err
	}

	return alephiumNode, port, walletFolder, nil
}
func tearDownContainer(ctx context.Context, c testcontainers.Container, dir string) {
	_ = os.RemoveAll(dir)
	_ = c.Terminate(ctx)
}

func TestCreateWalletE2E(t *testing.T) {

	log := logging.NewLogger()
	_ = logging.SetLogLevel(log, "debug")

	ctx, _ := context.WithCancel(context.Background())
	alephiumNode, port, walletFolder, err := setupContainer(ctx)
	defer tearDownContainer(ctx, alephiumNode, walletFolder)

	walletPassword := "dummy-password"
	walletName := "test-wallet"
	alephiumClient, err := New("http://localhost:"+port.Port(), log)
	assert.Nil(t, err)

	sync, err := alephiumClient.WaitUntilSyncedWithAtLeastOnePeer(context.Background())
	assert.Nil(t, err)
	assert.True(t, sync)


	genesisWalletName := "GenesisWallet-01"
	genesisWalletMnemonics := "convince crowd interest pen question tail curtain tenant buffalo advice mosquito position obey loyal gain local ecology tiger future turtle depend champion essence disorder"
	genesisWallet, err := alephiumClient.RestoreWallet(walletPassword, genesisWalletMnemonics, genesisWalletName, true, "")
	assert.Nil(t, err)

	genesisWalletAddresses, err := alephiumClient.GetWalletAddresses(genesisWallet.Name)
	assert.Nil(t, err)

	log.Printf("name: %s, activeAddress: %s, addresses: %s\n", genesisWallet.Name, genesisWalletAddresses.ActiveAddress, GetAddressesAsString(genesisWalletAddresses.Addresses))
	balance, err := alephiumClient.GetAddressBalance(genesisWalletAddresses.ActiveAddress, -1)
	assert.Equal(t, 1, balance.UtxoNum)
	log.Printf("Balance: %s\n", balance.Balance.PrettyString())

	newWallet, err := alephiumClient.CreateWallet(walletName, walletPassword, true, "")
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

	walletAddresses, err := alephiumClient.GetWalletAddresses(genesisWallet.Name)
	assert.Nil(t, err)

	log.Printf("name: %s, activeAddress: %s, addresses: %s\n", newWallet.Name, walletAddresses.ActiveAddress, GetAddressesAsString(walletAddresses.Addresses))

	walletBalances, err := alephiumClient.GetWalletBalances(newWallet.Name)
	assert.Nil(t, err)

	totalBalance, ok := walletBalances.GetTotalBalance()
	assert.True(t, ok)

	log.Printf("name: %s, total balance: %d\n", newWallet.Name, totalBalance)

	ok, err = alephiumClient.DeleteWallet(newWallet.Name, walletPassword)
	assert.True(t, ok)

	restoredWallet, err := alephiumClient.RestoreWallet(walletPassword, newWallet.Mnemonic, newWallet.Name, false, "")
	assert.Nil(t, err)

	walletExist, err := alephiumClient.CheckWalletExist(newWallet.Name)
	assert.Nil(t, err)
	assert.True(t, walletExist)

	walletAddresses, err = alephiumClient.GetWalletAddresses(restoredWallet.Name)
	assert.Nil(t, err)
	assert.Contains(t, GetAddressesAsString(walletAddresses.Addresses), walletAddresses.ActiveAddress)
	log.Printf("name: %s, activeAddress: %s, addresses: %s\n", restoredWallet.Name, walletAddresses.ActiveAddress, GetAddressesAsString(walletAddresses.Addresses))

	activeAddress, err := alephiumClient.ChangeActiveAddress(restoredWallet.Name, walletAddresses.ActiveAddress)
	assert.Nil(t, err)
	assert.True(t, activeAddress)
	log.Printf("name: %s, new activeAddress: %s, addresses: %s\n", restoredWallet.Name, walletAddresses.ActiveAddress, GetAddressesAsString(walletAddresses.Addresses))

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

	mnemonics, err := alephiumClient.RevealWalletMnemonic(restoredWallet.Name, walletPassword)
	assert.Nil(t, err)
	assert.Equal(t, newWallet.Mnemonic, mnemonics)

	encodedStr := hex.EncodeToString([]byte("some random data"))
	signature, err := alephiumClient.Sign(restoredWallet.Name, encodedStr)
	assert.Nil(t, err)
	log.Infof("%s", signature)
}

func TestJSONALF(t *testing.T) {
	amount, ok := ALPHFromALPHString("12.12")
	assert.True(t, ok)

	body := TransferRequest{
		Destinations: []TransferDestination{
			{Address: "1234", Amount: amount},
		},
	}

	b, err := json.Marshal(&body)
	assert.Nil(t, err)
	fmt.Printf("json: %s\n", string(b))
}
