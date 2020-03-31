package multiorg

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	contextAPI "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
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
)


type MultiorgContext struct {
	OrdererClientContext   contextAPI.ClientProvider
	Org1AdminClientContext contextAPI.ClientProvider
	Org2AdminClientContext contextAPI.ClientProvider
	ChClientOrg1User *channel.Client
	ChClientOrg2User *channel.Client
	Org1ResMgmt      *resmgmt.Client
	Org2ResMgmt      *resmgmt.Client
	Sdk              *fabsdk.FabricSDK
	ChannelID        string
}

func (mc *MultiorgContext) Init() {
	fabricConfigPath := "/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config-multiorg.yaml"
	sdk, err := fabsdk.New(config.FromFile(fabricConfigPath))
	if err != nil {
		panic(err)
	}
	mc.Sdk = sdk
	fmt.Println("sdk created successful")
	mc.OrdererClientContext = sdk.Context(fabsdk.WithUser(ordererAdminUser), fabsdk.WithOrg(ordererOrgName))
	mc.Org1AdminClientContext = sdk.Context(fabsdk.WithUser(org1AdminUser), fabsdk.WithOrg(org1))
	mc.Org2AdminClientContext = sdk.Context(fabsdk.WithUser(org2AdminUser), fabsdk.WithOrg(org2))
	org1ChannelClientContext := sdk.ChannelContext(mc.ChannelID, fabsdk.WithUser(org1User), fabsdk.WithOrg(org1))
	org2ChannelClientContext := sdk.ChannelContext(mc.ChannelID, fabsdk.WithUser(org2User), fabsdk.WithOrg(org2))
	err = createOrgsChannelClients(org1ChannelClientContext, org2ChannelClientContext, mc)
	if err != nil {
		panic(err)
	}
	org1ResMgmt, err := resmgmt.New(mc.Org1AdminClientContext)
	if err != nil {
		panic(err)
	}
	mc.Org1ResMgmt = org1ResMgmt
	org2ResMgmt, err := resmgmt.New(mc.Org2AdminClientContext)
	if err != nil {
		panic(err)
	}
	mc.Org2ResMgmt = org2ResMgmt
	fmt.Println("finish initializing the sdk")
}

func (mc *MultiorgContext) TearDown() {
	mc.Sdk.Close()
}

func createOrgsChannelClients(org1ChannelClientContext contextAPI.ChannelProvider, org2ChannelClientContext contextAPI.ChannelProvider,
	mc *MultiorgContext) error {
	// Org1 user connects to channel
	chClientOrg1User, err := channel.New(org1ChannelClientContext)
	if err != nil {
		return err
	}
	// Org2 user connects to channel
	chClientOrg2User, err := channel.New(org2ChannelClientContext)
	if err != nil {
		return err
	}
	mc.ChClientOrg1User = chClientOrg1User
	mc.ChClientOrg2User = chClientOrg2User
	return nil
}
