package types

import (
	"github.com/storyprotocol/iliad/contracts/bindings"
)

var (
	ipTokenSlashingABI = mustGetABI(bindings.IPTokenSlashingMetaData)
	UnjailEvent        = mustGetEvent(ipTokenSlashingABI, "Unjail")
)
