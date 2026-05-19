package utils

// Base64Bytes is a []byte that is base64-decoded from a string when populated
// from HTTP query parameters by QueryMapToVal. It exists so handlers can keep
// taking raw bytes for fields like cosmos-sdk PageRequest.Key, while the wire
// representation (set by amino MarshalJSON on []byte) stays as standard
// base64-encoded strings. Without this round-trip, page 2 requests submit the
// undecoded base64 ASCII as the iterator start key and silently return zero
// results.
type Base64Bytes []byte
