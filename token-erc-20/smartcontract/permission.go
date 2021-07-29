package smartcontract

import (
	"erc20/token-erc-20/entity/event"
	"erc20/token-erc-20/entity/permission"
	"fmt"
	"log"
	"strconv"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// Approve 用來准許
// 觸發Approval Event
func (s *SmartContract) Approve(ctx contractapi.TransactionContextInterface, input permission.ApproveInput) error {

	// 取得發出需求的client的ID
	owner, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return ErrFailedGetClientId(err)
	}

	// 建立許可鑰匙
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{owner, input.Spender})
	if err != nil {
		return ErrFailedCreateCompositeKey(fmt.Sprintf("%s: %v", allowancePrefix, err))
	}

	// 將許可的鑰匙與價值放進智能合約
	err = ctx.GetStub().PutState(allowanceKey, []byte(strconv.Itoa(input.Value)))
	if err != nil {
		return ErrFailedUpdateStateForKey(fmt.Sprintf("%s: %v", allowanceKey, err))
	}

	// 建立Approval Event
	err = event.SetEvent(ctx, "Approval", event.Event{From: owner, To: input.Spender, Value: input.Value})
	if err != nil {
		return err
	}

	log.Printf("client %s approved a withdrawal allowance of %d for spender %s", owner, input.Value, input.Spender)

	return nil
}

// Allowance回傳spender可以在owner帳戶所使用的餘額
func (s *SmartContract) Allowance(ctx contractapi.TransactionContextInterface, input permission.AllowanceInput) (int, error) {

	// 建立 allowanceKey
	allowanceKey, err := ctx.GetStub().CreateCompositeKey(allowancePrefix, []string{input.Owner, input.Spender})
	if err != nil {
		return 0, ErrFailedCreateCompositeKey(fmt.Sprintf("%s: %v", allowancePrefix, err))
	}

	// 讀取world state上的allowanceKey
	allowanceBytes, err := ctx.GetStub().GetState(allowanceKey)
	if err != nil {
		return 0, ErrFailedReadAllowance(fmt.Sprintf("%s: %v", allowanceKey, err))
	}

	var allowance int

	// 若無資料則初始化
	if allowanceBytes == nil {
		allowance = 0
	} else {
		allowance, _ = strconv.Atoi(string(allowanceBytes))
	}

	log.Printf("The allowance left for spender %s to withdraw from owner %s: %d", input.Spender, input.Owner, allowance)

	return allowance, nil
}
