package utils

import (
	"fmt"
	"github.com/kongyixueyuan.com/evote/sdkinit"
	"github.com/kongyixueyuan.com/evote/service"
	"os"
)

const (
	configFile = "config.yaml"
	initialized = false
	SimpleCC = "mycc"
)

func GetFab() *service.ServiceSetup {
	initInfo := &sdkinit.InitInfo{
		ChannelID: "mychannel",
		ChannelConfig: os.Getenv("GOPATH")+"/src/github.com/kongyixueyuan.com/evote/fixtures/artifacts/channel.tx",
		OrgAdmin: "Admin",
		OrgName: "Org1",
		OrdererOrgName: "orderer.example.com",
		ChaincodeID: SimpleCC,
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "github.com/kongyixueyuan.com/evote/chaincode/",
		UserName: "User1",
	}

	sdk,err:=sdkinit.SetupSDK(configFile,initialized)
	if err != nil {
		fmt.Println(err)
	}


	err = sdkinit.CreateChannel(sdk,initInfo)
	if err != nil {
		fmt.Println(err.Error())
	}

	channelClient,err:=sdkinit.InstallAndInstantiateCC(sdk,initInfo)
	if err != nil {
		fmt.Println("InstallAndInstantiateCC  error")
	}
	fmt.Println(channelClient)

	serviceSetup:=service.ServiceSetup{
		ChaincodeID: SimpleCC,
		Client: channelClient,
	}

	/*app:=controllers.Controller{
		Fabric: &serviceSetup,
	}
	return &app*/
	return &serviceSetup
}
