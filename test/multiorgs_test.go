package main

import (
	"fmt"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/status"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	fabAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	contextImpl "github.com/hyperledger/fabric-sdk-go/pkg/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/resource"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/test/metadata"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	org1             = "Org1"
	org2             = "Org2"
	ordererAdminUser = "Admin"
	ordererOrgName   = "OrdererOrg"
	org1AdminUser    = "Admin"
	org2AdminUser    = "Admin"
	org1User         = "User1"
	org2User         = "User1"
	channelID        = "mychannel"
	ccPath           = "github.com/mychaincode/go"
)

var (
	// SDK
	sdk *fabsdk.FabricSDK

	// Org MSP clients
	org1MspClient *mspclient.Client
	org2MspClient *mspclient.Client
	// Peers
	orgTestPeer0 fab.Peer
	orgTestPeer1 fab.Peer
	exampleCC    = "example_cc_e2e" + metadata.TestRunID
)

type multiorgContext struct {
	// client contexts
	ordererClientContext   contextAPI.ClientProvider
	org1AdminClientContext contextAPI.ClientProvider
	org2AdminClientContext contextAPI.ClientProvider
	org1ResMgmt            *resmgmt.Client
	org2ResMgmt            *resmgmt.Client
	ccName                 string
	ccVersion              string
	channelID              string
}

func Setup() error {
	var err error
	sdk, err = fabsdk.New(config.FromFile("/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config-multiorg.yaml"))
	if err != nil {
		return err
	}
	org1MspClient, err = mspclient.New(sdk.Context(), mspclient.WithOrg(org1))
	if err != nil {
		return err
	}

	org2MspClient, err = mspclient.New(sdk.Context(), mspclient.WithOrg(org2))
	if err != nil {
		return err
	}

	return nil
}

func teardown() {
	if sdk != nil {
		sdk.Close()
	}
}

func TestSetup(t *testing.T) {
	err := Setup()
	if err != nil {
		t.Fatal(err)
	}
	mc := multiorgContext{
		ordererClientContext:   sdk.Context(fabsdk.WithUser(ordererAdminUser), fabsdk.WithOrg(ordererOrgName)),
		org1AdminClientContext: sdk.Context(fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1)),
		org2AdminClientContext: sdk.Context(fabsdk.WithUser(org2AdminUser), fabsdk.WithOrg(org2)),
		ccName:                 "mycc",
		channelID:              "mychannel",
	}
	setupClientContextsAndChannel(t, sdk, &mc)
	testWithOrg1(t, sdk, &mc)
	teardown()
}

