package chaincode

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset

// 数据需求方调用这个方法，将数据请求到ledger

func (s *SmartContract) existRequest(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	requestJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, err
	}
	return requestJSON != nil, nil

}

// 用来购买字段的
func (s *SmartContract) SendRequest(ctx contractapi.TransactionContextInterface, id, TableID string, attributeID []int, proposal int, requetsTime string) error {
	// 判断是否存在这个id的请求
	isExist, err := s.existRequest(ctx, id)

	if err != nil {
		return err
	}

	if isExist {
		// 如果存在，结束函数
		return fmt.Errorf("Request %v already exists, please update your request ID", id)
	} else {
		// 如果不存在再做插入
		// 转账
		err = s.Transfer(ctx, "Org1MSP", ATTRIBUTE_FEE)
		if err != nil {
			return err
		}

		var request = Request{}
		demander, err := ctx.GetClientIdentity().GetMSPID()
		if err != nil {
			return err
		}

		request.Demander = demander
		request.AttributeID = attributeID // 前端进行处理
		request.TableID = TableID
		request.ID = id
		request.RequestTime = requetsTime
		request.Proposal = proposal // 这里是整数,注意在前端展示的时候修改
		request.Complete = false

		log.Println(request)
		requestJSON, err := json.Marshal(&request)
		if err != nil {
			return err
		}
		//err = ctx.GetStub().PutPrivateData(requestCollection, request.ID, requestJSON)
		err = ctx.GetStub().PutState(request.ID, requestJSON)
		if err != nil {
			return err
		}
		return nil
	}

}

// 用来购买原始表格的

// 数据提供方调用该方法，读取ledger中的请求， 也是HandleRequest的辅助函数
// 注意读取的时候只读取complete为false的请求
func (s *SmartContract) ReadRequest(ctx contractapi.TransactionContextInterface) ([]*Request, error) {
	// GetPrivateDataByRange()读取从[false, true)的值
	//requestIterator, err := ctx.GetStub().GetPrivateDataByRange(requestCollection, "", "")
	requestIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer requestIterator.Close()

	var proposals []*Request
	for requestIterator.HasNext() {
		requestJSON, err := requestIterator.Next()
		if err != nil {
			return nil, err
		}
		var temp Request
		err = json.Unmarshal(requestJSON.Value, &temp)
		proposals = append(proposals, &temp)
	}
	return proposals, nil
}
func (s *SmartContract) ReadRequestForView(ctx contractapi.TransactionContextInterface) ([]*Request, error) {
	// GetPrivateDataByRange()读取从[false, true)的值
	//requestIterator, err := ctx.GetStub().GetPrivateDataByRange(requestCollection, "", "")
	requestIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer requestIterator.Close()

	var proposals []*Request
	for requestIterator.HasNext() {
		requestJSON, err := requestIterator.Next()
		if err != nil {
			return nil, err
		}
		var temp Request
		err = json.Unmarshal(requestJSON.Value, &temp)
		if err != nil {
			return nil, err
		}
		response := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("ReadMyTableByID"), []byte(temp.ID)}, MYCHANNEL)
		if response.Status != 200 {
			return nil, fmt.Errorf("Failed to invoke chaincode.")
		}

		var targetTable []*Table
		err = json.Unmarshal([]byte(response.Payload), &targetTable)
		if err != nil {
			return nil, err
		}
		if targetTable == nil {
			return nil, fmt.Errorf("No table for id: %v", temp.TableID)
		}
		// 按理应该设置成两个struct，偷懒不设置了
		temp.TableID = targetTable[0].TableName

		proposals = append(proposals, &temp)
	}
	return proposals, nil
}

// service和serviceDescription可以用一个map[string]来实现

