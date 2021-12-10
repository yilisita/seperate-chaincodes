package chaincode

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 从私有数据库中取表格
func (s *SmartContract) FetchTable(ctx contractapi.TransactionContextInterface, tableID string) (*Table, error) {
	tableJSON, err := ctx.GetStub().GetPrivateData(myCollection, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if tableJSON == nil {
		return nil, fmt.Errorf("the table %s does not exist", tableID)
	}
	var table *Table
	// 注意是否有&
	err = json.Unmarshal(tableJSON, &table)
	if err != nil {
		return nil, err
	}

	// 用来存储加密了的字段的ID
	var index []int
	for i, v := range table.Flags {
		if v == 1 {
			index = append(index, i)
		}
	}

	// 解密表格
	for _, asset := range table.Assets {
		value := reflect.ValueOf(&asset).Elem()
		for _, i := range index {
			temp := decrypt(value.FieldByName(ATTRIBUTES[i]).Int())
			value.FieldByName(ATTRIBUTES[i]).Set(reflect.ValueOf(temp))
		}
	}

	return table, nil

}