func testWithOrg1(t *testing.T, sdk *fabsdk.FabricSDK, mc *multiorgContext) {

	//org1AdminChannelContext := sdk.ChannelContext(mc.channelID, fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))
	org1ChannelClientContext := sdk.ChannelContext(mc.channelID, fabsdk.WithUser(org1User), fabsdk.WithOrg(org1))
	org2ChannelClientContext := sdk.ChannelContext(mc.channelID, fabsdk.WithUser(org2User), fabsdk.WithOrg(org2))

	/*chClientOrg1User*/ _, chClientOrg2User := createOrgsChannelClients(org1ChannelClientContext, t, org2ChannelClientContext)

	ccPkg, err := packager.NewCCPackage(ccPath, os.Getenv("GOPATH"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ccPkg)
	createCC(t, mc, ccPkg, "testcccc", "1.0")
	cc := queryCC(chClientOrg2User, t, "testcccc")
	t.Log(string(cc))
	//transactionID := executeCC(chClientOrg2User, t, mc.ccName)
	//t.Log(transactionID)
}

func setupClientContextsAndChannel(t *testing.T, sdk *fabsdk.FabricSDK, mc *multiorgContext) {
	// Org1 resource management client (Org1 is default org)
	org1RMgmt, err := resmgmt.New(mc.org1AdminClientContext)
	require.NoError(t, err, "failed to create org1 resource management client")

	mc.org1ResMgmt = org1RMgmt

	// Org2 resource management client
	org2RMgmt, err := resmgmt.New(mc.org2AdminClientContext)
	require.NoError(t, err, "failed to create org2 resource management client")

	mc.org2ResMgmt = org2RMgmt
}


func createOrgsChannelClients(org1ChannelClientContext contextAPI.ChannelProvider, t *testing.T, org2ChannelClientContext contextAPI.ChannelProvider) (*channel.Client, *channel.Client) {
	// Org1 user connects to 'mychannel'
	chClientOrg1User, err := channel.New(org1ChannelClientContext)
	if err != nil {
		t.Fatalf("Failed to create new channel client for Org1 user: %s", err)
	}
	// Org2 user connects to 'mychannel'
	chClientOrg2User, err := channel.New(org2ChannelClientContext)
	if err != nil {
		t.Fatalf("Failed to create new channel client for Org2 user: %s", err)
	}
	return chClientOrg1User, chClientOrg2User
}

func queryCC(chClientOrg1User *channel.Client, t *testing.T, ccName string) []byte {
	response, err := chClientOrg1User.Query(channel.Request{ChaincodeID: ccName, Fcn: "query", Args: [][]byte{[]byte("b")}},
		channel.WithRetry(retry.DefaultChannelOpts))

	if err != nil {
		t.Fatal(err)
	}
	return response.Payload
}

func executeCC(chClientOrgUser *channel.Client, t *testing.T, ccName string) fab.TransactionID {
	response, err := chClientOrgUser.Execute(channel.Request{ChaincodeID: ccName, Fcn: "invoke", Args: [][]byte{[]byte("a"), []byte("b"), []byte("10")}}, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		t.Fatalf("Failed to move funds: %s", err)
	}
	if response.ChaincodeStatus == 0 {
		t.Fatal("Expected ChaincodeStatus")
	}
	if response.Responses[0].ChaincodeStatus != response.ChaincodeStatus {
		t.Fatal("Expected the chaincode status returned by successful Peer Endorsement to be same as Chaincode status for client response")
	}
	return response.TransactionID
}

func createCC(t *testing.T, mc *multiorgContext, ccPkg *resource.CCPackage, ccName, ccVersion string) {
	installCCReq := resmgmt.InstallCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Package: ccPkg}

	// Ensure that Gossip has propagated it's view of local peers before invoking
	// install since some peers may be missed if we call InstallCC too early
	org1Peers, err := DiscoverLocalPeers(mc.org1AdminClientContext, 1)
	if err != nil {
		t.Fatal(err)
	}
	org2Peers, err := DiscoverLocalPeers(mc.org2AdminClientContext, 1)
	if err != nil {
		t.Fatal(err)
	}

	// Install example cc to Org1 peers
	_, err = mc.org1ResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		t.Fatal(err)
	}

	// Install example cc to Org2 peers
	_, err = mc.org2ResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		t.Fatal(err)
	}

	//Ensure the CC is installed on all peers in both orgs
	installed := queryInstalledCC(t, "Org1", mc.org1ResMgmt, ccName, ccVersion, org1Peers)
	if !installed {
		t.Fatal("not installed on all peers in Org1")
	}
	t.Log("installed on all peers in org1")
	installed = queryInstalledCC(t, "Org2", mc.org2ResMgmt, ccName, ccVersion, org2Peers)
	if !installed {
		t.Fatal("not installed on all peers in Org2")
	}
	t.Log("installed on all peers in org2")
	instantiateCC(t, mc.org1ResMgmt, ccName, ccVersion, mc.channelID)

	// Ensure the CC is instantiated on all peers in both orgs
	found := queryInstantiatedCC(t, "Org1", mc.org1ResMgmt, mc.channelID, ccName, ccVersion, org1Peers)
	require.True(t, found, "Failed to find instantiated chaincode [%s:%s] in at least one peer in Org1 on channel [%s]", ccName, ccVersion, mc.channelID)

	found = queryInstantiatedCC(t, "Org2", mc.org2ResMgmt, mc.channelID, ccName, ccVersion, org2Peers)
	require.True(t, found, "Failed to find instantiated chaincode [%s:%s] in at least one peer in Org2 on channel [%s]", ccName, ccVersion, mc.channelID)
}

func DiscoverLocalPeers(ctxProvider contextAPI.ClientProvider, expectedPeers int) ([]fabAPI.Peer, error) {
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
				return nil, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Expecting %d peers but got %d", expectedPeers, len(peers)), nil)
			}
			return peers, nil
		},
	)
	if err != nil {
		return nil, err
	}

	return discoveredPeers.([]fabAPI.Peer), nil
}

