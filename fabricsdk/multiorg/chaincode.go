package multiorg

import (
	"fabric-sdk-go-test/util"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	packager "github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"github.com/pkg/errors"
	"os"
)

func (mc *MultiorgContext) QueryCC(ccName string, fcn string, args []string) (channel.Response, error) {
	response, err := mc.ChClientOrg1User.Query(channel.Request{ChaincodeID: ccName, Fcn: fcn, Args: util.PackageArgs(args)},
		channel.WithRetry(retry.DefaultChannelOpts))
	return response, err
}

func (mc *MultiorgContext) ExecuteCC(ccName string, fcn string, args []string) (string, error) {
	response, err := mc.ChClientOrg1User.Execute(channel.Request{ChaincodeID: ccName, Fcn: fcn, Args: util.PackageArgs(args)}, channel.WithRetry(retry.DefaultChannelOpts))
	if err != nil {
		return "", err
	}
	if response.ChaincodeStatus == 0 {
		return "", errors.Wrap(err, "response chaincodestatus != 0")
	}
	if response.Responses[0].ChaincodeStatus != response.ChaincodeStatus {
		return "", err
	}
	return string(response.TransactionID), nil
}

func (mc *MultiorgContext) InstallCC(ccPath, ccName, ccVersion string) (bool, error) {
	ccPkg, err := packager.NewCCPackage(ccPath, os.Getenv("GOPATH"))
	if err != nil {
		return false, errors.Wrap(err, "package chaincode failed")
	}
	installCCReq := resmgmt.InstallCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Package: ccPkg}
	org1Peers, err := discoverLocalPeers(mc.Org1AdminClientContext, 1)
	if err != nil {
		return false, errors.Wrap(err, "get org1's peers failed")
	}
	org2Peers, err := discoverLocalPeers(mc.Org2AdminClientContext, 1)
	if err != nil {
		return false, errors.Wrap(err, "get org2's peers failed")
	}
	_, err = mc.Org1ResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return false, errors.Wrap(err, "install chaincode on org1's peer failed")
	}
	_, err = mc.Org2ResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return false, errors.Wrap(err, "install chaincode on org2's peer failed")
	}
	installed := queryInstalledCC("Org1", mc.Org1ResMgmt, ccName, ccVersion, org1Peers)
	if !installed {
		return false, errors.New("not installed on all peers in Org1")
	}
	fmt.Println("installed on all peers in org1")
	installed = queryInstalledCC("Org2", mc.Org2ResMgmt, ccName, ccVersion, org2Peers)
	if !installed {
		return false, errors.New("not installed on all peers in Org2")
	}
	fmt.Println("installed on all peers in org2")
	return true, nil
}

func (mc *MultiorgContext) InstantiateCC(ccName, ccVersion, ccPath string, endorsementPolicy string, args []string) (bool, error) {
	org1Peers, err := discoverLocalPeers(mc.Org1AdminClientContext, 1)
	if err != nil {
		return false, errors.Wrap(err, "get org1's peers failed")
	}
	org2Peers, err := discoverLocalPeers(mc.Org2AdminClientContext, 1)
	if err != nil {
		return false, errors.Wrap(err, "get org2's peers failed")
	}
	_, err = instantiateChaincode(mc.Org1ResMgmt, mc.ChannelID, ccName, ccPath, ccVersion, endorsementPolicy,
		util.PackageArgs(args), nil)

	if err != nil {
		return false, errors.Wrap(err, "instantiate chaincode "+ccName+" failed")
	}
	found := queryInstantiatedCC("Org1", mc.Org1ResMgmt, mc.ChannelID, ccName, ccVersion, org1Peers)
	if !found {
		fmt.Printf("Failed to find instantiated chaincode [%s:%s] in at least one peer in Org1 on channel [%s]\n", ccName, ccVersion, mc.ChannelID)
	}

	found = queryInstantiatedCC("Org2", mc.Org2ResMgmt, mc.ChannelID, ccName, ccVersion, org2Peers)
	if !found {
		fmt.Printf("Failed to find instantiated chaincode [%s:%s] in at least one peer in Org2 on channel [%s]\n", ccName, ccVersion, mc.ChannelID)
	}
	return true, nil
}

func (mc *MultiorgContext) upgradeCC(ccName, ccPath, ccVersion, ep string, args []string) (bool, error) {

	ccPkg, err := packager.NewCCPackage(ccPath, os.Getenv("GOPATH"))
	if err != nil {
		return false, errors.Wrap(err, "package chaincode failed")
	}
	installCCReq := resmgmt.InstallCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Package: ccPkg}

	org1Peers, err := discoverLocalPeers(mc.Org1AdminClientContext, 2)

	org2Peers, err := discoverLocalPeers(mc.Org2AdminClientContext, 2)

	_, err = mc.Org1ResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return false, errors.Wrap(err, "install chaincode on org1's peer failed")
	}
	_, err = mc.Org2ResMgmt.InstallCC(installCCReq, resmgmt.WithRetry(retry.DefaultResMgmtOpts))
	if err != nil {
		return false, errors.Wrap(err, "install chaincode on org2's peer failed")
	}

	installed := queryInstalledCC("Org1", mc.Org1ResMgmt, ccName, ccVersion, org1Peers)
	if !installed {
		return false, errors.New("not installed on all peers in Org1")
	}
	installed = queryInstalledCC("Org2", mc.Org2ResMgmt, ccName, ccVersion, org2Peers)
	if !installed {
		return false, errors.New("not installed on all peers in Org2")
	}
	org1Andorg2Policy, err := cauthdsl.FromString(ep)

	_, err = mc.Org1ResMgmt.UpgradeCC(mc.ChannelID, resmgmt.UpgradeCCRequest{Name: ccName, Path: ccPath, Version: ccVersion, Args: util.PackageArgs(args), Policy: org1Andorg2Policy})
	if err != nil {
		return false, errors.Wrap(err, "upgrade chaincode failed")
	}
	return true, nil
}


