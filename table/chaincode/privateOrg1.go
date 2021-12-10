package chaincode

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"table/paillier"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset

// 用来加密数据
func encrypt(n int64) int64 {
	return paillier.Encryption(p, big.NewInt(n)).Int64()
}

// 解密数据
func decrypt(n int64) int64 {
	return paillier.Decryption(p, big.NewInt(n)).Int64()
}

// 选择字段进行加密, 可以进行多个字段的加密
func encryptAsset(asset Asset, attributes []int) Asset {
	var pointer = &asset
	value := reflect.ValueOf(pointer)
	for _, v := range attributes {
		var ciphertext = encrypt(value.Elem().FieldByName(ATTRIBUTES[v]).Int())
		value.Elem().FieldByName(ATTRIBUTES[v]).Set(reflect.ValueOf(ciphertext))
	}
	return *pointer
}

// 新建一个表格
func (s *SmartContract) CreateMyTable(ctx contractapi.TransactionContextInterface, tableStr string) error {
	var table Table

	err := json.Unmarshal([]byte(tableStr), &table)
	if err != nil {
		return err
	}
	isExist, err := s.MyTableExists(ctx, table.TableID)
	if err != nil {
		return err
	}

	if isExist {
		return fmt.Errorf("Table %v already exists", table.TableID)
	}

	// 获取要加密的属性值
	var attributes []int
	for k, v := range table.Flags {
		if v == 1 {
			attributes = append(attributes, k)
		}
	}

	for i, v := range table.Assets {
		table.Assets[i] = encryptAsset(v, attributes)
	}
	assetJSON, err := json.Marshal(table)
	if err != nil {
		return err
	}
	return ctx.GetStub().PutPrivateData(myCollection, table.TableID, assetJSON)

}

// 通过表格的ID来读取数据
func (s *SmartContract) ReadMyTableByID(ctx contractapi.TransactionContextInterface, tableID string) ([]*Table, error) {
	tableJSON, err := ctx.GetStub().GetPrivateData(myCollection, tableID)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if tableJSON == nil {
		return nil, fmt.Errorf("the table %s does not exist", tableID)
	}
	var res []*Table
	var table Table
	err = json.Unmarshal(tableJSON, &table)
	if err != nil {
		return nil, err
	}
	res = append(res, &table)
	return res, nil
}

// 检查表格是否存在
func (s *SmartContract) MyTableExists(ctx contractapi.TransactionContextInterface, tableID string) (bool, error) {
	tableJSON, err := ctx.GetStub().GetPrivateData(myCollection, tableID)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return tableJSON != nil, nil
}

// 用来获取ledger中一张指定的表中有多少记录