func queryInstalledCC(t *testing.T, orgID string, resMgmt *resmgmt.Client, ccName, ccVersion string, peers []fab.Peer) bool {
	installed, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			ok := isCCInstalled(t, orgID, resMgmt, ccName, ccVersion, peers)
			if !ok {
				return &ok, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Chaincode [%s:%s] is not installed on all peers in Org1", ccName, ccVersion), nil)
			}
			return &ok, nil
		},
	)

	require.NoErrorf(t, err, "Got error checking if chaincode was installed")
	return *(installed).(*bool)
}
func isCCInstalled(t *testing.T, orgID string, resMgmt *resmgmt.Client, ccName, ccVersion string, peers []fab.Peer) bool {
	t.Logf("Querying [%s] peers to see if chaincode [%s:%s] was installed", orgID, ccName, ccVersion)
	installedOnAllPeers := true
	for _, peer := range peers {
		t.Logf("Querying [%s] ...", peer.URL())
		resp, err := resMgmt.QueryInstalledChaincodes(resmgmt.WithTargets(peer))
		require.NoErrorf(t, err, "QueryInstalledChaincodes for peer [%s] failed", peer.URL())

		found := false
		for _, ccInfo := range resp.Chaincodes {
			t.Logf("... found chaincode [%s:%s]", ccInfo.Name, ccInfo.Version)
			if ccInfo.Name == ccName && ccInfo.Version == ccVersion {
				found = true
				break
			}
		}
		if !found {
			t.Logf("... chaincode [%s:%s] is not installed on peer [%s]", ccName, ccVersion, peer.URL())
			installedOnAllPeers = false
		}
	}
	return installedOnAllPeers
}

func instantiateCC(t *testing.T, resMgmt *resmgmt.Client, ccName, ccVersion string, channelID string) {
	instantiateResp, err := InstantiateChaincode(resMgmt, channelID, ccName, ccPath, ccVersion, "AND ('Org1MSP.peer','Org2MSP.peer')", [][]byte{[]byte("init"), []byte("a"), []byte("100"), []byte("b"), []byte("200")}, nil)
	require.NoError(t, err)
	require.NotEmpty(t, instantiateResp, "transaction response should be populated for instantateCC")
}

func InstantiateChaincode(resMgmt *resmgmt.Client, channelID, ccName, ccPath, ccVersion string, ccPolicyStr string, args [][]byte, collConfigs []*common.CollectionConfig) (resmgmt.InstantiateCCResponse, error) {
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

func queryInstantiatedCC(t *testing.T, orgID string, resMgmt *resmgmt.Client, channelID, ccName, ccVersion string, peers []fab.Peer) bool {
	require.Truef(t, len(peers) > 0, "Expecting one or more peers")
	t.Logf("Querying [%s] peers to see if chaincode [%s] was instantiated on channel [%s]", orgID, ccName, channelID)

	instantiated, err := retry.NewInvoker(retry.New(retry.TestRetryOpts)).Invoke(
		func() (interface{}, error) {
			ok := isCCInstantiated(t, resMgmt, channelID, ccName, ccVersion, peers)
			if !ok {
				return &ok, status.New(status.TestStatus, status.GenericTransient.ToInt32(), fmt.Sprintf("Did NOT find instantiated chaincode [%s:%s] on one or more peers in [%s].", ccName, ccVersion, orgID), nil)
			}
			return &ok, nil
		},
	)
	require.NoErrorf(t, err, "Got error checking if chaincode was instantiated")
	return *(instantiated).(*bool)
}

func isCCInstantiated(t *testing.T, resMgmt *resmgmt.Client, channelID, ccName, ccVersion string, peers []fab.Peer) bool {
	installedOnAllPeers := true
	for _, peer := range peers {
		t.Logf("Querying peer [%s] for instantiated chaincode [%s:%s]...", peer.URL(), ccName, ccVersion)
		chaincodeQueryResponse, err := resMgmt.QueryInstantiatedChaincodes(channelID, resmgmt.WithRetry(retry.DefaultResMgmtOpts), resmgmt.WithTargets(peer))
		require.NoError(t, err, "QueryInstantiatedChaincodes return error")

		t.Logf("Found %d instantiated chaincodes on peer [%s]:", len(chaincodeQueryResponse.Chaincodes), peer.URL())
		found := false
		for _, chaincode := range chaincodeQueryResponse.Chaincodes {
			t.Logf("Found instantiated chaincode Name: [%s], Version: [%s], Path: [%s] on peer [%s]", chaincode.Name, chaincode.Version, chaincode.Path, peer.URL())
			if chaincode.Name == ccName && chaincode.Version == ccVersion {
				found = true
				break
			}
		}
		if !found {
			t.Logf("... chaincode [%s:%s] is not instantiated on peer [%s]", ccName, ccVersion, peer.URL())
			installedOnAllPeers = false
		}
	}
	return installedOnAllPeers
}
