package simpleorg

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"testing"
)

var fm *FabricModel
var block *common.Block
var tx []byte
func init() {
	fm = &FabricModel{
		ConfigFile: "/home/fujitsu/IdeaProjects/com/fujitsu/fabric-sdk-go-test/config/connection-config.yaml",
		OrgAdmin:   "Admin",
		OrgName:    "Org1",
		UserName:   "User1",
		ChannelID:  "mychannel",
		HasInit:    false,
	}
	fm.Init()

	block, _ = fm.QueryBlockByNumber(4)
	tx = block.Data.Data[0]
}

func TestGetBlockNumber(t *testing.T) {
	number := GetBlockNumber(block)
	t.Log(number)
}

func TestGetDataHash(t *testing.T) {
	dataHash := GetDataHash(block)
	t.Log(dataHash)
}

func TestGetPreviousHash(t *testing.T) {
	hash := GetPreviousHash(block)
	t.Log(hash)
}

func TestGetTxCount(t *testing.T) {
	count := GetTxCount(block)
	t.Log(count)
}

func TestGetChannelVersion(t *testing.T) {
	config, err := fm.QueryConfig()
	if err != nil {
		t.Fatal(err)
	}
	version := GetChannelVersion(config)
	t.Log(version)
}

func TestGetCreateTime(t *testing.T) {
	time, err := GetCreateTime(block)

	if err != nil {
		t.Fatal(err)
	}
	t.Log(time)
}

func TestGenerateBlockHash(t *testing.T) {
	hash, err := GenerateBlockHash(block)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hash)
}

func TestGetInvalidTxCount(t *testing.T) {
	count, err := GetInvalidTxCount(fm.LedgerClient, block)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(count)
}

func TestGetTx(t *testing.T) {
	tx := GetTx(block, 0)
	t.Log(tx)
}

func TestGetTxFilter(t *testing.T) {
	filter := GetTxFilter(block, 0)
	t.Log(filter)
}
