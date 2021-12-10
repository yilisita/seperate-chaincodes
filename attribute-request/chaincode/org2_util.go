package chaincode

/*
*
*   这里是组织2，即买家可以调用的函数
*
 */

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// 不能直接调用
func (s *SmartContract) createResponse(ctx contractapi.TransactionContextInterface, response Response) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutPrivateData(org2Collection, response.ID, responseJSON)
}

// ReadAsset returns the asset stored in the world state with given id.

func (s *SmartContract) GetAllResponses(ctx contractapi.TransactionContextInterface) ([]*Response, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetPrivateDataByRange(org2Collection, "", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var responses []*Response
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var temp Response
		err = json.Unmarshal(queryResponse.Value, &temp)
		if err != nil {
			return nil, err
		}
		fmt.Println(temp)
		responses = append(responses, &temp)

	}

	return responses, nil
}
