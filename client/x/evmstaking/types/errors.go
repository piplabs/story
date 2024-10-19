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

func WrapErrWithCode(_ ErrCode, _ error) error {
	return nil
	/*
		if _, ok := codeToErr[code]; !ok {
			code = Unknown
		}

		return errors.Join(codeToErr[code], err)
	*/
}

func UnwrapErrCode(err error) ErrCode {
	switch {
	case errors.Is(err, ErrUnknown):
		return Unknown
	default:
		return Unknown
	}
}
