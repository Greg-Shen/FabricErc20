package smartcontract

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const totalSupplyKey = "totalSupply"

const allowancePrefix = "allowance"

type SmartContract struct {
	contractapi.Contract
}

func (s *SmartContract) Init(ctx contractapi.TransactionContextInterface) error {
	log.Printf("Into Init")
	return nil
}
