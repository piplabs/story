package types

import "github.com/piplabs/story/contracts/bindings"

var (
	cdrContractABI                            = mustGetABI(bindings.CDRMetaData)
	CDRVaultReadEvent                         = mustGetEvent(cdrContractABI, "VaultRead")
	CDREncryptedPartialDecryptionSubmittedEvent = mustGetEvent(cdrContractABI, "EncryptedPartialDecryptionSubmitted")
)
