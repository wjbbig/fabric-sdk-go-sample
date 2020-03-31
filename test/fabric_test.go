package main

import (
	"encoding/hex"
	"fabric-sdk-go-test/fabricsdk"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"os"
	"testing"
)

func TestLedger(t *testing.T) {
	fm := fabricsdk.FabricModel{
		ConfigFile: "/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml",
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		UserName:   "User1",
		ChannelID:  "mychannel",
		HasInit:    false,
	}

	fm.Init()
	//cc, _ := fm.QueryCC("mycc", "query", []string{"b"})
	//fmt.Println(string(cc.Payload))
	//fmt.Println(cc.TransactionID)
	//invokeCC, _ := fm.InvokeCC("mycc", "invoke", []string{"a", "b", "10"})
	//fmt.Println(invokeCC.TransactionID)
	//queryCC, _ := fm.QueryCC("mycc", "query", []string{"b"})
	//fmt.Println(string(queryCC.Payload))
	id, _ := fm.QueryBlockByTxID("c0a3db42ff2d040e1502c62259738a5bf6bf3a576168ad0d32213a195f059d6c")
	fmt.Println(hex.EncodeToString(id.Header.DataHash))
	number, _ := fm.QueryBlockByNumber(6)
	fmt.Println(hex.EncodeToString(number.Header.DataHash))
	info, _ := fm.QueryInfo()
	fmt.Println(info.BCI.Height)
	fmt.Println(info.Status)
	fmt.Println(info.Endorser)
	//hash, _ := fm.QueryBlockNumberByHash("9a3a29dd7be7110310db199643b46b7e186a88854516747692e4e768892c9ce3")
	//fmt.Println(hash.Header.Number)
	config, _ := fm.QueryConfig()
	//channel版本
	fmt.Println(config.Versions().Channel.Version)
}

func TestChannel(t *testing.T) {
	channelConfigPath := "/home/fujitsu/workspace/go/src/github.com/hyperledger/fabric/scripts/fabric-samples/first-network/channel-artifacts/"
	channelID := "mytest"
	fm := fabricsdk.FabricModel{
		ConfigFile:        "/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml",
		OrgAdmin:          "Admin",
		OrgName:           "Org1",
		OrgID:             "Org1MSP",
		UserName:          "User1",
		ChannelID:         channelID,
		ChannelConfigPath: channelConfigPath + channelID + ".tx",
		OrdererName:       "orderer.example.com",
		ChainCodeID:       "mytest",
		ChaincodeGoPath:   os.Getenv("GOPATH"),
		ChaincodePath:     "mychaincode/go",
	}
	fm.Init()
	channel, err := fm.ConstructChannel()
	//cc, err := fm.InstallCC()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(channel)
}

func TestRegister(t *testing.T) {
	sdk, _ := fabsdk.New(config.FromFile("/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml"))
	clientContext := sdk.Context()
	mspClient, _ := msp.New(clientContext)
	instance, s := getRegistrarEnrollmentCredentialsWithCAInstance(t, clientContext, "")
	fmt.Println(instance)
	fmt.Println(s)
	err := mspClient.Enroll(instance, msp.WithSecret(s))
	if err != nil {
		t.Fatalf("Enroll failed: %s", err)
	}

	var attributes []msp.Attribute
	attributes = append(attributes, msp.Attribute{Name: "test1", Value: "test2", ECert: true})
	attributes = append(attributes, msp.Attribute{Name: "test2", Value: "test3", ECert: true})
	register := msp.RegistrationRequest{
		Name:        "wjb",
		Type:        "user",
		Affiliation: "org1",
		Attributes:  attributes,
		//Secret:      "wjb12345",
	}
	secret, err := mspClient.Register(&register)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(secret)
}
func getRegistrarEnrollmentCredentialsWithCAInstance(t *testing.T, ctxProvider context.ClientProvider, caID string) (string, string) {

	ctx, err := ctxProvider()
	if err != nil {
		t.Fatalf("failed to get context: %s", err)
	}

	myOrg := ctx.IdentityConfig().Client().Organization

	if caID == "" {
		caID = ctx.EndpointConfig().NetworkConfig().Organizations[myOrg].CertificateAuthorities[0]
	}

	caConfig, ok := ctx.IdentityConfig().CAConfig(caID)
	if !ok {
		t.Fatal("CAConfig failed")
	}

	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
}

//func TestEnrollUser(t *testing.T) {
//	sdk, _ := fabsdk.New(config.FromFile("/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml"))
//	clientContext := sdk.Context()
//	mspClient, _ := msp.New(clientContext)
//	err := mspClient.Enroll("wjb", msp.WithSecret("NlgXFsGJaHEb"))
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	aa, err := mspClient.GetSigningIdentity("wjb")
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(aa)
//	identity.
//	channelContext := sdk.ChannelContext("mychannel", fabsdk.WithUser("wjb"))
//	channelClient, err := channel.New(channelContext)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("通道创建完成")
//	ledger, err := ledger.New(channelContext)
//}
