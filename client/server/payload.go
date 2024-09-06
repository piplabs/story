package server

import (
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type pagination struct {
	Key        string `mapstructure:"key"`
	Offset     uint64 `mapstructure:"offset"`
	Limit      uint64 `mapstructure:"limit"`
	CountTotal bool   `mapstructure:"count_total"`
	Reverse    bool   `mapstructure:"reverse"`
}

type getSigningInfosRequest struct {
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

type getAccountsRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getAccountAddressByIDRequest struct {
	AccountID uint64 `mapstructure:"account_id"`
}

type getSupplyRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getSupplyByDenomRequest struct {
	Denom string `mapstructure:"denom"`
}

type getBalancesByAddressRequest struct {
	ResolveDenom bool       `mapstructure:"resolve_denom"`
	Pagination   pagination `mapstructure:"pagination"`
}

type getBalancesByAddressDenomRequest struct {
	Denom string `mapstructure:"denom"`
}

type getDenomOwnersRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getDenomOwnersByQueryRequest struct {
	Denom      string     `mapstructure:"denom"`
	Pagination pagination `mapstructure:"pagination"`
}

type getDenomsMetadataRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getDenomMetadataByQueryStringRequest struct {
	Denom string `mapstructure:"denom"`
}

type getSendEnabledRequest struct {
	Denoms     []string   `mapstructure:"denoms"`
	Pagination pagination `mapstructure:"pagination"`
}
type getSpendableBalancesRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getSpendableBalanceByDenomRequest struct {
	Denom string `mapstructure:"denom"`
}

type getValidatorSlashesByValidatorAddressRequest struct {
	StartingHeight uint64     `mapstructure:"starting_height"`
	EndingHeight   uint64     `mapstructure:"ending_height"`
	Pagination     pagination `mapstructure:"pagination"`
}

type getComebftBlockEventsRequest struct {
	From            int64    `mapstructure:"from"`
	To              int64    `mapstructure:"to"`
	EventTypeFilter []string `mapstructure:"event_type_filter"`
}

type getComebftBlockEventsBlockResults struct {
	Height              int64        `json:"height"`
	FinalizeBlockEvents []abci.Event `json:"finalize_block_events"`
}

type getAllValidatorOutstandingRewardsRequest struct {
	From int64 `mapstructure:"from"`
	To   int64 `mapstructure:"to"`
}

type getAllValidatorOutstandingRewardsRequestBlockResults struct {
	Height     int64                   `json:"height"`
	Validators map[string]sdk.DecCoins `json:"validators"`
}

type getWithdrawalQueueRequest struct {
	Pagination pagination `mapstructure:"pagination"`
}

type getModuleVersionsRequest struct {
	ModuleName string `mapstructure:"module_name"`
}
