Alephium API Client for Golang
====

This project is a simple, yet expressive, API Client
for [Alephium](https://alephium.org/) blockchain,
written in [Golang](https://golang.org/).

# Getting started

The client wraps the API call in regular functions, ready to be used:

```
import (
	"github.com/sirupsen/logrus"
	"github.com/touilleio/alephium-go-client"
)

alephiumClient, err := alephium.New("http://localhost:12973", logrus.StandardLogger())

// Create a miner wallet
minerWallet, err := alephiumClient.CreateWallet("", "walletPassword", true, "")

// Get the addresses of the freshly created miner wallet
walletAddresses, err := alephiumClient.GetWalletAddresses(minerWallet.Name)

// Update miner wallet on the node (hint: you can set this in the config,
// see https://github.com/alephium/alephium/wiki/Miner-Guide for more details)
err = alephiumClient.UpdateMinersAddresses(walletAddresses.Addresses)

// Wait until the node is sync'ed with bootstrap nodes
alephiumClient.WaitUntilSyncedWithAtLeastOnePeer()

// Start mining
alephiumClient.StartMining()
```

# Hack

Build:

```
go build .
```

Test:

```
go test .
```

If you want to run your node manually,

```
docker run -it --rm -v ${PWD}/user-dev-standalone.conf:/alephium-home/.alephium/user.conf -p 12973:12973 alephium/alephium:v0.11.0
```
