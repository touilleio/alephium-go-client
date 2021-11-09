package alephium

import (
	"context"
	"time"
)

// GetSelfCliqueInfos gets the infos about the current clique
func (a *Client) GetSelfCliqueInfos() (SelfCliqueInfo, error) {
	var selfCliqueInfos SelfCliqueInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("infos/self-clique").
		Receive(&selfCliqueInfos, &errorDetail)
	return selfCliqueInfos, relevantError(err, errorDetail)
}

// GetInterCliquePeerInfos gets cliques about the other cliques connected to the current cllique
func (a *Client) GetInterCliquePeerInfos() ([]InterCliquePeerInfo, error) {
	var interCliquePeerInfos []InterCliquePeerInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("infos/inter-clique-peer-info").
		Receive(&interCliquePeerInfos, &errorDetail)

	return interCliquePeerInfos, relevantError(err, errorDetail)
}

// IsSyncedWithAtLeastOnePeer checks if the clique is connected with at least one clique
// or the list of peers is empty
func IsSyncedWithAtLeastOnePeer(peers []InterCliquePeerInfo) bool {
	atLeastOneSynced := false
	for _, peer := range peers {
		if peer.IsSynced {
			atLeastOneSynced = true
		}
	}
	return atLeastOneSynced || len(peers) == 0
}

// WaitUntilSyncedWithAtLeastOnePeer waits until the clique is connected to at least one clique
// or the context is done.
func (a *Client) WaitUntilSyncedWithAtLeastOnePeer(ctx context.Context) (bool, error) {
	for isSynced := false; ; {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:

		}
		var err error
		isSynced, err = a.IsSynced()
		if err != nil {
			return false, err
		}
		if isSynced {
			return true, nil
		} else {
			a.log.Debugf("Not sync'ed yet, sleeping %s", a.sleepTime)
			time.Sleep(a.sleepTime)
		}
	}
}

// IsSynced checks if the cilque is synced
func (a *Client) IsSynced() (bool, error) {
	peers, err := a.GetInterCliquePeerInfos()
	if err != nil {
		return false, err
	}
	isSynced := IsSyncedWithAtLeastOnePeer(peers)
	return isSynced, nil
}

// GetDiscoveredNeighbors gets the discovered neighbors
func (a *Client) GetDiscoveredNeighbors() ([]DiscoveredNeighbor, error) {
	var neighbors []DiscoveredNeighbor
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("infos/discovered-neighbors").
		Receive(&neighbors, &errorDetail)
	return neighbors, relevantError(err, errorDetail)
}

// GetMisbehaviors gets the misbehaving neighbors
func (a *Client) GetMisbehaviors() ([]Misbehavior, error) {
	var misbehaviors []Misbehavior
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("infos/misbehaviors").
		Receive(&misbehaviors, &errorDetail)
	return misbehaviors, relevantError(err, errorDetail)
}

type UnbanMisbehaviorsBodyParams struct {
	Type  string   `json:"type"`
	Peers []string `json:"peers"`
}

// UnbanMisbehaviors unbans misbehaving neighbors
func (a *Client) UnbanMisbehaviors(peers []string) (bool, error) {
	var errorDetail ErrorDetail
	params := UnbanMisbehaviorsBodyParams{
		Type:  "uban",
		Peers: peers,
	}
	_, err := a.slingClient.New().Post("infos/misbehaviors").
		BodyJSON(params).Receive(nil, &errorDetail)
	return true, relevantError(err, errorDetail)
}

// GetNodeInfos get the info of the node
func (a *Client) GetNodeInfos() (NodeInfo, error) {
	var nodeInfo NodeInfo
	var errorDetail ErrorDetail
	_, err := a.slingClient.New().Path("infos/node").
		Receive(&nodeInfo, &errorDetail)
	return nodeInfo, relevantError(err, errorDetail)
}
