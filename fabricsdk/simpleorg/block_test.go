package simpleorg

import (
	"encoding/hex"
	"fabric-sdk-go-test/util"
	"testing"
)

func TestGetCreateTime(t *testing.T) {
	fm := FabricModel{
		ConfigFile: "/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml",
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		UserName:   "User1",
		ChannelID:  "mychannel",
		HasInit:    false,
	}

	fm.Init()
	defer fm.Sdk.Close()
	block, err := fm.QueryBlockByNumber(4)
	hash, err := GenerateBlockHash(block)
	t.Log(hash)
	time, err := GetCreateTime(block)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time)
	bytes := util.Bytes(block.Header)
	blockHash := util.ComputeSHA256(bytes)
	t.Log(hex.EncodeToString(blockHash))
	block, err = fm.QueryBlockByHash(blockHash)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(block)

	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(time)
}
