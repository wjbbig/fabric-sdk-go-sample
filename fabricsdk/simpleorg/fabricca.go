package simpleorg

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/context"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

type FabricUser struct {
	Username string
	Secret string
	Affiliation string
	Type string
}

func RegisterUser(configPath string, username string, password string, affiliation string, attributes []msp.Attribute) (string, error) {
	sdk, _ := fabsdk.New(config.FromFile(configPath))
	clientContext := sdk.Context()
	mspClient, _ := msp.New(clientContext, msp.WithOrg("Org1"))
	instance, s := getRegistrarEnrollmentCredentialsWithCAInstance(clientContext, "")
	//登陆registrar
	err := mspClient.Enroll(instance, msp.WithSecret(s))
	if err != nil {
		return "", err
	}

	register := msp.RegistrationRequest{
		Name:        username,
		//使用Type可能导致ca注册的用户无法建立通道
		//Type:        "user",
		Affiliation: affiliation,
		Attributes:  attributes,
		Secret:      password,
	}
	//注册用户，secret为登陆密钥
	secret, err := mspClient.Register(&register)
	if err != nil {
		return "", err
	}
	return secret, nil
}

func getRegistrarEnrollmentCredentialsWithCAInstance( ctxProvider context.ClientProvider, caID string) (string, string) {
	ctx, err := ctxProvider()
	if err != nil {
		panic(err)
	}
	myOrg := ctx.IdentityConfig().Client().Organization
	if caID == "" {
		caID = ctx.EndpointConfig().NetworkConfig().Organizations[myOrg].CertificateAuthorities[0]
	}
	caConfig, ok := ctx.IdentityConfig().CAConfig(caID)
	if !ok {
		panic(err)
	}
	return caConfig.Registrar.EnrollID, caConfig.Registrar.EnrollSecret
}

