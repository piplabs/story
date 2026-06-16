//nolint:wrapcheck // Internal utils, don't need to wrap it.
package utils

import (
	"encoding/base64"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/piplabs/story/lib/errors"
)

var base64BytesType = reflect.TypeOf(Base64Bytes(nil))

// QueryMapToVal implements an all-in-one decoder to decode requests' query parameters to
// structured value.
func QueryMapToVal(query url.Values, val any) error {
	// Reject query keys nested deeper than maxQueryDepth instead of silently dropping them.
	for k := range query {
		if depth := strings.Count(k, "."); depth > maxQueryDepth {
			return fmt.Errorf("query key %q exceeds max nesting depth %d", k, maxQueryDepth)
		}
	}

	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		Result:           val,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			stringArrayToBase64Bytes(),
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			stringArrayToNative(),
			stringArrayToNativePtr(),
		),
	})
	if err != nil {
		return err
	}

	return decoder.Decode(buildMap(query))
}

// stringArrayToBase64Bytes decodes a single query-string value into Base64Bytes
// using standard base64 (matching how amino MarshalJSON serializes []byte in
// responses, e.g. pagination next_key). This lets callers round-trip the
// next_key field without manual decoding.
func stringArrayToBase64Bytes() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if t != base64BytesType {
			return data, nil
		}

		as, ok := data.([]string)
		if !ok || len(as) == 0 || as[0] == "" {
			return Base64Bytes(nil), nil
		}

		decoded, err := base64.StdEncoding.DecodeString(as[0])
		if err != nil {
			return nil, errors.Wrap(err, "invalid base64 value", "target_type", t.String())
		}

		return Base64Bytes(decoded), nil
	}
}

func stringArrayToNativePtr() mapstructure.DecodeHookFunc {
	return func(f reflect.Type, t reflect.Type, data any) (any, error) {
		if f.Kind() != reflect.Slice || t.Kind() != reflect.Ptr {
			return data, nil
		}

		as, ok := data.([]string)
		if !ok || len(as) == 0 {
			return data, nil
		}

		from := as[0]

		t = t.Elem()
		switch t.String() {
		case "big.Int":
			if i, ok := big.NewInt(0).SetString(from, 10); ok {
				return i, nil
			}
		default:
			return data, nil
		}

		return data, nil
	}
}

func stringArrayToNative() mapstructure.DecodeHookFunc {
	return func(f reflect.Kind, t reflect.Kind, data any) (any, error) {
		if f != reflect.Slice {
			return data, nil
		}

		as, ok := data.([]string)
		if !ok || len(as) == 0 {
			return data, nil
		}

		from := as[0]

		switch t {
		case reflect.Bool:
			return strconv.ParseBool(from)
		case reflect.Uint64:
			return strconv.ParseUint(from, 10, 64)
		case reflect.Int:
			return strconv.ParseInt(from, 10, 64)
		case reflect.Uint:
			parseUint, err := strconv.ParseUint(from, 10, 64)
			if parseUint > math.MaxUint {
				return 0, strconv.ErrRange
			}

			return uint(parseUint), err
		case reflect.Float32:
			return strconv.ParseFloat(from, 32)
		case reflect.Float64:
			return strconv.ParseFloat(from, 64)
		case reflect.String:
			return from, nil
		default:
			return data, nil
		}
	}
}

// maxQueryDepth limits the nesting depth of dot-separated query parameter keys
// to prevent O(D^2) CPU and memory consumption from deeply nested keys.
// Current max actual usage is depth 2 (e.g. "pagination.limit").
const maxQueryDepth = 5

func buildMap(query url.Values, prefix ...string) (ret map[string]any) {
	fullPrefix := strings.Join(prefix, ".")
	if len(fullPrefix) > 0 {
		fullPrefix += "."
	}

	ret = make(map[string]any)
	prefixMapped := make(map[string]bool)

	for k, vs := range query {
		if rk := strings.TrimPrefix(k, fullPrefix); fullPrefix == "" || rk != k {
			rk = strings.TrimSuffix(rk, "[]")
			if prefixIndex := strings.IndexByte(rk, '.'); prefixIndex != -1 {
				rk = rk[:prefixIndex]
				if !prefixMapped[rk] {
					prefixMapped[rk] = true
					ret[rk] = buildMap(query, append(slices.Clone(prefix), rk)...)
				}
			} else {
				ret[rk] = vs
			}
		}
	}

	return ret
}
