package service

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
)

func (t *ServiceSetup) SetVote(voter string, activity string, candidate string) (string, error) {
	eventID := "voteInfo"
	reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	defer t.Client.UnregisterChaincodeEvent(reg)
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "vote", Args: [][]byte{[]byte(voter), []byte(activity), []byte(candidate)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	err = eventResult(notifier, eventID)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) FindVotesbyName(candidate string, activity string) (string, error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "select", Args: [][]byte{[]byte(candidate), []byte(activity)}}

	response, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(response.Payload), nil
}

func (t *ServiceSetup) SetInfo(name, num string) (string, error) {

	//eventID := "eventSetInfo"
	//reg, notifier := regitserEvent(t.Client, t.ChaincodeID, eventID)
	//defer t.Client.UnregisterChaincodeEvent(reg)


	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "set", Args: [][]byte{[]byte(name), []byte(num)}}
	respone, err := t.Client.Execute(req)
	if err != nil {
		return "", err
	}

	return string(respone.TransactionID), nil
}

func (t *ServiceSetup) GetInfo(name string) (string, error){

	req := channel.Request{ChaincodeID: t.ChaincodeID, Fcn: "get", Args: [][]byte{[]byte(name)}}
	respone, err := t.Client.Query(req)
	if err != nil {
		return "", err
	}

	return string(respone.Payload), nil
}

func (t *ServiceSetup) Init_Activity_Candiate(activity,candidate1,candidate2,candidate3 string) (string,error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID,Fcn: "init_Activity_Candiate",Args: [][]byte{[]byte(activity),[]byte(candidate1),[]byte(candidate2),[]byte(candidate3)}}
	respone ,err := t.Client.Execute(req)
	if err != nil {
		return "",err
	}
	return string(respone.TransactionID),nil
}

func (t *ServiceSetup)Vote(VID1,CID1,AID1 string) (string,error) {
	req := channel.Request{ChaincodeID: t.ChaincodeID,Fcn: "Vote",Args: [][]byte{[]byte(AID1),[]byte(CID1),[]byte(VID1)}}
	respone,err:=t.Client.Execute(req)
	if err != nil {
		return "",err
	}
	return string(respone.TransactionID),nil
}

func (t *ServiceSetup) GetVotecandidate(CID1,AID1 string) (string,error) {
	fmt.Println("GetVotecandidate......................",CID1,AID1)
	req := channel.Request{ChaincodeID: t.ChaincodeID,Fcn: "GetVotecandidate",Args: [][]byte{[]byte(AID1),[]byte(CID1)}}
	respone,err := t.Client.Execute(req)
	fmt.Printf("service",respone.Payload[0])
	//fmt.Printf("service string",string(respone.Payload))
	if err != nil {
		return "",err
	}
	return string(respone.Payload),nil
}
