package alephium

import (
	"fmt"
	"math/big"
)

type WalletInfo struct {
	Wallet `json:",inline"`
	Locked bool `json:"locked"`
}

func (w WalletInfo) String() string {
	return fmt.Sprintf("%s, %v", w.Name, w.Locked)
}

type ErrorDetail struct {
	Resource string `json:"resource"`
	Detail   string `json:"detail"`
}

func (e ErrorDetail) Error() string {
	return e.Detail
}

type WalletBalances struct {
	TotalBalance string           `json:"totalBalance"`
	Balances     []AddressBalance `json:"balances"`
}

func (b WalletBalances) GetTotalBalance() (*big.Int, bool) {
	return new(big.Int).SetString(b.TotalBalance, 0)
}

type AddressBalance struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
}

func (b AddressBalance) GetBalance() (*big.Int, bool) {
	return new(big.Int).SetString(b.Balance, 0)
}

var WalletLockError = ErrorDetail{
	Detail: "WalletInfo is locked",
}

type WalletCreate struct {
	Wallet   `json:",inline"`
	Mnemonic string `json:"mnemonic"`
}

type WalletAddresses struct {
	ActiveAddress string   `json:"activeAddress"`
	Addresses     []string `json:"addresses"`
}

type Wallet struct {
	Name string `json:"walletName"`
}

type MinersAddresses struct {
	Addresses []string `json:"addresses"`
}

type InterCliquePeerInfo struct {
	CliqueId          string    `json:"cliqueId"`
	BrokerId          int       `json:"brokerId"`
	GroupNumPerBroker int       `json:"groupNumPerBroker"`
	Address           IPAndPort `json:"address"`
	IsSynced          bool      `json:"isSynced"`
}

type IPAndPort struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}

type SelfCliqueInfo struct {
	CliqueId              string        `json:"cliqueId"`
	NetworkType           string        `json:"networkType"`
	NumZerosAtLeastInHash int           `json:"numZerosAtLeastInHash"`
	Nodes                 []NodeAddress `json:"nodes"`
	Synced                bool          `json:"synced"`
	GroupNumPerBroker     int           `json:"groupNumPerBroker"`
	Groups                int           `json:"groups"`
}

type NodeAddress struct {
	Address  string `json:"address"`
	RestPort int    `json:"restPort"`
	WsPort   int    `json:"wsPort"`
}

type NodeInfo struct {
	IsMining bool `json:"isMining"`
}

type Transaction struct {
	TransactionId string `json:"txId"`
	FromGroup     int    `json:"fromGroup"`
	ToGroup       int    `json:"toGroup"`
}

type Address struct {
	Address string `json:"address"`
}

type TransactionStatus struct {
	Type                   string `json:"type"`
	BlockHash              string `json:"blockHash"`
	BlockIndex             int    `json:"blockIndex"`
	ChainConfirmations     int    `json:"chainConfirmations"`
	FromGroupConfirmations int    `json:"fromGroupConfirmations"`
	ToGroupConfirmations   int    `json:"toGroupConfirmations"`
}

type AddressUtxoBalance struct {
	Balance string `json:"balance"`
	LockedBalance string `json:"lockedBalance"`
	UtxoNum int `json:"utxoNum"`
}

type AddressGroup struct {
	Group int `json:"group"`
}

type DiscoveredNeighbor struct {
	CliqueId string `json:"cliqueId"`
	BrokerId int `json:"brokerId"`
	GroupNumPerBroker int `json:"groupNumPerBroker"`
}

type Misbehavior struct {
	Peer string `json:"peer"`
	Status MisbehaviorStatus `json:"status"`
}

type MisbehaviorStatus struct {
	Type string `json:"type"`
	Value int `json:"value"`
}