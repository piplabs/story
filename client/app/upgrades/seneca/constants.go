package seneca

import "github.com/piplabs/story/lib/netconf"

const UpgradeName = netconf.Seneca

// NewMaxValidators is the target active validator set size after the seneca
// upgrade. See SIP-00011 for the rationale of reducing from 80 to 21.
var NewMaxValidators uint32 = 21
