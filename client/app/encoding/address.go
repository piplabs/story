package encoding

import sdk "github.com/cosmos/cosmos-sdk/types"

// Bech32HRP is the human-readable-part of the Bech32 address format.
const (
	Bech32HRP = "story"

	AccountAddressPrefix   = Bech32HRP
	AccountPubKeyPrefix    = Bech32HRP + sdk.PrefixPublic
	ValidatorAddressPrefix = Bech32HRP + sdk.PrefixValidator + sdk.PrefixOperator
	ValidatorPubKeyPrefix  = Bech32HRP + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	ConsNodeAddressPrefix  = Bech32HRP + sdk.PrefixValidator + sdk.PrefixConsensus
	ConsNodePubKeyPrefix   = Bech32HRP + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)
