package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (s *SmartContract) ReadPurchasedTable(ctx contractapi.TransactionContextInterface) ([]*Table, error) {
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(tableCollection, "", "")
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	defer resultsIterator.Close()

	var tables []*Table
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var temp Table
		err = json.Unmarshal(queryResponse.Value, &temp)
		if err != nil {
			return nil, err
		}

		tables = append(tables, &temp)

	}

	return tables, nil
}
