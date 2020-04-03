package simpleorg

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/ledger"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/resmgmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type FabricModel struct {
	ConfigFile        string            //sdk的配置文件路径
	ChainCodeID       string            // 链码名称
	ChaincodePath     string            // 链码在工程中的存放目录
	ChaincodeGoPath   string            // GOPATH
	OrgAdmin          string            // 组织的管理员用户
	OrgName           string            // 组织名称
	OrgID             string            // 组织id
	UserName          string            // 组织的普通用户
	ChannelID         string            // 通道id
	ChannelConfigPath string            //组织的通道文件路径
	OrdererName       string            // 将组织添加到通道时候使用!
	LedgerClient      *ledger.Client    //账本客户端
	Sdk               *fabsdk.FabricSDK // 保存实例化后的sdk
	ResMgmtClient     *resmgmt.Client   // 资源管理客户端,也需要在安装链码时候的使用
	ChannelClient     *channel.Client   // 通道客户端
	HasInit           bool              // 是否已经初始化了sdk
}

func (fm *FabricModel) Init() {
	if fm.HasInit {
		panic("sdk已初始化")
	}
	sdk, err := fabsdk.New(config.FromFile(fm.ConfigFile))
	if err != nil {
		panic(err)
	}
	fmt.Println("sdk创建完成")
	fm.Sdk = sdk
	//defer sdk.Close()
	clientContext := sdk.Context(fabsdk.WithUser(fm.OrgAdmin), fabsdk.WithOrg(fm.OrgName))
	if clientContext == nil {
		panic("根据指定的组织名称与管理员创建资源管理客户端Context失败")
	}
	resMgmtClient, err := resmgmt.New(clientContext)
	if err != nil {
		panic(err)
	}
	fmt.Println("资源管理器创建完毕")
	fm.ResMgmtClient = resMgmtClient
	if fm.ChannelConfigPath == "" {
		channelContext := sdk.ChannelContext(fm.ChannelID, fabsdk.WithUser(fm.UserName))
		channelClient, err := channel.New(channelContext)
		if err != nil {
			panic(err)
		}
		fmt.Println("通道创建完成")
		ledger, err := ledger.New(channelContext)
		if err != nil {
			panic(err)
		}
		fm.LedgerClient = ledger
		fm.ChannelClient = channelClient
		fm.HasInit = true
	}
	fmt.Println("sdk初始化完毕")
}
