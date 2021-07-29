package smartcontract

import (
	"erc20/token-erc-20/entity/event"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Mint 建立新的token並新增到Minter的account
// 觸發一個Transfer Event
func (s *SmartContract) Mint(ctx contractapi.TransactionContextInterface, amount int) error {

	// 確認minter的身份，這邊以Org1當做唯一可以製作token的身份

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return ErrFailedGetMSPID(err)
	}
	if clientMSPID != "Org1MSP" {
		return ErrClientIsNotAuthorizedToMint(nil)
	}

	// 確認為Org1後，確認製作token者的身份
	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return ErrFailedGetClientId(err)
	}

	// 若要製作的token價值小於0則無法製作token
	if amount <= 0 {
		return ErrMintAmountMustBePositive(nil)
	}

	// 取得minter的餘額
	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return ErrFailedReadMinterAccount(fmt.Sprintf("%s: %v", minter, err))
	}

	var currentBalance int

	// 若minter帳戶不存在，則初始化此minter帳戶
	if currentBalanceBytes == nil {
		currentBalance = 0
	} else {
		currentBalance, _ = strconv.Atoi(string(currentBalanceBytes))
	}

	// 將minter餘額新增token價值
	updatedBalance := currentBalance + amount
	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return err
	}

	// 將totalSupply更新
	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return ErrFailedRetrieveTotalSupply(err)
	}

	var totalSupply int

	// 若原本沒有totalSupply則初始化totalSupply
	if totalSupplyBytes == nil {
		totalSupply = 0
	} else {
		totalSupply, _ = strconv.Atoi(string(totalSupplyBytes))
	}

	// 將totalSupply加上amount
	totalSupply += amount
	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return err
	}

	// 建立Transfer Event
	err = event.SetEvent(ctx, "Transfer", event.Event{From: "0x0", To: minter, Value: amount})
	if err != nil {
		return err
	}

	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)

	return nil
}

// Burn刪減Minter balance
// 觸發Transfer Event
func (s *SmartContract) Burn(ctx contractapi.TransactionContextInterface, amount int) error {

	clientMSPID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return ErrFailedGetMSPID(err)
	}
	if clientMSPID != "Org1MSP" {
		return ErrClientIsNotAuthorizedToMint(nil)
	}

	minter, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return ErrFailedGetClientId(err)
	}

	if amount <= 0 {
		return ErrBurnAmountMustBePositive(nil)
	}

	currentBalanceBytes, err := ctx.GetStub().GetState(minter)
	if err != nil {
		return ErrFailedReadMinterAccount(fmt.Sprintf("%s: %v", minter, err))
	}

	var currentBalance int

	if currentBalanceBytes == nil {
		return ErrBalanceNotExist(nil)
	}

	currentBalance, _ = strconv.Atoi(string(currentBalanceBytes))

	updatedBalance := currentBalance - amount

	err = ctx.GetStub().PutState(minter, []byte(strconv.Itoa(updatedBalance)))
	if err != nil {
		return err
	}

	totalSupplyBytes, err := ctx.GetStub().GetState(totalSupplyKey)
	if err != nil {
		return ErrFailedRetrieveTotalSupply(err)
	}

	if totalSupplyBytes == nil {
		return ErrTotalSupplyNotExist(nil)
	}

	totalSupply, _ := strconv.Atoi(string(totalSupplyBytes))

	totalSupply -= amount
	err = ctx.GetStub().PutState(totalSupplyKey, []byte(strconv.Itoa(totalSupply)))
	if err != nil {
		return err
	}

	err = event.SetEvent(ctx, "Transfer", event.Event{From: minter, To: "0x0", Value: amount})
	if err != nil {
		return err
	}

	log.Printf("minter account %s balance updated from %d to %d", minter, currentBalance, updatedBalance)

	return nil
}
