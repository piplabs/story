package types

import (
	"errors"
)

type ErrCode uint32

const (
	Unspecified              ErrCode = 0
	InvalidUncmpPubKey       ErrCode = 1
	ValidatorNotFound        ErrCode = 2
	ValidatorAlreadyExists   ErrCode = 3
	InvalidTokenType         ErrCode = 4
	InvalidPeriodType        ErrCode = 5
	InvalidOperator          ErrCode = 6
	InvalidCommissionRate    ErrCode = 7
	InvalidMinSelfDelegation ErrCode = 8
	InvalidDelegationAmount  ErrCode = 9
	PeriodDelegationNotFound ErrCode = 10
)

var (
	ErrUnspecified              = errors.New("unspecified")
	ErrInvalidUncmpPubKey       = errors.New("invalid_uncompressed_pubkey")
	ErrValidatorNotFound        = errors.New("validator_not_found")
	ErrValidatorAlreadyExists   = errors.New("validator_already_exists")
	ErrInvalidTokenType         = errors.New("invalid_token_type")
	ErrInvalidPeriodType        = errors.New("invalid_period_type")
	ErrInvalidOperator          = errors.New("invalid_operator")
	ErrInvalidCommissionRate    = errors.New("invalid_commission_rate")
	ErrInvalidMinSelfDelegation = errors.New("invalid_min_self_delegation")
	ErrInvalidDelegationAmount  = errors.New("invalid_delegation_amount")
	ErrPeriodDelegationNotFound = errors.New("period_delegation_not_found")
)

var codeToErr = map[ErrCode]error{
	Unspecified:              ErrUnspecified,
	InvalidUncmpPubKey:       ErrInvalidUncmpPubKey,
	ValidatorNotFound:        ErrValidatorNotFound,
	ValidatorAlreadyExists:   ErrValidatorAlreadyExists,
	InvalidTokenType:         ErrInvalidTokenType,
	InvalidPeriodType:        ErrInvalidPeriodType,
	InvalidOperator:          ErrInvalidOperator,
	InvalidCommissionRate:    ErrInvalidCommissionRate,
	InvalidMinSelfDelegation: ErrInvalidMinSelfDelegation,
	InvalidDelegationAmount:  ErrInvalidDelegationAmount,
	PeriodDelegationNotFound: ErrPeriodDelegationNotFound,
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

	return errors.Join(codeToErr[code], err)
}

func UnwrapErrCode(err error) ErrCode {
	switch {
	case errors.Is(err, ErrUnspecified):
		return Unspecified
	case errors.Is(err, ErrInvalidUncmpPubKey):
		return InvalidUncmpPubKey
	case errors.Is(err, ErrValidatorNotFound):
		return ValidatorNotFound
	case errors.Is(err, ErrValidatorAlreadyExists):
		return ValidatorAlreadyExists
	case errors.Is(err, ErrInvalidTokenType):
		return InvalidTokenType
	case errors.Is(err, ErrInvalidPeriodType):
		return InvalidPeriodType
	case errors.Is(err, ErrInvalidOperator):
		return InvalidOperator
	case errors.Is(err, ErrInvalidCommissionRate):
		return InvalidCommissionRate
	case errors.Is(err, ErrInvalidMinSelfDelegation):
		return InvalidMinSelfDelegation
	case errors.Is(err, ErrInvalidDelegationAmount):
		return InvalidDelegationAmount
	case errors.Is(err, ErrPeriodDelegationNotFound):
		return PeriodDelegationNotFound
	default:
		return Unspecified
	}
}
