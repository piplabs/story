package types

import (
	"github.com/piplabs/story/contracts/bindings"
)

var (
	dkgContractABI                            = mustGetABI(bindings.DKGMetaData)
	DKGInitializedEvent                       = mustGetEvent(dkgContractABI, "DKGInitialized")
	DKGNetworkSetEvent                        = mustGetEvent(dkgContractABI, "DKGNetworkSet")
	DKGFinalizedEvent                         = mustGetEvent(dkgContractABI, "DKGFinalized")
	DKGUpgradeScheduledEvent                  = mustGetEvent(dkgContractABI, "UpgradeScheduled")
	DKGRemoteAttestationProcessedOnChainEvent = mustGetEvent(dkgContractABI, "RemoteAttestationProcessedOnChain")
	DKGDealComplaintsSubmittedEvent           = mustGetEvent(dkgContractABI, "DealComplaintsSubmitted")
	DKGDealVerifiedEvent                      = mustGetEvent(dkgContractABI, "DealVerified")
	DKGInvalidDealEvent                       = mustGetEvent(dkgContractABI, "InvalidDeal")
	DKGThresholdDecryptRequestedEvent         = mustGetEvent(dkgContractABI, "ThresholdDecryptRequested")
)
