package errors

import (
	stderrors "errors"
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
	InvalidRequest           ErrCode = 11
	UnexpectedCondition      ErrCode = 12
)

var (
	ErrUnspecified              = stderrors.New("unspecified")
	ErrInvalidUncmpPubKey       = stderrors.New("invalid_uncompressed_pubkey")
	ErrValidatorNotFound        = stderrors.New("validator_not_found")
	ErrValidatorAlreadyExists   = stderrors.New("validator_already_exists")
	ErrInvalidTokenType         = stderrors.New("invalid_token_type")
	ErrInvalidPeriodType        = stderrors.New("invalid_period_type")
	ErrInvalidOperator          = stderrors.New("invalid_operator")
	ErrInvalidCommissionRate    = stderrors.New("invalid_commission_rate")
	ErrInvalidMinSelfDelegation = stderrors.New("invalid_min_self_delegation")
	ErrInvalidDelegationAmount  = stderrors.New("invalid_delegation_amount")
	ErrPeriodDelegationNotFound = stderrors.New("period_delegation_not_found")
	ErrInvalidRequest           = stderrors.New("invalid_request")
	ErrUnexpectedCondition      = stderrors.New("unexpected_condition")
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
	InvalidRequest:           ErrInvalidRequest,
	UnexpectedCondition:      ErrUnexpectedCondition,
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
	case stderrors.Is(err, ErrInvalidUncmpPubKey):
		return InvalidUncmpPubKey
	case stderrors.Is(err, ErrValidatorNotFound):
		return ValidatorNotFound
	case stderrors.Is(err, ErrValidatorAlreadyExists):
		return ValidatorAlreadyExists
	case stderrors.Is(err, ErrInvalidTokenType):
		return InvalidTokenType
	case stderrors.Is(err, ErrInvalidPeriodType):
		return InvalidPeriodType
	case stderrors.Is(err, ErrInvalidOperator):
		return InvalidOperator
	case stderrors.Is(err, ErrInvalidCommissionRate):
		return InvalidCommissionRate
	case stderrors.Is(err, ErrInvalidMinSelfDelegation):
		return InvalidMinSelfDelegation
	case stderrors.Is(err, ErrInvalidDelegationAmount):
		return InvalidDelegationAmount
	case stderrors.Is(err, ErrPeriodDelegationNotFound):
		return PeriodDelegationNotFound
	case stderrors.Is(err, ErrInvalidRequest):
		return InvalidRequest
	case stderrors.Is(err, ErrUnexpectedCondition):
		return UnexpectedCondition
	default:
		return Unspecified
	}
}
