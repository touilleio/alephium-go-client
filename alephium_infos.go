package alephium

import (
	"fmt"
	"time"
)

// GetSelfCliqueInfos
func (a *AlephiumClient) GetSelfCliqueInfos() (SelfCliqueInfo, error) {
	var selfCliqueInfos SelfCliqueInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("infos/self-clique").
		Receive(&selfCliqueInfos, &errorDetail)
	return selfCliqueInfos, relevantError(err, errorDetail)
}

// GetInterCliquePeerInfos
func (a *AlephiumClient) GetInterCliquePeerInfos() ([]InterCliquePeerInfo, error) {
	var interCliquePeerInfos []InterCliquePeerInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("infos/inter-clique-peer-info").
		Receive(&interCliquePeerInfos, &errorDetail)

	return interCliquePeerInfos, relevantError(err, errorDetail)
}

func IsSyncedWithAtLeastOnePeer(peers []InterCliquePeerInfo) bool {
	atLeastOneSynced := false
	for _, peer := range peers {
		if peer.IsSynced {
			atLeastOneSynced = true
		}
	}
	return atLeastOneSynced
}

func (a *AlephiumClient) WaitUntilSyncedWithAtLeastOnePeer() error {
	for isSynced := false; !isSynced; {
		var err error
		isSynced, err = a.IsSynced()
		if err != nil {
			return err
		}
		if !isSynced {
			a.log.Debugf("Not sync'ed yet, sleeping 10s")
			time.Sleep(10 * time.Second)
		}
	}
	return nil
}

func (a *AlephiumClient) IsSynced() (bool, error) {
	peers, err := a.GetInterCliquePeerInfos()
	if err != nil {
		return false, err
	}
	isSynced := IsSyncedWithAtLeastOnePeer(peers)
	return isSynced, nil
}

// GetDiscoveredNeighbors
func (a *AlephiumClient) GetDiscoveredNeighbors() error {
	return fmt.Errorf("not implemented yet")
}

// GetMisbehaviors
func (a *AlephiumClient) GetMisbehaviors() error {
	return fmt.Errorf("not implemented yet")
}
