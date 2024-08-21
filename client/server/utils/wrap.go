package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type SimpleWrapFunc func(r *http.Request) (resp any, err error)
type AutoSimpleInterfaceWrapFunc[T any] func(typ *T, r *http.Request) (resp any, err error)
type RespFunc func(w http.ResponseWriter, r *http.Request)

type Response struct {
	Code  int             `json:"code"`
	Msg   json.RawMessage `json:"msg"`
	Error string          `json:"error"`
}

// SimpleWrap implements a response encoder for requests without any query paramenter.
//
//nolint:nestif // Complicated scenarios.
func SimpleWrap(codec *codec.LegacyAmino, api SimpleWrapFunc) RespFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := api(r)
		resp := &Response{}
		if err != nil {
			resp.Code = 500
			resp.Error = err.Error()
		} else {
			var msgByte []byte
			if codec != nil {
				msgByte, err = codec.MarshalJSON(msg)
			} else {
				msgByte, err = json.Marshal(msg)
			}

			if err == nil {
				resp.Code = 200
				resp.Msg = msgByte
			} else {
				resp.Code = 502
				resp.Error = err.Error()
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.Code)
		_ = json.NewEncoder(w).Encode(resp)
	}
}

// AutoWrap implements a response encoder for requests with query paramenters.
func AutoWrap[T any](codec *codec.LegacyAmino, api AutoSimpleInterfaceWrapFunc[T]) RespFunc {
	typ := reflect.TypeOf(new(T)).Elem()
	return SimpleWrap(codec, func(r *http.Request) (resp any, err error) {
		val := reflect.New(typ).Interface()

		if r.URL.RawQuery != "" {
			err = QueryMapToVal(r.URL.Query(), val)
			if err != nil {
				return nil, NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("decode `%s` query err: %v", typ.String(), err))
			}
		}

		if r.ContentLength > 0 {
			err = json.NewDecoder(r.Body).Decode(val)
			if err != nil {
				return nil, NewHTTPError(http.StatusUnprocessableEntity, fmt.Sprintf("decode `%s` body err: %v", typ.String(), err))
			}
		}

		err = validate.Struct(val)
		if err != nil {
			return nil, WrapHTTPErrorWithCode(http.StatusUnprocessableEntity, err)
		}

		valT, ok := val.(*T)
		if !ok {
			return nil, WrapHTTPErrorWithCode(http.StatusUnprocessableEntity, errors.New("type assertion failed"))
		}

		return api(valT, r)
	})
}
