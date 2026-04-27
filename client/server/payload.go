package server

import (
	"encoding/hex"
	stdmath "math"
	"strings"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/piplabs/story/lib/errors"
)

// requesterPubKeyByteLen is the length of an uncompressed secp256k1 public key.
const requesterPubKeyByteLen = 65

// maxCDRUuid is the maximum vault uuid the CDR contract can allocate. The
// contract enforces `require($.uuid < type(uint32).max)` before assigning a
// new uuid, so valid uuids fall in [0, math.MaxUint32-1].
const maxCDRUuid = stdmath.MaxUint32 - 1

type pagination struct {
	Key        string `mapstructure:"key"`
	Offset     uint64 `mapstructure:"offset"`
	Limit      uint64 `mapstructure:"limit"`
	CountTotal bool   `mapstructure:"count_total"`
	Reverse    bool   `mapstructure:"reverse"`
}

type getSupplyByDenomRequest struct {
	Denom string `mapstructure:"denom"`
}

type getBalancesByAddressDenomRequest struct {
	Denom string `mapstructure:"denom"`
}

type getSpendableBalancesByAddressDenomRequest struct {
	Denom string `mapstructure:"denom"`
}

type getValidatorSlashesByValidatorAddressRequest struct {
	StartingHeight uint64     `mapstructure:"starting_height"`
	EndingHeight   uint64     `mapstructure:"ending_height"`
	Pagination     pagination `mapstructure:"pagination"`
}

type getWithdrawalQueueRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getRewardWithdrawalQueueRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getValidatorsRequest struct {
	Status     string     `mapstructure:"status"`
	Pagination pagination `mapstructure:"pagination"`
}

type getValidatorDelegationsByValidatorAddressRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getValidatorUnbondingDelegationsRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getDelegationsByDelegatorAddressRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getRedelegationsByDelegatorAddressRequest struct {
	SrcValidatorAddr string     `mapstructure:"src_validator_addr"`
	DstValidatorAddr string     `mapstructure:"dst_validator_addr"`
	Pagination       pagination `mapstructure:"pagination"`
}

type getUnbondingDelegationsByDelegatorAddressRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getValidatorsByDelegatorAddressRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getPeriodDelegationsRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getModuleVersionsRequest struct {
	ModuleName string `mapstructure:"module_name"`
}

type getStakedTokenByDelegatorAddressRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getRewardsTokenByDelegatorAddressRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type QueryTotalDelegationsCountResponse struct {
	Total int `json:"total"`
}

type QueryTotalStakedTokenResponse struct {
	TotalStakedToken math.Int `json:"total_staked_token"`
}

type DelegationStakedToken struct {
	ValidatorOperatorAddress string         `json:"validator_operator_address"`
	StakedToken              math.LegacyDec `json:"staked_token"`
}

type QueryStakedTokenByDelegatorAddressResponse struct {
	DelegationStakedToken []DelegationStakedToken `json:"delegation_staked_token"`
	Pagination            *query.PageResponse     `json:"pagination"`
}

type QueryTotalStakedTokenByDelegatorAddressResponse struct {
	StakedToken math.LegacyDec `json:"staked_token"`
}

type DelegationRewardsToken struct {
	ValidatorOperatorAddress string         `json:"validator_operator_address"`
	RewardsToken             math.LegacyDec `json:"rewards_token"`
}

type QueryRewardsTokenByDelegatorAddressResponse struct {
	DelegationRewardsToken []DelegationRewardsToken `json:"delegation_rewards_token"`
	Pagination             *query.PageResponse      `json:"pagination"`
}

type QueryTotalRewardsTokenByDelegatorAddressResponse struct {
	RewardsToken math.LegacyDec `json:"rewards_token"`
}

type getDKGNetworkRequest struct {
	Round uint32 `mapstructure:"round"`
	// Deprecated: accepted but ignored — the keeper scopes results only by round.
	CodeCommitmentHex string `mapstructure:"code_commitment_hex"`
}

type getAllDKGRegistrationsRequest struct {
	Round uint32 `mapstructure:"round"`
}

type getVerifiedDKGRegistrationsRequest struct {
	Round uint32 `mapstructure:"round"`
	// Deprecated: accepted but ignored — the keeper scopes results only by round.
	CodeCommitmentHex string `mapstructure:"code_commitment_hex"`
}

type getCDRPartialsRequest struct {
	Uuid               uint32 `mapstructure:"uuid"`
	RequesterPubKeyHex string `mapstructure:"requester_pub_key_hex"`
}

func (r *getCDRPartialsRequest) validate() error {
	// Accept both `0x`-prefixed and bare hex; normalize the field in place so
	// the keeper sees the canonical (no-prefix) form.
	if len(r.RequesterPubKeyHex) >= 2 && strings.EqualFold(r.RequesterPubKeyHex[:2], "0x") {
		r.RequesterPubKeyHex = r.RequesterPubKeyHex[2:]
	}
	if r.RequesterPubKeyHex == "" {
		return errors.New("requester_pub_key_hex is required")
	}
	pubKey, err := hex.DecodeString(r.RequesterPubKeyHex)
	if err != nil {
		return errors.Wrap(err, "requester_pub_key_hex must be valid hex")
	}
	if len(pubKey) != requesterPubKeyByteLen {
		return errors.New("requester_pub_key_hex must decode to an uncompressed secp256k1 public key (65 bytes)")
	}
	if r.Uuid > maxCDRUuid {
		return errors.New("uuid out of range")
	}

	return nil
}

type QueryDKGGlobalPublicKeyResponse struct {
	PublicKeyHex string `json:"public_key"`
}
