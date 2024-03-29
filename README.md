Alephium API Client for Golang
====

> **Warning**
> This project is no longer maintained, please use the official go-sdk from Alephium: [https://github.com/alephium/go-sdk](https://github.com/alephium/go-sdk)

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

// Wait until the node is sync'ed with bootstrap nodes
alephiumClient.WaitUntilSyncedWithAtLeastOnePeer()
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
docker run -it --rm -v ${PWD}/user-dev-standalone.conf:/alephium-home/.alephium/user.conf -p 12973:12973 alephium/alephium:v1.1.13
```
