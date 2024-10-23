package server

type DelegatorBaseInfo struct {
	DelegatorAddr   string `json:"delegator_addr"`
	WithdrawAddress string `json:"withdraw_address"`
	RewardAddress   string `json:"reward_address"`
	OperatorAddress string `json:"operator_address"`
}
