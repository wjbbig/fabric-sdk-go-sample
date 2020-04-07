package simpleorg

import (
	"encoding/hex"
	"fabric-sdk-go-test/util"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/pkg/errors"
	"time"
)

func GetBlockNumber(block *common.Block) uint64 {
	return block.Header.Number
}

func GetDataHash(block *common.Block) string {
	return hex.EncodeToString(block.Header.DataHash)
}

func GetPreviousHash(block *common.Block) string {
	return hex.EncodeToString(block.Header.PreviousHash)
}

func GetTxCount(block *common.Block) int {
	return len(block.Data.Data)
}

func GetChannelVersion(config fab.ChannelCfg) uint64 {
	return config.Versions().Channel.Version
}

func GetCreateTime(block *common.Block) (string, error) {
	firstTx := block.Data.Data[0]
	channelHeader, err := util.GetChannelHeader(firstTx)
	if err != nil {
		return "", errors.Wrap(err, "error getting channelHeader")
	}
	timestamp := channelHeader.Timestamp
	return time.Unix(timestamp.Seconds, 0).Format("2006-01-02 15:04:05"), nil
}

func GenerateBlockHash(block *common.Block) (string, error) {
	bytes := util.Bytes(block.Header)
	blockHash := util.ComputeSHA256(bytes)
	return hex.EncodeToString(blockHash), nil
}

func GetInvalidTxCount(lc *ledger.Client, block *common.Block) (uint64, error) {
	var count uint64
	for _, tx := range block.Data.Data {
		channelHeader, err := util.GetChannelHeader(tx)
		if err != nil {
			return 0, errors.Wrap(err, "error getting channelHeader")
		}
		txId := fab.TransactionID(channelHeader.TxId)
		transaction, err := lc.QueryTransaction(txId)
		if err != nil {
			return 0, errors.Wrap(err, "error getting transaction")
		}
		if transaction.ValidationCode != 0 {
			count++
		}
	}
	return count, nil
}

func GetTx(block *common.Block, i uint64) []byte {
	dataBytes := block.Data.Data[i]
	return dataBytes
}

func GetTxFilter(block *common.Block, i uint64) int {
	return int(block.Metadata.Metadata[2][i])
}
