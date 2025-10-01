package server

import (
	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/types/query"
)

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
