package fabricsdk

import (
	"fabric-sdk-go-test/util"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
)

func (fm *FabricModel) QueryCC(chaincodeID string, fcn string, args []string) (channel.Response, error) {
	queryRequest := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         fcn,
		Args:        util.PackageArgs(args),
	}

	response, err := fm.ChannelClient.Query(queryRequest)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (fm *FabricModel) InvokeCC(chaincodeID string, fcn string, args []string) (channel.Response, error) {
	invokeRequest := channel.Request{
		ChaincodeID: chaincodeID,
		Fcn:         fcn,
		Args:        util.PackageArgs(args),
	}

	response, err := fm.ChannelClient.Execute(invokeRequest)
	if err != nil {
		return response, err
	}
	return response, nil
}

func (fm *FabricModel) ConstructChannel() (fab.TransactionID, error) {

	mspClient, err := mspclient.New(fm.Sdk.Context(), mspclient.WithOrg(fm.OrgName))
	if err != nil {
		return "", err
	}
	adminIdentity, err := mspClient.GetSigningIdentity(fm.OrgAdmin)
	if err != nil {
		return "", err
	}

	channelRequest := resmgmt.SaveChannelRequest{
		ChannelID:         fm.ChannelID,
		ChannelConfigPath: fm.ChannelConfigPath,
		SigningIdentities: []msp.SigningIdentity{adminIdentity},
	}

	response, err := fm.ResMgmtClient.SaveChannel(channelRequest)
	if err != nil {
		return "", err
	}
	err = fm.ResMgmtClient.JoinChannel(
		fm.ChannelID,
		resmgmt.WithRetry(retry.DefaultResMgmtOpts),
		resmgmt.WithOrdererEndpoint(fm.OrdererName),
	)
	return response.TransactionID, err
}

func (fm *FabricModel) InstallCC() (bool, error) {
	ccp, err := gopackager.NewCCPackage(fm.ChaincodePath, fm.ChaincodeGoPath)
	if err != nil {
		return false, err
	}
	installCCRequest := resmgmt.InstallCCRequest{
		Name:    fm.ChainCodeID,   // 链码名称
		Path:    fm.ChaincodePath, //链码在工程中的路径
		Version: "1.0",
		Package: ccp,
	}

	response, err := fm.ResMgmtClient.InstallCC(installCCRequest)
	if response[0].Status == 200 && err == nil {
		return true, nil
	}
	return false, err
}
