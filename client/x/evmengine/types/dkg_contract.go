package types

import (
	"github.com/piplabs/story/contracts/bindings"
)

var (
	dkgContractABI                            = mustGetABI(bindings.DKGMetaData)
	DKGInitializedEvent                       = mustGetEvent(dkgContractABI, "DKGInitialized")
	DKGCommitmentsUpdatedEvent                = mustGetEvent(dkgContractABI, "DKGCommitmentsUpdated")
	DKGFinalizedEvent                         = mustGetEvent(dkgContractABI, "DKGFinalized")
	DKGUpgradeScheduledEvent                  = mustGetEvent(dkgContractABI, "UpgradeScheduled")
	DKGRegistrationChallengedEvent            = mustGetEvent(dkgContractABI, "RegistrationChallenged")
	DKGInvalidDKGInitializationEvent          = mustGetEvent(dkgContractABI, "InvalidDKGInitialization")
	DKGRemoteAttestationProcessedOnChainEvent = mustGetEvent(dkgContractABI, "RemoteAttestationProcessedOnChain")
	DKGDealComplaintsSubmittedEvent           = mustGetEvent(dkgContractABI, "DealComplaintsSubmitted")
	DKGDealVerifiedEvent                      = mustGetEvent(dkgContractABI, "DealVerified")
	DKGInvalidDealEvent                       = mustGetEvent(dkgContractABI, "InvalidDeal")
)
