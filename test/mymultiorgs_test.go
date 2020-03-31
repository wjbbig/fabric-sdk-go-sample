package main

import (
	"fabric-sdk-go-test/fabricsdk/multiorg"
	"testing"
)

func TestCreateSdk(t *testing.T) {
	var mc multiorg.MultiorgContext
	mc.ChannelID = "mychannel"
	mc.Init()
	cc, err := mc.QueryCC("testcccc", "query", []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cc.Payload))
	s, err := mc.ExecuteCC("testcccc", "invoke", []string{"a", "b", "10"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(s)
	cc1, err := mc.QueryCC("testcccc", "query", []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cc1.Payload))
	mc.TearDown()

}

func TestInstallAndInstantiateCC(t *testing.T) {
	var mc multiorg.MultiorgContext
	mc.ChannelID = "mychannel"
	mc.Init()
	ccPath := "github.com/mychaincode/go"
	ccVersion := "1.0"
	ccName := "testmycc"
	_, err := mc.InstallCC(ccPath, ccName, ccVersion)
	if err != nil {
		t.Fatal(err)
	}
	ep := "AND ('Org1MSP.peer','Org2MSP.peer')"
	args := []string{"Init", "a", "100", "b", "200"}
	_, err = mc.InstantiateCC(ccName, ccVersion, ccPath, ep, args)
	if err != nil {
		t.Fatal(err)
	}
	cc, err := mc.QueryCC(ccName, "query", []string{"a"})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cc.Payload))
	mc.TearDown()
}


