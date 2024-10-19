package types

import (
	"errors"
)

type ErrCode uint32

const (
	Unknown ErrCode = 0
)

var (
	ErrUnknown = errors.New("unknown")
)

var codeToErr = map[ErrCode]error{
	Unknown: ErrUnknown,
}

func (c ErrCode) String() string {
	if _, ok := codeToErr[c]; !ok {
		return ErrUnknown.Error()
	}

	return codeToErr[c].Error()
}

//nolint:wrapcheck // we are wrapping the error with the code
func WrapErrWithCode(code ErrCode, err error) error {
	if _, ok := codeToErr[code]; !ok {
		code = Unknown
	}

	return errors.Join(codeToErr[code], err)
}

func UnwrapErrCode(err error) ErrCode {
	switch {
	case errors.Is(err, ErrUnknown):
		return Unknown
	default:
		return Unknown
	}
}
