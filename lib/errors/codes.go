package errors

import (
	stderrors "errors"
)

type ErrCode uint32

const (
	Unspecified              ErrCode = 0
	UnexpectedCondition      ErrCode = 1
	InvalidCmpPubKey         ErrCode = 2
	ValidatorNotFound        ErrCode = 3
	ValidatorAlreadyExists   ErrCode = 4
	InvalidTokenType         ErrCode = 5
	InvalidOperator          ErrCode = 6
	InvalidCommissionRate    ErrCode = 7
	InvalidMinSelfDelegation ErrCode = 8
	InvalidPeriodType        ErrCode = 9
	InvalidDelegationAmount  ErrCode = 10
	DelegationNotFound       ErrCode = 11
	PeriodDelegationNotFound ErrCode = 12
	InvalidRequest           ErrCode = 13
	SelfRedelegation         ErrCode = 14
	TokenTypeMismatch        ErrCode = 15
	MissingSelfDelegation    ErrCode = 16
	ValidatorNotJailed       ErrCode = 17
	ValidatorStillJailed     ErrCode = 18
	PendingUpgradeExists     ErrCode = 19
)

var (
	ErrUnspecified              = stderrors.New("unspecified")
	ErrUnexpectedCondition      = stderrors.New("unexpected_condition")
	ErrInvalidCmpPubKey         = stderrors.New("invalid_compressed_pubkey")
	ErrValidatorNotFound        = stderrors.New("validator_not_found")
	ErrValidatorAlreadyExists   = stderrors.New("validator_already_exists")
	ErrInvalidTokenType         = stderrors.New("invalid_token_type")
	ErrInvalidOperator          = stderrors.New("invalid_operator")
	ErrInvalidCommissionRate    = stderrors.New("invalid_commission_rate")
	ErrInvalidMinSelfDelegation = stderrors.New("invalid_min_self_delegation")
	ErrInvalidPeriodType        = stderrors.New("invalid_period_type")
	ErrInvalidDelegationAmount  = stderrors.New("invalid_delegation_amount")
	ErrDelegationNotFound       = stderrors.New("delegation_not_found")
	ErrPeriodDelegationNotFound = stderrors.New("period_delegation_not_found")
	ErrInvalidRequest           = stderrors.New("invalid_request")
	ErrSelfRedelegation         = stderrors.New("self_redelegation")
	ErrTokenTypeMismatch        = stderrors.New("token_type_mismatch")
	ErrMissingSelfDelegation    = stderrors.New("missing_self_delegation")
	ErrValidatorNotJailed       = stderrors.New("validator_not_jailed")
	ErrValidatorStillJailed     = stderrors.New("validator_still_jailed")
	ErrPendingUpgradeExists     = stderrors.New("pending_upgrade_exists")
)

var codeToErr = map[ErrCode]error{
	Unspecified:              ErrUnspecified,
	UnexpectedCondition:      ErrUnexpectedCondition,
	InvalidCmpPubKey:         ErrInvalidCmpPubKey,
	ValidatorNotFound:        ErrValidatorNotFound,
	ValidatorAlreadyExists:   ErrValidatorAlreadyExists,
	InvalidTokenType:         ErrInvalidTokenType,
	InvalidOperator:          ErrInvalidOperator,
	InvalidCommissionRate:    ErrInvalidCommissionRate,
	InvalidMinSelfDelegation: ErrInvalidMinSelfDelegation,
	InvalidPeriodType:        ErrInvalidPeriodType,
	InvalidDelegationAmount:  ErrInvalidDelegationAmount,
	DelegationNotFound:       ErrDelegationNotFound,
	PeriodDelegationNotFound: ErrPeriodDelegationNotFound,
	InvalidRequest:           ErrInvalidRequest,
	SelfRedelegation:         ErrSelfRedelegation,
	TokenTypeMismatch:        ErrTokenTypeMismatch,
	MissingSelfDelegation:    ErrMissingSelfDelegation,
	ValidatorNotJailed:       ErrValidatorNotJailed,
	ValidatorStillJailed:     ErrValidatorStillJailed,
	PendingUpgradeExists:     ErrPendingUpgradeExists,
}

func (c ErrCode) String() string {
	if _, ok := codeToErr[c]; !ok {
		return ErrUnspecified.Error()
	}

	return codeToErr[c].Error()
}

//nolint:wrapcheck // we are wrapping the error with the code
func WrapErrWithCode(code ErrCode, err error) error {
	if _, ok := codeToErr[code]; !ok {
		code = Unspecified
	}

	return stderrors.Join(codeToErr[code], err)
}

func UnwrapErrCode(err error) ErrCode {
	switch {
	case stderrors.Is(err, ErrUnspecified):
		return Unspecified
	case stderrors.Is(err, ErrUnexpectedCondition):
		return UnexpectedCondition
	case stderrors.Is(err, ErrInvalidCmpPubKey):
		return InvalidCmpPubKey
	case stderrors.Is(err, ErrValidatorNotFound):
		return ValidatorNotFound
	case stderrors.Is(err, ErrValidatorAlreadyExists):
		return ValidatorAlreadyExists
	case stderrors.Is(err, ErrInvalidTokenType):
		return InvalidTokenType
	case stderrors.Is(err, ErrInvalidOperator):
		return InvalidOperator
	case stderrors.Is(err, ErrInvalidCommissionRate):
		return InvalidCommissionRate
	case stderrors.Is(err, ErrInvalidMinSelfDelegation):
		return InvalidMinSelfDelegation
	case stderrors.Is(err, ErrInvalidPeriodType):
		return InvalidPeriodType
	case stderrors.Is(err, ErrInvalidDelegationAmount):
		return InvalidDelegationAmount
	case stderrors.Is(err, ErrDelegationNotFound):
		return DelegationNotFound
	case stderrors.Is(err, ErrPeriodDelegationNotFound):
		return PeriodDelegationNotFound
	case stderrors.Is(err, ErrInvalidRequest):
		return InvalidRequest
	case stderrors.Is(err, ErrSelfRedelegation):
		return SelfRedelegation
	case stderrors.Is(err, ErrTokenTypeMismatch):
		return TokenTypeMismatch
	case stderrors.Is(err, ErrMissingSelfDelegation):
		return MissingSelfDelegation
	case stderrors.Is(err, ErrValidatorNotJailed):
		return ValidatorNotJailed
	case stderrors.Is(err, ErrValidatorStillJailed):
		return ValidatorStillJailed
	case stderrors.Is(err, ErrPendingUpgradeExists):
		return PendingUpgradeExists
	default:
		return Unspecified
	}
}
