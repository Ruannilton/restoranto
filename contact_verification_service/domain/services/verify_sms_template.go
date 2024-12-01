package services

import messagecontracts "github.com/Ruannilton/go-msg-contracts/pkg/message_contracts"

func CreateVerifySMSContent(costumerInfo messagecontracts.VerifySMSMessage, verificationCode string) string {
	if costumerInfo.NewCostumer {
		return newUserSMSVerificationTemplate(costumerInfo.CostumerName, verificationCode)
	}
	return alreadyExistingUserSMSVerificationTemplate(costumerInfo.CostumerName, verificationCode)
}

func newUserSMSVerificationTemplate(costumerName string, verificationCode string) string {
	return ""
}

func alreadyExistingUserSMSVerificationTemplate(costumerName string, verificationCode string) string {
	return ""
}
