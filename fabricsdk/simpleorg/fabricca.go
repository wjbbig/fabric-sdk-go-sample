package fabricsdk

import (
	"github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
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
	sdk, err := fabsdk.New(config.FromFile(configPath))
	if err != nil {
		return "", err
	}
	clientContext := sdk.Context()
	mspClient, err := msp.New(clientContext)
	if err != nil {
		return "", err
	}
	register := msp.RegistrationRequest{
		Name:        username,
		Type:        "user",
		Affiliation: affiliation,
		Attributes:  attributes,
		Secret:      password,
	}
	secret, err := mspClient.Register(&register)
	return secret, err
}


