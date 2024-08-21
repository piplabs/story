package static

import _ "embed"

// original: anvil-state.json
//
//go:embed anvil-state.json
var anvilState []byte

func GetDevnetAnvilState() []byte {
	return anvilState
}
