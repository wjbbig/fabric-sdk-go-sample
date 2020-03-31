package multiorg

import (
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
)

func queryInstalledCC(orgID string, resMgmt *resmgmt.Client, ccName, ccVersion string, peers []fab.Peer) bool {
	installed, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			ok := isCCInstalled(orgID, resMgmt, ccName, ccVersion, peers)
			if !ok {
				return &ok, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Chaincode [%s:%s] is not installed on all peers in Org1\n", ccName, ccVersion), nil)
			}
			return &ok, nil
		},
	)
	if err != nil {
		panic(err)
	}
	return *(installed).(*bool)
}
func isCCInstalled(orgID string, resMgmt *resmgmt.Client, ccName, ccVersion string, peers []fab.Peer) bool {
	fmt.Printf("Querying [%s] peers to see if chaincode [%s:%s] was installed\n", orgID, ccName, ccVersion)
	installedOnAllPeers := true
	for _, peer := range peers {
		fmt.Printf("Querying [%s] ...\n", peer.URL())
		resp, err := resMgmt.QueryInstalledChaincodes(resmgmt.WithTargets(peer))
		if err != nil {
			fmt.Printf("QueryInstalledChaincodes for peer [%s] failed\n", peer.URL())
		}

		found := false
		for _, ccInfo := range resp.Chaincodes {
			fmt.Printf("... found chaincode [%s:%s]\n", ccInfo.Name, ccInfo.Version)
			if ccInfo.Name == ccName && ccInfo.Version == ccVersion {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("... chaincode [%s:%s] is not installed on peer [%s]\n", ccName, ccVersion, peer.URL())
			installedOnAllPeers = false
		}
	}
	return installedOnAllPeers
}

func discoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int) ([]fab.Peer, error) {
	ctx, err := contextImpl.NewLocal(ctxProvider)
	if err != nil {
		return nil, errors.Wrap(err, "error creating local context")
	}
	discoveredPeers, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			peers, serviceErr := ctx.LocalDiscoveryService().GetPeers()
			if serviceErr != nil {
				return nil, errors.Wrapf(serviceErr, "error getting peers for MSP [%s]", ctx.Identifier().MSPID)
			}
			if len(peers) < expectedPeers {
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d\n", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}
	return discoveredPeers.([]fab.Peer), nil
}

func instantiateChaincode(resMgmt *resmgmt.Client, channelID, ccName, ccPath, ccVersion string, ccPolicyStr string, args [][]byte, collConfigs []*common.CollectionConfig) (resmgmt.InstantiateCCResponse, error) {
	ccPolicy, err := cauthdsl.FromString(ccPolicyStr)
	if err != nil {
		return resmgmt.InstantiateCCResponse{}, errors.Wrapf(err, "error creating CC policy [%s]", ccPolicyStr)
	}

	return resMgmt.InstantiateCC(
		channelID,
		resmgmt.InstantiateCCRequest{
			Name:       ccName,
			Path:       ccPath,
			Version:    ccVersion,
			Args:       args,
			Policy:     ccPolicy,
			CollConfig: collConfigs,
		},
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
	)
}


func queryInstantiatedCC(orgID string, resMgmt *resmgmt.Client, channelID, ccName, ccVersion string, peers []fab.Peer) bool {
	if len(peers) == 0 {
		fmt.Printf( "Expecting one or more peers\n")
		return false
	}
	fmt.Printf("Querying [%s] peers to see if chaincode [%s] was instantiated on channel [%s]\n", orgID, ccName, channelID)

	instantiated, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			ok := isCCInstantiated(resMgmt, channelID, ccName, ccVersion, peers)
			if !ok {
				return &ok, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Did NOT find instantiated chaincode [%s:%s] on one or more peers in [%s].", ccName, ccVersion, orgID), nil)
			}
			return &ok, nil
		},
	)
	if err != nil {
		panic(err)
	}
	return *(instantiated).(*bool)
}

func isCCInstantiated(resMgmt *resmgmt.Client, channelID, ccName, ccVersion string, peers []fab.Peer) bool {
	installedOnAllPeers := true
	for _, peer := range peers {
		fmt.Printf("Querying peer [%s] for instantiated chaincode [%s:%s]...\n", peer.URL(), ccName, ccVersion)
		chaincodeQueryResponse, err := resMgmt.QueryInstantiatedChaincodes(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargets(peer))
		if err != nil {
			fmt.Println("QueryInstantiatedChaincodes return error")
		}

		fmt.Printf("Found %d instantiated chaincodes on peer [%s]:\n", len(chaincodeQueryResponse.Chaincodes), peer.URL())
		found := false
		for _, chaincode := range chaincodeQueryResponse.Chaincodes {
			fmt.Printf("Found instantiated chaincode Name: [%s], Version: [%s], Path: [%s] on peer [%s]\n", chaincode.Name, chaincode.Version, chaincode.Path, peer.URL())
			if chaincode.Name == ccName && chaincode.Version == ccVersion {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("... chaincode [%s:%s] is not instantiated on peer [%s]\n", ccName, ccVersion, peer.URL())
			installedOnAllPeers = false
		}
	}
	return installedOnAllPeers
}
