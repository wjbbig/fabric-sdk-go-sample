package simpleorg

import (
	"encoding/hex"
	"fabric-sdk-go-test/util"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset"
	"github.com/hyperledger/fabric-protos-go/ledger/rwset/kvrwset"
	"github.com/hyperledger/fabric-protos-go/msp"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

func GetTxHash(tx []byte) (string, error) {
	channelHeader, err := util.GetChannelHeader(tx)
	if err != nil {
		return "", errors.Wrap(err, "error getting channel header")
	}
	return channelHeader.TxId, nil
}

func GetChannelId(tx []byte) (string, error) {
	channelHeader, err := util.GetChannelHeader(tx)
	if err != nil {
		return "", errors.Wrap(err, "error getting channel header")
	}
	return channelHeader.ChannelId, nil
}

func GetTxType(tx []byte) (int32, error) {
	channelHeader, err := util.GetChannelHeader(tx)
	if err != nil {
		return -1, errors.Wrap(err, "error getting channel header")
	}
	return channelHeader.Type, nil
}

func GetCreatTime(tx []byte) (string, error) {
	channelHeader, err := util.GetChannelHeader(tx)
	if err != nil {
		return "", errors.Wrap(err, "error getting channel header")
	}
	return time.Unix(channelHeader.Timestamp.Seconds, 0).Format("2006-01-02 15:04:05"), nil
}

func GetChaincodeName(tx []byte) (string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return "", errors.Wrap(err, "error getting chaincode action")
	}
	return chaincodeAction.ChaincodeId.Name, nil
}

func GetResponseStatus(tx []byte) (int32, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return -1, errors.Wrap(err, "error getting chaincode action")
	}
	return chaincodeAction.Response.Status, nil
}

func GetCreatorMSPId(tx []byte) (string, error) {
	env, err := util.GetEnvelopeFromBlock(tx)
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
	signatureHeader := &common.SignatureHeader{}
	err = proto.Unmarshal(payload.Header.SignatureHeader, signatureHeader)
	if err != nil {
		return "", errors.Wrap(err, "error extracting signature header")
	}
	mspContent := &msp.SerializedIdentity{}
	err = proto.Unmarshal(signatureHeader.Creator, mspContent)
	if err != nil {
		return "", err
	}
	return mspContent.Mspid, nil
}

func GetEndorserMSPId(tx []byte) ([]string, error) {
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	var endorserMSPIds []string
	for _, endorsement := range chaincodeActionPayload.Action.Endorsements {
		mspContent := &msp.SerializedIdentity{}
		err = proto.Unmarshal(endorsement.Endorser, mspContent)
		endorserMSPIds = append(endorserMSPIds, mspContent.Mspid)
	}
	return endorserMSPIds, nil
}

func GetReadSet(tx []byte) ([]map[string]string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var readSets []map[string]string
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" {
			readSet := make(map[string]string)
			readSet["namespace"] = set.Namespace
			kvrwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, kvrwset)
			if err != nil {
				return nil, err
			}
			var setStr []string
			for _, read := range kvrwset.Reads {
				setStr = append(setStr, read.String())
			}
			readSet["set"] = strings.Join(setStr, "")
			readSets = append(readSets, readSet)
		}
	}
	return readSets, nil
}

func GetReadKeyList(tx []byte) ([]string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var readSet []string
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" {
			kvrwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, kvrwset)
			if err != nil {
				return nil, err
			}
			for _, read := range kvrwset.Reads {
				if read.Version == nil {
					readSet = append(readSet, read.Key+"!#null")
				} else {
					readSet = append(readSet, read.Key+"!#"+strconv.Itoa(int(read.Version.BlockNum))+"_"+
						strconv.Itoa(int(read.Version.TxNum)))
				}
			}
		}
	}
	return readSet, nil
}

func GetWriteSet(tx []byte) ([]map[string]string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var writeSets []map[string]string
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" {
			writeSet := make(map[string]string)
			writeSet["namespace"] = set.Namespace
			kvrwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, kvrwset)
			if err != nil {
				return nil, err
			}
			var setStr []string
			for _, write := range kvrwset.Writes {
				setStr = append(setStr, write.String())
			}
			writeSet["set"] = strings.Join(setStr, "")
			writeSets = append(writeSets, writeSet)
		}
	}
	return writeSets, nil
}

