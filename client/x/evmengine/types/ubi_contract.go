package types

import "github.com/piplabs/story/contracts/bindings"

var (
	ubiPoolABI            = mustGetABI(bindings.UBIPoolMetaData)
	UBIPercentageSetEvent = mustGetEvent(ubiPoolABI, "UBIPercentageSet")
)
