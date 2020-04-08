package simpleorg

import (
	"testing"
)

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

func TestGetCreatorMSPId(t *testing.T) {
	id, err := GetCreatorMSPId(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}

func TestGetEndorserMSPId(t *testing.T) {
	id, err := GetEndorserMSPId(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(id)
}

func TestGetCreatTime(t *testing.T) {
	time, err := GetCreatTime(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(time)
}

func TestGetChaincodeName(t *testing.T) {
	name, err := GetChaincodeName(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(name)
}

func TestGetResponseStatus(t *testing.T) {
	status, err := GetResponseStatus(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(status)
}

func TestGetEndorserSignature(t *testing.T) {
	signature, err := GetEndorserSignature(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(signature)
}

func TestGetReadSet(t *testing.T) {
	set, err := GetReadSet(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestGetReadKeyList(t *testing.T) {
	list, err := GetReadKeyList(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}

func TestGetWriteSet(t *testing.T) {
	set, err := GetWriteSet(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(set)
}

func TestGetWriteKeyList(t *testing.T) {
	list, err := GetWriteKeyList(tx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(list)
}
