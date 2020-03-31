package main

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	_ "github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
)

const (
	orgAdmin  = "Admin"
	orgName   = "Org1"
	channelID = "mychannel"
	username  = "User1"
)

func main() {
	//sdk, err := fabsdk.New(config.FromFile("/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml"))
	//if err != nil {
	//	panic(err)
	//}
	//
	//defer sdk.Close()
	//clientContext := sdk.Context(fabsdk.WithUser(orgAdmin), fabsdk.WithOrg(orgName))
	//if clientContext == nil {
	//	panic("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	//}
	////resMgmtClient, err := resmgmt.New(clientContext)
	////if err != nil {
	////	panic(err)
	////}
	////mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	////if err != nil {
	////	panic(err)
	////}
	////adminIdentity, err := mspClient.GetSigningIdentity(orgAdmin)
	////if err != nil {
	////	panic(err)
	////}
	//
	//channelContext := sdk.ChannelContext(channelID, fabsdk.WithUser(username))
	//channelClient, err := channel.New(channelContext)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Channel client created")
	//queryRequest := channel.Request{
	//	ChaincodeID:     "mycc",
	//	Fcn:             "query",
	//	Args:            packArgs("a"),
	//}
	//response, err := channelClient.Query(queryRequest)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(string(response.Payload))


	//ledger, err := ledger.New(channelContext)
	//if err != nil {
	//	panic(err)
	//}
	//info, err := ledger.QueryBlockByTxID("ef7946a5ae5f180b3f45d1a81eaf2fab739517e6ef10dd4c0b5333c7b6d41bbd")
	//fmt.Println(info.Header.Number)

	//invokeRequest := channel.Request{
	//	ChaincodeID: "mycc",
	//	Fcn:         "invoke",
	//	Args:        packArgs("b", "a", "20"),
	//}
	//execute, err := channelClient.Execute(invokeRequest)
	////channel.
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(execute.TransactionID)



}

func packArgs(args ...string) [][]byte {
	var argSlice [][]byte
	for i := 0; i < len(args); i++ {
		argSlice = append(argSlice, []byte(args[i]))
	}

	return argSlice
}


func QueryInfo(channelContext context.ChannelProvider) {
	ledger, err := ledger.New(channelContext)
	if err != nil {
		panic(err)
	}
	id, err := ledger.QueryBlockByTxID("6")
	fmt.Println(id.Data.Data)
}