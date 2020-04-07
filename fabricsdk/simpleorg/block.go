package simpleorg

import (
	"encoding/hex"
	"fabric-sdk-go-test/util"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
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
	env, err := util.GetEnvelopeFromBlock(firstTx)
	if err != nil {
		return "", err
	}
	if env == nil {
		return "", errors.New("nil envelope")
	}
	payload, err := util.GetPayload(env)
	if err != nil {
		return "", errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	channelHeaderBytes := payload.Header.ChannelHeader
	channelHeader := &common.ChannelHeader{}
	err = proto.Unmarshal(channelHeaderBytes, channelHeader)
	if err != nil {
		return "", errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	fmt.Println(channelHeader.Type)
	timestamp := channelHeader.Timestamp
	return time.Unix(timestamp.Seconds, 0).Format("2006-01-02 15:04:05"), nil
}

func GenerateBlockHash(block *common.Block) (string, error) {
	bytes := util.Bytes(block.Header)
	blockHash := util.ComputeSHA256(bytes)
	return hex.EncodeToString(blockHash), nil
}

func GetInvalidTxCount(block *common.Block) {

}

func GetTx(block *common.Block, i uint64) []byte {
	dataBytes := block.Data.Data[i]
	return dataBytes
}

func GetTxFilter(block *common.Block, i uint64) byte {
	return block.Metadata.Metadata[2][i]
}
