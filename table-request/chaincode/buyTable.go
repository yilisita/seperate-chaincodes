/*
*	我们使用InvokeChaincode，主要是想要跨链码调用返回的结果，不对另一份chaincode做更新操作
*	除了：account链码之外;
*
 */
package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 购买一张表格
// 这个表格请求是发送到本链码的ledger中
func (s *SmartContract) BuyOneTable(ctx contractapi.TransactionContextInterface, id, tableID, requestTime string) error {
	err := s.Transfer(ctx, "Org1MSP", TABLE_FEE)
	if err != nil {
		return err
	}
	var buyTable BuyTable
	demander, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return err
	}

	fmt.Println(demander)
	buyTable.Demander = "buyer"
	// 前端进行处理
	buyTable.TableID = tableID
	buyTable.ID = id
	buyTable.RequestTime = requestTime
	buyTable.Complete = false

	buyTableJSON, err := json.Marshal(&buyTable)
	if err != nil {
		return err
	}
	//err = ctx.GetStub().PutPrivateData(requestCollection, request.ID, requestJSON)
	err = ctx.GetStub().PutState(buyTable.ID, buyTableJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *SmartContract) ReadTablePurchase(ctx contractapi.TransactionContextInterface) ([]*BuyTable, error) {
	requestIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer requestIterator.Close()

	var proposals []*BuyTable
	for requestIterator.HasNext() {
		buyTableJSON, err := requestIterator.Next()
		if err != nil {
			return nil, err
		}
		var temp BuyTable
		err = json.Unmarshal(buyTableJSON.Value, &temp)
		proposals = append(proposals, &temp)
	}
	return proposals, nil
}

// 查看购表请求,前端视图的，需要交互链码:table
func (s *SmartContract) ReadTablePurchaseForView(ctx contractapi.TransactionContextInterface) ([]*BuyTable, error) {
	requestIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer requestIterator.Close()

	var proposals []*BuyTable
	for requestIterator.HasNext() {
		buyTableJSON, err := requestIterator.Next()
		if err != nil {
			return nil, err
		}
		var temp BuyTable
		err = json.Unmarshal(buyTableJSON.Value, &temp)
		// 这里需要获取到表格的名称，所以要和table链码进行交互, []*Table
		tempArry := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("ReadMyTableByID"), []byte(temp.TableID)}, MYCHANNEL)
		if tempArry.Status != 200 {
			return nil, fmt.Errorf(string(tempArry.Payload))
		}
		var tempTable []*Table
		// tempArry.Payload []*Table
		err = json.Unmarshal([]byte(tempArry.Payload), &tempTable)
		if err != nil {
			return nil, err
		}
		temp.TableID = tempTable[0].TableName
		proposals = append(proposals, &temp)
	}
	return proposals, nil
}

// 处理所有请求
func (s *SmartContract) HandleAllTablePurchase(ctx contractapi.TransactionContextInterface) ([]*Table, error) {
	buyTables, err := s.ReadTablePurchase(ctx)
	if err != nil {
		return nil, err
	}

	var res []*Table
	for _, buy := range buyTables {
		if buy.Complete == true {
			return nil, fmt.Errorf("Already handled.")
		}
		temp, err := s.FetchTable(ctx, buy.TableID)
		if err != nil {
			return nil, err
		}
		s.Tablerequest = append(s.Tablerequest, *buy)
		res = append(res, temp)
	}

	return res, nil
}

// 处理单个请求
func (s *SmartContract) HandleSingleTablePurchase(ctx contractapi.TransactionContextInterface, buyTableID string) ([]*Table, error) {
	buyTableJSON, err := ctx.GetStub().GetState(buyTableID)
	if err != nil {
		return nil, err
	}
	if buyTableJSON == nil {
		return nil, fmt.Errorf("No such buyTable.")
	}

	var buyTable BuyTable
	err = json.Unmarshal(buyTableJSON, &buyTable)
	if err != nil {
		return nil, err
	}

	if buyTable.Complete == true {
		return nil, fmt.Errorf("Already handled.")
	}
	var res []*Table
	temp, err := s.FetchTable(ctx, buyTable.TableID)
	if err != nil {
		return nil, err
	}
	s.Tablerequest = append(s.Tablerequest, buyTable)
	res = append(res, temp)
	return res, nil
}

// endorsing: org2, 还需要标记它为已经处理了的
// 发送表格
func (s *SmartContract) SendTable(ctx contractapi.TransactionContextInterface, tableStr string) error {
	var table []*Table
	err := json.Unmarshal([]byte(tableStr), &table)
	if err != nil {
		return err
	}

	for _, t := range table {
		tableJSON, err := json.Marshal(t)
		if err != nil {
			return err
		}
		err = ctx.GetStub().PutPrivateData(tableCollection, t.TableID, tableJSON)
		if err != nil {
			return err
		}
	}

	for _, id := range s.Tablerequest {
		id.Complete = true
		idJSON, err := json.Marshal(id)
		if err != nil {
			return err
		}
		// 修改ledger中的状态
		err = ctx.GetStub().PutState(id.TableID, idJSON)
		if err != nil {
			return err
		}

	}
	// 处理完了请求之后将它清空
	s.Tablerequest = s.Tablerequest[0:0]
	return nil

}

// 从私有数据库中取表格
// FetchTable就相当于读取这个表格，并解密原来加密了的数据。
func (s *SmartContract) FetchTable(ctx contractapi.TransactionContextInterface, tableID string) (*Table, error) {
	// 从table链码中获取私有数据,response是响应结果
	response := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("FetchTable"), []byte(tableID)}, MYCHANNEL)

	// 表格数据在response.Payload中 *Table
	if response.Status != 200 {
		return nil, fmt.Errorf(string(response.Payload))
	}
	var table *Table
	err := json.Unmarshal(response.Payload, &table)
	if err != nil {
		return nil, err
	}
	return table, nil
}

func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, to string, amount float64) error {
	response := ctx.GetStub().InvokeChaincode(ACCOUNT, [][]byte{[]byte("Transfer"), []byte(to), []byte(strconv.FormatFloat(amount, 'f', 2, 64))}, MYCHANNEL)
	if response.Status != 200 {
		return fmt.Errorf(string(response.Payload))
	}
	return nil
}

// 增加一个获取最大ID的函数，然后将这个ID值返回给我们的服务器；
