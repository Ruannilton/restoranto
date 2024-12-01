package services

import (
	messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"
)

func CreateVerifyEmailContent(costumerInfo messagecontracts.VerifyEmailMessage, verificationLink string) string {
	if costumerInfo.NewCostumer {
		return newUserEmailVerificationTemplate(costumerInfo.CostumerName, verificationLink)
	}
	return alreadyExistingUserEmailVerificationTemplate(costumerInfo.CostumerName, verificationLink)
}

func newUserEmailVerificationTemplate(costumerName string, verificationLink string) string {
	return ""
}

func alreadyExistingUserEmailVerificationTemplate(costumerName string, verificationLink string) string {
	return ""
}
