package chaincode

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 链码1的名称
var PRIVATE = "private"

// 通道名称
var CHANNEL = "mychannel"

// 账户结构体
type Account struct {
	Owner  string  `json:"所有者"`
	Amount float64 `json:"余额"`
}

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

// OKAY!
func (s *SmartContract) GetRequestFromPrivate1(ctx contractapi.TransactionContextInterface) string {
	fmt.Println("GetGetRequestFromPrivate")
	response := ctx.GetStub().InvokeChaincode(PRIVATE, [][]byte{[]byte("ReadRequest")}, CHANNEL)
	var res = response.GetMessage()
	fmt.Println(res)
	return string(response.Payload)
}

// OKAY!
func (s *SmartContract) GetRequestFromPrivate2(ctx contractapi.TransactionContextInterface) string {
	response := ctx.GetStub().InvokeChaincode(PRIVATE, [][]byte{[]byte("ReadRequest")}, CHANNEL)
	var res = response.GetMessage()
	fmt.Println(res)
	return response.String()
}

// OKAY!
func (s *SmartContract) GetRequestFromPrivate3(ctx contractapi.TransactionContextInterface) int32 {
	response := ctx.GetStub().InvokeChaincode(PRIVATE, [][]byte{[]byte("ReadRequest")}, CHANNEL)
	var res = response.GetMessage()
	fmt.Println(res)
	return response.Status
}

// 初始化账户;
func (s *SmartContract) CreateAccount(ctx contractapi.TransactionContextInterface) error {
	owner, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		fmt.Println(err)
		return err
	}

	isExist, err := s.MyAccountExist(ctx)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("账户已存在.")
	}
	fmt.Println("初始化账户......")
	fmt.Println("账户名称:", owner)

	var account = Account{
		Owner:  owner,
		Amount: 0,
	}
	accountJSON, err := json.Marshal(account)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return ctx.GetStub().PutState(owner, accountJSON)
}

// 检测账户是否存在
func (s *SmartContract) MyAccountExist(ctx contractapi.TransactionContextInterface) (bool, error) {
	owner, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	accountJSON, err := ctx.GetStub().GetState(owner)
	if err != nil {
		return false, err
	}
	return accountJSON != nil, nil
}

// 充值金额
func (s *SmartContract) Recharge(ctx contractapi.TransactionContextInterface, amount float64) error {
	isExist, err := s.MyAccountExist(ctx)
	if err != nil {
		return err
	}
	if isExist != true {
		return fmt.Errorf("请先创建账户.")
	}

	owner, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}
	accountJSON, err := ctx.GetStub().GetState(owner)
	if err != nil {
		return err
	}
	var account Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return err
	}
	account.Amount += amount
	accountJSON, err = json.Marshal(account)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutState(owner, accountJSON)
}

// 转账
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, to string, amount float64) error {

	// 检测转账人的账户是否存在
	isFromExist, err := s.MyAccountExist(ctx)
	if err != nil {
		return err
	}
	if !isFromExist {
		return fmt.Errorf("请先创建账户")
	}

	// 获取转账人和收账人的账户信息
	from, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return err
	}
	fromAccountJSON, err := ctx.GetStub().GetState(from)
	if err != nil {
		return err
	}

	toAccountJSON, err := ctx.GetStub().GetState(to)
	if err != nil {
		return err
	}

	var (
		fromAccount Account
		toAccount   Account
	)
	err = json.Unmarshal(fromAccountJSON, &fromAccount)
	if err != nil {
		return err
	}
	err = json.Unmarshal(toAccountJSON, &toAccount)
	if err != nil {
		return err
	}

	// 检查转账人的余额是否足够
	if fromAccount.Amount < amount {
		return fmt.Errorf("余额不足，请先充值.")
	}

	// 转账
	fromAccount.Amount -= amount
	toAccount.Amount += amount

	fromAccountJSON, err = json.Marshal(fromAccount)
	if err != nil {
		return err
	}
	toAccountJSON, err = json.Marshal(toAccount)
	if err != nil {
		return err
	}

	// 更新ledger中的数据
	err = ctx.GetStub().PutState(from, fromAccountJSON)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(to, toAccountJSON)
	if err != nil {
		return err
	}
	return nil
}

