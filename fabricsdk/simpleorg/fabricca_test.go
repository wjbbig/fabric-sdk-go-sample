package simpleorg

import (
	"fabric-sdk-go-test/util"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"testing"
)

const (
	configPath = "/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml"
	username = "aa"
	password = "abcaa"
	affiliation = "org1.department1"
)

func TestRegisterUser(t *testing.T) {

	secret, err := RegisterUser(configPath, username, password, affiliation, nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(secret)
}

func TestEnrollUser(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		t.Fatal(err.Error())
	}
	clientContext := sdk.Context()
	mspClient, err := msp.New(clientContext, msp.WithOrg("Org1"))
	if err != nil {
		t.Fatal(err.Error())
	}
	err = mspClient.Enroll(username, msp.WithSecret(password))
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(mspClient.GetSigningIdentity(username))
}

func TestUserNewCaUserQueryLedger(t *testing.T) {
	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		t.Fatal(err.Error())
	}
	clientContext := sdk.Context()
	mspClient, err := msp.New(clientContext, msp.WithOrg("Org1"))
	if err != nil {
		t.Fatal(err.Error())
	}
	err = mspClient.Enroll(username, msp.WithSecret(password))
	if err != nil {
		t.Fatal(err.Error())
	}

	identity, err := mspClient.GetSigningIdentity(username)
	if err != nil {
		t.Fatal(err.Error())
	}
	channelContext := sdk.ChannelContext("mychannel", fabsdk.WithUser(username), fabsdk.WithIdentity(identity))
	channelClient, err := channel.New(channelContext)
	if err != nil {
		t.Fatal(err.Error())
	}
	queryRequest := channel.Request{
		ChaincodeID: "mycc",
		Fcn:         "query",
		Args:        util.PackageArgs([]string{"a"}),
	}

	response, err := channelClient.Query(queryRequest)
	if err != nil {
		t.Fatal(err.Error())
	}
	t.Log(string(response.Payload))
}
