package smartcontract

import (
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// ClientAccountID回傳有提出要求的client的ID
func (s *SmartContract) ClientAccountID(ctx contractapi.TransactionContextInterface) (string, error) {

	// 取得要查詢的client的ID
	clientAccountID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return "", ErrFailedGetClientId(err)
	}

	return clientAccountID, nil
}

// BalanceOf 回傳所輸入帳戶名稱的balance
func (s *SmartContract) BalanceOf(ctx contractapi.TransactionContextInterface, account string) (int, error) {
	balanceBytes, err := ctx.GetStub().GetState(account)
	if err != nil {
		return 0, ErrFailedReadWorldState(err)
	}
	if balanceBytes == nil {
		return 0, ErrAccountNotExist(account)
	}

	balance, _ := strconv.Atoi(string(balanceBytes))

	return balance, nil
}

// ClientAccountBalance 回傳有提出要求的client的balance
func (s *SmartContract) ClientAccountBalance(ctx contractapi.TransactionContextInterface) (int, error) {

	// 取得client的ID
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return 0, ErrFailedGetClientId(err)
	}

	balanceBytes, err := ctx.GetStub().GetState(clientID)
	if err != nil {
		return 0, ErrFailedReadWorldState(err)
	}
	if balanceBytes == nil {
		return 0, ErrAccountNotExist(clientID)
	}

	balance, _ := strconv.Atoi(string(balanceBytes))

	return balance, nil
}

// TotalSupply 回傳Token Supply
func (s *SmartContract) TotalSupply(ctx contractapi.TransactionContextInterface) (int, error) {

	// 取得total supply
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return 0, ErrFailedRetrieveTotalSupply(err)
	}

	var totalSupply int

	// 若無則回傳0
	if totalSupplyBytes == nil {
		totalSupply = 0
	} else {
		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes))
	}

	log.Printf("TotalSupply: %d tokens", totalSupply)

	return totalSupply, nil
}
