package types

import (
	"errors"
)

type ErrCode uint32

const (
	Internal               ErrCode = 0
	InvalidUncmpPubKey     ErrCode = 1
	ValidatorNotFound      ErrCode = 2
	ValidatorAlreadyExists ErrCode = 3
	InvalidTokenType       ErrCode = 4
	InvalidPeriodType      ErrCode = 5
	InvalidOperator        ErrCode = 6
)

var (
	ErrInternal               = errors.New("internal_error")
	ErrInvalidUncmpPubKey     = errors.New("invalid_uncompressed_pubkey")
	ErrValidatorNotFound      = errors.New("validator_not_found")
	ErrValidatorAlreadyExists = errors.New("validator_already_exists")
	ErrInvalidTokenType       = errors.New("invalid_token_type")
	ErrInvalidPeriodType      = errors.New("invalid_period_type")
	ErrInvalidOperator        = errors.New("invalid_operator")
)

var codeToErr = map[ErrCode]error{
	Internal:               ErrInternal,
	InvalidUncmpPubKey:     ErrInvalidUncmpPubKey,
	ValidatorNotFound:      ErrValidatorNotFound,
	ValidatorAlreadyExists: ErrValidatorAlreadyExists,
	InvalidTokenType:       ErrInvalidTokenType,
	InvalidPeriodType:      ErrInvalidPeriodType,
	InvalidOperator:        ErrInvalidOperator,
}

func (c ErrCode) String() string {
	if _, ok := codeToErr[c]; !ok {
		return ErrInternal.Error()
	}

	return codeToErr[c].Error()
}

//nolint:wrapcheck // we are wrapping the error with the code
func WrapErrWithCode(code ErrCode, err error) error {
	if _, ok := codeToErr[code]; !ok {
		code = Internal
	}

	return errors.Join(codeToErr[code], err)
}

func UnwrapErrCode(err error) ErrCode {
	switch {
	case errors.Is(err, ErrInternal):
		return Internal
	default:
		return Internal
	}
}
