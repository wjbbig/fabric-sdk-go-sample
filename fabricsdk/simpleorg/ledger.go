package simpleorg

import (
	"github.com/hyperledger/fabric-protos-go/common"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

func (fm *FabricModel) QueryBlockByTxID(txID string) (*common.Block, error) {
	return fm.LedgerClient.QueryBlockByTxID(fab.TransactionID(txID))
}

//QueryBlockByNumber 通过区块号查询区块
func (fm *FabricModel) QueryBlockByNumber(blockNum uint64) (*common.Block, error) {
	return fm.LedgerClient.QueryBlock(blockNum)
}

func (fm *FabricModel) QueryInfo() (*fab.BlockchainInfoResponse, error) {
	return fm.LedgerClient.QueryInfo()
}

func (fm *FabricModel) QueryBlockByHash(hash []byte) (*common.Block, error) {
	return fm.LedgerClient.QueryBlockByHash(hash)
}

func (fm *FabricModel) QueryConfig() (fab.ChannelCfg, error) {
	return fm.LedgerClient.QueryConfig()
}

func (fm *FabricModel) QueryTransaction(transactionID fab.TransactionID) (*pb.ProcessedTransaction, error) {
	return fm.LedgerClient.QueryTransaction(transactionID)
}