func GetWriteKeyList(tx []byte) ([]string, error) {
	chaincodeAction, err := getChaincodeAction(tx)
	if err != nil {
		return nil, err
	}
	result := &rwset.TxReadWriteSet{}
	err = proto.Unmarshal(chaincodeAction.Results, result)
	if err != nil {
		return nil, err
	}
	var writeSet []string
	for _, set := range result.NsRwset {
		if set.Namespace != "lscc" {

			kvrwset := &kvrwset.KVRWSet{}
			err := proto.Unmarshal(set.Rwset, kvrwset)
			if err != nil {
				return nil, err
			}
			for _, write := range kvrwset.Writes {
				writeSet = append(writeSet, write.Key)
			}
		}
	}
	return writeSet, nil
}

func GetEndorserSignature(tx []byte) ([]map[string]string, error) {
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	var endorserSignatures []map[string]string
	for _, endorsement := range chaincodeActionPayload.Action.Endorsements {
		endorserSignature := make(map[string]string)
		endorserSignature["signature"] = hex.EncodeToString(endorsement.Signature)
		mspContent := &msp.SerializedIdentity{}
		err = proto.Unmarshal(endorsement.Endorser, mspContent)
		endorserSignature["msp_id"] = mspContent.Mspid
		endorserSignature["cerficate"] = string(mspContent.IdBytes)
		endorserSignatures = append(endorserSignatures, endorserSignature)
	}

	return endorserSignatures, nil
}

func GetChaincodeFunction(tx []byte) (string, error) {
	invokeSpec, err := getTxAllArgs(tx)
	if err != nil {
		return "", err
	}
	return string(invokeSpec.ChaincodeSpec.Input.Args[0]), nil
}

func GetFunctionParameters(tx []byte) ([]string, error) {
	invokeSpec, err := getTxAllArgs(tx)
	if err != nil {
		return nil, err
	}
	var args []string
	for i := 1; i < len(invokeSpec.ChaincodeSpec.Input.Args); i++ {
		args = append(args, string(invokeSpec.ChaincodeSpec.Input.Args[i]))
	}
	return args, err
}

func getTxAllArgs(tx []byte) (*pb.ChaincodeInvocationSpec, error) {

	propPayload := &pb.ChaincodeProposalPayload{}
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	if err := proto.Unmarshal(chaincodeActionPayload.ChaincodeProposalPayload, propPayload); err != nil {
		return nil, errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	invokeSpec := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(propPayload.Input, invokeSpec)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	return invokeSpec, nil
}

func getChaincodeActionPayload(tx []byte) (*pb.ChaincodeActionPayload, error) {
	env, err := util.GetEnvelopeFromBlock(tx)
	if err != nil {
		return nil, err
	}
	if env == nil {
		return nil, errors.New("nil envelope")
	}
	payload, err := util.GetPayload(env)
	if err != nil {
		return nil, errors.Wrap(err, "error extracting ChannelHeader from payload")
	}
	transaction, err := util.GetTransaction(payload.Data)
	if err != nil {
		return nil, errors.Wrap(err, "error getting transaction")
	}
	chaincodeActionPayload, err := util.GetChaincodeActionPayload(transaction.Actions[0].Payload)
	if err != nil {
		return nil, errors.Wrap(err, "error getting chaincodeActionPayload")
	}
	return chaincodeActionPayload, err
}

func getChaincodeAction(tx []byte) (*pb.ChaincodeAction, error) {
	chaincodeActionPayload, err := getChaincodeActionPayload(tx)
	if err != nil {
		return nil, err
	}
	proposalResponsePayload := &pb.ProposalResponsePayload{}
	err = proto.Unmarshal(chaincodeActionPayload.Action.ProposalResponsePayload, proposalResponsePayload)
	if err != nil {
		return nil, err
	}
	chaincodeAction := &pb.ChaincodeAction{}
	err = proto.Unmarshal(proposalResponsePayload.Extension, chaincodeAction)
	if err != nil {
		return nil, err
	}
	return chaincodeAction, nil
}
