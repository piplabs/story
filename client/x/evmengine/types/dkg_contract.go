package types

import (
	"github.com/piplabs/story/contracts/bindings"
)

var (
	dkgContractABI                            = mustGetABI(bindings.DKGMetaData)
	DKGRegisteredEvent                        = mustGetEvent(dkgContractABI, "Registered")
	DKGFinalizedEvent                         = mustGetEvent(dkgContractABI, "Finalized")
	DKGUpgradeScheduledEvent                  = mustGetEvent(dkgContractABI, "UpgradeScheduled")
	DKGRemoteAttestationProcessedOnChainEvent = mustGetEvent(dkgContractABI, "RemoteAttestationProcessedOnChain")
	DKGDealComplaintsSubmittedEvent           = mustGetEvent(dkgContractABI, "DealComplaintsSubmitted")
	DKGDealVerifiedEvent                      = mustGetEvent(dkgContractABI, "DealVerified")
	DKGInvalidDealEvent                       = mustGetEvent(dkgContractABI, "InvalidDeal")
	DKGThresholdDecryptRequestedEvent         = mustGetEvent(dkgContractABI, "ThresholdDecryptRequested")
	DKGMinReqRegisteredParticipantsSetEvent   = mustGetEvent(dkgContractABI, "MinReqRegisteredParticipantsSet")
	DKGMinReqFinalizedParticipantsSetEvent    = mustGetEvent(dkgContractABI, "MinReqFinalizedParticipantsSet")
	DKGOperationalThresholdSetEvent           = mustGetEvent(dkgContractABI, "OperationalThresholdSet")
)
