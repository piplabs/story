package types

import "github.com/piplabs/story/contracts/bindings"

const (
	CDRFeeTypeAllocate      uint8 = 0
	CDRFeeTypeWrite         uint8 = 1
	CDRFeeTypeRead          uint8 = 2
	CDRFeeTypeSubmitPartial uint8 = 3
)

var (
	cdrContractABI                              = mustGetABI(bindings.CDRMetaData)
	CDRVaultReadEvent                           = mustGetEvent(cdrContractABI, "VaultRead")
	CDREncryptedPartialDecryptionSubmittedEvent = mustGetEvent(cdrContractABI, "EncryptedPartialDecryptionSubmitted")
	CDRFeeCollectedEvent                        = mustGetEvent(cdrContractABI, "FeeCollected")
)
