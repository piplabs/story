package types

import (
	"github.com/piplabs/story/contracts/bindings"
)

var (
	dkgContractABI                          = mustGetABI(bindings.DKGMetaData)
	DKGRegisteredEvent                      = mustGetEvent(dkgContractABI, "Registered")
	DKGFinalizedEvent                       = mustGetEvent(dkgContractABI, "Finalized")
	DKGUpgradeScheduledEvent                = mustGetEvent(dkgContractABI, "UpgradeScheduled")
	DKGUpgradeCancelledEvent                = mustGetEvent(dkgContractABI, "UpgradeCancelled")
	DKGMinReqRegisteredParticipantsSetEvent = mustGetEvent(dkgContractABI, "MinReqRegisteredParticipantsSet")
	DKGMinReqFinalizedParticipantsSetEvent  = mustGetEvent(dkgContractABI, "MinReqFinalizedParticipantsSet")
	DKGOperationalThresholdSetEvent         = mustGetEvent(dkgContractABI, "OperationalThresholdSet")
	DKGEnclaveTypeWhitelistedEvent          = mustGetEvent(dkgContractABI, "EnclaveTypeWhitelisted")
)
