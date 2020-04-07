package util

import (
	"encoding/asn1"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/factory"
	"github.com/pkg/errors"
	"math"
)

func GetEnvelopeFromBlock(data []byte) (*common.Envelope, error) {
	// Block always begins with an envelope
	var err error
	env := &common.Envelope{}
	if err = proto.Unmarshal(data, env); err != nil {
		return nil, errors.Wrap(err, "error unmarshaling Envelope")
	}
	return env, nil
}

func GetPayload(e *common.Envelope) (*common.Payload, error) {
	payload := &common.Payload{}
	err := proto.Unmarshal(e.Payload, payload)
	return payload, errors.Wrap(err, "error unmarshaling Payload")
}

func ComputeSHA256(data []byte) (hash []byte) {
	hash, err := factory.GetDefault().Hash(data, &bccsp.SHA256Opts{})
	if err != nil {
		panic(fmt.Errorf("failed computing SHA256 on [% x]", data))
	}
	return hash
}

type asn1Header struct {
	Number       int64
	PreviousHash []byte
	DataHash     []byte
}

func Bytes(b *common.BlockHeader) []byte {
	asn1Header := asn1Header{
		PreviousHash: b.PreviousHash,
		DataHash:     b.DataHash,
	}
	if b.Number > uint64(math.MaxInt64) {
		panic(fmt.Errorf("golang does not currently support encoding uint64 to asn1"))
	} else {
		asn1Header.Number = int64(b.Number)
	}
	result, err := asn1.Marshal(asn1Header)
	if err != nil {
		panic(err)
	}
	return result
}

func GetTransaction(txBytes []byte) (*peer.Transaction, error) {
	tx := &peer.Transaction{}
	err := proto.Unmarshal(txBytes, tx)
	return tx, errors.Wrap(err, "error unmarshaling Transaction")
}
