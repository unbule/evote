//package main
//
//import (
//	"fmt"
//	"github.com/hyperledger/fabric/core/chaincode/shim"
//	"github.com/hyperledger/fabric/protos/peer"
//)
//
//type SimpleChaincode struct {
//}
//
//func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
//
//	return shim.Success(nil)
//}
//
//func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
//	fun, args := stub.GetFunctionAndParameters()
//
//	var result string
//	var err error
//	if fun == "set" {
//		result, err = set(stub, args)
//	} else {
//		result, err = get(stub, args)
//	}
//	if err != nil {
//		return shim.Error(err.Error())
//	}
//	return shim.Success([]byte(result))
//}
//
//func set(stub shim.ChaincodeStubInterface, args []string) (string, error) {
//
//	if len(args) != 2{
//		return "", fmt.Errorf("给定的参数错误")
//	}
//
//	err := stub.PutState(args[0], []byte(args[1]))
//	if err != nil{
//		return "", fmt.Errorf(err.Error())
//	}
//
//	return string(args[0]), nil
//}
//
//func get(stub shim.ChaincodeStubInterface, args []string) (string, error) {
//	if len(args) != 1 {
//		return "", fmt.Errorf("给定的参数错误")
//	}
//	result, err := stub.GetState(args[0])
//
//	if err != nil {
//		return "", fmt.Errorf("获取数据发生错误")
//	}
//	if result == nil {
//		return "", fmt.Errorf("根据 %s 没有获取到相应的数据", args[0])
//	}
//	return string(result), nil
//
//}
//
//func main() {
//	err := shim.Start(new(SimpleChaincode))
//	if err != nil {
//		fmt.Printf("启动SimpleChaincode时发生错误: %s", err)
//	}
//}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type  VoteChaincode struct {
}
type Voter struct {
	VID string //VID 为键
	AIDlist []string  //添加活动ID,表明该VID参加的活动
}
type Activity struct {
	AID string //键
	CIDlist []string  //候选人ID列表

}
type Candiater struct {
	CIDandAID string  //键
	Voterecieved int
	Voterlist []string  //存放投票人ID列表 ，某候选人被投的投票人列表
}
//初始化应该先规定好AID和相应的CID
func (t *VoteChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("\nVotingApp Is Starting Up,Init() args\n")
	return shim.Success(nil)
}

func (t *VoteChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fn, args := stub.GetFunctionAndParameters()
	fmt.Println(fn, args)
	var result string
	var err error
	if fn=="init_Activity_Candiate"{
		result,err = init_Activity_Candiate(stub,args)
	} else if fn == "Vote" {
		return Vote(stub, args)
	} else if fn == "GetVotecandidate" {
		result,err =  GetVotecandidate(stub, args)
	}
	if err != nil{
		return shim.Error(err.Error())
	}
	return shim.Success([]byte(result))
}

//初始化 活动AID 候选人CID
//args[0] VID CID AID
func init_Activity_Candiate(stub shim.ChaincodeStubInterface,args []string) (string,error){
	_, args = stub.GetFunctionAndParameters() //字符串转化为切片字符串
	var a  string
	a = args[0]
	var b []string
	for _,v := range args[1:]{ //1包含
		b = append(b, v)
	}
	Activity1 := Activity{}
	Activity1.AID = a
	Activity1.CIDlist = b
	AIDAsBytes, _ := json.Marshal(Activity1)
	err :=stub.PutState(Activity1.AID, AIDAsBytes)
	if err != nil {
		return "",fmt.Errorf(err.Error())
	}
	return "chaincode : init success",nil
}

//AID CID VID
//活动AID，候选人CID，投票人VID
func Vote(stub shim.ChaincodeStubInterface,args []string) peer.Response{
	fmt.Println("投票")
	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}
	var VID1,CID1,AID1,AIDandCID1 string
	VID1 =args[2]
	CID1 =args[1]
	AID1 =args[0]
	AIDandCID1 = AID1+CID1
	Vote1 :=Voter{}
	Candiate1 :=Candiater{}

	//先用VID查找有无AID
	VIDAsBytes, err := stub.GetState(VID1)
	if err !=nil{
		fmt.Println("err!=nil 出错")
		return shim.Error(err.Error())
	}
	err = json.Unmarshal(VIDAsBytes, &Vote1) //反序列化
	if err != nil {
		shim.Error(err.Error())
	}
	for _,v :=range Vote1.AIDlist{
		if v==AID1{
			fmt.Println("该投票人已参加过该活动的投票，不允许投票")
			return shim.Error("该投票人已参加该活动的投票，不允许投票")
		}
	}
	//没有参加过投票时的处理
	//1、将VID 放入AIDlist 表明该VID参加过该活动了
	Vote1.VID=VID1
	Vote1.AIDlist=append(Vote1.AIDlist,AID1)
	VIDAsBytes, _ = json.Marshal(Vote1)
	//更新数据
	err = stub.PutState(VID1, VIDAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}
	//2、将VID加入Voterlist列表里,修改票数
	CIDandAIDAsBytes,err :=stub.GetState(AIDandCID1)
	if err !=nil{
		fmt.Println("err!=nil 出错")
		return shim.Error(err.Error())
	}
	fmt.Println("该活动AID和候选人ID数据在账本中已存在，在此上修改票数")
	err = json.Unmarshal(CIDandAIDAsBytes, &Candiate1) //反序列化
	if err != nil {
		shim.Error(err.Error())
	}
	Candiate1.CIDandAID=AIDandCID1
	Candiate1.Voterlist = append(Candiate1.Voterlist,VID1)
	Candiate1.Voterecieved = Candiate1.Voterecieved +1
	CIDandAIDAsBytes, _ = json.Marshal(Candiate1)
	//更新数据
	err = stub.PutState(AIDandCID1, CIDandAIDAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

//根据AID,CID获取所有候选人所得票数
func GetVotecandidate(stub shim.ChaincodeStubInterface ,args []string) (string,error){
	//Activities := Activity{}
	Candiaters :=Candiater{}
	var A ,C,AandC string
	A = args[0]
	C = args[1]
	AandC = A+C

	AIDAsBytes,_:=stub.GetState(AandC)
	err := json.Unmarshal(AIDAsBytes, &Candiaters)
	if err !=nil{
		return "",fmt.Errorf(err.Error())
	}
	fmt.Println(Candiaters)
	votenumber :=Candiaters.Voterecieved
	fmt.Println(votenumber)
	fmt.Println(string(votenumber))
	//将查询的结果以json字符串的形式写入buffer 此处需要改善
	//fmt.Println(votenumber)
	//VoteNumBytes, _ := json.Marshal(votenumber)
	return string(votenumber),nil
}//

func main() {
	//链码开始
	err := shim.Start(new(VoteChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

