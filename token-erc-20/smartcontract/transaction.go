package smartcontract

import (
	"erc20/token-erc-20/entity/event"
	"erc20/token-erc-20/entity/transaction"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Transfer 把token從傳送者傳到收件者
// 觸發Transfer Event
func (s *SmartContract) Transfer(ctx contractapi.TransactionContextInterface, input transaction.TransferInput) error {

	// 取得傳送者ID
	clientID, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return ErrFailedGetClientId(err)
	}

	err = transaction.TransferBasic(ctx, clientID, input.Recipient, input.Amount)
	if err != nil {
		return ErrFailedTransfer(err)
	}

	// 建立 Transfer Event
	err = event.SetEvent(ctx, "Transfer", event.Event{From: clientID, To: input.Recipient, Value: input.Amount})
	if err != nil {
		return err
	}

	return nil
}

// TransferFrom 讓人把amount從別的傳送位址傳到收件位址
// 觸發Transfer Event
func (s *SmartContract) TransferFrom(ctx contractapi.TransactionContextInterface, input transaction.TransferFromInput) error {

	// 取得spenderID
	spender, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return ErrFailedGetClientId(err)
	}

	// CreateCompositeKey為chaincode interface中function。
	// 生成的鑰匙可直接被PutState()使用.
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{input.From, spender})
	if err != nil {
		return ErrFailedCreateCompositeKey(fmt.Sprintf("%s: %v", allowancePrefix, err))
	}

	// 取得此鑰匙所儲存的價值
	currentAllowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return ErrFailedReadAllowance(fmt.Sprintf("%s: %v", allowanceKey, err))
	}

	var currentAllowance int
	currentAllowance, _ = strconv.Atoi(string(currentAllowanceBytes))

	// 若許可值小於要傳送值
	if currentAllowance < input.Value {
		return ErrSpenderNotHaveEnoughAllowance(nil)
	}

	// 使用TransferBasic進行transfer
	err = transaction.TransferBasic(ctx, input.From, input.To, input.Value)
	if err != nil {
		return ErrFailedTransfer(err)
	}

	// 刪減spender的Allowance
	updatedAllowance := currentAllowance - input.Value
	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(updatedAllowance)))
	if err != nil {
		return err
	}

	// 建立Transfer Event
	err = event.SetEvent(ctx, "Transfer", event.Event{From: input.From, To: input.To, Value: input.Value})
	if err != nil {
		return err
	}

	log.Printf("spender %s allowance updated from %d to %d", spender, currentAllowance, updatedAllowance)

	return nil
}
