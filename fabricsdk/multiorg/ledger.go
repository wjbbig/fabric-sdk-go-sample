package multiorg

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func (mc *MultiorgContext) QueryBlockByTxID(txID string) (*common.Block, error) {
	return mc.LedgerClient.QueryBlockByTxID(fab.TransactionID(txID))
}

//QueryBlockByNumber 通过区块号查询区块
func (mc *MultiorgContext) QueryBlockByNumber(blockNum uint64) (*common.Block, error) {
	return mc.LedgerClient.QueryBlock(blockNum)
}

func (mc *MultiorgContext) QueryInfo() (*fab.BlockchainInfoResponse, error) {
	return mc.LedgerClient.QueryInfo()
}

func (mc *MultiorgContext) QueryBlockNumberByHash(hash string) (*common.Block, error) {
	return mc.LedgerClient.QueryBlockByHash([]byte(hash))
}

func (mc *MultiorgContext) QueryConfig() (fab.ChannelCfg, error) {
	return mc.LedgerClient.QueryConfig()
}