// 数据提供方调用这个方法来处理ledger中的请求
func (s *SmartContract) HandleAll(ctx contractapi.TransactionContextInterface) ([]*Response, error) {
	// 数据提供方能够对外提供的服务
	// 获取请求数据库中的所有请求
	var requests, err = s.ReadRequest(ctx)
	if err != nil {
		return nil, err
	}

	// 获取反射值，开始处理请求
	value := reflect.ValueOf(s)
	// 要返回的response切片
	var res []*Response
	for _, request := range requests {
		// 要是该请求完成状态是true说明已经处理过了，跳过
		if request.Complete == true {
			continue
		}

		// 根据这个请求的proposal来获取对应的智能合约函数
		f := value.MethodByName(PROPOSALS[request.Proposal])

		// temp是待获得的查询结果
		var temp float64
		// 调用proposal对应的函数,并传入其需要的参数
		if len(request.AttributeID) == 2 {
			temp, err = s.GetPearson(ctx, request.TableID, request.AttributeID)
		} else {
			temp = f.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(request.TableID), reflect.ValueOf(request.AttributeID[0])})[0].Float()
		}

		fmt.Println(temp)

		//var detailedProposal = serviceDescription[request.Proposal]

		targetTable, err := s.ReadMyTableByID(ctx, request.TableID)
		if err != nil {
			return nil, err
		}
		if targetTable == nil {
			return nil, fmt.Errorf("No table for id: %v", request.TableID)
		}
		var response_attributes []string
		for _, v := range request.AttributeID {
			response_attributes = append(response_attributes, ATTRIBUTES_STRING[v])
		}
		var response = Response{
			ID:          request.ID,
			Attribute:   response_attributes,
			TableName:   targetTable[0].TableName,
			ProposalStr: PROPOSALS_STRINGS[request.Proposal],
			RequestTime: request.RequestTime,
			Amount:      temp,
		}
		fmt.Println("hello")
		res = append(res, &response)
	}
	return res, nil
}

func (s *SmartContract) HandleSingle(ctx contractapi.TransactionContextInterface, requestID string) ([]*Response, error) {
	requestJSON, err := ctx.GetStub().GetState(requestID)
	if err != nil {
		return nil, err
	}
	if requestJSON == nil {
		return nil, fmt.Errorf("No such request.")
	}

	var request Request
	err = json.Unmarshal(requestJSON, &request)
	if err != nil {
		return nil, err
	}

	if request.Complete == true {
		return nil, fmt.Errorf("Already handled.")
	}

	isExist, err := s.MyTableExists(ctx, request.TableID)
	if err != nil {
		return nil, fmt.Errorf("Failed to read from world state %v", err)
	}
	if !isExist {
		return nil, fmt.Errorf("Table dose not exist.")
	}

	targetTable, err := s.ReadMyTableByID(ctx, request.TableID)
	if err != nil {
		return nil, err
	}
	if targetTable == nil {
		return nil, fmt.Errorf("No table for id: %v", request.TableID)
	}

	var response_attributes []string
	for _, v := range request.AttributeID {
		response_attributes = append(response_attributes, ATTRIBUTES_STRING[v])
	}
	var (
		v    = reflect.ValueOf(s)
		temp float64
		f    = v.MethodByName(PROPOSALS[request.Proposal])
	)
	if len(request.AttributeID) == 2 {
		temp, err = s.GetPearson(ctx, request.TableID, request.AttributeID)
	} else {
		temp = f.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(request.TableID), reflect.ValueOf(request.AttributeID[0])})[0].Float()
	}
	var (
		res      []*Response
		response = Response{
			ID:          request.ID,
			TableName:   targetTable[0].TableName,
			Attribute:   response_attributes,
			ProposalStr: PROPOSALS_STRINGS[request.Proposal],
			RequestTime: request.RequestTime,
			Amount:      temp,
		}
	)
	res = append(res, &response)
	// 如果编号是3，说明就是做整表的购买；
	if request.Proposal == 3 {

	}
	return res, nil
}

func (s *SmartContract) SendResponse(ctx contractapi.TransactionContextInterface, responseStr string) error {
	// 在private data 1上做了range query就不能在private data 2上做range query，所有handleRequest和sendResponse必须分开定义
	var responses []*Response
	err := json.Unmarshal([]byte(responseStr), &responses)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	for _, response := range responses {
		err := s.createResponse(ctx, *response)
		if err != nil {
			return err
		}

		// 将处理过的请求标记为已经完成
		//requestJSON, err := ctx.GetStub().GetPrivateData(requestCollection, response.ID)
		requestJSON, err := ctx.GetStub().GetState(response.ID)
		if err != nil {
			return err
		}
		var request *Request
		err = json.Unmarshal(requestJSON, &request)
		if err != nil {
			return err
		}
		request.Complete = true

		requestJSON, err = json.Marshal(request)
		if err != nil {
			return err
		}

		//err = ctx.GetStub().PutPrivateData(requestCollection, request.ID, requestJSON)
		err = ctx.GetStub().PutState(request.ID, requestJSON)
		if err != nil {
			return err
		}

		// err := ctx.GetStub().PutPrivateData(requestCollection, response.RequestTime, )
	}
	return nil
}
