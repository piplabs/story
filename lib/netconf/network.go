package netconf

import (
	"slices"

	"github.com/piplabs/story/lib/errors"
)

// ID is a network identifier.
type ID string

// Static returns the static config and data for the network.
func (i ID) Static() Static {
	return statics[i]
}

func (i ID) Verify() error {
	if !supported[i] {
		return errors.New("unsupported network", "network", i)
	}

	return nil
}

func (i ID) String() string {
	return string(i)
}

func (i ID) Version() string {
	return i.Static().Version
}

const (
	// Iliad is the official Story Protocol public testnet.
	Iliad ID = "iliad"

	// Local is for local network testing.
	Local ID = "local"

	// Odyssey is the official Story Protocol public testnet.
	Odyssey ID = "odyssey"

	// Aeneid is the official Story Protocol public testnet.
	Aeneid ID = "aeneid"

	// Story is the official Story Protocol mainnet.
	Story ID = "story"
)

// supported is a map of supported networks.
//
//nolint:gochecknoglobals // Global state here is fine.
var supported = map[ID]bool{
	Iliad:   true,
	Odyssey: true,
	Aeneid:  true,
	Local:   true,
	Story:   true,
}

// IsAny returns true if the `ID` matches any of the provided targets.
func IsAny(id ID, targets ...ID) bool {
	return slices.Contains(targets, id)
}

// All returns all the supported network IDs.
func All() []ID {
	var resp []ID

	for id, ok := range supported {
		if ok {
			resp = append(resp, id)
		}
	}

	return resp
}
