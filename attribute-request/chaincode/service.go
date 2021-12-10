package chaincode

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

/*******************************************************************
**																  **
**                   TABLE CHAINCODE                              **
********************************************************************/
// 获取本月发电量总额 get total genAmountThisMonth
func (s *SmartContract) GetAttributeTotal(ctx contractapi.TransactionContextInterface, tableID string, attributeID int) (float64, error) {
	// 需要invoke的时候，将int类型的参数转化为string类型的，然后[]byte强制转换
	response := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("GetAttributeTotal"), []byte(tableID), []byte(strconv.FormatInt(int64(attributeID), 10))}, MYCHANNEL)

	// 表格数据在response.Payload中 *Table
	if response.Status != 200 {
		return -1, fmt.Errorf(string(response.Payload))
	}
	var res float64
	res, err := strconv.ParseFloat(string(response.Payload), 64)
	if err != nil {
		return -1, err
	}
	return res, nil
}

// 获取本月发电量的平均值
func (s *SmartContract) GetAttributeAve(ctx contractapi.TransactionContextInterface, tableName string, attributeID int) (float64, error) {
	// 需要invoke的时候，将int类型的参数转化为string类型的，然后[]byte强制转换
	response := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("GetAttributeAve"), []byte(tableName), []byte(strconv.FormatInt(int64(attributeID), 10))}, MYCHANNEL)

	// 表格数据在response.Payload中 *Table
	if response.Status != 200 {
		return -1, fmt.Errorf(string(response.Payload))
	}
	var res float64
	res, err := strconv.ParseFloat(string(response.Payload), 64)
	if err != nil {
		return -1, err
	}
	return res, nil
}

func (s *SmartContract) GetPearson(ctx contractapi.TransactionContextInterface, tableID string, attributes []int) (float64, error) {
	// 需要invoke的时候，将int类型的参数转化为string类型的，然后[]byte强制转换
	// 数组可以直接转化为[]byte，利用json.Marshal
	attributesBYTES, err := json.Marshal(attributes)
	if err != nil {
		return -2, err
	}
	response := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("GetPearson"), []byte(tableID), (attributesBYTES)}, MYCHANNEL)

	// 表格数据在response.Payload中 *Table
	if response.Status != 200 {
		return -1, fmt.Errorf(string(response.Payload))
	}
	var res float64
	res, err = strconv.ParseFloat(string(response.Payload), 64)
	if err != nil {
		return -1, err
	}
	return res, nil

}

func (s *SmartContract) ReadMyTableByID(ctx contractapi.TransactionContextInterface, tableID string) ([]*Table, error) {
	response := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("ReadMyTableByID"), []byte(tableID)}, MYCHANNEL)

	// 表格数据在response.Payload中 *Table
	if response.Status != 200 {
		return nil, fmt.Errorf(string(response.Payload))
	}
	var res []*Table
	err := json.Unmarshal((response.Payload), &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *SmartContract) MyTableExists(ctx contractapi.TransactionContextInterface, tableID string) (bool, error) {
	response := ctx.GetStub().InvokeChaincode(TABLE, [][]byte{[]byte("MyTableExists"), []byte(tableID)}, MYCHANNEL)
	if response.Status != 200 {
		return false, fmt.Errorf(string(response.Payload))
	}
	var res bool
	res, err := strconv.ParseBool(string((response.Payload)))
	if err != nil {
		return false, err
	}
	return res, nil
}

/*******************************************************************
**																  **
**                   ACCOUNT CHAINCODE                            **
********************************************************************/
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, to string, amount float64) error {
	response := ctx.GetStub().InvokeChaincode(ACCOUNT, [][]byte{[]byte("Transfer"), []byte(to), []byte(strconv.FormatFloat(amount, 'f', 2, 64))}, MYCHANNEL)
	fmt.Println(response)
	if response.Status != 200 {
		return fmt.Errorf(string(response.Payload))
	}
	return nil
}