// 查看账户
func (s *SmartContract) ViewMyAccount(ctx contractapi.TransactionContextInterface) (*Account, error) {
	isExist, err := s.MyAccountExist(ctx)
	if err != nil {
		return nil, nil
	}
	if !isExist {
		return nil, fmt.Errorf("账户不存在")
	}

	owner, err := s.GetSubmittingClientIdentity(ctx)
	if err != nil {
		return nil, err
	}

	accountJSON, err := ctx.GetStub().GetState(owner)
	if err != nil {
		return nil, err
	}

	var account *Account
	err = json.Unmarshal(accountJSON, &account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *SmartContract) GetSubmittingClientIdentity(ctx contractapi.TransactionContextInterface) (string, error) {

	b64ID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", fmt.Errorf("Failed to read clientID: %v", err)
	}
	decodeID, err := base64.StdEncoding.DecodeString(b64ID)
	if err != nil {
		return "", fmt.Errorf("failed to base64 decode clientID: %v", err)
	}
	s1 := ctx.GetClientIdentity()
	s2, err := ctx.GetClientIdentity().GetX509Certificate()
	s3, err := ctx.GetClientIdentity().GetMSPID()
	fmt.Println("ctx.GetClientIdentity().GetMSPID():", s3)
	fmt.Println("-------------------------------------------------")
	fmt.Println("ctx.GetClientIdentity().GetX509Certificate():", s2)
	fmt.Println("-------------------------------------------------")
	fmt.Println("ctx.GetClientIdentity():", s1)
	fmt.Println("-------------------------------------------------")
	fmt.Println("ctx.GetClientIdentity().GetID():", string(decodeID))
	fmt.Println("-------------------------------------------------")
	return s3, nil
}

// IntToBytes 将int类型的数转化为字节并以小端存储
func IntToBytes(intNum int) []byte {
	uint16Num := uint16(intNum)
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.LittleEndian, uint16Num)
	return buf.Bytes()
}

// BytesToInt 将以小端存储的长为1/2字节的数转化成int类型的数
func BytesToInt(bytesArr []byte) int {
	var intNum int
	if len(bytesArr) == 1 {
		bytesArr = append(bytesArr, byte(0))
		intNum = int(binary.LittleEndian.Uint16(bytesArr))
	} else if len(bytesArr) == 2 {
		intNum = int(binary.LittleEndian.Uint16(bytesArr))
	}

	return intNum
}

func (s *SmartContract) Test1(ctx contractapi.TransactionContextInterface) string {
	fmt.Println("Test")
	response := ctx.GetStub().InvokeChaincode(PRIVATE, [][]byte{[]byte("Test"), []byte("187")}, CHANNEL)
	var res = response.GetMessage()
	fmt.Println(res)
	return string(response.Payload)
}

func (s *SmartContract) Test2(ctx contractapi.TransactionContextInterface) string {
	fmt.Println("Test")
	response := ctx.GetStub().InvokeChaincode(PRIVATE, [][]byte{[]byte("Test"), IntToBytes(187)}, CHANNEL)
	var res = response.GetMessage()
	fmt.Println(res)
	return string(response.Payload)
}

func (s *SmartContract) Test3(ctx contractapi.TransactionContextInterface) string {
	fmt.Println("Test")
	var myArray = []int{1, 2, 3, 4, 4, 3, 2, 1}
	myArrayJSON, err := json.Marshal(myArray)
	if err != nil {
		return "error"
	}
	response := ctx.GetStub().InvokeChaincode(PRIVATE, [][]byte{[]byte("TestArray"), myArrayJSON}, CHANNEL)
	var res = response.GetMessage()
	fmt.Println(res)
	return string(response.Payload)
}
