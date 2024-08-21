package types

import (
	"github.com/piplabs/story/contracts/bindings"
)

var (
	ipTokenSlashingABI = mustGetABI(bindings.IPTokenSlashingMetaData)
	UnjailEvent        = mustGetEvent(ipTokenSlashingABI, "Unjail")
)
