package simpleorg

import (
	"testing"
)

var fm1 *FabricModel
var tx []byte

func init() {
	fm1 = &FabricModel{
		ConfigFile: "/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml",
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		UserName:   "User1",
		ChannelID:  "mychannel",
		HasInit:    false,
	}
	fm1.Init()

	block1, _ := fm.QueryBlockByNumber(4)
	tx = block1.Data.Data[0]
}

func TestGetTxHash(t *testing.T) {
	hash, err := GetTxHash(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}

func TestGetChannelId(t *testing.T) {
	id, err := GetChannelId(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}

func TestGetTxType(t *testing.T) {
	txType, err := GetTxType(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(txType)
}

func TestGetChaincodeFunction(t *testing.T) {
	function, err := GetChaincodeFunction(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(function)
}

func TestGetFunctionParameters(t *testing.T) {
	parameters, err := GetFunctionParameters(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(parameters)
}
